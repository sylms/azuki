package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"github.com/sylms/azuki/infrastructure/persistence"
	"github.com/sylms/azuki/interface/handler"
	"github.com/sylms/azuki/usecase"
)

const (
	envSylmsPostgresDBKey       = "SYLMS_POSTGRES_DB"
	envSylmsPostgresUserKey     = "SYLMS_POSTGRES_USER"
	envSylmsPostgresPasswordKey = "SYLMS_POSTGRES_PASSWORD"
	envSylmsPostgresHostKey     = "SYLMS_POSTGRES_HOST"
	envSylmsPostgresPortKey     = "SYLMS_POSTGRES_PORT"
	envSylmsPort                = "SYLMS_PORT"
	envSecretKey                = "SECRET_KEY"
)

func main() {
	envKeys := []string{envSylmsPostgresDBKey, envSylmsPostgresUserKey, envSylmsPostgresPasswordKey, envSylmsPostgresHostKey, envSylmsPostgresPortKey, envSylmsPort, envSecretKey}
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

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", postgresHost, postgresPort, postgresUser, postgresPassword, postgresDb))
	if err != nil {
		log.Fatalf("%+v", err)
	}

	persistence := persistence.NewCoursePersistence(db)
	useCase := usecase.NewCourseUseCase(persistence)
	handler := handler.NewCourseHandler(useCase, os.Getenv(envSecretKey))

	r := mux.NewRouter()
	r.HandleFunc("/course", handler.Search).Methods("POST")
	r.HandleFunc("/facet", handler.Facet).Methods("POST")
	r.HandleFunc("/csv", handler.Csv).Methods("POST")
	r.HandleFunc("/update", handler.Update).Methods("POST")
	c := cors.Default().Handler(r)
	log.Printf("Listen Port: %s", portStr)
	err = http.ListenAndServe(fmt.Sprintf(":%s", portStr), c)
	if err != nil {
		log.Fatalln(err)
	}
}
