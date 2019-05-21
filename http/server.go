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

package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/openairtech/apiserver/db"
	v1 "github.com/openairtech/apiserver/http/handler/v1"
)

type Server struct {
	http *http.Server
}

func NewServer(addr string, db *db.Db) *Server {
	var router = mux.NewRouter()

	var v1Api = router.PathPrefix("/v1").Subrouter()
	v1Api.Handle("/feeder", v1.FeederHandler(db)).Methods("POST")

	s := &Server{
		http: &http.Server{
			Addr:         addr,
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      router,
		},
	}

	return s
}

func (s *Server) Run() error {
	if err := s.http.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
