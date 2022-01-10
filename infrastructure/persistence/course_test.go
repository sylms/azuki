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
			name: "CourseName で一意になる検索をしてそれが得られる",
			fields: coursePersistence{
				db: db,
			},
			args: args{
				query: domain.CourseQuery{
					CourseName:           "情報社会と法制度",
					CourseNameFilterType: "and",
					Limit:                50,
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
		{
			name: "CourseOverview で一意になる検索をしてそれが得られる",
			fields: coursePersistence{
				db: db,
			},
			args: args{
				query: domain.CourseQuery{
					CourseOverview:           "知的財産に関する法制度を主要な概念や法理に基づいて学ぶ。著作権法、特許法を中心に、不正競争防止法、商標法など、知的財産諸法についての基礎的な知識を身につけ、知的財産法の法技術的な特色を踏まえた上で、情報化社会における望ましい制度のあり方について考察し、情報の保護と利用についてのバランス感覚や、問題解決能力を身につけることを目的とする。",
					CourseOverviewFilterType: "and",
					Limit:                    50,
				},
			},
			want: []*domain.Course{
				{
					ID:                       18011,
					CourseNumber:             "GA10201",
					CourseName:               "知的財産概論",
					InstructionalType:        1,
					Credits:                  "2.0",
					StandardRegistrationYear: []string{"2"},
					Term:                     []int{4, 5},
					Period:                   []string{"金5", "金6"},
					Classroom:                "",
					Instructor:               []string{"村井 麻衣子"},
					CourseOverview:           "知的財産に関する法制度を主要な概念や法理に基づいて学ぶ。著作権法、特許法を中心に、不正競争防止法、商標法など、知的財産諸法についての基礎的な知識を身につけ、知的財産法の法技術的な特色を踏まえた上で、情報化社会における望ましい制度のあり方について考察し、情報の保護と利用についてのバランス感覚や、問題解決能力を身につけることを目的とする。",
					Remarks:                  "オンライン(オンデマンド型)",
					CreditedAuditors:         0,
					ApplicationConditions:    "正規生に対しても受講制限をしているため",
					AltCourseName:            "Introduction to Intellectual Property",
					CourseCode:               "GA10201",
					CourseCodeName:           "知的財産概論",
					Year:                     2021,
				},
			},
		},
		{
			name: "CourseNumber で一意になる検索をしてそれが得られる",
			fields: coursePersistence{
				db: db,
			},
			args: args{
				query: domain.CourseQuery{
					CourseNumber: "GA10101",
					Limit:        50,
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
		{
			name: "Period で検索",
			fields: coursePersistence{
				db: db,
			},
			args: args{
				query: domain.CourseQuery{
					// csv2sql/kdb で "木3" と "木4" に分割される
					// "木3"、"木4" のどちらかが含まれている科目を検索する
					// -> 空いている時間割から科目を検索したい需要
					Period: "木3木4",
					Limit:  50,
				},
			},
			want: []*domain.Course{
				{
					ID:                       18020,
					CourseNumber:             "GA14201",
					CourseName:               "知識情報システム概説",
					InstructionalType:        1,
					Credits:                  "1.0",
					StandardRegistrationYear: []string{"1"},
					Term:                     []int{2, 3},
					Period:                   []string{"木4"},
					Classroom:                "",
					Instructor:               []string{"高久 雅生", "佐藤 哲司", "阪口 哲男", "鈴木 伸崇"},
					CourseOverview:           "ネットワーク社会における知識の構造化、提供、共有のための枠組みについて講義する。",
					Remarks:                  "専門導入科目(事前登録対象) オンライン(オンデマンド型)",
					CreditedAuditors:         0,
					ApplicationConditions:    "正規生に対しても受講制限をしているため",
					AltCourseName:            "Foundations of Knowledge Information Systems",
					CourseCode:               "GA14201",
					CourseCodeName:           "知識情報システム概説",
					Year:                     2021,
				},
				{
					ID:                       18021,
					CourseNumber:             "GA14301",
					CourseName:               "図書館概論",
					InstructionalType:        1,
					Credits:                  "2.0",
					StandardRegistrationYear: []string{"1"},
					Term:                     []int{4, 5},
					Period:                   []string{"木3", "木4"},
					Classroom:                "",
					Instructor:               []string{"吉田 右子"},
					CourseOverview:           "図書館とは何かについて概説し、これからの図書館の在り方を考える。図書館の歴史と現状、機能と社会的意義、館種別図書館と利用者、図書館職員、類縁機関と関係団体、図書館の課題と展望等について幅広く学ぶ。",
					Remarks:                  "専門導入科目(事前登録対象) オンライン(オンデマンド型) GE22001「図書館概論」を修得済みの者は履修不可。",
					CreditedAuditors:         1,
					ApplicationConditions:    "本学(学群・大学院)卒業・修了者又は本学の大学院在学者で司書・司書教諭資格希望者に限る",
					AltCourseName:            "Introduction to Librarianship",
					CourseCode:               "GA14301",
					CourseCodeName:           "図書館概論",
					Year:                     2021,
				},
				{
					ID:                       18060,
					CourseNumber:             "GB11404",
					CourseName:               "電磁気学",
					InstructionalType:        4,
					Credits:                  "2.0",
					StandardRegistrationYear: []string{"2"},
					Term:                     []int{4, 5},
					Period:                   []string{"木3", "木4"},
					Classroom:                "3A306",
					Instructor:               []string{"安永 守利"},
					CourseOverview:           "集積回路(IC)やハードディスク,タッチパネルや無線LANなど,我々の身の回りの情報通信機器は,電磁現象を原理として動作している.本講義では,これらの電磁現象の基礎を解説する.講義の前半では,「電荷」からスタートして「電場」,「電位」という場の概念とポテンシャルの概念を解説する.また,これらの現象を利用した応用事例も紹介する.後半では,はじめに磁気現象の本質は電流であることを説明し,「磁場」の概念,および「電磁誘導」等の電流と磁気現象の関係を解説する.また,磁気現象を利用した応用事例も紹介する.最後に,「電場」と「磁場」がマクスウェル方程式としてまとめられることを示し,「電磁波」の導出とその応用事例について言及する.",
					Remarks:                  "オンライン(オンデマンド型)",
					CreditedAuditors:         2,
					ApplicationConditions:    "",
					AltCourseName:            "Electromagnetics",
					CourseCode:               "GB11404",
					CourseCodeName:           "電磁気学",
					Year:                     2021,
				},
			},
		},
		{
			name: "複数の Term で検索したらそれらすべてを開講時期と指定している科目を返す",
			fields: coursePersistence{
				db: db,
			},
			args: args{
				query: domain.CourseQuery{
					// TODO: 検索需要わからなくなった確認する
					// csv2sql/kdb で "春A" と "春B" と "春C" に分割される
					// クエリの開講時期のすべてが含まれている科目を検索する
					// -> 余裕のある開講時期をいくつか指定して科目を検索したい需要
					Term:  "春A春B",
					Limit: 50,
				},
			},
			want: []*domain.Course{
				{
					ID:                       18022,
					CourseNumber:             "GA15111",
					CourseName:               "情報数学A",
					InstructionalType:        1,
					Credits:                  "2.0",
					StandardRegistrationYear: []string{"1"},
					Term:                     []int{1, 2},
					Period:                   []string{"木5", "木6"},
					Classroom:                "3A203",
					Instructor:               []string{"西出 隆志", "亀山 幸義"},
					CourseOverview:           "本授業では,情報学の基礎となる数学的概念について学ぶ.その中でも特に重要な概念である集合,論理,写像,関係,グラフ等を取りあげ,その基礎的な事項について講義する.また,講義内容に対する理解を深めるため,演習も行う.",
					Remarks:                  "平成31年度以降入学の者に限る。情報科学類生は1・2クラスを対象とする。 オンライン(オンデマンド型) 定員を超過した場合は履修調整をする場合がある（情報科学類生および総合学域群生(情報科学類への移行希望者・学籍番号の下一桁が奇数)優先）。 ",
					ApplicationConditions:    "正規生に対しても受講制限をしているため",
					AltCourseName:            "Mathematics for Informatics A",
					CourseCode:               "GA15101",
					CourseCodeName:           "情報数学A",
					Year:                     2021,
				},
				{
					ID:                       18047,
					CourseNumber:             "GB10244",
					CourseName:               "線形代数B",
					InstructionalType:        4,
					Credits:                  "2.0",
					StandardRegistrationYear: []string{"2"},
					Term:                     []int{1, 2},
					Period:                   []string{"月1", "月2"},
					Classroom:                "3A207",
					Instructor:               []string{"山田 武志"},
					CourseOverview:           "線形代数の基礎。 内容:ベクトル空間,1次写像,核と像,内積空間,固有値・固有ベクトルと対角化",
					Remarks:                  "情報科学類3・4クラス対象 オンライン(オンデマンド型) 対面",
					CreditedAuditors:         2,
					AltCourseName:            "Linear Algebra B",
					CourseCode:               "GB10244",
					CourseCodeName:           "線形代数B",
					Year:                     2021,
				},
			},
		},
		{
			name: "CourseName の部分一致検索が動作する",
			fields: coursePersistence{
				db: db,
			},
			args: args{
				query: domain.CourseQuery{
					CourseName:           "分積",
					CourseNameFilterType: "and",
					Limit:                50,
				},
			},
			want: []*domain.Course{
				{
					ID:                       18030,
					CourseNumber:             "GA15311",
					CourseName:               "微分積分A",
					InstructionalType:        1,
					Credits:                  "2.0",
					StandardRegistrationYear: []string{"1"},
					Term:                     []int{4, 5},
					Period:                   []string{"金3", "金4"},
					Classroom:                "3B302",
					Instructor:               []string{"町田 文雄", "堀江 和正"},
					CourseOverview:           "解析学の基礎として,実数,関数,数列ならびに連続性や極限などの基本概念と,1変数関数の微分法および積分法について講義を行う。",
					Remarks:                  "情報科学類生は1・2クラスを対象とする。定員を超過した場合は履修調整をする場合がある（情報科学類生および総合学域 群生(情報科学類への移行希望者・学籍番号の下一桁が奇数)優先）。履修申請期 限は9月21日(火)まで。 オンライン(オンデマンド型) 平成30年度までに開設された「解析学I」(GB10314,GB10324)の単位を修得した者 の履修は認めない。",
					ApplicationConditions:    "正規生に対しても受講制限をしているため",
					AltCourseName:            "Calculus A",
					CourseCode:               "GA15301",
					CourseCodeName:           "微分積分A",
					Year:                     2021,
				},
				{
					ID:                       18033,
					CourseNumber:             "GA15341",
					CourseName:               "微分積分A",
					InstructionalType:        1,
					Credits:                  "2.0",
					StandardRegistrationYear: []string{"1"},
					Term:                     []int{4, 5},
					Period:                   []string{"金3", "金4"},
					Instructor:               []string{"加藤 誠"},
					CourseOverview:           "解析学の基礎として,実数,関数,数列ならびに連続性や極限などの基本概念と,1変数関数の微分法および積分法について講義を行う。",
					Remarks:                  "知識学類生および総合学域群生（知識学類への移行希望者）優先。 履修申請期限は9月21日(火)まで。 定員を超過した場合は履修調整をする場合がある 。 オンライン(オンデマンド型)",
					ApplicationConditions:    "正規生に対しても受講制限をしているため",
					AltCourseName:            "Calculus A",
					CourseCode:               "GA15301",
					CourseCodeName:           "微分積分A",
					Year:                     2021,
				},
			},
		},
		{
			name: "CourseOverview の部分一致検索が動作する",
			fields: coursePersistence{
				db: db,
			},
			args: args{
				query: domain.CourseQuery{
					CourseOverview:           "アルゴリズム",
					CourseOverviewFilterType: "and",
					Limit:                    50,
				},
			},
			want: []*domain.Course{
				{
					ID:                       18066,
					CourseNumber:             "GB11931",
					CourseName:               "データ構造とアルゴリズム",
					InstructionalType:        1,
					Credits:                  "3.0",
					StandardRegistrationYear: []string{"2"},
					Term:                     []int{4, 5, 6},
					Period:                   []string{"月1", "月2"},
					Classroom:                "3B402",
					Instructor:               []string{"天笠 俊之", "長谷部 浩二", "藤田 典久"},
					CourseOverview:           "ソフトウェアを書く上で基本となるデータ構造とアルゴリズムの考え方について学ぶ。線形構造,木構造,グラフ構造,データ整列,データ探索について学習する。",
					Remarks:                  "平成25年度までに開設された「データ構造とアルゴリズム」(GB11911, GB11921)の単位を修得した者の履修は認めない。 オンライン(同時双方向型)",
					CreditedAuditors:         2,
					AltCourseName:            "Data Structures and Algorithms",
					CourseCode:               "GB11931",
					CourseCodeName:           "データ構造とアルゴリズム",
					Year:                     2021,
				},
				{
					ID:                       18067,
					CourseNumber:             "GB11956",
					CourseName:               "データ構造とアルゴリズム実験",
					InstructionalType:        6,
					Credits:                  "2.0",
					StandardRegistrationYear: []string{"2"},
					Term:                     []int{4, 5, 6},
					Period:                   []string{"月3", "月4", "月5", "月3", "月4"},
					Classroom:                "3C113,3C205",
					Instructor:               []string{"天笠 俊之"},
					CourseOverview:           "データ構造とアルゴリズムに関して,実際にJava言語を用いてプログラムを作成し,そのプログラムが稼働することを確認する。プログラムは,毎週,あるいは隔週に一個の割合で作成する。",
					Remarks:                  "1・2クラス オンライン(同時双方向型) 令和2年度までに開設された「データ構造とアルゴリズム実験」(GB11936,GB11946)または平成26年度までに開設された「データ構造とアルゴリズム実験」(GB11916, GB11926)の単位を修得した者の履修は認めない。",
					ApplicationConditions:    "施設設備の許容量上の制約と学類生に対する良質の少人数教育を行うため",
					AltCourseName:            "Data Structures and Algorithms Laboratory",
					CourseCode:               "GB11956",
					CourseCodeName:           "データ構造とアルゴリズム実験",
					Year:                     2021,
				},
			},
		},
		{
			name: "スペース区切りで複数 CourseName を与え CourseNameFilterType = and がきちんと働く",
			fields: coursePersistence{
				db: db,
			},
			args: args{
				query: domain.CourseQuery{
					CourseName:           "システム 情報科学",
					CourseNameFilterType: "and",
					Limit:                50,
				},
			},
			want: []*domain.Course{
				{
					ID:                       18014,
					CourseNumber:             "GA12301",
					CourseName:               "システムと情報科学",
					InstructionalType:        1,
					Credits:                  "1.0",
					StandardRegistrationYear: []string{"1"},
					Term:                     []int{5},
					Period:                   []string{"火5", "火6"},
					Instructor:               []string{"山際 伸一", "山口 佳樹", "佐藤 聡", "西出 隆志", "大山 恵弘"},
					CourseOverview:           "情報科学への導入となる基礎理論から応用までを概説し、専門的科目への導入としての基礎知識を習得する。本科目は特に、システムを中心に専門性を習得する上での事前知識となる原理や技術、理論について説明する。",
					Remarks:                  "専門導入科目(事前登録対象) オンライン(オンデマンド型)",
					ApplicationConditions:    "正規生に対しても受講制限をしているため",
					AltCourseName:            "Introduction to Information Science:Information Systems",
					CourseCode:               "GA12301",
					CourseCodeName:           "システムと情報科学",
					Year:                     2021,
				},
			},
		},
		{
			name: "スペース区切りで複数 CourseName を与え CourseNameFilterType = or がきちんと働く",
			fields: coursePersistence{
				db: db,
			},
			args: args{
				query: domain.CourseQuery{
					CourseName:           "システム 情報科学",
					CourseNameFilterType: "or",
					Limit:                50,
				},
			},
			want: []*domain.Course{
				{
					ID:                       18014,
					CourseNumber:             "GA12301",
					CourseName:               "システムと情報科学",
					InstructionalType:        1,
					Credits:                  "1.0",
					StandardRegistrationYear: []string{"1"},
					Term:                     []int{5},
					Period:                   []string{"火5", "火6"},
					Instructor:               []string{"山際 伸一", "山口 佳樹", "佐藤 聡", "西出 隆志", "大山 恵弘"},
					CourseOverview:           "情報科学への導入となる基礎理論から応用までを概説し、専門的科目への導入としての基礎知識を習得する。本科目は特に、システムを中心に専門性を習得する上での事前知識となる原理や技術、理論について説明する。",
					Remarks:                  "専門導入科目(事前登録対象) オンライン(オンデマンド型)",
					ApplicationConditions:    "正規生に対しても受講制限をしているため",
					AltCourseName:            "Introduction to Information Science:Information Systems",
					CourseCode:               "GA12301",
					CourseCodeName:           "システムと情報科学",
					Year:                     2021,
				},
				{
					ID:                       18020,
					CourseNumber:             "GA14201",
					CourseName:               "知識情報システム概説",
					InstructionalType:        1,
					Credits:                  "1.0",
					StandardRegistrationYear: []string{"1"},
					Term:                     []int{2, 3},
					Period:                   []string{"木4"},
					Instructor:               []string{"高久 雅生", "佐藤 哲司", "阪口 哲男", "鈴木 伸崇"},
					CourseOverview:           "ネットワーク社会における知識の構造化、提供、共有のための枠組みについて講義する。",
					Remarks:                  "専門導入科目(事前登録対象) オンライン(オンデマンド型)",
					ApplicationConditions:    "正規生に対しても受講制限をしているため",
					AltCourseName:            "Foundations of Knowledge Information Systems",
					CourseCode:               "GA14201",
					CourseCodeName:           "知識情報システム概説",
					Year:                     2021,
				},
			},
		},
		// TODO: CourseOverview + CourseOverviewFilterType = {and,or} の確認
		// TODO: CourseName, CourseOverview たちをまとめる FilterType の確認
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

func Test_coursePersistence_Facet(t *testing.T) {
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
		want    []*domain.Facet
		wantErr bool
	}{
		{
			name: "temp",
			fields: coursePersistence{
				db: db,
			},
			args: args{
				query: domain.CourseQuery{
					CourseName:           "情報",
					CourseNameFilterType: "and",
					Limit:                50,
				},
			},
			want: []*domain.Facet{
				{
					Term:      1,
					TermCount: 1,
				},
				{
					Term:      5,
					TermCount: 2,
				},
				{
					Term:      2,
					TermCount: 2,
				},
				{
					Term:      4,
					TermCount: 1,
				},
				{
					Term:      3,
					TermCount: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields
			got, err := p.Facet(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("coursePersistence.Facet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("coursePersistence.Facet() mismatch: (-got +want)\n%s", diff)
			}
		})
	}
}
