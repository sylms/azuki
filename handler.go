package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func courseSimpleSearchHandler(w http.ResponseWriter, r *http.Request) {

	//Validate request
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	//parse json
	var query CourseQuery
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	options, err := validateSearchCourseOptions(query)
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
	coursesJSON := []CourseJSON{}
	for _, c := range courses {

		var term []int
		for _, i := range c.Term {
			term = append(term, int(i))
		}

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
	_, err = w.Write(j)
	if err != nil {
		log.Printf("%+v", err)
	}
}
