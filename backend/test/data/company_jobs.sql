-- テストデータ用シードファイル
-- 既存のデータをクリア
DELETE FROM job_custom_fields;
DELETE FROM job_postings;
DELETE FROM company_custom_fields;
DELETE FROM companies;

-- 企業データのシード
INSERT INTO companies (name, business_description) VALUES 
('株式会社テックイノベーション', 'AIと機械学習を活用した革新的なソリューションを提供する企業です。クラウドサービス、データ分析、自然言語処理など、最先端技術を駆使したサービスを展開しています。'),
('グローバルコンサルティング株式会社', '世界各国の企業に対して、経営戦略、デジタルトランスフォーメーション、組織改革などのコンサルティングサービスを提供しています。'),
('未来フィンテック株式会社', 'ブロックチェーン技術を活用した次世代の金融サービスを開発。個人向けおよび法人向けの革新的な決済ソリューションを提供しています。'),
('エコテクノロジー株式会社', '再生可能エネルギーとスマートグリッド技術を組み合わせた環境配慮型のエネルギーマネジメントシステムを開発・提供しています。'),
('ヘルスケアソリューションズ株式会社', 'IoTとAIを活用した遠隔医療プラットフォームの開発・運営。予防医療から治療後のケアまで、包括的な医療サービスを提供しています。'),
('デジタルエデュケーション株式会社', 'オンライン教育プラットフォームの開発・運営。個別最適化された学習体験を提供し、生涯学習をサポートしています。');

-- 最後に挿入したIDの取得用変数
SET @tech_id = LAST_INSERT_ID();
SET @consult_id = @tech_id + 1;
SET @fintech_id = @tech_id + 2;
SET @eco_id = @tech_id + 3;
SET @health_id = @tech_id + 4;
SET @edu_id = @tech_id + 5;

-- 企業カスタムフィールドのシード
INSERT INTO company_custom_fields (company_id, field_name, content) VALUES
(@tech_id, '業界', 'IT・テクノロジー'),
(@tech_id, '従業員数', '150名'),
(@tech_id, '本社所在地', '東京都渋谷区'),
(@consult_id, '業界', 'コンサルティング'),
(@consult_id, '従業員数', '300名'),
(@consult_id, '本社所在地', '東京都千代田区'),
(@fintech_id, '業界', 'フィンテック'),
(@fintech_id, '従業員数', '80名'),
(@fintech_id, '本社所在地', '東京都港区'),
(@eco_id, '業界', 'エネルギー・環境'),
(@eco_id, '従業員数', '120名'),
(@eco_id, '本社所在地', '大阪府大阪市'),
(@health_id, '業界', 'ヘルスケア'),
(@health_id, '従業員数', '200名'),
(@health_id, '本社所在地', '福岡県福岡市'),
(@edu_id, '業界', 'エドテック'),
(@edu_id, '従業員数', '90名'),
(@edu_id, '本社所在地', '東京都新宿区');

-- 求人データのシード
INSERT INTO job_postings (company_id, title, description) VALUES
(@tech_id, 'AIエンジニア', 'AIと機械学習を活用したソリューション開発を担当していただきます。機械学習の実務経験3年以上、Pythonでの開発経験必須。勤務地：東京（リモート可）、給与：600万円〜900万円'),
(@tech_id, 'フルスタックエンジニア', 'フロントエンドからバックエンドまで幅広く開発を担当していただきます。JavaScript/TypeScript, React, Node.js, GoまたはPythonの経験。勤務地：東京（リモート可）、給与：500万円〜800万円'),
(@consult_id, 'ITコンサルタント', 'クライアント企業のDX推進を支援していただきます。コンサルティングまたはIT業界での実務経験3年以上。勤務地：東京・大阪、給与：700万円〜1000万円'),
(@fintech_id, 'ブロックチェーンエンジニア', 'ブロックチェーン技術を活用した金融サービスの開発を担当していただきます。ブロックチェーン開発経験、Solidityの知識。勤務地：東京、給与：800万円〜1200万円'),
(@eco_id, 'エネルギーシステムエンジニア', 'スマートグリッドシステムの開発・実装を担当していただきます。組み込みシステム開発経験、IoTデバイス連携の知識。勤務地：横浜、給与：600万円〜900万円'),
(@health_id, 'AIヘルスケアエンジニア', '医療データ分析と予測モデル開発を担当していただきます。機械学習の知識、医療データ分析の経験があれば尚可。勤務地：東京・福岡、給与：650万円〜950万円'),
(@edu_id, 'エドテックプロダクトマネージャー', 'オンライン教育プラットフォームの企画・開発を担当していただきます。プロダクトマネジメント経験、教育分野の知識があれば尚可。勤務地：東京（リモート可）、給与：550万円〜850万円');

-- 最後に挿入したIDの取得用変数
SET @ai_eng_id = LAST_INSERT_ID();
SET @fullstack_id = @ai_eng_id + 1;
SET @consultant_id = @ai_eng_id + 2;
SET @blockchain_id = @ai_eng_id + 3;
SET @energy_id = @ai_eng_id + 4;
SET @healthcare_id = @ai_eng_id + 5;
SET @edtech_id = @ai_eng_id + 6;

-- 求人カスタムフィールドのシード
INSERT INTO job_custom_fields (job_id, field_name, content) VALUES
(@ai_eng_id, '雇用形態', '正社員'),
(@ai_eng_id, '勤務時間', '9:00-18:00（フレックスタイム制）'),
(@ai_eng_id, '休日休暇', '完全週休2日制、年間休日125日'),
(@fullstack_id, '雇用形態', '正社員'),
(@fullstack_id, '勤務時間', '9:00-18:00（フレックスタイム制）'),
(@fullstack_id, '休日休暇', '完全週休2日制、年間休日125日'),
(@consultant_id, '雇用形態', '正社員'),
(@consultant_id, '勤務時間', '9:00-18:00（フレックスタイム制）'),
(@consultant_id, '休日休暇', '完全週休2日制、年間休日120日'),
(@blockchain_id, '雇用形態', '正社員'),
(@blockchain_id, '勤務時間', '10:00-19:00（フレックスタイム制）'),
(@blockchain_id, '休日休暇', '完全週休2日制、年間休日125日'),
(@energy_id, '雇用形態', '正社員'),
(@energy_id, '勤務時間', '9:00-18:00（フレックスタイム制）'),
(@energy_id, '休日休暇', '完全週休2日制、年間休日125日'),
(@healthcare_id, '雇用形態', '正社員'),
(@healthcare_id, '勤務時間', '9:00-18:00（フレックスタイム制）'),
(@healthcare_id, '休日休暇', '完全週休2日制、年間休日125日'),
(@edtech_id, '雇用形態', '正社員'),
(@edtech_id, '勤務時間', '10:00-19:00（フレックスタイム制）'),
(@edtech_id, '休日休暇', '完全週休2日制、年間休日125日'); 