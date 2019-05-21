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
	"database/sql/driver"
	"time"

	"github.com/cridenour/go-postgis"
)

type Station struct {
	Id          int
	TokenId     string `db:"token_id"`
	Description string
	Created     time.Time
	Location    postgis.PointS
}

func (s Station) Value() (driver.Value, error) {
	return int64(s.Id), nil
}

type Measurement struct {
	Id          int
	Station     Station   `db:"station_id"`
	Timestamp   time.Time `db:"tstamp"`
	Temperature sql.NullFloat64
	Humidity    sql.NullFloat64
	Pressure    sql.NullFloat64
	Pm25        sql.NullFloat64
	Pm10        sql.NullFloat64
	Aqi         sql.NullInt64
}
