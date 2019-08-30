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
	return api.Station{
		Id:              &sId,
		Created:         api.UnixTime(s.Created),
		Description:     s.Description.String,
		Longitude:       s.Location.X,
		Latitude:        s.Location.Y,
		LastMeasurement: m,
	}
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
