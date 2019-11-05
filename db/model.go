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
	"database/sql"
	"fmt"
	"time"

	"github.com/cridenour/go-postgis"

	"github.com/openairtech/api"
)

type Station struct {
	Id          int
	TokenId     string `db:"token_id"`
	Description sql.NullString
	Version     sql.NullString
	Created     time.Time
	Seen        *time.Time
	Location    postgis.PointS
	Measurement `db:"m"`
}

func (s Station) Copy() Station {
	c := Station{}
	c = s
	return c
}

func (s Station) ApiStation() api.Station {
	sId := s.Id
	var m *api.Measurement = nil
	if s.Measurement.Id.Valid {
		am := s.Measurement.ApiMeasurement()
		m = &am
	}
	as := api.Station{
		Id:              &sId,
		Created:         api.UnixTime(s.Created),
		Description:     s.Description.String,
		Longitude:       s.Location.X,
		Latitude:        s.Location.Y,
		LastMeasurement: m,
	}
	if s.Seen != nil {
		ass := api.UnixTime(*s.Seen)
		as.Seen = &ass
	}
	return as
}

type Measurement struct {
	Id          sql.NullInt64
	StationId   sql.NullInt64 `db:"station_id"`
	Timestamp   *time.Time    `db:"tstamp"`
	Temperature sql.NullFloat64
	Humidity    sql.NullFloat64
	Pressure    sql.NullFloat64
	Pm25        sql.NullFloat64
	Pm10        sql.NullFloat64
	Aqi         sql.NullInt64
}

func (m Measurement) ApiMeasurement() api.Measurement {
	var ts *api.UnixTime
	if m.Timestamp != nil {
		ats := api.UnixTime(*m.Timestamp)
		ts = &ats
	}
	return api.Measurement{
		Timestamp:   ts,
		Temperature: fromNullFloat64(m.Temperature),
		Humidity:    fromNullFloat64(m.Humidity),
		Pressure:    fromNullFloat64(m.Pressure),
		Pm25:        fromNullFloat64(m.Pm25),
		Pm10:        fromNullFloat64(m.Pm10),
		Aqi:         fromNullInt64(m.Aqi),
	}
}

func MeasurementDbColumns(amv []string) ([]interface{}, error) {
	s := make(map[interface{}]struct{})
	// FIXME Use reflection with API/DB measurement structs
	for _, v := range amv {
		switch v {
		case "timestamp":
			s["tstamp"] = struct{}{}
		case "temperature":
			s["temperature"] = struct{}{}
		case "humidity":
			s["humidity"] = struct{}{}
		case "pressure":
			s["pressure"] = struct{}{}
		case "pm25":
			s["pm25"] = struct{}{}
		case "pm10":
			s["pm10"] = struct{}{}
		case "aqi":
			s["aqi"] = struct{}{}
		default:
			return nil, fmt.Errorf("unknown variable: %s", v)
		}
	}
	c := make([]interface{}, 0, len(s))
	for v := range s {
		c = append(c, v)
	}
	return c, nil
}
