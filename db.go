package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sylms/azuki/util"
)

const (
	filterTypeAnd = "and"
	filterTypeOr  = "or"
)

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

// buildSearchCourseQuery からの切り分け
// 空白区切りで分割し指定されたフィルタータイプで繋いだクエリを生成する。
// 与えられたプレースホルダーカウントの値から順にプレースホルダーに整数を割り当てていく
// TODO : create test
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

// buildSearchCourseQuery からの切り分け
// 単なるテキスト配列用
// TODO : create test
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
			resQuery += fmt.Sprintf(`]::varchar[] <@ %s`, dbColumnName)
		}
		if dbColumnName == "term" {
			resQuery += fmt.Sprintf(`]::int[] <@ %s`, dbColumnName)
		}
	}
	return resQuery, placeholderCount, selectArgs
}

// buildSearchCourseQuery からの切り分け
// TODO : 汎用的にしたい、str の配列などで
// TODO : create test
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

// validateSearchCourseOptions() の返り値の searchCourseOptions を元に DB へ投げるクエリ文字列とそれら引数を作成する
func buildSearchCourseQuery(options CourseQuery) (string, []interface{}, error) {
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

	const queryHead = `select * from courses where `
	return queryHead + queryWhere + queryOrderBy + queryLimit + queryOffset, selectArgs, nil
}

// validateSearchCourseOptions() の返り値の searchCourseOptions を元に DB へ投げるクエリ文字列とそれら引数を作成する
func buildGetFacetQuery(options CourseQuery) (string, []interface{}, error) {
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
	queryCourseName, _, selectArgs := buildSimpleQuery(options.CourseName, options.CourseNameFilterType, "course_name", selectArgs, placeholderCount)
	queryLists = append(queryLists, queryCourseName)
	queryCourseOverview, _, selectArgs := buildSimpleQuery(options.CourseOverview, options.CourseOverviewFilterType, "course_overview", selectArgs, placeholderCount)
	queryLists = append(queryLists, queryCourseOverview)
	queryCourseNumber, _, selectArgs := buildSimpleQuery(options.CourseNumber, options.CourseOverviewFilterType, "course_number", selectArgs, placeholderCount)
	queryLists = append(queryLists, queryCourseNumber)
	queryPeriod, _, selectArgs := buildArrayQuery(options.Period, options.CourseOverviewFilterType, "period_", selectArgs, placeholderCount)
	queryLists = append(queryLists, queryPeriod)
	queryTerm, _, selectArgs := buildArrayQuery(options.Term, options.CourseOverviewFilterType, "term", selectArgs, placeholderCount)
	queryLists = append(queryLists, queryTerm)

	// カラムごとに生成されたクエリを接続
	queryWhere := connectEachSimpleQuery(queryLists, options.FilterType)

	// order by
	// const queryOrderBy = "order by id asc "

	// limit 部分を構築
	// queryLimit := fmt.Sprintf(`limit $%d `, placeholderCount)
	// placeholderCount++
	// selectArgs = append(selectArgs, strconv.Itoa(options.Limit))

	// offset 部分を構築
	// queryOffset := fmt.Sprintf(`offset $%d`, placeholderCount)
	// selectArgs = append(selectArgs, strconv.Itoa(options.Offset))

	const queryHead = `select unnest(term) as term from courses where `
	return `select term, count(term) as term_count from(` + queryHead + queryWhere + `) as s1 group by term`, selectArgs, nil
}

func searchCourse(query string, args []interface{}) ([]CoursesDB, error) {
	var result []CoursesDB
	err := db.Select(&result, query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func getFacet(query string, args []interface{}) ([]FacetDB, error) {
	var result []FacetDB
	err := db.Select(&result, query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

func validateSearchCourseOptions(query CourseQuery) error {

	allowedFilterType := []string{filterTypeAnd, filterTypeOr}
	emptyQuery := true
	if !util.Contains(allowedFilterType, query.FilterType) {
		return fmt.Errorf("FilterType error: %s, %+v", query.FilterType, allowedFilterType)
	}
	if query.CourseNumber != "" {
		emptyQuery = false
	}
	if query.CourseName != "" {
		if !util.Contains(allowedFilterType, query.CourseNameFilterType) {
			return fmt.Errorf("CourseNameFilterType error: %s, %+v", query.CourseNameFilterType, allowedFilterType)
		}
		emptyQuery = false
	}
	if query.CourseOverview != "" {
		if !util.Contains(allowedFilterType, query.CourseOverviewFilterType) {
			return fmt.Errorf("CourseOverviewFilterType error: %s, %+v", query.CourseOverviewFilterType, allowedFilterType)
		}
		emptyQuery = false
	}
	if query.Period != "" {
		emptyQuery = false
	}

	if query.Term != "" {
		emptyQuery = false
	}

	// どのカラムも検索対象としていなければ検索そのものが実行できないので、不正なリクエストである
	if emptyQuery {
		return errors.New("all parameter is empty")
	}

	if query.Limit < 0 {
		return errors.New("limit is negative")
	}

	if query.Offset < 0 {
		return errors.New("offset is negative")
	}

	return nil
}

// csv2sql からのコピペいい感じにモジュールとして入れたい
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
