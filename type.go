package main

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type nullString struct {
	sql.NullString
}

type CoursesDB struct {
	ID                       int            `db:"id"`
	CourseNumber             string         `db:"course_number"`
	CourseName               string         `db:"course_name"`
	InstructionalType        int            `db:"instructional_type"`
	Credits                  string         `db:"credits"`
	StandardRegistrationYear pq.StringArray `db:"standard_registration_year"`
	Term                     pq.Int64Array  `db:"term"`
	Period                   pq.StringArray `db:"period_"`
	Classroom                sql.NullString `db:"classroom"`
	Instructor               pq.StringArray `db:"instructor"`
	CourseOverview           sql.NullString `db:"course_overview"`
	Remarks                  sql.NullString `db:"remarks"`
	CreditedAuditors         int            `db:"credited_auditors"`
	ApplicationConditions    sql.NullString `db:"application_conditions"`
	AltCourseName            sql.NullString `db:"alt_course_name"`
	CourseCode               sql.NullString `db:"course_code"`
	CourseCodeName           sql.NullString `db:"course_code_name"`
	CSVUpdatedAt             time.Time      `db:"csv_updated_at"`
	Year                     int            `db:"year"`
	CreatedAt                time.Time      `db:"created_at"`
	UpdatedAt                time.Time      `db:"updated_at"`
}

type FacetDB struct {
	Term      int `db:"term"`
	TermCount int `db:"term_count"`
}

type CourseJSON struct {
	ID                       int        `json:"id"`
	CourseNumber             string     `json:"course_number"`
	CourseName               string     `json:"course_name"`
	InstructionalType        int        `json:"instructional_type"`
	Credits                  string     `json:"credits"`
	StandardRegistrationYear []string   `json:"standard_registration_year"`
	Term                     []int      `json:"term"`
	Period                   []string   `json:"period"`
	Classroom                nullString `json:"classroom"`
	Instructor               []string   `json:"instructor"`
	CourseOverview           nullString `json:"course_overview"`
	Remarks                  nullString `json:"remarks"`
	CreditedAuditors         int        `json:"credited_auditors"`
	ApplicationConditions    nullString `json:"application_conditions"`
	AltCourseName            nullString `json:"alt_course_name"`
	CourseCode               nullString `json:"course_code"`
	CourseCodeName           nullString `json:"course_code_name"`
	CSVUpdatedAt             time.Time  `json:"csv_updated_at"`
	Year                     int        `json:"year"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

type CourseQuery struct {
	CourseNumber             string `json:"course_number"`
	CourseName               string `json:"course_name"`
	InstructionalType        int    `json:"instructional_type"`
	Credits                  string `json:"credits"`
	StandardRegistrationYear int    `json:"standard_registration_year"`
	Term                     string `json:"term"`
	Period                   string `json:"period"`
	Classroom                string `json:"classroom"`
	Instructor               string `json:"instructor"`
	CourseOverview           string `json:"course_overview"`
	Remarks                  string `json:"remarks"`
	CourseNameFilterType     string `json:"course_name_filter_type"`
	CourseOverviewFilterType string `json:"course_overview_filter_type"`
	FilterType               string `json:"filter_type"`
	Limit                    int    `json:"limit"`
	Offset                   int    `json:"offset"`
}

type FacetJSON struct {
	TermFacet map[int]int `json:"term_facet"`
}
