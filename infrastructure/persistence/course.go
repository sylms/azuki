package persistence

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sylms/azuki/domain"
	"github.com/sylms/azuki/util"
	"github.com/sylms/csv2sql/kdb"
)

type CoursesPostgresql struct {
	ID                       int            `db:"id"`
	CourseNumber             string         `db:"course_number"`
	CourseName               string         `db:"course_name"`
	InstructionalType        int            `db:"instructional_type"`
	Credits                  string         `db:"credits"`
	StandardRegistrationYear pq.StringArray `db:"standard_registration_year"`
	Term                     pq.Int64Array  `db:"term"`
	Period                   pq.StringArray `db:"period_"`
	Classroom                string         `db:"classroom"`
	Instructor               pq.StringArray `db:"instructor"`
	CourseOverview           string         `db:"course_overview"`
	Remarks                  string         `db:"remarks"`
	CreditedAuditors         int            `db:"credited_auditors"`
	ApplicationConditions    string         `db:"application_conditions"`
	AltCourseName            string         `db:"alt_course_name"`
	CourseCode               string         `db:"course_code"`
	CourseCodeName           string         `db:"course_code_name"`
	CSVUpdatedAt             time.Time      `db:"csv_updated_at"`
	Year                     int            `db:"year"`
	CreatedAt                time.Time      `db:"created_at"`
	UpdatedAt                time.Time      `db:"updated_at"`
}

type FacetPostgresql struct {
	Term      int `db:"term"`
	TermCount int `db:"term_count"`
}

type coursePersistence struct {
	db *sqlx.DB
}

func NewCoursePersistence(db *sqlx.DB) domain.CourseRepository {
	return &coursePersistence{
		db: db,
	}
}

func (p *coursePersistence) Search(query domain.CourseQuery) ([]*domain.Course, error) {
	queryStr, queryArgs, err := buildSearchCourseQuery(query)
	if err != nil {
		return nil, err
	}

	// とりあえず具体的な PostgreSQL と指定
	// TODO: これはもっと抽象にするべき？調査
	coursesDb, err := p.selectPostgresql(queryStr, queryArgs)
	if err != nil {
		return nil, err
	}

	var courses []*domain.Course
	for _, courseDb := range coursesDb {
		course := courseDb.toCourse()
		courses = append(courses, &course)
	}

	return courses, nil
}

func (p *coursePersistence) Facet(query domain.CourseQuery) ([]*domain.Facet, error) {
	queryStr, queryArgs, err := buildGetFacetQuery(query)
	if err != nil {
		return nil, err
	}

	// とりあえず具体的な PostgreSQL と指定
	// TODO: これはもっと抽象にするべき？調査
	facetsDb, err := p.selectPostgresqlFacet(queryStr, queryArgs)
	if err != nil {
		return nil, err
	}

	// 無駄な気がする
	var facets []*domain.Facet
	for _, facetDb := range facetsDb {
		facet := &domain.Facet{
			Term:      facetDb.Term,
			TermCount: facetDb.TermCount,
		}
		facets = append(facets, facet)
	}

	return facets, nil
}

// PostgreSQL へ SELECT を実行する
func (p *coursePersistence) selectPostgresql(query string, args []interface{}) ([]*CoursesPostgresql, error) {
	var result []*CoursesPostgresql
	err := p.db.Select(&result, query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

// selectPostgresql とひとまとめにしたい
func (p *coursePersistence) selectPostgresqlFacet(query string, args []interface{}) ([]*FacetPostgresql, error) {
	var result []*FacetPostgresql
	err := p.db.Select(&result, query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

// domain.Course に変換
// pq パッケージに依存しているところを整形する
func (c *CoursesPostgresql) toCourse() domain.Course {
	var term []int
	for _, i := range c.Term {
		term = append(term, int(i))
	}

	course := domain.Course{
		ID:                       c.ID,
		CourseNumber:             c.CourseNumber,
		CourseName:               c.CourseName,
		InstructionalType:        c.InstructionalType,
		Credits:                  c.Credits,
		StandardRegistrationYear: c.StandardRegistrationYear,
		Term:                     term,
		Period:                   c.Period,
		Classroom:                c.Classroom,
		Instructor:               c.Instructor,
		CourseOverview:           c.CourseOverview,
		Remarks:                  c.Remarks,
		CreditedAuditors:         c.CreditedAuditors,
		ApplicationConditions:    c.ApplicationConditions,
		AltCourseName:            c.AltCourseName,
		CourseCode:               c.CourseCode,
		CourseCodeName:           c.CourseCodeName,
		CSVUpdatedAt:             c.CSVUpdatedAt,
		Year:                     c.Year,
		CreatedAt:                c.CreatedAt,
		UpdatedAt:                c.UpdatedAt,
	}
	return course
}

func buildSearchCourseQuery(options domain.CourseQuery) (string, []interface{}, error) {
	// それぞれのカラムに対してカラム内検索の AND/OR が指定されている場合はそれで構築を行なう
	// それぞれのカラムに対して検索文字列を構築したらそれぞれの間を FilterType で埋める
	// 全体に対して offset, limit を指定する

	// PostgreSQL へ渡す $1, $2 プレースホルダーのインクリメントのカウンタ
	placeholderCount := 1

	// PostgreSQL へ渡す select 文のプレースホルダーに割り当てる変数を格納
	selectArgs := []interface{}{}

	// それぞれのカラムに対する小さなクエリの集合
	queryLists := []string{}

	// where 部分を構築
	queryCourseName, placeholderCount, selectArgs := buildSimpleQuery(options.CourseName, options.CourseNameFilterType, "course_name", selectArgs, placeholderCount)
	queryLists = append(queryLists, queryCourseName)
	queryCourseOverview, placeholderCount, selectArgs := buildSimpleQuery(options.CourseOverview, options.CourseOverviewFilterType, "course_overview", selectArgs, placeholderCount)
	queryLists = append(queryLists, queryCourseOverview)
	queryCourseNumber, placeholderCount, selectArgs := buildSimpleQuery(options.CourseNumber, options.CourseOverviewFilterType, "course_number", selectArgs, placeholderCount)
	queryLists = append(queryLists, queryCourseNumber)
	queryPeriod, placeholderCount, selectArgs := buildArrayQuery(options.Period, options.CourseOverviewFilterType, "period_", selectArgs, placeholderCount)
	queryLists = append(queryLists, queryPeriod)
	queryTerm, placeholderCount, selectArgs := buildArrayQuery(options.Term, options.CourseOverviewFilterType, "term", selectArgs, placeholderCount)
	queryLists = append(queryLists, queryTerm)

	// カラムごとに生成されたクエリを接続
	queryWhere := connectEachSimpleQuery(queryLists, options.FilterType)

	// order by
	const queryOrderBy = "order by id asc "

	// limit 部分を構築
	queryLimit := fmt.Sprintf(`limit $%d `, placeholderCount)
	placeholderCount++
	selectArgs = append(selectArgs, strconv.Itoa(options.Limit))

	// offset 部分を構築
	queryOffset := fmt.Sprintf(`offset $%d`, placeholderCount)
	selectArgs = append(selectArgs, strconv.Itoa(options.Offset))

	const queryHead = `select * from courses `
	queryWhere = "where " + queryWhere
	if queryWhere == "where ()" {
		queryWhere = ""
	}
	return queryHead + queryWhere + queryOrderBy + queryLimit + queryOffset, selectArgs, nil
}

func buildSimpleQuery(rawStr string, filterType string, dbColumnName string, selectArgs []interface{}, placeholderCount int) (string, int, []interface{}) {
	separatedStrList := util.SplitSpace(rawStr)
	resQuery := ""
	for count, separseparatedStr := range separatedStrList {
		if count == 0 {
			resQuery += fmt.Sprintf(`%s like $%d `, dbColumnName, placeholderCount)
		} else {
			resQuery += fmt.Sprintf(`%s %s like $%d `, filterType, dbColumnName, placeholderCount)
		}
		placeholderCount++
		// 現時点では、キーワードを含むものを検索
		selectArgs = append(selectArgs, "%"+separseparatedStr+"%")
	}
	return resQuery, placeholderCount, selectArgs
}

func connectEachSimpleQuery(queryLists []string, filterType string) string {
	resStr := ""
	for _, query := range queryLists {
		if query != "" {
			if resStr != "" {
				resStr += filterType
			}
			resStr += "(" + query + ")"
		}
	}
	return "(" + resStr + ")"
}

func buildArrayQuery(rawStr string, filterType string, dbColumnName string, selectArgs []interface{}, placeholderCount int) (string, int, []interface{}) {
	var separatedStrList []string
	if dbColumnName == "period_" {
		separatedStrList, _ = kdb.PeriodParser(rawStr)
	}
	if dbColumnName == "term" {
		terms := kdb.TermParser(rawStr)
		termsInt := []int{}
		for _, term := range terms {
			termInt, _ := kdb.TermStrToInt(term)
			termsInt = append(termsInt, termInt)
		}
		for _, termInt := range termsInt {
			separatedStrList = append(separatedStrList, strconv.Itoa(termInt))
		}
	}
	resQuery := ""
	if len(separatedStrList) != 0 {
		resQuery += "array["
		for count, separseparatedStr := range separatedStrList {
			if count == 0 {
				resQuery += fmt.Sprintf(`$%d`, placeholderCount)
			} else {
				resQuery += fmt.Sprintf(`, $%d`, placeholderCount)
			}
			placeholderCount++
			selectArgs = append(selectArgs, separseparatedStr)
		}
		if dbColumnName == "period_" {
			resQuery += fmt.Sprintf(`]::varchar[] @> %s and array[]::varchar[] <> %s`, dbColumnName, dbColumnName)
		}
		if dbColumnName == "term" {
			resQuery += fmt.Sprintf(`]::int[] @> %s and array[]::int[] <> %s`, dbColumnName, dbColumnName)
		}
	}
	return resQuery, placeholderCount, selectArgs
}

// TODO: buildSearchCourseQuery とほぼ同じなところを抜き出す
func buildGetFacetQuery(options domain.CourseQuery) (string, []interface{}, error) {
	// それぞれのカラムに対してカラム内検索の AND/OR が指定されている場合はそれで構築を行なう
	// それぞれのカラムに対して検索文字列を構築したらそれぞれの間を FilterType で埋める
	// 全体に対して offset, limit を指定する

	// PostgreSQL へ渡す $1, $2 プレースホルダーのインクリメントのカウンタ
	placeholderCount := 1

	// PostgreSQL へ渡す select 文のプレースホルダーに割り当てる変数を格納
	selectArgs := []interface{}{}

	// それぞれのカラムに対する小さなクエリの集合
	queryLists := []string{}

	// where 部分を構築
	queryCourseName, placeholderCount, selectArgs := buildSimpleQuery(options.CourseName, options.CourseNameFilterType, "course_name", selectArgs, placeholderCount)
	queryLists = append(queryLists, queryCourseName)
	queryCourseOverview, placeholderCount, selectArgs := buildSimpleQuery(options.CourseOverview, options.CourseOverviewFilterType, "course_overview", selectArgs, placeholderCount)
	queryLists = append(queryLists, queryCourseOverview)
	queryCourseNumber, placeholderCount, selectArgs := buildSimpleQuery(options.CourseNumber, options.CourseOverviewFilterType, "course_number", selectArgs, placeholderCount)
	queryLists = append(queryLists, queryCourseNumber)
	queryPeriod, placeholderCount, selectArgs := buildArrayQuery(options.Period, options.CourseOverviewFilterType, "period_", selectArgs, placeholderCount)
	queryLists = append(queryLists, queryPeriod)
	queryTerm, _, selectArgs := buildArrayQuery(options.Term, options.CourseOverviewFilterType, "term", selectArgs, placeholderCount)
	queryLists = append(queryLists, queryTerm)

	// カラムごとに生成されたクエリを接続
	queryWhere := connectEachSimpleQuery(queryLists, options.FilterType)

	const queryHead = `select unnest(term) as term from courses `
	queryWhere = "where " + queryWhere
	if queryWhere == "where ()" {
		queryWhere = ""
	}
	return `select term, count(term) as term_count from(` + queryHead + queryWhere + `) as s1 group by term`, selectArgs, nil
}
