package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/sylms/azuki/domain"
	"github.com/sylms/azuki/usecase"
	"github.com/sylms/azuki/util"
	"github.com/sylms/csv2sql/kdb"
)

type CourseHandler interface {
	Search(http.ResponseWriter, *http.Request)
	Csv(http.ResponseWriter, *http.Request)
	Facet(http.ResponseWriter, *http.Request)
}

type courseHandler struct {
	uc usecase.CourseUseCase
}

func NewCourseHandler(uc usecase.CourseUseCase) CourseHandler {
	return &courseHandler{
		uc: uc,
	}
}

type CourseJSON struct {
	ID                       int       `json:"id"`
	CourseNumber             string    `json:"course_number"`
	CourseName               string    `json:"course_name"`
	InstructionalType        int       `json:"instructional_type"`
	Credits                  string    `json:"credits"`
	StandardRegistrationYear []string  `json:"standard_registration_year"`
	Term                     []int     `json:"term"`
	Period                   []string  `json:"period"`
	Classroom                string    `json:"classroom"`
	Instructor               []string  `json:"instructor"`
	CourseOverview           string    `json:"course_overview"`
	Remarks                  string    `json:"remarks"`
	CreditedAuditors         int       `json:"credited_auditors"`
	ApplicationConditions    string    `json:"application_conditions"`
	AltCourseName            string    `json:"alt_course_name"`
	CourseCode               string    `json:"course_code"`
	CourseCodeName           string    `json:"course_code_name"`
	CSVUpdatedAt             time.Time `json:"csv_updated_at"`
	Year                     int       `json:"year"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type CourseCSV struct {
	CourseNumber             string    `csv:"科目番号"`
	CourseName               string    `csv:"科目名"`
	InstructionalType        int       `csv:"授業方法"`
	Credits                  string    `csv:"単位数"`
	StandardRegistrationYear string    `csv:"標準履修年次"`
	Term                     string    `csv:"実施学期"`
	Period                   string    `csv:"曜時限"`
	Classroom                string    `csv:"教室"`
	Instructor               string    `csv:"担当教員"`
	CourseOverview           string    `csv:"授業概要"`
	Remarks                  string    `csv:"備考"`
	CreditedAuditors         int       `csv:"科目等履修生申請可否"`
	ApplicationConditions    string    `csv:"申請条件"`
	AltCourseName            string    `csv:"英語(日本語)科目名"`
	CourseCode               string    `csv:"科目コード"`
	CourseCodeName           string    `csv:"要件科目名"`
	UpdatedAt                time.Time `csv:"データ更新日"`
}

type FacetJSON struct {
	TermFacet map[int]int `json:"term_facet"`
}

func (h *courseHandler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// TODO: domain が依存しているが良いか？
	var query domain.CourseQuery
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validateSearchCourseQuery(query)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	courses, err := h.uc.Search(query)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var coursesJson []CourseJSON
	for _, course := range courses {
		courseJson := CourseJSON(*course)
		coursesJson = append(coursesJson, courseJson)
	}

	resJson, err := json.Marshal(coursesJson)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resJson)
	if err != nil {
		log.Printf("%+v", err)
	}
}

// TODO: これは interface or usecase ？
func validateSearchCourseQuery(query domain.CourseQuery) error {
	allowedFilterType := []string{"and", "or"}
	if !util.Contains(allowedFilterType, query.FilterType) {
		return fmt.Errorf("FilterType error: %s, %+v", query.FilterType, allowedFilterType)
	}
	if query.CourseName != "" {
		if !util.Contains(allowedFilterType, query.CourseNameFilterType) {
			return fmt.Errorf("CourseNameFilterType error: %s, %+v", query.CourseNameFilterType, allowedFilterType)
		}
	}
	if query.CourseOverview != "" {
		if !util.Contains(allowedFilterType, query.CourseOverviewFilterType) {
			return fmt.Errorf("CourseOverviewFilterType error: %s, %+v", query.CourseOverviewFilterType, allowedFilterType)
		}
	}

	if query.Period != "" {
		_, err := kdb.PeriodParser(query.Period)
		if err != nil {
			return fmt.Errorf("'period' parse error: %+v", err)
		}
	}

	if query.Term != "" {
		terms := kdb.TermParser(query.Term)
		if len(terms) == 0 {
			// Term に何か与えられているもののパースした結果どの開講時期でも無いので与えられた文字列がおかしい
			// "春Aははは" みたいな、きちんとした開講時期とおかしな文字列の両方が含まれる場合については、とりあえず考えないこととする
			return fmt.Errorf("'term' parse error")
		}
	}

	if query.Limit < 0 {
		return errors.New("limit is negative")
	}

	if query.Offset < 0 {
		return errors.New("offset is negative")
	}

	return nil
}

func (h *courseHandler) Csv(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// TODO: domain が依存しているが良いか？
	var query domain.CourseQuery
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validateSearchCourseQuery(query)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: 無理矢理書き換えないようにする
	query.Offset = 0
	query.Limit = 10000000

	courses, err := h.uc.Search(query)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var coursesCsv []CourseCSV
	for _, course := range courses {
		// TODO: Term をカンマ区切りで結合する
		term := ""
		for _, termIndex := range course.Term {
			termStr, err := decodeTerm(termIndex)
			if err != nil {
				log.Printf("%+v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			term += termStr
		}

		courseCsv := CourseCSV{
			CourseNumber:             course.CourseNumber,
			CourseName:               course.CourseName,
			InstructionalType:        course.InstructionalType,
			Credits:                  course.Credits,
			StandardRegistrationYear: strings.Join(course.StandardRegistrationYear, ","),
			Term:                     term,
			Period:                   strings.Join(course.Period, ","),
			Classroom:                course.Classroom,
			Instructor:               strings.Join(course.Instructor, ","),
			CourseOverview:           course.CourseOverview,
			Remarks:                  course.Remarks,
			CreditedAuditors:         course.CreditedAuditors,
			ApplicationConditions:    course.ApplicationConditions,
			AltCourseName:            course.AltCourseName,
			CourseCode:               course.CourseCode,
			CourseCodeName:           course.CourseCodeName,
			UpdatedAt:                course.UpdatedAt,
		}
		coursesCsv = append(coursesCsv, courseCsv)
	}

	csvStr, err := gocsv.MarshalString(coursesCsv)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(csvStr))
	if err != nil {
		log.Printf("%+v", err)
	}
}

// 開講時期を数値から文字列に変換
func decodeTerm(index int) (string, error) {
	terms := []string{
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

	if index < 0 || len(terms)-1 < index {
		return "", fmt.Errorf("index range error: 0 - %d", len(terms)-1)
	}

	return terms[index], nil
}

func (h *courseHandler) Facet(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// TODO: domain が依存しているが良いか？
	var query domain.CourseQuery
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validateSearchCourseQuery(query)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	facets, err := h.uc.Facet(query)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	termFacet := map[int]int{}
	for _, facet := range facets {
		termFacet[facet.Term] = facet.TermCount
	}
	facetJson := FacetJSON{
		TermFacet: termFacet,
	}

	j, err := json.Marshal(facetJson)
	if err != nil {
		log.Printf("%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(j)
	if err != nil {
		log.Printf("%+v", err)
	}
}
