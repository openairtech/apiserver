// Copyright © 2019 Victor Antonovich <victor@antonovich.me>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"fmt"
	"time"

	"github.com/openairtech/apiserver/util"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Db struct {
	sqlx *sqlx.DB
}

func NewDb(host string, port int, user, password, name string, maxConn int) (*Db, error) {
	var db *sqlx.DB
	var err error

	if db, err = sqlx.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s sslmode=disable", host, port, user, password, name)); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxConn)

	return &Db{
		sqlx: db,
	}, nil
}

func (db *Db) Close() {
	_ = db.sqlx.Close()
}

func (db *Db) StationByTokenId(tokenId string) (*Station, error) {
	s := Station{}
	if err := db.sqlx.Get(&s, "SELECT * FROM Stations WHERE token_id = $1", tokenId); err != nil {
		return nil, err
	}
	return &s, nil
}
func (db *Db) Stations(bbox []float64, mlast time.Duration) ([]Station, error) {
	var s []Station

	query := `SELECT DISTINCT ON (s.id) s.*, m.id "m.id", m.tstamp "m.tstamp", m.temperature "m.temperature", 
                          m.pressure "m.pressure", m.humidity "m.humidity", m.pm25 "m.pm25", m.pm10 "m.pm10", 
                          m.aqi "m.aqi"
		FROM Stations s LEFT JOIN Measurements m ON s.id = m.station_id `

	if mlast > 0 {
		query += fmt.Sprintf("AND m.tstamp > NOW() - INTERVAL '%d SECONDS' ", int(mlast.Seconds()))
	}

	if len(bbox) == 4 {
		query += fmt.Sprintf("WHERE s.location @ ST_MakeEnvelope(%f, %f, %f, %f) ",
			bbox[0], bbox[1], bbox[2], bbox[3])
	}

	query += "ORDER BY s.id, m.tstamp DESC"

	if err := db.sqlx.Select(&s, query); err != nil {
		return nil, err
	}

	return s, nil
}

func (db *Db) AddMeasurement(station *Station, timestamp time.Time,
	temperature, humidity, pressure, pm25, pm10 *float32, aqi *int) (*Measurement, error) {

	m := Measurement{
		StationId:   toNullInt64(&station.Id),
		Temperature: toNullFloat64(temperature),
		Humidity:    toNullFloat64(humidity),
		Pressure:    toNullFloat64(pressure),
		Pm25:        toNullFloat64(pm25),
		Pm10:        toNullFloat64(pm10),
		Aqi:         toNullInt64(aqi),
		Timestamp:   &timestamp,
	}

	// TODO Add `ON CONFLICT` clause and its handling logic
	query := `INSERT INTO Measurements(station_id, tstamp, temperature, humidity, pressure, pm25, pm10, aqi)
		VALUES (:station_id, :tstamp, :temperature, :humidity, :pressure, :pm25, :pm10, :aqi) 
		RETURNING id`
	rows, err := db.sqlx.NamedQuery(query, m)

	if err != nil {
		return nil, err
	}

	defer util.CloseQuietly(rows)

	if rows.Next() {
		if err := rows.Scan(&m.Id); err != nil {
			return nil, err
		}
	} else {
		fmt.Printf("No rows!\n")
	}

	return &m, nil
}
