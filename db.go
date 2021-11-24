package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sylms/azuki/util"
)

const (
	filterTypeAnd = "and"
	filterTypeOr  = "or"
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
// 空白区切りで分割し指定されたフィルタータイプで繋いだクエリを生成する。
// 与えられたプレースホルダーカウントの値から順にプレースホルダーに整数を割り当てていく
// TODO : create test
func buildArrayQuery(rawStr string, filterType string, dbColumnName string, selectArgs []interface{}, placeholderCount int) (string, int, []interface{}) {
	separatedStrList, _ := periodParser(rawStr)
	resQuery := ""
	for count, separseparatedStr := range separatedStrList {
		if count == 0 {
			resQuery += fmt.Sprintf(`$%d = ANY(%s) `, placeholderCount, dbColumnName)
		} else {
			resQuery += fmt.Sprintf(`%s $%d = ANY(%s) `, filterType, placeholderCount, dbColumnName)
		}
		placeholderCount++
		// 現時点では、キーワードを含むものを検索
		selectArgs = append(selectArgs, separseparatedStr)
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
	fmt.Printf("===%s===\n", options.Period)
	queryPeriod, placeholderCount, selectArgs := buildArrayQuery(options.Period, options.CourseOverviewFilterType, "period_", selectArgs, placeholderCount)
	queryLists = append(queryLists, queryPeriod)

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

func searchCourse(query string, args []interface{}) ([]CoursesDB, error) {
	var result []CoursesDB
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
