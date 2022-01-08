package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gocarina/gocsv"
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

	err = validateSearchCourseOptions(query)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// SQL クエリ文字列を構築
	queryStr, queryArgs, err := buildSearchCourseQuery(query)
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

func courseCSVHandler(w http.ResponseWriter, r *http.Request) {

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
	query.Offset = 0
	query.Limit = 10000000

	err = validateSearchCourseOptions(query)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// SQL クエリ文字列を構築
	queryStr, queryArgs, err := buildSearchCourseQuery(query)
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

	// CourseDB -> CourseCSV
	// TODO: わざわざ2つ定義しているのは面倒なのでひとつにしたい
	coursesCSV := []CourseCSV{}
	for _, c := range courses {

		term := ""
		for _, i := range c.Term {
			// term = append(term, int(i))
			if 0 <= i && i < 12 {
				termLs := [...]string{
					"春A",
					"春B",
					"春C",
					"秋A",
					"秋B",
					"秋C",
					"夏季休業中",
					"春季休業中",
					"通年",
					"春学期",
					"秋学期",
				}
				term += termLs[i-1]
			} else {
				term += "その他"
			}
			// decodeTerm: function (num: number) {
			// 	return ls[num - 1];
			//   },
		}

		courseCSV := CourseCSV{
			CourseNumber:             c.CourseNumber,
			CourseName:               c.CourseName,
			InstructionalType:        c.InstructionalType,
			Credits:                  c.Credits,
			StandardRegistrationYear: strings.Join(c.StandardRegistrationYear, ","),
			Term:                     term,
			Period:                   strings.Join(c.Period, ","),
			Classroom:                c.Classroom.String,
			Instructor:               strings.Join(c.Instructor, ","),
			CourseOverview:           c.CourseOverview.String,
			Remarks:                  c.Remarks.String,
			CreditedAuditors:         c.CreditedAuditors,
			ApplicationConditions:    c.ApplicationConditions.String,
			AltCourseName:            c.AltCourseName.String,
			CourseCode:               c.CourseCode.String,
			CourseCodeName:           c.CourseCodeName.String,
			UpdatedAt:                c.UpdatedAt,
		}
		coursesCSV = append(coursesCSV, courseCSV)
	}

	// j, err := json.Marshal(coursesCSV)
	// if err != nil {
	// 	log.Printf("%+v", err)
	// }

	// for _, c := range courses {
	// 	fmt.Println(c)
	// }

	csvStr, err := gocsv.MarshalString(&coursesCSV)
	if err != nil {
		log.Printf("%+v", err)
	}

	w.Header().Set("Content-Type", "text/csv; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	csvStrByte := []byte(csvStr)
	_, err = w.Write(csvStrByte)
	if err != nil {
		log.Printf("%+v", err)
	}
}

func courseFacetSearchHandler(w http.ResponseWriter, r *http.Request) {

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

	err = validateSearchCourseOptions(query)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// SQL クエリ文字列を構築
	queryStr, queryArgs, err := buildGetFacetQuery(query)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// DB へクエリを投げ結果を取得
	courses, err := getFacet(queryStr, queryArgs)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// if err != nil {
	// 	log.Printf("%+v", err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// CourseDB -> FacetJSON
	// TODO: わざわざ2つ定義しているのは面倒なのでひとつにしたい
	// for _, c := range courses {

	// 	var term []int
	// 	for _, i := range c.Term {
	// 		term = append(term, int(i))
	// 	}

	// 	FacetJSON := FacetJSON{
	// 		// ID: c.ID,
	// 		ID: 12345,
	// 	}
	// 	facetsJSON = append(facetsJSON, FacetJSON)
	// }
	termFacet := make(map[int]int)
	for _, c := range courses {
		termFacet[c.Term] = c.TermCount
	}
	facetJSON := FacetJSON{
		TermFacet: termFacet,
	}
	j, err := json.Marshal(facetJSON)
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
