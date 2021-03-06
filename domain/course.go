package domain

import "time"

type Course struct {
	ID                       int
	CourseNumber             string
	CourseName               string
	InstructionalType        int
	Credits                  string
	StandardRegistrationYear []string
	Term                     []int
	Period                   []string
	Classroom                string
	Instructor               []string
	CourseOverview           string
	Remarks                  string
	CreditedAuditors         int
	ApplicationConditions    string
	AltCourseName            string
	CourseCode               string
	CourseCodeName           string
	CSVUpdatedAt             time.Time
	Year                     int
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type CourseQuery struct {
	CourseNumber string `json:"course_number"`
	// スペース区切り
	// CourseNameFilterType も指定する
	CourseName               string `json:"course_name"`
	InstructionalType        int    `json:"instructional_type"`
	Credits                  string `json:"credits"`
	StandardRegistrationYear int    `json:"standard_registration_year"`
	Term                     string `json:"term"`
	// csv2sql/kdb の PeriodParser で認識できる形式であれば良い
	Period     string `json:"period"`
	Classroom  string `json:"classroom"`
	Instructor string `json:"instructor"`
	// スペース区切り
	// CourseOverviewFilterType も指定する
	CourseOverview           string `json:"course_overview"`
	Remarks                  string `json:"remarks"`
	CourseNameFilterType     string `json:"course_name_filter_type"`
	CourseOverviewFilterType string `json:"course_overview_filter_type"`
	FilterType               string `json:"filter_type"`
	// 必須
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Facet struct {
	Term      int
	TermCount int
}

type CourseRepository interface {
	Search(CourseQuery) ([]*Course, error)
	Facet(CourseQuery) ([]*Facet, error)
}
