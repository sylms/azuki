package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sylms/azuki/domain"
)

type courseUseCaseMock struct {
	domain.Course
	FakeSearch func(domain.CourseQuery) ([]*domain.Course, error)
	FakeFacet  func(domain.CourseQuery) ([]*domain.Facet, error)
}

func (uc *courseUseCaseMock) Search(query domain.CourseQuery) ([]*domain.Course, error) {
	return uc.FakeSearch(query)
}

func (uc *courseUseCaseMock) Facet(query domain.CourseQuery) ([]*domain.Facet, error) {
	return uc.FakeFacet(query)
}

func Test_courseHandler_Search(t *testing.T) {
	type fakeSearch struct {
		Search func(domain.CourseQuery) ([]*domain.Course, error)
	}
	tests := []struct {
		name                 string
		fakeSearch           fakeSearch
		reqContentTypeHeader string
		reqBody              string
		wantResStatusCode    int
		wantResBody          string
	}{
		{
			name: "normal",
			fakeSearch: fakeSearch{
				Search: func(cq domain.CourseQuery) ([]*domain.Course, error) {
					courses := []*domain.Course{
						{
							ID:                       18010,
							CourseNumber:             "GA10101",
							CourseName:               "情報社会と法制度",
							InstructionalType:        1,
							Credits:                  "2.0",
							StandardRegistrationYear: []string{"2"},
							Term:                     []int{4, 5},
							Period:                   []string{"月5", "月6"},
							Classroom:                "",
							Instructor:               []string{"髙良 幸哉"},
							CourseOverview:           "情報化社会における法制度や情報モラル向上に必要な基礎知識を習得することを目指すため、現行の我が国の法制度の基礎を学び、ネットワーク社会における法整備の現状について講義する。",
							Remarks:                  "オンライン(オンデマンド型)",
							CreditedAuditors:         0,
							ApplicationConditions:    "正規生に対しても受講制限をしているため",
							AltCourseName:            "Information Society Law",
							CourseCode:               "GA10101",
							CourseCodeName:           "情報社会と法制度",
							Year:                     2021,
						},
					}
					return courses, nil
				},
			},
			reqContentTypeHeader: "application/json",
			reqBody: `{
		    "course_number": "GA10101",
		    "course_name": "情報社会と法制度",
		    "instructional_type": -1,
		    "credits": "",
		    "standard_registration_year": -1,
		    "term": "",
		    "period": "",
		    "classroom": "",
		    "instructor": "",
		    "course_overview": "",
		    "remarks": "",
		    "course_name_filter_type": "and",
		    "course_overview_filter_type": "and",
		    "filter_type": "and",
		    "limit": 20,
		    "offset": 0
		}`,
			wantResStatusCode: http.StatusOK,
			wantResBody:       `[{"id":18010,"course_number":"GA10101","course_name":"情報社会と法制度","instructional_type":1,"credits":"2.0","standard_registration_year":["2"],"term":[4,5],"period":["月5","月6"],"classroom":"","instructor":["髙良 幸哉"],"course_overview":"情報化社会における法制度や情報モラル向上に必要な基礎知識を習得することを目指すため、現行の我が国の法制度の基礎を学び、ネットワーク社会における法整備の現状について講義する。","remarks":"オンライン(オンデマンド型)","credited_auditors":0,"application_conditions":"正規生に対しても受講制限をしているため","alt_course_name":"Information Society Law","course_code":"GA10101","course_code_name":"情報社会と法制度","csv_updated_at":"0001-01-01T00:00:00Z","year":2021,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name: "検索条件に該当する科目が存在しないときに空配列を表す JSON 文字列を返す",
			fakeSearch: fakeSearch{
				Search: func(cq domain.CourseQuery) ([]*domain.Course, error) {
					courses := []*domain.Course{}
					return courses, nil
				},
			},
			reqContentTypeHeader: "application/json",
			reqBody: `{
		    "course_number": "GA10101",
		    "course_name": "情報社会と法制度",
		    "instructional_type": -1,
		    "credits": "",
		    "standard_registration_year": -1,
		    "term": "",
		    "period": "",
		    "classroom": "",
		    "instructor": "",
		    "course_overview": "",
		    "remarks": "",
		    "course_name_filter_type": "and",
		    "course_overview_filter_type": "and",
		    "filter_type": "and",
		    "limit": 20,
		    "offset": 0
		}`,
			wantResStatusCode: http.StatusOK,
			wantResBody:       `[]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/courses", bytes.NewBufferString(tt.reqBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", tt.reqContentTypeHeader)

			res := httptest.NewRecorder()

			h := &courseHandler{
				uc: &courseUseCaseMock{
					FakeSearch: tt.fakeSearch.Search,
				},
			}

			h.Search(res, req)

			resBodyGot := res.Body.String()
			if resBodyGot != tt.wantResBody {
				t.Errorf("response mismatch:\ngot: %s\nwant: %s", resBodyGot, tt.wantResBody)
			}

			statusCodeGot := res.Code
			if statusCodeGot != tt.wantResStatusCode {
				t.Errorf("response status code mismatch:\ngot: %d\nwant: %d", statusCodeGot, tt.wantResStatusCode)
			}
		})
	}
}

func Test_courseHandler_Csv(t *testing.T) {
	type fakeSearch struct {
		Search func(domain.CourseQuery) ([]*domain.Course, error)
	}
	tests := []struct {
		name                 string
		fakeSearch           fakeSearch
		reqContentTypeHeader string
		reqBody              string
		wantResStatusCode    int
		wantResBody          string
	}{
		{
			name: "temp",
			fakeSearch: fakeSearch{
				Search: func(cq domain.CourseQuery) ([]*domain.Course, error) {
					courses := []*domain.Course{
						{
							ID:                       18010,
							CourseNumber:             "GA10101",
							CourseName:               "情報社会と法制度",
							InstructionalType:        1,
							Credits:                  "2.0",
							StandardRegistrationYear: []string{"2"},
							Term:                     []int{4, 5},
							Period:                   []string{"月5", "月6"},
							Classroom:                "",
							Instructor:               []string{"髙良 幸哉"},
							CourseOverview:           "情報化社会における法制度や情報モラル向上に必要な基礎知識を習得することを目指すため、現行の我が国の法制度の基礎を学び、ネットワーク社会における法整備の現状について講義する。",
							Remarks:                  "オンライン(オンデマンド型)",
							CreditedAuditors:         0,
							ApplicationConditions:    "正規生に対しても受講制限をしているため",
							AltCourseName:            "Information Society Law",
							CourseCode:               "GA10101",
							CourseCodeName:           "情報社会と法制度",
							Year:                     2021,
						},
					}
					return courses, nil
				},
			},
			reqContentTypeHeader: "application/json",
			reqBody: `{
		    "course_number": "GA10101",
		    "course_name": "情報社会と法制度",
		    "instructional_type": -1,
		    "credits": "",
		    "standard_registration_year": -1,
		    "term": "",
		    "period": "",
		    "classroom": "",
		    "instructor": "",
		    "course_overview": "",
		    "remarks": "",
		    "course_name_filter_type": "and",
		    "course_overview_filter_type": "and",
		    "filter_type": "and",
		    "limit": 20,
		    "offset": 0
		}`,
			wantResStatusCode: http.StatusOK,
			wantResBody: `科目番号,科目名,授業方法,単位数,標準履修年次,実施学期,曜時限,教室,担当教員,授業概要,備考,科目等履修生申請可否,申請条件,英語(日本語)科目名,科目コード,要件科目名,データ更新日
GA10101,情報社会と法制度,1,2.0,2,秋A秋B,"月5,月6",,髙良 幸哉,情報化社会における法制度や情報モラル向上に必要な基礎知識を習得することを目指すため、現行の我が国の法制度の基礎を学び、ネットワーク社会における法整備の現状について講義する。,オンライン(オンデマンド型),0,正規生に対しても受講制限をしているため,Information Society Law,GA10101,情報社会と法制度,0001-01-01T00:00:00Z
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/csv", bytes.NewBufferString(tt.reqBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", tt.reqContentTypeHeader)

			res := httptest.NewRecorder()

			h := &courseHandler{
				uc: &courseUseCaseMock{
					FakeSearch: tt.fakeSearch.Search,
				},
			}

			h.Csv(res, req)

			resBodyGot := res.Body.String()
			if resBodyGot != tt.wantResBody {
				t.Errorf("response mismatch:\ngot: %s\nwant: %s", resBodyGot, tt.wantResBody)
			}

			statusCodeGot := res.Code
			if statusCodeGot != tt.wantResStatusCode {
				t.Errorf("response status code mismatch:\ngot: %d\nwant: %d", statusCodeGot, tt.wantResStatusCode)
			}
		})
	}
}

func Test_courseHandler_Facet(t *testing.T) {
	t.Run("temp", func(t *testing.T) {
		want := `{"term_facet":{"1":111,"2":222}}`
		reqBody := `{
		    "course_number": "GA10101",
		    "course_name": "情報社会と法制度",
		    "instructional_type": -1,
		    "credits": "",
		    "standard_registration_year": -1,
		    "term": "",
		    "period": "",
		    "classroom": "",
		    "instructor": "",
		    "course_overview": "",
		    "remarks": "",
		    "course_name_filter_type": "and",
		    "course_overview_filter_type": "and",
		    "filter_type": "and",
		    "limit": 20,
		    "offset": 0
		}`

		req, err := http.NewRequest(http.MethodPost, "/facet", bytes.NewBufferString(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		res := httptest.NewRecorder()

		h := &courseHandler{
			uc: &courseUseCaseMock{
				FakeFacet: func(cq domain.CourseQuery) ([]*domain.Facet, error) {
					courses := []*domain.Facet{
						{
							Term:      1,
							TermCount: 111,
						},
						{
							Term:      2,
							TermCount: 222,
						},
					}
					return courses, nil
				},
			},
		}
		h.Facet(res, req)

		got := res.Body.String()
		if got != want {
			t.Errorf("response mismatch:\ngot:\n%s\nwant:\n%s", got, want)
		}
	})
}

func Test_validateSearchCourseQuery(t *testing.T) {
	type args struct {
		query domain.CourseQuery
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "exist all parameter",
			args: args{
				query: domain.CourseQuery{
					CourseName:               "情報",
					CourseNameFilterType:     "and",
					CourseOverview:           "科学",
					CourseOverviewFilterType: "and",
					FilterType:               "and",
					Limit:                    100,
					Offset:                   50,
				},
			},
			wantErr: false,
		},
		{
			name: "一部パラメーターのみ存在する(つまり (全て存在する または 全て存在しない)ではない)",
			args: args{
				query: domain.CourseQuery{
					CourseName:               "",
					CourseNameFilterType:     "and",
					CourseOverview:           "科学",
					CourseOverviewFilterType: "and",
					FilterType:               "and",
					Limit:                    100,
					Offset:                   50,
				},
			},
			wantErr: false,
		},
		{
			name: "一部パラメーターのみ存在する(つまり (全て存在する または 全て存在しない)ではない)",
			args: args{
				query: domain.CourseQuery{
					CourseName:               "情報",
					CourseNameFilterType:     "and",
					CourseOverview:           "",
					CourseOverviewFilterType: "and",
					FilterType:               "and",
					Limit:                    100,
					Offset:                   50,
				},
			},
			wantErr: false,
		},
		{
			name: "パラメーターが存在していないもののフィルタータイプが異常(エラーは発生しない)",
			args: args{
				query: domain.CourseQuery{
					CourseName:               "",
					CourseNameFilterType:     "fake",
					CourseOverview:           "科学",
					CourseOverviewFilterType: "and",
					FilterType:               "and",
					Limit:                    100,
					Offset:                   50,
				},
			},
			wantErr: false,
		},
		{
			name: "パラメーターが存在していないもののフィルタータイプが異常(エラーは発生しない)",
			args: args{
				query: domain.CourseQuery{
					CourseName:               "情報",
					CourseNameFilterType:     "and",
					CourseOverview:           "",
					CourseOverviewFilterType: "fake",
					FilterType:               "and",
					Limit:                    100,
					Offset:                   50,
				},
			},
			wantErr: false,
		},
		{
			name: "all str is empty",
			args: args{
				query: domain.CourseQuery{
					CourseName:               "",
					CourseNameFilterType:     "and",
					CourseOverview:           "",
					CourseOverviewFilterType: "and",
					FilterType:               "and",
					Limit:                    100,
					Offset:                   50,
				},
			},
			wantErr: false,
		},
		{
			name: "cause FilterType error",
			args: args{
				query: domain.CourseQuery{
					CourseName:               "情報",
					CourseNameFilterType:     "and",
					CourseOverview:           "科学",
					CourseOverviewFilterType: "and",
					FilterType:               "andor",
					Limit:                    100,
					Offset:                   50,
				},
			},
			wantErr: true,
		},
		{
			name: "cause FilterType is empty error",
			args: args{
				query: domain.CourseQuery{
					CourseName:               "情報",
					CourseNameFilterType:     "and",
					CourseOverview:           "科学",
					CourseOverviewFilterType: "and",
					FilterType:               "",
					Limit:                    100,
					Offset:                   50,
				},
			},
			wantErr: true,
		},
		{
			name: "cause CourseNameFilterType error",
			args: args{
				query: domain.CourseQuery{
					CourseName:               "情報",
					CourseNameFilterType:     "andor",
					CourseOverview:           "科学",
					CourseOverviewFilterType: "and",
					FilterType:               "and",
					Limit:                    100,
					Offset:                   50,
				},
			},
			wantErr: true,
		},
		{
			name: "cause CourseNameFilterType is empty error",
			args: args{
				query: domain.CourseQuery{
					CourseName:               "情報",
					CourseNameFilterType:     "",
					CourseOverview:           "科学",
					CourseOverviewFilterType: "and",
					FilterType:               "and",
					Limit:                    100,
					Offset:                   50,
				},
			},
			wantErr: true,
		},
		{
			name: "cause CourseOverviewFilterType error",
			args: args{
				query: domain.CourseQuery{
					CourseName:               "情報",
					CourseNameFilterType:     "and",
					CourseOverview:           "科学",
					CourseOverviewFilterType: "andor",
					FilterType:               "and",
					Limit:                    100,
					Offset:                   50,
				},
			},
			wantErr: true,
		},
		{
			name: "cause CourseOverviewFilterType is empty error",
			args: args{
				query: domain.CourseQuery{
					CourseName:               "情報",
					CourseNameFilterType:     "and",
					CourseOverview:           "科学",
					CourseOverviewFilterType: "",
					FilterType:               "and",
					Limit:                    100,
					Offset:                   50,
				},
			},
			wantErr: true,
		},
		{
			name: "limit is negative",
			args: args{
				query: domain.CourseQuery{
					CourseName:               "情報",
					CourseNameFilterType:     "and",
					CourseOverview:           "科学",
					CourseOverviewFilterType: "and",
					FilterType:               "and",
					Limit:                    -10,
					Offset:                   50,
				},
			},
			wantErr: true,
		},
		{
			name: "offset is negative",
			args: args{
				query: domain.CourseQuery{
					CourseName:               "情報",
					CourseNameFilterType:     "and",
					CourseOverview:           "科学",
					CourseOverviewFilterType: "and",
					FilterType:               "and",
					Limit:                    100,
					Offset:                   -100,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSearchCourseQuery(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateSearchCourseQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
