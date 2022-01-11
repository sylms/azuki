--
-- PostgreSQL database dump
--

-- Dumped from database version 12.7
-- Dumped by pg_dump version 12.7

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: credited_auditors; Type: TYPE; Schema: public; Owner: sylms
--

CREATE TYPE public.credited_auditors AS ENUM (
    '0',
    '1',
    '2'
);


ALTER TYPE public.credited_auditors OWNER TO sylms;

--
-- Name: instructional_type; Type: TYPE; Schema: public; Owner: sylms
--

CREATE TYPE public.instructional_type AS ENUM (
    '0',
    '1',
    '2',
    '3',
    '4',
    '5',
    '6',
    '7',
    '8'
);


ALTER TYPE public.instructional_type OWNER TO sylms;

--
-- Name: standard_registration_year; Type: TYPE; Schema: public; Owner: sylms
--

CREATE TYPE public.standard_registration_year AS ENUM (
    '?',
    '1',
    '2',
    '3',
    '4',
    '5',
    '6'
);


ALTER TYPE public.standard_registration_year OWNER TO sylms;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: courses; Type: TABLE; Schema: public; Owner: sylms
--

CREATE TABLE public.courses (
    id integer NOT NULL,
    course_number character varying(16) NOT NULL,
    course_name character varying(256) NOT NULL,
    instructional_type public.instructional_type NOT NULL,
    credits character varying(8) NOT NULL,
    standard_registration_year public.standard_registration_year[] NOT NULL,
    term integer[] NOT NULL,
    period_ character varying(16)[] NOT NULL,
    classroom character varying(256) NOT NULL,
    instructor character varying(256)[] NOT NULL,
    course_overview text NOT NULL,
    remarks text NOT NULL,
    credited_auditors public.credited_auditors NOT NULL,
    application_conditions character varying(256) NOT NULL,
    alt_course_name character varying(256) NOT NULL,
    course_code character varying(16) NOT NULL,
    course_code_name character varying(256) NOT NULL,
    csv_updated_at timestamp with time zone NOT NULL,
    year integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


ALTER TABLE public.courses OWNER TO sylms;

--
-- Name: courses_id_seq; Type: SEQUENCE; Schema: public; Owner: sylms
--

CREATE SEQUENCE public.courses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.courses_id_seq OWNER TO sylms;

--
-- Name: courses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: sylms
--

ALTER SEQUENCE public.courses_id_seq OWNED BY public.courses.id;


--
-- Name: gorp_migrations; Type: TABLE; Schema: public; Owner: sylms
--

CREATE TABLE public.gorp_migrations (
    id text NOT NULL,
    applied_at timestamp with time zone
);


ALTER TABLE public.gorp_migrations OWNER TO sylms;

--
-- Name: courses id; Type: DEFAULT; Schema: public; Owner: sylms
--

ALTER TABLE ONLY public.courses ALTER COLUMN id SET DEFAULT nextval('public.courses_id_seq'::regclass);


--
-- Data for Name: courses; Type: TABLE DATA; Schema: public; Owner: sylms
--

COPY public.courses (id, course_number, course_name, instructional_type, credits, standard_registration_year, term, period_, classroom, instructor, course_overview, remarks, credited_auditors, application_conditions, alt_course_name, course_code, course_code_name, csv_updated_at, year, created_at, updated_at) FROM stdin;
18010	GA10101	情報社会と法制度	1	2.0	{2}	{4,5}	{月5,月6}		{"髙良 幸哉"}	情報化社会における法制度や情報モラル向上に必要な基礎知識を習得することを目指すため、現行の我が国の法制度の基礎を学び、ネットワーク社会における法整備の現状について講義する。	オンライン(オンデマンド型)	0	正規生に対しても受講制限をしているため	Information Society Law	GA10101	情報社会と法制度	2021-03-01 16:18:19+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18011	GA10201	知的財産概論	1	2.0	{2}	{4,5}	{金5,金6}		{"村井 麻衣子"}	知的財産に関する法制度を主要な概念や法理に基づいて学ぶ。著作権法、特許法を中心に、不正競争防止法、商標法など、知的財産諸法についての基礎的な知識を身につけ、知的財産法の法技術的な特色を踏まえた上で、情報化社会における望ましい制度のあり方について考察し、情報の保護と利用についてのバランス感覚や、問題解決能力を身につけることを目的とする。	オンライン(オンデマンド型)	0	正規生に対しても受講制限をしているため	Introduction to Intellectual Property	GA10201	知的財産概論	2021-03-01 16:18:19+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18014	GA12301	システムと情報科学	1	1.0	{1}	{5}	{火5,火6}		{"山際 伸一","山口 佳樹","佐藤 聡","西出 隆志","大山 恵弘"}	情報科学への導入となる基礎理論から応用までを概説し、専門的科目への導入としての基礎知識を習得する。本科目は特に、システムを中心に専門性を習得する上での事前知識となる原理や技術、理論について説明する。	専門導入科目(事前登録対象) オンライン(オンデマンド型)	0	正規生に対しても受講制限をしているため	Introduction to Information Science:Information Systems	GA12301	システムと情報科学	2021-04-14 14:23:34+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18047	GB10244	線形代数B	4	2.0	{2}	{1,2}	{月1,月2}	3A207	{"山田 武志"}	線形代数の基礎。 内容:ベクトル空間,1次写像,核と像,内積空間,固有値・固有ベクトルと対角化	情報科学類3・4クラス対象 オンライン(オンデマンド型) 対面	2		Linear Algebra B	GB10244	線形代数B	2021-03-23 11:43:45+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18020	GA14201	知識情報システム概説	1	1.0	{1}	{2,3}	{木4}		{"高久 雅生","佐藤 哲司","阪口 哲男","鈴木 伸崇"}	ネットワーク社会における知識の構造化、提供、共有のための枠組みについて講義する。	専門導入科目(事前登録対象) オンライン(オンデマンド型)	0	正規生に対しても受講制限をしているため	Foundations of Knowledge Information Systems	GA14201	知識情報システム概説	2021-03-07 20:25:51+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18021	GA14301	図書館概論	1	2.0	{1}	{4,5}	{木3,木4}		{"吉田 右子"}	図書館とは何かについて概説し、これからの図書館の在り方を考える。図書館の歴史と現状、機能と社会的意義、館種別図書館と利用者、図書館職員、類縁機関と関係団体、図書館の課題と展望等について幅広く学ぶ。	専門導入科目(事前登録対象) オンライン(オンデマンド型) GE22001「図書館概論」を修得済みの者は履修不可。	1	本学(学群・大学院)卒業・修了者又は本学の大学院在学者で司書・司書教諭資格希望者に限る	Introduction to Librarianship	GA14301	図書館概論	2021-03-01 16:18:19+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18022	GA15111	情報数学A	1	2.0	{1}	{1,2}	{木5,木6}	3A203	{"西出 隆志","亀山 幸義"}	本授業では,情報学の基礎となる数学的概念について学ぶ.その中でも特に重要な概念である集合,論理,写像,関係,グラフ等を取りあげ,その基礎的な事項について講義する.また,講義内容に対する理解を深めるため,演習も行う.	平成31年度以降入学の者に限る。情報科学類生は1・2クラスを対象とする。 オンライン(オンデマンド型) 定員を超過した場合は履修調整をする場合がある（情報科学類生および総合学域群生(情報科学類への移行希望者・学籍番号の下一桁が奇数)優先）。 	0	正規生に対しても受講制限をしているため	Mathematics for Informatics A	GA15101	情報数学A	2021-03-15 01:06:18+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18029	GA15241	線形代数A	1	2.0	{1}	{2,3}	{金3,金4}		{"長谷川 秀彦"}	行列の基礎概念を学び、それを基に行列演算、連立1次方程式の解法、行列式の性質や展開について講義と演習を行なう。	知識情報・図書館学類生および総合学域群生（知識情報・図書館学類への移行希望者）優先。 履修申請期限は5月11日(火)まで。 定員を超過した場合は履修調整をする場合がある 。 期末試験は対面で実施予定 オンライン(オンデマンド型)	0	正規生に対しても受講制限をしているため	Linear Algebra A	GA15201	線形代数A	2021-04-28 16:21:44+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18048	GB10414	解析学II	4	2.0	{3,4}	{4,5}	{水1,水2}	3A308	{"片岸 一起"}	1変数関数の積分と多変数関数の微分を中心に講義を行う。 	平成30年度以前入学者対象 オンライン(オンデマンド型)	2		Analysis II	GB10414	解析学II	2021-03-23 11:43:40+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18030	GA15311	微分積分A	1	2.0	{1}	{4,5}	{金3,金4}	3B302	{"町田 文雄","堀江 和正"}	解析学の基礎として,実数,関数,数列ならびに連続性や極限などの基本概念と,1変数関数の微分法および積分法について講義を行う。	情報科学類生は1・2クラスを対象とする。定員を超過した場合は履修調整をする場合がある（情報科学類生および総合学域 群生(情報科学類への移行希望者・学籍番号の下一桁が奇数)優先）。履修申請期 限は9月21日(火)まで。 オンライン(オンデマンド型) 平成30年度までに開設された「解析学I」(GB10314,GB10324)の単位を修得した者 の履修は認めない。	0	正規生に対しても受講制限をしているため	Calculus A	GA15301	微分積分A	2021-04-14 14:27:18+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18033	GA15341	微分積分A	1	2.0	{1}	{4,5}	{金3,金4}		{"加藤 誠"}	解析学の基礎として,実数,関数,数列ならびに連続性や極限などの基本概念と,1変数関数の微分法および積分法について講義を行う。	知識学類生および総合学域群生（知識学類への移行希望者）優先。 履修申請期限は9月21日(火)まで。 定員を超過した場合は履修調整をする場合がある 。 オンライン(オンデマンド型)	0	正規生に対しても受講制限をしているため	Calculus A	GA15301	微分積分A	2021-03-01 16:18:17+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18034	GA18212	プログラミング入門A	2	2.0	{1}	{4,5}	{木5,木6}	3A402	{アランニャ," クラウス","新城 靖"}	プログラミングの有用性と必要性を理解し、単純な処理を行うプログラムを書けるようになることを目指す。	情報科学類生および総合学域群生(情報科学類への移行希望者）優先。定員を超過した場合は履修調整をする場合がある。履修申請期限は9月14日(火) まで。原則的に「プログラミング入門B」（GA18312）と同一年度に履修すること。 その他の実施形態 令和2年度までに開設された「プログラミング入門」(GA18112)または平成30年度 までに開設された「プログラミング入門A・B」(GB10664,GB10684)の単位を修得 した者の履修は認めない。	0	正規生に対しても受講制限をしているため	Introduction to Programming A	GA18202	プログラミング入門A	2021-04-23 08:54:49+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18054	GB10524	微分方程式	4	2.0	{2}	{4,5}	{水1,水2}	3B405	{"國廣 昇"}	自然現象を数理モデル化する手段の一つとして微分方程式は有用である.この講義では,線形微分方程式の解法を中心に,微分方程式全般について講義する.	「解析学III」(GB10504)の単位を修得した者の履修は認めない。 オンライン(オンデマンド型)	2			GB10524	微分方程式	2021-03-23 11:43:41+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18060	GB11404	電磁気学	4	2.0	{2}	{4,5}	{木3,木4}	3A306	{"安永 守利"}	集積回路(IC)やハードディスク,タッチパネルや無線LANなど,我々の身の回りの情報通信機器は,電磁現象を原理として動作している.本講義では,これらの電磁現象の基礎を解説する.講義の前半では,「電荷」からスタートして「電場」,「電位」という場の概念とポテンシャルの概念を解説する.また,これらの現象を利用した応用事例も紹介する.後半では,はじめに磁気現象の本質は電流であることを説明し,「磁場」の概念,および「電磁誘導」等の電流と磁気現象の関係を解説する.また,磁気現象を利用した応用事例も紹介する.最後に,「電場」と「磁場」がマクスウェル方程式としてまとめられることを示し,「電磁波」の導出とその応用事例について言及する.	オンライン(オンデマンド型)	2		Electromagnetics	GB11404	電磁気学	2021-03-23 11:43:43+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18061	GB11514	シミュレーション物理	4	1.0	{2}	{6}	{木1,木2}	3A311	{"狩野 均"}	計算機を用いた物理実験について,実験方法から結果のまとめ方まで,演習を交えて系統的に学ぶ。	オンライン(オンデマンド型) 対面	0	計算機の台数制限のため	Computer Simulation Methods in Physics	GB11514	シミュレーション物理	2021-03-23 11:43:44+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18062	GB11601	確率論	1	2.0	{2}	{4,5}	{火5,火6}	3A402	{"馬場 雪乃"}	確率論の基礎。 内容:確率の公理,確率空間,確率変数,分布関数,期待値,特性関数,極限定理など	オンライン(オンデマンド型) 「確率・統計」(GB11611)の単位を修得した者の履修は認めない。	2		Probability Theory	GB11601	確率論	2021-04-14 12:40:03+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18063	GB11611	確率・統計	1	2.0	{2}	{4,5}	{火5,火6}	3A402	{"馬場 雪乃"}	確率論の基礎。 内容:確率の公理,確率空間,確率変数,分布関数,期待値,特性関数,極限定理など	教員免許取得希望者対象。 オンライン(オンデマンド型) 「確率論」(GB11601)の単位を修得した者の履修は認めない。	2		Probability Theory and Statistics	GB11611	確率・統計	2021-04-14 12:40:40+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18064	GB11621	統計学	1	2.0	{2}	{4,5}	{木5,木6}	3A416	{"秋本 洋平"}	数理統計学(統計的推定,仮説検定)ならびに分散分析の基礎と応用(ヒューマンインタフェース評価実験の計画と解析)。理論構成の理解を深めるために,コンピュータを利用した演習を実施。	「確率論」(または同等科目)の履修を前提とする。 オンライン(オンデマンド型) 情報科学類生は2019年度以前の入学生に限る。「統計学」(GB41204)の単位を修得した者の履修は認めない。	2		Statistics	GB11621	統計学	2021-03-23 11:43:43+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18066	GB11931	データ構造とアルゴリズム	1	3.0	{2}	{4,5,6}	{月1,月2}	3B402	{"天笠 俊之","長谷部 浩二","藤田 典久"}	ソフトウェアを書く上で基本となるデータ構造とアルゴリズムの考え方について学ぶ。線形構造,木構造,グラフ構造,データ整列,データ探索について学習する。	平成25年度までに開設された「データ構造とアルゴリズム」(GB11911, GB11921)の単位を修得した者の履修は認めない。 オンライン(同時双方向型)	2		Data Structures and Algorithms	GB11931	データ構造とアルゴリズム	2021-03-23 11:43:37+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
18067	GB11956	データ構造とアルゴリズム実験	6	2.0	{2}	{4,5,6}	{月3,月4,月5,月3,月4}	3C113,3C205	{"天笠 俊之"}	データ構造とアルゴリズムに関して,実際にJava言語を用いてプログラムを作成し,そのプログラムが稼働することを確認する。プログラムは,毎週,あるいは隔週に一個の割合で作成する。	1・2クラス オンライン(同時双方向型) 令和2年度までに開設された「データ構造とアルゴリズム実験」(GB11936,GB11946)または平成26年度までに開設された「データ構造とアルゴリズム実験」(GB11916, GB11926)の単位を修得した者の履修は認めない。	0	施設設備の許容量上の制約と学類生に対する良質の少人数教育を行うため	Data Structures and Algorithms Laboratory	GB11956	データ構造とアルゴリズム実験	2021-03-23 11:43:36+09	2021	2022-01-08 23:59:18.11643+09	2022-01-08 23:59:18.11643+09
\.


--
-- Data for Name: gorp_migrations; Type: TABLE DATA; Schema: public; Owner: sylms
--

COPY public.gorp_migrations (id, applied_at) FROM stdin;
20210619141018-init.sql	2022-01-08 23:59:18.158106+09
\.


--
-- Name: courses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: sylms
--

SELECT pg_catalog.setval('public.courses_id_seq', 19796, true);


--
-- Name: courses courses_pkey; Type: CONSTRAINT; Schema: public; Owner: sylms
--

ALTER TABLE ONLY public.courses
    ADD CONSTRAINT courses_pkey PRIMARY KEY (id);


--
-- Name: gorp_migrations gorp_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: sylms
--

ALTER TABLE ONLY public.gorp_migrations
    ADD CONSTRAINT gorp_migrations_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

