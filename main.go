package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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
)

func main() {
	envKeys := []string{envSylmsPostgresDBKey, envSylmsPostgresUserKey, envSylmsPostgresPasswordKey, envSylmsPostgresHostKey, envSylmsPostgresPortKey}
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

	var err error
	db, err = sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", postgresHost, postgresPort, postgresUser, postgresPassword, postgresDb))
	if err != nil {
		log.Fatalf("%+v", err)
	}

	r := mux.NewRouter()
	// とりあえず科目名と授業概要で検索できるように
	r.HandleFunc("/course", courseSimpleSearchHandler).Queries(
		"course_name", "{course_name}",
		"course_name_filter_type", "{course_name_filter_type}",
		"course_overview", "{course_overview}",
		"course_overview_filter_type", "{course_overview_filter_type}",
		"filter_type", "{filter_type}",
		"limit", "{limit}",
	).Methods("GET")
	c := cors.Default().Handler(r)
	http.ListenAndServe(":8080", c)
}

func courseSimpleSearchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseName := vars["course_name"]
	courseNameFilterType := vars["course_name_filter_type"]
	courseOverview := vars["course_overview"]
	courseOverviewFilterType := vars["course_overview_filter_type"]
	filterType := vars["filter_type"]
	limit := vars["limit"]

	// if !(filterType == "and" || filterType == "or") {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// TODO: デフォルト値を設定する
	if limit == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	courses, err := searchCourse(searchCourseOptions{
		courseName:               courseName,
		courseNameFilterType:     courseNameFilterType,
		courseOverview:           courseOverview,
		courseOverviewFilterType: courseOverviewFilterType,
		filterType:               filterType,
		limit:                    limitInt,
	})

	if err != nil {
		log.Fatalf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// CourseDB -> CourseJSON
	// TODO: わざわざ2つ定義しているのは面倒なのでひとつにしたい
	var coursesJSON []CourseJSON
	for _, c := range courses {
		var sry []string
		c.StandardRegistrationYear.Scan(&sry)

		var term []int
		c.Term.Scan(&term)

		var period []string
		c.Period.Scan(&period)

		var instructor []string
		c.Instructor.Scan(&instructor)

		courseJSON := CourseJSON{
			ID:                       c.ID,
			CourseNumber:             c.CourseNumber,
			CourseName:               c.CourseName,
			InstructionalType:        c.InstructionalType,
			Credits:                  c.Credits,
			StandardRegistrationYear: sry,
			Term:                     term,
			Period:                   period,
			Classroom:                newNullString(c.Classroom.String, c.Classroom.Valid),
			Instructor:               instructor,
			CourseOverview:           newNullString(c.CourseOverview.String, c.CourseOverview.Valid),
			Remarks:                  newNullString(c.Remarks.String, c.Remarks.Valid),
			CreditedAuditors:         c.CreditedAuditors,
			ApplicationConditions:    newNullString(c.ApplicationConditions.String, c.ApplicationConditions.Valid),
			AltCourseName:            newNullString(c.AltCourseName.String, c.AltCourseName.Valid),
			CourseCode:               newNullString(c.CourseCode.String, c.CourseCode.Valid),
			CourseCodeName:           newNullString(c.CourseCodeName.String, c.CourseCodeName.Valid),
			CSVUpdatedAt:             c.CSVUpdatedAt,
			Year:                     c.Year,
			CreatedAt:                c.CreatedAt,
			UpdatedAt:                c.UpdatedAt,
		}
		coursesJSON = append(coursesJSON, courseJSON)
	}

	j, err := json.Marshal(coursesJSON)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
