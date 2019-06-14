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

package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/openairtech/api"
	"github.com/openairtech/apiserver/aqi"
	"github.com/openairtech/apiserver/db"
)

func FeederHandler(db *db.Db) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var f api.FeederData
		err := decoder.Decode(&f)
		if err != nil {
			m := fmt.Sprintf("invalid request: %v", err)
			writeResult(w, api.StatusBadRequest, m)
			return
		}

		s, err := db.StationByTokenId(f.TokenId)
		if err != nil {
			m := fmt.Sprintf("can't get station by token id [%s]: %v", f.TokenId, err)
			writeResult(w, api.StatusBadRequest, m)
			log.Error(m)
			return
		}

		for i, fm := range f.Measurements {
			// Check for timestamp presence in measurement
			if fm.Timestamp == nil {
				// Skip all measurements without timestamp except last
				if i < len(f.Measurements)-1 {
					log.Warnf("skipped measurement without timestamp: %+v", fm)
					continue
				}
				// Set timestamp to the last measurement
				now := api.UnixTime(time.Now())
				fm.Timestamp = &now
			}
			ts := time.Time(*fm.Timestamp)

			// Use provided AQI value or compute it from PM values
			if fm.Aqi == nil && fm.Pm10 != nil && fm.Pm25 != nil {
				pm := aqi.PM{Pm25: *fm.Pm25, Pm10: *fm.Pm10}
				ac := pm.Aqi()
				fm.Aqi = &ac
			}

			// Store measurement data to db
			m, err := db.AddMeasurement(s, ts, fm.Temperature, fm.Humidity, fm.Pressure,
				fm.Pm25, fm.Pm10, fm.Aqi)

			if err != nil {
				m := fmt.Sprintf("can't store measurement: %v", err)
				writeResult(w, api.StatusServerError, m)
				return
			}

			log.Debugf("added measurement: %+v", m)
		}

		writeResult(w, api.StatusOk, "")
	})
}
