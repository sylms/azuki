package main

import (
	"reflect"
	"testing"

	_ "github.com/lib/pq"
)

func Test_validateSearchCourseOptions(t *testing.T) {
	type args struct {
		courseName               string
		courseNameFilterType     string
		courseOverview           string
		courseOverviewFilterType string
		filterType               string
		limit                    string
		offset                   string
	}
	tests := []struct {
		name    string
		args    args
		want    searchCourseOptions
		wantErr bool
	}{
		{
			name: "exist all parameter",
			args: args{
				courseName:               "情報",
				courseNameFilterType:     "and",
				courseOverview:           "科学",
				courseOverviewFilterType: "and",
				filterType:               "and",
				limit:                    "100",
				offset:                   "50",
			},
			want: searchCourseOptions{
				courseName:               "情報",
				courseNameFilterType:     "and",
				courseOverview:           "科学",
				courseOverviewFilterType: "and",
				filterType:               "and",
				limit:                    100,
				offset:                   50,
			},
			wantErr: false,
		},
		{
			name: "courseName and courseOverview are empty",
			args: args{
				courseName:               "",
				courseNameFilterType:     "and",
				courseOverview:           "",
				courseOverviewFilterType: "and",
				filterType:               "and",
				limit:                    "100",
				offset:                   "50",
			},
			want:    searchCourseOptions{},
			wantErr: true,
		},
		{
			name: "courseNameFilterType: invalid",
			args: args{
				courseName:               "情報",
				courseNameFilterType:     "andandand",
				courseOverview:           "科学",
				courseOverviewFilterType: "and",
				filterType:               "and",
				limit:                    "100",
				offset:                   "50",
			},
			want:    searchCourseOptions{},
			wantErr: true,
		},
		{
			name: "courseOverviewFilterType: invalid",
			args: args{
				courseName:               "情報",
				courseNameFilterType:     "and",
				courseOverview:           "科学",
				courseOverviewFilterType: "andandand",
				filterType:               "and",
				limit:                    "100",
				offset:                   "50",
			},
			want:    searchCourseOptions{},
			wantErr: true,
		},
		{
			name: "filterType: invalid",
			args: args{
				courseName:               "情報",
				courseNameFilterType:     "and",
				courseOverview:           "科学",
				courseOverviewFilterType: "and",
				filterType:               "andandand",
				limit:                    "100",
				offset:                   "50",
			},
			want:    searchCourseOptions{},
			wantErr: true,
		},
		{
			name: "limit is text",
			args: args{
				courseName:               "情報",
				courseNameFilterType:     "and",
				courseOverview:           "科学",
				courseOverviewFilterType: "and",
				filterType:               "andandand",
				limit:                    "fake",
				offset:                   "50",
			},
			want:    searchCourseOptions{},
			wantErr: true,
		},
		{
			name: "limit < 0",
			args: args{
				courseName:               "情報",
				courseNameFilterType:     "and",
				courseOverview:           "科学",
				courseOverviewFilterType: "and",
				filterType:               "andandand",
				limit:                    "-50",
				offset:                   "50",
			},
			want:    searchCourseOptions{},
			wantErr: true,
		},
		{
			name: "offset is text",
			args: args{
				courseName:               "情報",
				courseNameFilterType:     "and",
				courseOverview:           "科学",
				courseOverviewFilterType: "and",
				filterType:               "andandand",
				limit:                    "100",
				offset:                   "fake",
			},
			want:    searchCourseOptions{},
			wantErr: true,
		},
		{
			name: "offset < 0",
			args: args{
				courseName:               "情報",
				courseNameFilterType:     "and",
				courseOverview:           "科学",
				courseOverviewFilterType: "and",
				filterType:               "andandand",
				limit:                    "100",
				offset:                   "-50",
			},
			want:    searchCourseOptions{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateSearchCourseOptions(tt.args.courseName, tt.args.courseNameFilterType, tt.args.courseOverview, tt.args.courseOverviewFilterType, tt.args.filterType, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateSearchCourseOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateSearchCourseOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
