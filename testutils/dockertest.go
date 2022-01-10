package testutils

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
)

func CreateDB() (*sqlx.DB, error) {
	pwd, _ := os.Getwd()

	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}

	runOptions := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.7-alpine",
		Env: []string{
			"TZ=Asia/Tokyo",
			"POSTGRES_DB=courses",
			"POSTGRES_USER=sylms",
			"POSTGRES_PASSWORD=sylms",
		},
		Mounts: []string{
			pwd + "/testdata/testdata1.sql:/docker-entrypoint-initdb.d/testdata1.sql",
		},
	}

	resource, err := pool.RunWithOptions(runOptions)
	if err != nil {
		return nil, fmt.Errorf("could not start resource: %s", err)
	}

	const waitMaxSeconds = 120
	err = resource.Expire(waitMaxSeconds)
	if err != nil {
		return nil, err
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://sylms:sylms@%s/courses?sslmode=disable", hostAndPort)

	var db *sqlx.DB
	pool.MaxWait = waitMaxSeconds * time.Second
	if err = pool.Retry(func() error {
		db, err = sqlx.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		return nil, fmt.Errorf("could not connect to docker: %s", err)
	}

	return db, nil
}
