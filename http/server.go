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

	"github.com/gorilla/handlers"
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

	v1Api.NotFoundHandler = http.HandlerFunc(v1.ErrorNotFoundHandler)
	v1Api.MethodNotAllowedHandler = http.HandlerFunc(v1.ErrorMethodNotAllowedHandler)

	v1Api.Handle("/feeder", v1.FeederHandler(db)).Methods("POST")

	sgh := v1.StationsGetHandler(db)
	v1Api.Handle("/stations", sgh).Methods("GET")
	//v1Api.Handle("/stations", sgh).
	//	Queries("offset", "{offset:[0-9]+}", "limit", "{limit:[0-9]+}").
	//	Methods("GET")

	originsOk := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	s := &Server{
		http: &http.Server{
			Addr:         addr,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
			IdleTimeout:  60 * time.Second,
			Handler:      handlers.CORS(originsOk, headersOk, methodsOk)(router),
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
