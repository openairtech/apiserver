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

	log "github.com/sirupsen/logrus"

	"github.com/openairtech/api"
	"github.com/openairtech/apiserver/db"
	httputil "github.com/openairtech/apiserver/http/util"
	"github.com/openairtech/apiserver/util"
)

func StationsGetHandler(db *db.Db) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bbox, err := util.ParseBBox(r.URL.Query().Get("bbox"))
		if err != nil {
			writeResult(w, api.StatusBadRequest, fmt.Sprint(err))
			return
		}

		mlast, err := util.ParseDuration(r.URL.Query().Get("mlast"))
		if err != nil {
			writeResult(w, api.StatusBadRequest, fmt.Sprint(err))
			return
		}

		dss, err := db.Stations(bbox, mlast)
		if err != nil {
			m := fmt.Sprintf("can't get stations: %v", err)
			writeResult(w, api.StatusServerError, m)
			log.Error(m)
			return
		}

		var as []api.Station
		for _, ds := range dss {
			as = append(as, ds.ApiStation())
		}

		httputil.WriteJsonResponse(w, api.StationsResult{
			Result:   api.Result{Status: api.StatusOk},
			Stations: as,
		})
	})
}
