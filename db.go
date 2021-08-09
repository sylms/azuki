package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sylms/azuki/util"
)

const (
	filterTypeAnd = "and"
	filterTypeOr  = "or"
)

type searchCourseOptions struct {
	courseName     string
	courseOverview string
	filterType     string
	limit          int
}

func searchCourse(options searchCourseOptions) ([]CoursesDB, error) {
	allowedFilterType := []string{filterTypeAnd, filterTypeOr}
	if !util.Contains(allowedFilterType, options.filterType) {
		return []CoursesDB{}, fmt.Errorf("filterType error: %s, %+v", options.filterType, allowedFilterType)
	}
	var result []CoursesDB
	const queryFormat = `select * from courses where course_name like $1 %s course_overview like $2 limit $3`
	queryFilterApplied := fmt.Sprintf(queryFormat, options.filterType)
	err := db.Select(&result, queryFilterApplied, "%"+options.courseName+"%", "%"+options.courseOverview+"%", options.limit)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}
