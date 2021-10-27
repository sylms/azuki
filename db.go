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

type searchCourseOptions struct {
	courseName               string
	courseNameFilterType     string
	courseOverview           string
	courseOverviewFilterType string
	filterType               string
	limit                    int
	offset                   int
}

// validateSearchCourseOptions() の返り値の searchCourseOptions を元に DB へ投げるクエリ文字列とそれら引数を作成する
func buildSearchCourseQuery(options searchCourseOptions) (string, []interface{}, error) {
	// PostgreSQL へ渡す $1, $2 プレースホルダーのインクリメント
	placeholderCount := 1

	// PostgreSQL へ渡す select 文のプレースホルダーに割り当てる変数を格納
	selectArgs := []interface{}{}

	// スペース区切りとみなして単語を分割
	courseNames := util.SplitSpace(options.courseName)
	courseOverviews := util.SplitSpace(options.courseOverview)

	// where 部分を構築
	queryCourseName := ""
	for count, courseName := range courseNames {
		if count == 0 {
			queryCourseName += fmt.Sprintf(`course_name like $%d `, placeholderCount)
		} else {
			queryCourseName += fmt.Sprintf(`%s course_name like $%d `, options.courseNameFilterType, placeholderCount)
		}
		placeholderCount++
		// 現時点では、キーワードを含むものを検索
		selectArgs = append(selectArgs, "%"+courseName+"%")
	}

	queryCourseOverview := ""
	for count, courseOverview := range courseOverviews {
		if count == 0 {
			queryCourseOverview += fmt.Sprintf(`course_overview like $%d `, placeholderCount)
		} else {
			queryCourseOverview += fmt.Sprintf(`%s course_overview like $%d `, options.courseOverviewFilterType, placeholderCount)
		}
		placeholderCount++
		// 現時点では、キーワードを含むものを検索
		selectArgs = append(selectArgs, "%"+courseOverview+"%")
	}

	// 若干無理矢理な気もするのできれいにしたい
	queryWhere := ""
	if queryCourseName != "" {
		queryWhere += "( " + queryCourseName + ") "
	}
	if queryCourseOverview != "" {
		if queryWhere == "" {
			queryWhere = "( " + queryCourseOverview + ") "
		} else {
			queryWhere += options.filterType + " ( " + queryCourseOverview + ") "
		}
	}

	// order by
	const queryOrderBy = "order by id asc "

	// limit 部分を構築
	queryLimit := fmt.Sprintf(`limit $%d `, placeholderCount)
	placeholderCount++
	selectArgs = append(selectArgs, strconv.Itoa(options.limit))

	// offset 部分を構築
	queryOffset := fmt.Sprintf(`offset $%d`, placeholderCount)
	selectArgs = append(selectArgs, strconv.Itoa(options.offset))

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

// 各パラメーターに問題がないかを確認し、問題なければ整形したものを返す
func validateSearchCourseOptions(query CourseQuery) (searchCourseOptions, error) {

	allowedFilterType := []string{filterTypeAnd, filterTypeOr}
	emptyQuert := true
	if !util.Contains(allowedFilterType, query.FilterType) {
		return searchCourseOptions{}, fmt.Errorf("FilterType error: %s, %+v", query.FilterType, allowedFilterType)
	}
	if query.CourseName != "" {
		if !util.Contains(allowedFilterType, query.CourseNameFilterType) {
			return searchCourseOptions{}, fmt.Errorf("CourseNameFilterType error: %s, %+v", query.CourseNameFilterType, allowedFilterType)
		}
		emptyQuert = false
	}
	if query.CourseOverview != "" {
		if !util.Contains(allowedFilterType, query.CourseOverviewFilterType) {
			return searchCourseOptions{}, fmt.Errorf("CourseOverviewFilterType error: %s, %+v", query.CourseOverviewFilterType, allowedFilterType)
		}
		emptyQuert = false
	}

	// どのカラムも検索対象としていなければ検索そのものが実行できないので、不正なリクエストである
	if emptyQuert {
		return searchCourseOptions{}, errors.New("all query string is empty")
	}

	if query.Limit < 0 {
		return searchCourseOptions{}, errors.New("limit is negative")
	}

	if query.Offset < 0 {
		return searchCourseOptions{}, errors.New("offset is negative")
	}

	return searchCourseOptions{
		courseName:               query.CourseName,
		courseNameFilterType:     query.CourseNameFilterType,
		courseOverview:           query.CourseOverview,
		courseOverviewFilterType: query.CourseOverviewFilterType,
		filterType:               query.FilterType,
		limit:                    query.Limit,
		offset:                   query.Offset,
	}, nil
}
