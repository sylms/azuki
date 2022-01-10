package persistence

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sylms/azuki/domain"
	"github.com/sylms/azuki/util"
)

// TODO: csv2sql からパッケージとして読み込む
const (
	// 開講時期
	_               = iota
	termSpringACode // 春A: 1
	termSpringBCode
	termSpringCCode
	termFallACode
	termFallBCode
	termFallCCode
	termSummerVacationCode
	termSpringVacationCode
	termAllCode
	termSpringCode
	termFallCode
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
		separatedStrList, _ = periodParser(rawStr)
	}
	if dbColumnName == "term" {
		term := termParser(rawStr)
		separatedStrList, _ = termStrToInt(term)
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
			// 現時点では、キーワードを含むものを検索
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

// === csv2sql から ===
// TODO: csv2sql からのコピペいい感じにモジュールとして入れたい
func periodParser(periodString string) ([]string, error) {
	period := []string{}
	periodString = strings.Replace(periodString, " ", "", -1)
	periodString = strings.Replace(periodString, "　", "", -1)
	periodString = strings.Replace(periodString, "ー", "-", -1)
	periodString = strings.Replace(periodString, "・", "", -1)
	periodString = strings.Replace(periodString, ",", "", -1)
	periodString = strings.Replace(periodString, "集中", "集0", -1)
	periodString = strings.Replace(periodString, "応談", "応0", -1)
	periodString = strings.Replace(periodString, "随時", "随0", -1)

	for i := 1; i <= 8; i++ {
		listPeriod := strconv.Itoa(i)
		for j := i + 1; j <= 8; j++ {
			listPeriod = listPeriod + strconv.Itoa(j)
			spanPeriod := strconv.Itoa(i) + "-" + strconv.Itoa(j)
			periodString = strings.Replace(periodString, spanPeriod, listPeriod, -1)
		}
	}

	for i := 0; i <= 8; i++ {
		for _, dayOfWeek := range []string{"月", "火", "水", "木", "金", "土", "日", "応", "随", "集"} {
			beforeStr1 := strconv.Itoa(i) + dayOfWeek
			beforeStr2 := dayOfWeek + strconv.Itoa(i)
			afterStr1 := strconv.Itoa(i) + "," + dayOfWeek
			afterStr2 := dayOfWeek + ":" + strconv.Itoa(i)
			periodString = strings.Replace(periodString, beforeStr1, afterStr1, -1)
			periodString = strings.Replace(periodString, beforeStr2, afterStr2, -1)
		}
	}
	if len(periodString) == 0 {
		return period, nil
	}
	strList := strings.Split(periodString, ",")

	for _, str := range strList {
		strList2 := strings.Split(str, ":")
		if len(strList2) != 2 {
			fmt.Println("-" + periodString + "-")
			return nil, errors.New("unexpected period input : " + str)
		} else {
			dayOfWeek := strList2[0]
			timeTimetable := strList2[1]
			for i := 0; i < len([]rune(dayOfWeek)); i++ {
				for j := 0; j < len([]rune(timeTimetable)); j++ {
					inputStr := string([]rune(dayOfWeek)[i]) + string([]rune(timeTimetable)[j])
					inputStr = strings.Replace(inputStr, "集0", "集", -1)
					inputStr = strings.Replace(inputStr, "集", "集中", -1)
					inputStr = strings.Replace(inputStr, "随0", "随", -1)
					inputStr = strings.Replace(inputStr, "随", "随時", -1)
					inputStr = strings.Replace(inputStr, "応0", "応", -1)
					inputStr = strings.Replace(inputStr, "応", "応談", -1)
					period = append(period, inputStr)
				}
			}
		}
	}

	return period, nil
}

func termStrToInt(term []string) ([]string, error) {
	res := []string{}
	for _, t := range term {
		switch t {
		case "春A":
			res = append(res, strconv.Itoa(termSpringACode))
		case "春B":
			res = append(res, strconv.Itoa(termSpringBCode))
		case "春C":
			res = append(res, strconv.Itoa(termSpringCCode))
		case "秋A":
			res = append(res, strconv.Itoa(termFallACode))
		case "秋B":
			res = append(res, strconv.Itoa(termFallBCode))
		case "秋C":
			res = append(res, strconv.Itoa(termFallCCode))
		case "夏季休業中":
			res = append(res, strconv.Itoa(termSummerVacationCode))
		case "春季休業中":
			res = append(res, strconv.Itoa(termSpringVacationCode))
		case "通年":
			res = append(res, strconv.Itoa(termAllCode))
		case "春学期":
			res = append(res, strconv.Itoa(termSpringCode))
		case "秋学期":
			res = append(res, strconv.Itoa(termFallCode))
		default:
			return nil, fmt.Errorf("invalid term string: %s", t)
		}
	}
	return res, nil
}

func termParser(termString string) []string {
	res := []string{}
	if termString == "" {
		return []string{}
	}
	var re *regexp.Regexp
	re = regexp.MustCompile(`(春A|春AA|春AA|春AB|春BA|春AC|春CA|春ABC)`)
	if re.MatchString(termString) {
		res = append(res, "春A")
	}
	re = regexp.MustCompile(`(春B|春BA|春AB|春BB|春BB|春BC|春CB|春ABC)`)
	if re.MatchString(termString) {
		res = append(res, "春B")
	}
	re = regexp.MustCompile(`(春C|春CA|春AC|春CB|春BC|春CC|春CC|春ABC)`)
	if re.MatchString(termString) {
		res = append(res, "春C")
	}
	re = regexp.MustCompile(`(秋A|秋AA|秋AA|秋AB|秋BA|秋AC|秋CA|秋ABC)`)
	if re.MatchString(termString) {
		res = append(res, "秋A")
	}
	re = regexp.MustCompile(`(秋B|秋BA|秋AB|秋BB|秋BB|秋BC|秋CB|秋ABC)`)
	if re.MatchString(termString) {
		res = append(res, "秋B")
	}
	re = regexp.MustCompile(`(秋C|秋CA|秋AC|秋CB|秋BC|秋CC|秋CC|秋ABC)`)
	if re.MatchString(termString) {
		res = append(res, "秋C")
	}
	re = regexp.MustCompile(`(夏季休業中)`)
	if re.MatchString(termString) {
		res = append(res, "夏季休業中")
	}
	re = regexp.MustCompile(`(春季休業中)`)
	if re.MatchString(termString) {
		res = append(res, "春季休業中")
	}
	re = regexp.MustCompile(`(通年)`)
	if re.MatchString(termString) {
		res = append(res, "通年")
	}
	re = regexp.MustCompile(`(春学期)`)
	if re.MatchString(termString) {
		res = append(res, "春学期")
	}
	re = regexp.MustCompile(`(秋学期)`)
	if re.MatchString(termString) {
		res = append(res, "秋学期")
	}
	return res
}

// ==========
