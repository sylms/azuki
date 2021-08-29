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
func validateSearchCourseOptions(courseName string, courseNameFilterType string, courseOverview string, courseOverviewFilterType string, filterType string, limit string, offset string) (searchCourseOptions, error) {
	allowedFilterType := []string{filterTypeAnd, filterTypeOr}
	if !util.Contains(allowedFilterType, filterType) {
		return searchCourseOptions{}, fmt.Errorf("filterType error: %s, %+v", filterType, allowedFilterType)
	}
	if courseName != "" && !util.Contains(allowedFilterType, courseNameFilterType) {
		return searchCourseOptions{}, fmt.Errorf("courseNameFilterType error: %s, %+v", courseNameFilterType, allowedFilterType)
	}
	if courseOverview != "" && !util.Contains(allowedFilterType, courseOverviewFilterType) {
		return searchCourseOptions{}, fmt.Errorf("courseOverviewFilterType error: %s, %+v", courseOverviewFilterType, allowedFilterType)
	}

	// どのカラムも検索対象としていなければ検索そのものが実行できないので、不正なリクエストである
	if courseName == "" && courseOverview == "" {
		return searchCourseOptions{}, errors.New("course_name and course_overview are empty")
	}

	var limitInt int
	if limit == "" {
		limitInt = searchQueryDefaultLimit
	} else {
		var err error
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			return searchCourseOptions{}, errors.New("limit is not int")
		}
		if limitInt < 0 {
			return searchCourseOptions{}, errors.New("limit is negative")
		}
	}

	var offsetInt int
	if offset == "" {
		// offset = 0 であれば offset を指定しないときと同じ結果を得られる
		offsetInt = 0
	} else {
		var err error
		offsetInt, err = strconv.Atoi(offset)
		if err != nil {
			return searchCourseOptions{}, errors.New("offset is not int")
		}
		if offsetInt < 0 {
			return searchCourseOptions{}, errors.New("offset is negative")
		}
	}

	return searchCourseOptions{
		courseName:               courseName,
		courseNameFilterType:     courseNameFilterType,
		courseOverview:           courseOverview,
		courseOverviewFilterType: courseOverviewFilterType,
		filterType:               filterType,
		limit:                    limitInt,
		offset:                   offsetInt,
	}, nil
}
