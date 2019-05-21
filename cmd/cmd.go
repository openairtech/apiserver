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

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	dbpkg "github.com/openairtech/apiserver/db"
	"github.com/openairtech/apiserver/http"
)

const (
	FlagVersion         = "version"
	FlagGracefulTimeout = "graceful-timeout"
	FlagDebug           = "debug"

	FlagDbHost     = "db-host"
	FlagDbPort     = "db-port"
	FlagDbUser     = "db-user"
	FlagDbPassword = "db-pass"
	FlagDbName     = "db-name"
	FlagDbMaxConn  = "db-max-conn"

	FlagHttpHost = "http-host"
	FlagHttpPort = "http-port"
)

var (
	BuildVersion   = "unknown"
	BuildTimestamp = "unknown"
)

var (
	debug                              bool
	gracefulTimeout                    time.Duration
	dbHost, dbUser, dbPassword, dbName string
	httpHost                           string
	dbPort, dbMaxConn, httpPort        int
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "openair-apiserver",
		Long: "OpenAir API server.",
		Run:  runCmd,
	}
	initCmd(cmd)
	return cmd
}

func initCmd(cmd *cobra.Command) {
	// Command-related flags set
	f := cmd.Flags()
	f.BoolP(FlagVersion, "V", false, "display the build number and timestamp")
	f.DurationVarP(&gracefulTimeout, FlagGracefulTimeout, "T", time.Second*15, "graceful shutdown timeout")
	f.BoolVarP(&debug, FlagDebug, "d", false, "enable debug logging")

	f.StringVarP(&dbHost, FlagDbHost, "H", "localhost", "database server host")
	f.IntVarP(&dbPort, FlagDbPort, "P", 5432, "database server port")
	f.StringVarP(&dbUser, FlagDbUser, "U", "openair", "database user name")
	f.StringVarP(&dbPassword, FlagDbPassword, "W", "openair", "database user password")
	f.StringVarP(&dbName, FlagDbName, "D", "openair", "database name to connect to")
	f.IntVarP(&dbMaxConn, FlagDbMaxConn, "M", 0, "database maximum number of open connections (0 for unlimited)")

	f.StringVarP(&httpHost, FlagHttpHost, "s", "localhost", "HTTP server host")
	f.IntVarP(&httpPort, FlagHttpPort, "p", 8081, "HTTP server port")
}

func runCmd(cmd *cobra.Command, _ []string) {
	if f, _ := cmd.Flags().GetBool(FlagVersion); f {
		fmt.Printf("Build version: %s\n", BuildVersion)
		fmt.Printf("Build timestamp: %s\n", BuildTimestamp)
		return
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Debug("connecting to database...")
	db, err := dbpkg.NewDb(dbHost, dbPort, dbUser, dbPassword, dbName, dbMaxConn)
	if err != nil {
		log.Errorf("can't connect to database: %v", err)
		return
	}
	defer db.Close()

	s := http.NewServer(fmt.Sprintf("%s:%d", httpHost, httpPort), db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Info("starting server")
		if err := s.Run(); err != nil {
			log.Errorf("server error: %v", err)
		}
		cancel()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	select {
	case <-c:
		log.Info("stopping server")
		ctx, _ := context.WithTimeout(ctx, gracefulTimeout)
		if err := s.Shutdown(ctx); err != nil {
			log.Errorf("can't shutdown server: %v", err)
		}
	case <-ctx.Done():
		break
	}

	log.Info("server stopped")
}
