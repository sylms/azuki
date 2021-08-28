package main

import (
	"reflect"
	"testing"
)

func Test_buildSearchCourseQuery(t *testing.T) {
	type args struct {
		options searchCourseOptions
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []interface{}
		wantErr bool
	}{
		{
			name: "exist all parameter",
			args: args{
				options: searchCourseOptions{
					courseName:               "情報",
					courseNameFilterType:     "and",
					courseOverview:           "科学",
					courseOverviewFilterType: "and",
					filterType:               "and",
					limit:                    100,
					offset:                   50,
				},
			},
			want: `select * from courses where course_name like $1 or course_overview like $2 limit $3 offset $4`,
			want: `select * from courses where course_name like $1 and course_overview like $2 order by id asc limit $3 offset $4`,
			want1: []interface{}{
				`%情報%`,
				`%科学%`,
				"100",
				"50",
			},
			wantErr: false,
		},
		{
			name: "empty: courseName",
			args: args{
				options: searchCourseOptions{
					courseName:               "",
					courseNameFilterType:     "and",
					courseOverview:           "科学",
					courseOverviewFilterType: "and",
					filterType:               "and",
					limit:                    100,
					offset:                   50,
				},
			},
			want: `select * from courses where course_overview like $1 order by id asc limit $2 offset $3`,
			want1: []interface{}{
				`%科学%`,
				"100",
				"50",
			},
			wantErr: false,
		},
		{
			name: "empty: course_overview",
			args: args{
				options: searchCourseOptions{
					courseName:               "情報",
					courseNameFilterType:     "and",
					courseOverview:           "",
					courseOverviewFilterType: "and",
					filterType:               "and",
					limit:                    100,
					offset:                   50,
				},
			},
			want: `select * from courses where course_name like $1 order by id asc limit $2 offset $3`,
			want1: []interface{}{
				`%情報%`,
				"100",
				"50",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := buildSearchCourseQuery(tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildSearchCourseQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("buildSearchCourseQuery() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("buildSearchCourseQuery() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
