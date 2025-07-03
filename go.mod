module github.com/openairtech/apiserver

// replace github.com/openairtech/api v0.1.0 => ../openair-api

require (
	github.com/cridenour/go-postgis v1.0.1
	github.com/doug-martin/goqu/v7 v7.4.0
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.10.9
	github.com/openairtech/api v0.1.0
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.9.1
)

require (
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	golang.org/x/sys v0.33.0 // indirect
)

go 1.23.0

toolchain go1.24.4
