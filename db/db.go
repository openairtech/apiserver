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
	"errors"
	"fmt"
	"time"

	gq "github.com/doug-martin/goqu/v7"
	_ "github.com/doug-martin/goqu/v7/dialect/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/openairtech/apiserver/util"
)

var d = gq.Dialect("postgres")

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

// StationByTokenId finds station by its tokenId.
// It returns reference to Station struct or nil if no station with given token ID was found,
// or error if something went wrong.
func (db *Db) StationByTokenId(tokenId string) (*Station, error) {
	s := Station{}
	if err := db.sqlx.Get(&s, "SELECT * FROM stations WHERE token_id = $1", tokenId); err != nil {
		return nil, err
	}
	return &s, nil
}

// Stations gets slice of stations with their last measurements according to given parameters.
// bbox, if not empty, defines a bounding box [min_long, min_lat, max_long, max_lat] to get stations within it.
// mfrom specifies the upper time limit for last measurement to include in result, now() if nil.
// mlast, if not nil, defines the lower time limit for last measurement to include in result,
// and it is computed as (mfrom - mlast).
func (db *Db) Stations(bbox []float64, mfrom *time.Time, mlast *time.Duration) ([]Station, error) {
	var s []Station

	lj := []gq.Expression{gq.I("s.id").Eq(gq.I("m.station_id"))}

	if mfrom != nil {
		lj = append(lj, gq.L("m.tstamp <= ?", mfrom))
	}

	if mlast != nil {
		if mfrom != nil {
			lj = append(lj, gq.L("m.tstamp > ?::TIMESTAMP - ? * INTERVAL '1 SECOND'", mfrom, int(mlast.Seconds())))
		} else {
			lj = append(lj, gq.L("m.tstamp > NOW() - ? * INTERVAL '1 SECOND'", int(mlast.Seconds())))
		}
	}

	var w []gq.Expression

	if len(bbox) == 4 {
		w = append(w, gq.L("s.location @ ST_MakeEnvelope(?, ?, ?, ?)",
			bbox[0], bbox[1], bbox[2], bbox[3]))
	}

	q := d.From(gq.T("stations").As("s")).
		Select(gq.L(`DISTINCT ON (s.id) s.*, m.id "m.id", m.tstamp "m.tstamp", m.temperature "m.temperature", 
			m.pressure "m.pressure", m.humidity "m.humidity", m.pm25 "m.pm25", m.pm10 "m.pm10", m.aqi "m.aqi"`)).
		LeftJoin(gq.T("measurements").As("m"), gq.On(lj...))

	if len(w) > 0 {
		q = q.Where(w...)
	}

	q = q.Order(gq.I("s.id").Asc(), gq.I("m.tstamp").Desc())

	query, args, err := q.Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}

	if err := db.sqlx.Select(&s, query, args...); err != nil {
		return nil, err
	}

	return s, nil
}

// UpdateStation updates station s data by the differences found while comparing it with updated data su
func (db *Db) UpdateStation(s, su *Station) error {
	if s == su {
		// No fields to update
		return nil
	}

	r := make(gq.Record)
	if s.Id != su.Id {
		return errors.New(fmt.Sprintf("station id %d change to %d is not allowed", s.Id, su.Id))
	}
	if s.TokenId != su.TokenId {
		r["token_id"] = su.TokenId
	}
	if s.Description != su.Description {
		r["description"] = su.Description
	}
	if s.Version != su.Version {
		r["version"] = su.Version
	}
	if s.Created != su.Created {
		r["created"] = su.Created
	}
	if s.Seen != su.Seen {
		r["seen"] = su.Seen
	}
	if s.Location != su.Location {
		r["location"] = su.Location
	}

	if len(r) == 0 {
		return errors.New(fmt.Sprintf("station objects are different "+
			"but no differences found:\ninitial: %+v,\nupdated: %+v", s, su))
	}

	query, args, err := d.From("stations").Prepared(true).Where(gq.C("id").Eq(s.Id)).ToUpdateSQL(r)
	if err != nil {
		return err
	}

	_, err = db.sqlx.Exec(query, args...)

	return err
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
	query := `INSERT INTO measurements(station_id, tstamp, temperature, humidity, pressure, pm25, pm10, aqi)
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

// Measurements gets slice of station measurements sorted by timestamp according to given time interval.
// stationId is identifier of station to get measurements.
// timeFrom specifies the start time of interval to get measurements.
// timeEnd specifies the end time of interval to get measurements.
func (db *Db) Measurements(stationId int, timeFrom time.Time, timeTo time.Time) ([]Measurement, error) {
	var m []Measurement
	if timeFrom.After(timeTo) {
		timeFrom, timeTo = timeTo, timeFrom
	}

	q := d.From("measurements")
	q = q.Where(gq.C("station_id").Eq(stationId))
	q = q.Where(gq.C("tstamp").Between(gq.Range(timeFrom, timeTo)))
	q = q.Order(gq.I("tstamp").Asc())

	query, args, err := q.Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}

	if err := db.sqlx.Select(&m, query, args...); err != nil {
		return nil, err
	}

	return m, nil
}
