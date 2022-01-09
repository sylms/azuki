package persistence

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sylms/azuki/domain"
	"github.com/sylms/azuki/testutils"
)

func Test_coursePersistence_Search(t *testing.T) {
	db, err := testutils.CreateDB()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		query domain.CourseQuery
	}
	tests := []struct {
		name    string
		fields  coursePersistence
		args    args
		want    []*domain.Course
		wantErr bool
	}{
		{
			name: "temp",
			fields: coursePersistence{
				db: db,
			},
			args: args{
				query: domain.CourseQuery{
					CourseName: "情報社会と法制度",
					Limit:      50,
				},
			},
			want: []*domain.Course{
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields
			got, err := p.Search(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("coursePersistence.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want, cmpopts.IgnoreFields(domain.Course{}, "CSVUpdatedAt", "CreatedAt", "UpdatedAt")); diff != "" {
				t.Errorf("coursePersistence.Search() mismatch: (-got +want)\n%s", diff)
			}
		})
	}
}
