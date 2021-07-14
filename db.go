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

func getSyllabusByNameAndOverview(courseName string, filterType string, courseOverview string, limit int) ([]CoursesDB, error) {
	allowedFilterType := []string{filterTypeAnd, filterTypeOr}
	if !util.Contains(allowedFilterType, filterType) {
		return []CoursesDB{}, fmt.Errorf("filterType error: %s, %+v", filterType, allowedFilterType)
	}
	var result []CoursesDB
	const queryFormat = `select * from courses where course_name like $1 %s course_overview like $2 limit $3`
	queryFilterApplied := fmt.Sprintf(queryFormat, filterType)
	err := db.Select(&result, queryFilterApplied, "%"+courseName+"%", "%"+courseOverview+"%", limit)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}
