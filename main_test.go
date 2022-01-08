package main

import (
	"testing"

	_ "github.com/lib/pq"
)

func Test_validateSearchCourseOptions(t *testing.T) {
	tests := []struct {
		name    string
		args    CourseQuery
		wantErr bool
	}{
		{
			name: "exist all parameter",
			args: CourseQuery{
				CourseName:               "情報",
				CourseNameFilterType:     "and",
				CourseOverview:           "科学",
				CourseOverviewFilterType: "and",
				FilterType:               "and",
				Limit:                    100,
				Offset:                   50,
			},
			wantErr: false,
		},
		{
			name: "一部パラメーターのみ存在する(つまり (全て存在する または 全て存在しない)ではない)",
			args: CourseQuery{
				CourseName:               "",
				CourseNameFilterType:     "and",
				CourseOverview:           "科学",
				CourseOverviewFilterType: "and",
				FilterType:               "and",
				Limit:                    100,
				Offset:                   50,
			},
			wantErr: false,
		},
		{
			name: "一部パラメーターのみ存在する(つまり (全て存在する または 全て存在しない)ではない)",
			args: CourseQuery{
				CourseName:               "情報",
				CourseNameFilterType:     "and",
				CourseOverview:           "",
				CourseOverviewFilterType: "and",
				FilterType:               "and",
				Limit:                    100,
				Offset:                   50,
			},
			wantErr: false,
		},
		{
			name: "パラメーターが存在していないもののフィルタータイプが異常(エラーは発生しない)",
			args: CourseQuery{
				CourseName:               "",
				CourseNameFilterType:     "fake",
				CourseOverview:           "科学",
				CourseOverviewFilterType: "and",
				FilterType:               "and",
				Limit:                    100,
				Offset:                   50,
			},
			wantErr: false,
		},
		{
			name: "パラメーターが存在していないもののフィルタータイプが異常(エラーは発生しない)",
			args: CourseQuery{
				CourseName:               "情報",
				CourseNameFilterType:     "and",
				CourseOverview:           "",
				CourseOverviewFilterType: "fake",
				FilterType:               "and",
				Limit:                    100,
				Offset:                   50,
			},
			wantErr: false,
		},
		{
			name: "all str is empty",
			args: CourseQuery{
				CourseName:               "",
				CourseNameFilterType:     "and",
				CourseOverview:           "",
				CourseOverviewFilterType: "and",
				FilterType:               "and",
				Limit:                    100,
				Offset:                   50,
			},
			wantErr: false,
		},
		{
			name: "cause FilterType error",
			args: CourseQuery{
				CourseName:               "情報",
				CourseNameFilterType:     "and",
				CourseOverview:           "科学",
				CourseOverviewFilterType: "and",
				FilterType:               "andor",
				Limit:                    100,
				Offset:                   50,
			},
			wantErr: true,
		},
		{
			name: "cause FilterType is empty error",
			args: CourseQuery{
				CourseName:               "情報",
				CourseNameFilterType:     "and",
				CourseOverview:           "科学",
				CourseOverviewFilterType: "and",
				FilterType:               "",
				Limit:                    100,
				Offset:                   50,
			},
			wantErr: true,
		},
		{
			name: "cause CourseNameFilterType error",
			args: CourseQuery{
				CourseName:               "情報",
				CourseNameFilterType:     "andor",
				CourseOverview:           "科学",
				CourseOverviewFilterType: "and",
				FilterType:               "and",
				Limit:                    100,
				Offset:                   50,
			},
			wantErr: true,
		},
		{
			name: "cause CourseNameFilterType is empty error",
			args: CourseQuery{
				CourseName:               "情報",
				CourseNameFilterType:     "",
				CourseOverview:           "科学",
				CourseOverviewFilterType: "and",
				FilterType:               "and",
				Limit:                    100,
				Offset:                   50,
			},
			wantErr: true,
		},
		{
			name: "cause CourseOverviewFilterType error",
			args: CourseQuery{
				CourseName:               "情報",
				CourseNameFilterType:     "and",
				CourseOverview:           "科学",
				CourseOverviewFilterType: "andor",
				FilterType:               "and",
				Limit:                    100,
				Offset:                   50,
			},
			wantErr: true,
		},
		{
			name: "cause CourseOverviewFilterType is empty error",
			args: CourseQuery{
				CourseName:               "情報",
				CourseNameFilterType:     "and",
				CourseOverview:           "科学",
				CourseOverviewFilterType: "",
				FilterType:               "and",
				Limit:                    100,
				Offset:                   50,
			},
			wantErr: true,
		},
		{
			name: "limit is negative",
			args: CourseQuery{
				CourseName:               "情報",
				CourseNameFilterType:     "and",
				CourseOverview:           "科学",
				CourseOverviewFilterType: "and",
				FilterType:               "and",
				Limit:                    -10,
				Offset:                   50,
			},
			wantErr: true,
		},
		{
			name: "offset is negative",
			args: CourseQuery{
				CourseName:               "情報",
				CourseNameFilterType:     "and",
				CourseOverview:           "科学",
				CourseOverviewFilterType: "and",
				FilterType:               "and",
				Limit:                    100,
				Offset:                   -100,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSearchCourseOptions(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateSearchCourseOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
