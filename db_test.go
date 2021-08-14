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
				},
			},
			want: `select * from courses where course_name like $1 and course_overview like $2 limit $3`,
			want1: []interface{}{
				`%情報%`,
				`%科学%`,
				"100",
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
				},
			},
			want: `select * from courses where course_name like $1 and course_overview like $2 limit $3`,
			want1: []interface{}{
				`%%`,
				`%科学%`,
				"100",
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
				},
			},
			want: `select * from courses where course_name like $1 and course_overview like $2 limit $3`,
			want1: []interface{}{
				`%情報%`,
				`%%`,
				"100",
			},
			wantErr: false,
		},
		{
			name: "limit ok?",
			args: args{
				options: searchCourseOptions{
					courseName:               "情報",
					courseNameFilterType:     "and",
					courseOverview:           "科学",
					courseOverviewFilterType: "and",
					filterType:               "and",
					limit:                    1234,
				},
			},
			want: `select * from courses where course_name like $1 and course_overview like $2 limit $3`,
			want1: []interface{}{
				`%情報%`,
				`%科学%`,
				"1234",
			},
			wantErr: false,
		},
		{
			name: "courseNameFilterType: invalid",
			args: args{
				options: searchCourseOptions{
					courseName:               "情報",
					courseNameFilterType:     "andandand",
					courseOverview:           "科学",
					courseOverviewFilterType: "and",
					filterType:               "and",
					limit:                    100,
				},
			},
			want:    ``,
			want1:   nil,
			wantErr: true,
		},
		{
			name: "courseOverviewFilterType: invalid",
			args: args{
				options: searchCourseOptions{
					courseName:               "情報",
					courseNameFilterType:     "and",
					courseOverview:           "科学",
					courseOverviewFilterType: "andandand",
					filterType:               "and",
					limit:                    100,
				},
			},
			want:    ``,
			want1:   nil,
			wantErr: true,
		},
		{
			name: "filterType: invalid",
			args: args{
				options: searchCourseOptions{
					courseName:               "情報",
					courseNameFilterType:     "and",
					courseOverview:           "科学",
					courseOverviewFilterType: "and",
					filterType:               "andand",
					limit:                    100,
				},
			},
			want:    ``,
			want1:   nil,
			wantErr: true,
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
