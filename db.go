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

// searchCourseOptions を元に DB へ投げるクエリ文字列とそれら引数を作成する
func buildSearchCourseQuery(options searchCourseOptions) (string, []interface{}, error) {
	allowedFilterType := []string{filterTypeAnd, filterTypeOr}
	if !util.Contains(allowedFilterType, options.filterType) {
		return "", nil, fmt.Errorf("filterType error: %s, %+v", options.filterType, allowedFilterType)
	}
	if !util.Contains(allowedFilterType, options.courseNameFilterType) {
		return "", nil, fmt.Errorf("courseNameFilterType error: %s, %+v", options.filterType, allowedFilterType)
	}
	if !util.Contains(allowedFilterType, options.courseOverviewFilterType) {
		return "", nil, fmt.Errorf("courseOverviewFilterType error: %s, %+v", options.filterType, allowedFilterType)
	}

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

	// とりあえず各カラムの検索結果は AND でつなげるように
	// 若干無理矢理な気もするのできれいにしたい
	queryWhere := ""
	if queryCourseName != "" {
		queryWhere += queryCourseName
	}
	if queryCourseOverview != "" {
		if queryWhere == "" {
			queryWhere = queryCourseOverview
		} else {
			queryWhere += "and " + queryCourseOverview
		}
	}

	// limit 部分を構築
	queryLimit := fmt.Sprintf(`limit $%d`, placeholderCount)
	placeholderCount++
	selectArgs = append(selectArgs, strconv.Itoa(options.limit))

	// offset 部分を構築
	queryOffset := fmt.Sprintf(`offset $%d`, placeholderCount)
	selectArgs = append(selectArgs, strconv.Itoa(options.offset))

	const queryHead = `select * from courses where `
	return queryHead + queryWhere + queryLimit + queryOffset, selectArgs, nil
}

func searchCourse(query string, args []interface{}) ([]CoursesDB, error) {
	var result []CoursesDB
	err := db.Select(&result, query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}
