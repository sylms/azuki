package main

import (
	"fmt"
	"strconv"

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
func buildSimpleQuery(rawStr string, filterType string, dbColumnName string, selectArgs []interface{}, placeholderCount int) (string, int, []interface{}, error) {
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
	return resQuery, placeholderCount, selectArgs, nil
}

// buildSearchCourseQuery からの切り分け
// TODO : 汎用的にしたい、str の配列などで
func connectEachSimpleQuery(queryCourseName string, queryCourseOverview string, filterType string) (string, error) {
	resStr := ""
	if queryCourseName != "" {
		if resStr != "" {
			resStr += filterType
		}
		resStr += "(" + queryCourseName + ")"
	}
	if queryCourseOverview != "" {
		if resStr != "" {
			resStr += filterType
		}
		resStr += "(" + queryCourseOverview + ")"
	}
	return "(" + resStr + ")", nil
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

	// where 部分を構築
	queryCourseName, placeholderCount, selectArgs, _ := buildSimpleQuery(options.CourseName, options.CourseNameFilterType, "course_name", selectArgs, placeholderCount)
	queryCourseOverview, placeholderCount, selectArgs, _ := buildSimpleQuery(options.CourseOverview, options.CourseOverviewFilterType, "course_overview", selectArgs, placeholderCount)

	// カラムごとに生成されたクエリを接続
	queryWhere, _ := connectEachSimpleQuery(queryCourseName, queryCourseOverview, options.FilterType)

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
