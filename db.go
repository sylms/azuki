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

type searchCourseOptions struct {
	courseName               string
	courseNameFilterType     string
	courseOverview           string
	courseOverviewFilterType string
	filterType               string
	limit                    int
}

func searchCourse(options searchCourseOptions) ([]CoursesDB, error) {
	allowedFilterType := []string{filterTypeAnd, filterTypeOr}
	if !util.Contains(allowedFilterType, options.filterType) {
		return []CoursesDB{}, fmt.Errorf("filterType error: %s, %+v", options.filterType, allowedFilterType)
	}
	if !util.Contains(allowedFilterType, options.courseNameFilterType) {
		return []CoursesDB{}, fmt.Errorf("courseNameFilterType error: %s, %+v", options.filterType, allowedFilterType)
	}
	if !util.Contains(allowedFilterType, options.courseOverviewFilterType) {
		return []CoursesDB{}, fmt.Errorf("courseOverviewFilterType error: %s, %+v", options.filterType, allowedFilterType)
	}

	// PostgreSQL へ渡す $1, $2 プレースホルダーのインクリメント
	placeholderCount := 1

	selectArgs := []interface{}{}

	// TODO: 半角だけでなく全角にも対応
	// スペース区切りとみなして単語を分割
	courseNames := strings.Split(options.courseName, " ")
	courseOverviews := strings.Split(options.courseOverview, " ")

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
	queryWhere := queryCourseName + "and " + queryCourseOverview
	queryLimit := fmt.Sprintf(`limit $%d`, placeholderCount)
	selectArgs = append(selectArgs, strconv.Itoa(options.limit))

	var result []CoursesDB
	const queryHead = `select * from courses where `
	err := db.Select(&result, queryHead+queryWhere+queryLimit, selectArgs...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}
