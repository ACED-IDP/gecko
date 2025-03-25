package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ACED-IDP/gecko/gecko"
	"github.com/jmoiron/sqlx"
	"github.com/uc-cdis/go-authutils/authutils"
)

func main() {
	var jwkEndpointEnv string = os.Getenv("JWKS_ENDPOINT")

	var port *uint = flag.Uint("port", 8080, "port on which to expose the API")
	var jwkEndpoint *string = flag.String(
		"jwks",
		jwkEndpointEnv,
		"endpoint from which the application can fetch a JWKS",
	)

	if *jwkEndpoint == "" {
		print("WARNING: no $JWKS_ENDPOINT or --jwks specified; endpoints requiring JWT validation will error\n")
	}

	var dbUrl *string = flag.String(
		"db",
		"",
		"URL to connect to database: postgresql://user:password@netloc:port/dbname\n"+
			"can also be specified through the postgres\n"+
			"environment variables. If using the commandline argument, add\n"+
			"?sslmode=disable",
	)
	flag.Parse()

	db, err := sqlx.Open("postgres", *dbUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("DB ping failed: %v", err)
		panic(err)
	}
	defer db.Close()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	jwtApp := authutils.NewJWTApplication(*jwkEndpoint)
	log.Printf("JWT APP: %#v\n", jwtApp.Keys)
	log.Printf("LOGGER: %#v\n", logger)

	geckoServer, err := gecko.NewServer().
		WithLogger(logger).
		WithJWTApp(jwtApp).
		WithDB(db).
		Init()
	if err != nil {
		log.Fatalf("Failed to initialize gecko server: %v", err)
	}

	app := geckoServer.MakeRouter()

	// Configure Iris logger to output to your httpLogger
	httpLogger := log.New(os.Stdout, "", log.LstdFlags)
	app.Logger().SetOutput(httpLogger.Writer())

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     httpLogger,
		Handler:      app,
	}

	httpLogger.Println("gecko serving at", httpServer.Addr)
	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
