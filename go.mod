module github.com/openairtech/apiserver

replace github.com/openairtech/api v0.0.0 => ../openair-api

require (
	github.com/cridenour/go-postgis v1.0.0
	github.com/doug-martin/goqu/v7 v7.3.1
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.1
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.1.1
	github.com/openairtech/api v0.0.0
	github.com/sirupsen/logrus v1.4.1
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3 // indirect
)

go 1.13
