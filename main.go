package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/sylms/azuki/util"
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
	r.HandleFunc("/course", courseSimpleSearchHandler).Methods("GET")
	c := cors.Default().Handler(r)
	log.Printf("Listen Port: %s", portStr)
	err = http.ListenAndServe(fmt.Sprintf(":%s", portStr), c)
	if err != nil {
		log.Fatalln(err)
	}
}

func courseSimpleSearchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	courseName := q.Get("course_name")
	courseNameFilterType := q.Get("course_name_filter_type")
	courseOverview := q.Get("course_overview")
	courseOverviewFilterType := q.Get("course_overview_filter_type")
	filterType := q.Get("filter_type")
	limit := q.Get("limit")
	offset := q.Get("offset")

	// if !(filterType == "and" || filterType == "or") {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	options, err := validateSearchCourseOptions(courseName, courseNameFilterType, courseOverview, courseOverviewFilterType, filterType, limit, offset)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// SQL クエリ文字列を構築
	queryStr, queryArgs, err := buildSearchCourseQuery(options)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// DB へクエリを投げ結果を取得
	courses, err := searchCourse(queryStr, queryArgs)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// CourseDB -> CourseJSON
	// TODO: わざわざ2つ定義しているのは面倒なのでひとつにしたい
	var coursesJSON []CourseJSON
	for _, c := range courses {

		var term []int
		c.Term.Scan(&term)

		courseJSON := CourseJSON{
			ID:                       c.ID,
			CourseNumber:             c.CourseNumber,
			CourseName:               c.CourseName,
			InstructionalType:        c.InstructionalType,
			Credits:                  c.Credits,
			StandardRegistrationYear: c.StandardRegistrationYear,
			Term:                     term,
			Period:                   c.Period,
			Classroom:                newNullString(c.Classroom.String, c.Classroom.Valid),
			Instructor:               c.Instructor,
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
		log.Printf("%+v", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// 各パラメーターに問題がないかを確認し、問題なければ整形したものを返す
func validateSearchCourseOptions(courseName string, courseNameFilterType string, courseOverview string, courseOverviewFilterType string, filterType string, limit string, offset string) (searchCourseOptions, error) {
	allowedFilterType := []string{filterTypeAnd, filterTypeOr}
	if !util.Contains(allowedFilterType, filterType) {
		return searchCourseOptions{}, fmt.Errorf("filterType error: %s, %+v", filterType, allowedFilterType)
	}
	if !util.Contains(allowedFilterType, courseNameFilterType) {
		return searchCourseOptions{}, fmt.Errorf("courseNameFilterType error: %s, %+v", courseNameFilterType, allowedFilterType)
	}
	if !util.Contains(allowedFilterType, courseOverviewFilterType) {
		return searchCourseOptions{}, fmt.Errorf("courseOverviewFilterType error: %s, %+v", courseOverviewFilterType, allowedFilterType)
	}

	// どのカラムも検索対象としていなければ検索そのものが実行できないので、不正なリクエストである
	if courseName == "" && courseOverview == "" {
		return searchCourseOptions{}, errors.New("course_name and course_overview are empty")
	}

	var limitInt int
	if limit == "" {
		limitInt = searchQueryDefaultLimit
	} else {
		var err error
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			return searchCourseOptions{}, errors.New("limit is not int")
		}
		if limitInt < 0 {
			return searchCourseOptions{}, errors.New("limit is negative")
		}
	}

	var offsetInt int
	if offset == "" {
		// offset = 0 であれば offset を指定しないときと同じ結果を得られる
		offsetInt = 0
	} else {
		var err error
		offsetInt, err = strconv.Atoi(offset)
		if err != nil {
			return searchCourseOptions{}, errors.New("offset is not int")
		}
		if offsetInt < 0 {
			return searchCourseOptions{}, errors.New("offset is negative")
		}
	}

	return searchCourseOptions{
		courseName:               courseName,
		courseNameFilterType:     courseNameFilterType,
		courseOverview:           courseOverview,
		courseOverviewFilterType: courseOverviewFilterType,
		filterType:               filterType,
		limit:                    limitInt,
		offset:                   offsetInt,
	}, nil
}
