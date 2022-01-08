package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

var (
	db *sqlx.DB
)

const (
	envSylmsPostgresDBKey       = "SYLMS_POSTGRES_DB"
	envSylmsPostgresUserKey     = "SYLMS_POSTGRES_USER"
	envSylmsPostgresPasswordKey = "SYLMS_POSTGRES_PASSWORD"
	envSylmsPostgresHostKey     = "SYLMS_POSTGRES_HOST"
	envSylmsPostgresPortKey     = "SYLMS_POSTGRES_PORT"
	envSylmsPort                = "SYLMS_PORT"
)

const (
	searchQueryDefaultLimit = 50
)

func main() {
	envKeys := []string{envSylmsPostgresDBKey, envSylmsPostgresUserKey, envSylmsPostgresPasswordKey, envSylmsPostgresHostKey, envSylmsPostgresPortKey, envSylmsPort}
	for _, key := range envKeys {
		val, ok := os.LookupEnv(key)
		if !ok || val == "" {
			log.Fatalf("%s is not set or empty\n", key)
		}
	}

	postgresDb := os.Getenv(envSylmsPostgresDBKey)
	postgresUser := os.Getenv(envSylmsPostgresUserKey)
	postgresPassword := os.Getenv(envSylmsPostgresPasswordKey)
	postgresHost := os.Getenv(envSylmsPostgresHostKey)
	postgresPort := os.Getenv(envSylmsPostgresPortKey)
	portStr := os.Getenv(envSylmsPort)

	var err error
	db, err = sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", postgresHost, postgresPort, postgresUser, postgresPassword, postgresDb))
	if err != nil {
		log.Fatalf("%+v", err)
	}

	r := mux.NewRouter()
	// とりあえず科目名と授業概要で検索できるように
	// TODO: course_name や course_overview を指定しない検索方法に対応
	r.HandleFunc("/course", courseSimpleSearchHandler).Methods("POST")
	r.HandleFunc("/facet", courseFacetSearchHandler).Methods("POST")
	r.HandleFunc("/csv", courseCSVHandler).Methods("POST")
	c := cors.Default().Handler(r)
	log.Printf("Listen Port: %s", portStr)
	err = http.ListenAndServe(fmt.Sprintf(":%s", portStr), c)
	if err != nil {
		log.Fatalln(err)
	}
}
