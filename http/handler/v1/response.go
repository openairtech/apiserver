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
	"net/http"

	"github.com/openairtech/api"
)

func writeResult(w http.ResponseWriter, sc api.StatusCode, m string) {
	w.Header().Set("Content-Type", "application/json")
	r := api.Result{
		Status:  sc,
		Message: m,
	}
	_ = json.NewEncoder(w).Encode(r)
}
