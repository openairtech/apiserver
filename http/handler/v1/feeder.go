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
	"database/sql"
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
			em := fmt.Sprintf("invalid request: %v", err)
			writeResult(w, api.StatusBadRequest, em)
			return
		}

		s, err := db.StationByTokenId(f.TokenId)
		if err != nil {
			em := fmt.Sprintf("can't get station by token id [%s]: %v", f.TokenId, err)
			writeResult(w, api.StatusBadRequest, em)
			log.Error(em)
			return
		}

		ms := stationDbMeasurements(s, f)
		if err := db.AddMeasurements(s, ms); err != nil {
			em := fmt.Sprintf("station [%d]: can't add %d measurement(s): %v", s.Id, len(ms), err)
			writeResult(w, api.StatusServerError, em)
			return
		}

		m := fmt.Sprintf("station [%d]: added %d measurement(s)", s.Id, len(ms))
		if len(ms) > 1 {
			log.Info(m)
		} else {
			log.Debug(m)
		}

		// Update station data
		su := s.Copy()
		seen := time.Now()
		su.Seen = &seen
		su.Version = sql.NullString{String: f.Version, Valid: len(f.Version) > 0}
		if err := db.UpdateStation(s, &su); err != nil {
			em := fmt.Sprintf("station [%d]: can't update station data: %v", s.Id, err)
			writeResult(w, api.StatusServerError, em)
			return
		}

		writeResult(w, api.StatusOk, "")
	})
}

func stationDbMeasurements(station *db.Station, f api.FeederData) []db.Measurement {
	var ms []db.Measurement

	for i, am := range f.Measurements {
		// Check for timestamp presence in measurement
		if am.Timestamp == nil {
			// Skip all measurements without timestamp except last
			if i < len(f.Measurements)-1 {
				log.Warnf("skipped measurement without timestamp: %+v", am)
				continue
			}
			// Set timestamp to the last measurement
			now := api.UnixTime(time.Now())
			am.Timestamp = &now
		}

		// Use provided AQI value or compute it from PM values
		if am.Aqi == nil && am.Pm10 != nil && am.Pm25 != nil {
			pm := aqi.PM{Pm25: *am.Pm25, Pm10: *am.Pm10}
			ac := pm.Aqi()
			am.Aqi = &ac
		}

		ms = append(ms, db.NewMeasurement(station, am))
	}

	return ms
}
