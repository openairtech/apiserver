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
	"fmt"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/openairtech/api"
	"github.com/openairtech/apiserver/db"
	httputil "github.com/openairtech/apiserver/http/util"
	"github.com/openairtech/apiserver/util"
)

func MeasurementsGetHandler(db *db.Db) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ss := r.URL.Query().Get("station")
		if ss == "" {
			writeResult(w, api.StatusBadRequest, fmt.Sprint("'station' parameter not set"))
			return
		}
		s, err := strconv.ParseInt(ss, 10, 32)
		if err != nil {
			writeResult(w, api.StatusBadRequest, fmt.Sprintf("can't parse station id: %v", err))
			return
		}

		from, err := util.ParseUnixTime(r.URL.Query().Get("from"))
		if err != nil {
			writeResult(w, api.StatusBadRequest, fmt.Sprint(err))
			return
		}
		if from == nil {
			writeResult(w, api.StatusBadRequest, fmt.Sprint("'from' parameter not set"))
			return
		}

		to, err := util.ParseUnixTime(r.URL.Query().Get("to"))
		if err != nil {
			writeResult(w, api.StatusBadRequest, fmt.Sprint(err))
			return
		}
		if to == nil {
			writeResult(w, api.StatusBadRequest, fmt.Sprint("'to' parameter not set"))
			return
		}

		var vars []string
		v := r.URL.Query().Get("v")
		if v != "" {
			vars = strings.Split(v, ",")
		}

		dms, err := db.Measurements(int(s), *from, *to, vars)
		if err != nil {
			m := fmt.Sprintf("can't get measurements: %v", err)
			writeResult(w, api.StatusServerError, m)
			log.Error(m)
			return
		}

		var ms []api.Measurement
		for _, dm := range dms {
			ms = append(ms, dm.ApiMeasurement())
		}

		httputil.WriteJsonResponse(w, api.MeasurementsResult{
			Result:       api.Result{Status: api.StatusOk},
			Measurements: ms,
		})
	})
}
