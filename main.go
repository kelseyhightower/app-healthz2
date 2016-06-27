package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/app-healthz2/healthz"
)

var version = "2.0.0"

func main() {
	log.Println("Starting app...")

	httpAddr := os.Getenv("HTTP_ADDR")
	databaseHost := os.Getenv("DATABASE_HOST")
	vaultAddr := os.Getenv("VAULT_ADDR")
	vaultToken := os.Getenv("VAULT_TOKEN")

	vc, err := newVaultClient(vaultAddr, vaultToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Getting database credentials...")
	username, password, err := vc.getDatabaseCredentials("mysql/creds/app")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initializing database connection pool...")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/app", username, password, databaseHost)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	hc := &healthz.Config{
		Hostname: hostname,
		Database: healthz.DatabaseConfig{
			DataSourceName: dsn,
		},
		Vault: healthz.VaultConfig{
			Address: vaultAddr,
		},
	}

	healthzHandler, err := healthz.Handler(hc)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/healthz", healthzHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, html, hostname, version)
	})

	log.Printf("HTTP service listening on %s", httpAddr)
	http.ListenAndServe(httpAddr, nil)
}
