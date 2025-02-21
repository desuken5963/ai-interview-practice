# AI面接評価フィードバック機能

## 1. 概要
AI面接練習後の評価・フィードバック機能の基本設計を記載する。

## 2. 画面設計

### 2.1 面接完了画面 (/interview-sessions/{session_id}/complete)

#### 2.1.1 面接完了メッセージ
- **完了通知**
  - 面接練習が完了した旨のメッセージを表示
  - フィードバック生成中のローディング表示

#### 2.1.2 面接サマリー
- **面接情報**
  - 実施日時
  - 企業名（設定時）
  - 求人情報（設定時）
  - 面接フェーズ
  - 質問数

#### 2.1.3 アクション
- **フィードバック確認ボタン**
  - フィードバック生成完了後に表示
  - クリックでフィードバック詳細画面へ遷移
- **フィードバック一覧へのリンク**
  - 過去の面接履歴一覧画面へ遷移

### 2.2 フィードバック一覧画面 (/feedback)

#### 2.2.1 フィルター・検索
- **期間フィルター**
  - 週間/月間/年間の切り替え
  - カスタム期間の指定
- **企業・求人フィルター**
  - 企業名での絞り込み
  - 求人情報での絞り込み
- **面接フェーズフィルター**
  - 面接フェーズでの絞り込み

#### 2.2.2 面接履歴リスト
- **各面接のサマリーカード**
  - 実施日時
  - 企業名・求人情報
  - 面接フェーズ
  - 総合ランク（A-E）
  - 総合スコア（0-100点）
  - レーダーチャートのサムネイル（6軸評価の概要）

#### 2.2.3 統計情報
- **全体の成長推移グラフ**
  - 総合スコアの推移
  - 各評価軸のスコア推移
- **強み・弱み分析**
  - 評価軸ごとの平均スコア
  - 特に高評価な項目
  - 改善が必要な項目

### 2.3 フィードバック詳細画面 (/feedback/{session_id})

#### 2.3.1 面接情報ヘッダー
- **基本情報**
  - 実施日時
  - 企業名・求人情報
  - 面接フェーズ
  - 質問数

#### 2.3.2 総合評価セクション
- **総合評価バッジ**
  - ランク（A-E）を大きく表示
  - ランクに応じたカラーリング
  - ランク判定基準：
    - A: 90-100点
    - B: 80-89点
    - C: 70-79点
    - D: 60-69点
    - E: 0-59点
- **総合スコア**
  - 0-100点のスコアを表示
  - 前回からの変化（＋/-）
- **総評コメント**
  - 面接全体の評価コメント
  - 特に評価できる点
  - 改善を推奨する点

#### 2.3.3 詳細評価セクション
- **レーダーチャート**
  - 6つの評価軸でのスコアを図示
    1. 論理的思考力（ストーリー構成、論理展開）
    2. コミュニケーション力（説明の明確さ、対話力）
    3. 技術力（専門知識、スキルの深さ）
    4. 問題解決能力（課題分析、解決アプローチ）
    5. 志望度・意欲（熱意、モチベーション）
    6. カルチャーフィット（企業文化との適合性）
- **各評価軸の詳細**
  - スコア（0-100点）
  - 評価コメント
  - 改善アドバイス

#### 2.3.4 質問別評価セクション
- **タイムライン形式で表示**
  - 質問内容
  - 回答内容
  - 個別評価
    - スコア（0-100点）
    - 良かった点（箇条書き）
    - 改善点（箇条書き）
    - 具体的なアドバイス

#### 2.3.5 アクション
- **フィードバックの共有**
  - PDFでのエクスポート
  - リンクでの共有
- **新規面接練習の開始**
  - 同じ設定での再挑戦
  - 設定を変更しての挑戦

## 3. データ構造

### 3.1 面接評価テーブル (interview_evaluations)
| カラム名 | 型 | 説明 | NULL |
|---------|-----|------|------|
| id | UUID | 主キー | NO |
| session_id | UUID | 面接セッションID (FK) | NO |
| total_rank | ENUM('A','B','C','D','E') | 総合ランク | NO |
| total_score | INTEGER | 総合スコア（0-100） | NO |
| overall_comment | TEXT | 総評コメント | NO |
| logical_score | INTEGER | 論理的思考力スコア（0-100） | NO |
| communication_score | INTEGER | コミュニケーション力スコア（0-100） | NO |
| technical_score | INTEGER | 技術力スコア（0-100） | NO |
| problem_solving_score | INTEGER | 問題解決能力スコア（0-100） | NO |
| motivation_score | INTEGER | 志望度・意欲スコア（0-100） | NO |
| culture_fit_score | INTEGER | カルチャーフィットスコア（0-100） | NO |
| created_at | TIMESTAMP | 作成日時 | NO |

### 3.2 回答評価テーブル (answer_evaluations)
| カラム名 | 型 | 説明 | NULL |
|---------|-----|------|------|
| id | UUID | 主キー | NO |
| session_id | UUID | 面接セッションID (FK) | NO |
| question_id | UUID | 面接質問ID (FK) | NO |
| logical_score | INTEGER | 論理的思考力スコア（0-100） | NO |
| communication_score | INTEGER | コミュニケーション力スコア（0-100） | NO |
| technical_score | INTEGER | 技術力スコア（0-100） | NO |
| problem_solving_score | INTEGER | 問題解決能力スコア（0-100） | NO |
| motivation_score | INTEGER | 志望度・意欲スコア（0-100） | NO |
| culture_fit_score | INTEGER | カルチャーフィットスコア（0-100） | NO |
| question_comment | TEXT | 質問に対するコメント | NO |
| strengths | TEXT[] | 良かった点（配列） | NO |
| improvements | TEXT[] | 改善点（配列） | NO |
| created_at | TIMESTAMP | 作成日時 | NO |

## 4. API設計

### 4.1 評価生成API

#### POST /api/v1/interview-sessions/{session_id}/evaluate
面接セッション全体の評価を生成

##### 評価生成プロセス
1. **回答単位の評価生成**
   - 面接セッション情報（企業、求人、フェーズ）を取得
   - 各質問・回答ペアに対して個別に評価を実施
   - 評価結果を`answer_evaluations`テーブルに保存
   ```
   # 回答評価プロンプト
   企業：{company_name}
   求人：{job_posting_details}
   面接フェーズ：{interview_phase}
   質問：{question_content}
   回答：{answer_content}

   # 評価項目
   1. 論理的思考力（ストーリー構成、論理展開）
   2. コミュニケーション力（説明の明確さ、対話力）
   3. 技術力（専門知識、スキルの深さ）
   4. 問題解決能力（課題分析、解決アプローチ）
   5. 志望度・意欲（熱意、モチベーション）
   6. カルチャーフィット（企業文化との適合性）

   # 出力形式
   {
       "logical_score": 0-100,
       "communication_score": 0-100,
       "technical_score": 0-100,
       "problem_solving_score": 0-100,
       "motivation_score": 0-100,
       "culture_fit_score": 0-100,
       "question_comment": "回答に対する詳細なコメント",
       "strengths": ["良かった点の配列"],
       "improvements": ["改善点の配列"]
   }
   ```

2. **セッション全体の評価生成**
   - `answer_evaluations`のデータを基に`interview_evaluations`を作成
   - 以下の項目を`answer_evaluations`の対応項目の単純平均から計算
     - 各評価軸のスコア（logical_score, communication_score, technical_score, problem_solving_score, motivation_score, culture_fit_score）
     - 総合スコア（total_score）：6つの評価軸の単純平均
     - 総合ランク（total_rank）：total_scoreに基づき判定
   - 総評コメント（overall_comment）のみ以下のプロンプトで生成
   ```
   # 総評生成プロンプト
   企業：{company_name}
   求人：{job_posting_details}
   面接フェーズ：{interview_phase}

   # 回答評価結果一覧
   [
     {
       "question": "質問内容",
       "scores": {
         "logical_score": 90,
         "communication_score": 85,
         ...
       },
       "strengths": ["..."],
       "improvements": ["..."]
     },
     ...
   ]

   # 計算済み評価スコア
   - 論理的思考力: {logical_score}
   - コミュニケーション力: {communication_score}
   - 技術力: {technical_score}
   - 問題解決能力: {problem_solving_score}
   - 志望度・意欲: {motivation_score}
   - カルチャーフィット: {culture_fit_score}
   - 総合スコア: {total_score}
   - 総合ランク: {total_rank}

   # 出力形式
   {
       "overall_comment": "面接全体の詳細な評価コメント。候補者の強みと改善点を含めた具体的なフィードバック。"
   }
   ```

- リクエストボディ: なし（セッションIDからデータを取得）

- レスポンス
```json
{
    "evaluation": {
        "id": "uuid",
        "total_rank": "A-E",
        "total_score": 85,
        "overall_comment": "評価コメント",
        "scores": {
            "logical_score": 90,
            "communication_score": 85,
            "technical_score": 88,
            "problem_solving_score": 82,
            "motivation_score": 95,
            "culture_fit_score": 87
        },
        "answer_evaluations": [
            {
                "question_id": "uuid",
                "logical_score": 90,
                "communication_score": 85,
                "technical_score": 88,
                "problem_solving_score": 82,
                "motivation_score": 95,
                "culture_fit_score": 87,
                "question_comment": "回答に対する詳細なコメント",
                "strengths": ["良かった点1", "良かった点2"],
                "improvements": ["改善点1", "改善点2"]
            }
        ]
    }
}
```

### 4.2 評価履歴API

#### GET /api/v1/users/{user_id}/interview-evaluations
ユーザーの面接評価履歴を取得

- クエリパラメータ
  - `period`: 取得期間（week/month/year）
  - `limit`: 取得件数（デフォルト10件）
  - `offset`: 開始位置（ページネーション用）

- レスポンス
```json
{
    "total_count": 100,
    "evaluations": [
        {
            "id": "uuid",
            "session_id": "uuid",
            "company_name": "企業名",
            "job_posting_title": "求人タイトル",
            "total_rank": "A-E",
            "total_score": 85,
            "created_at": "2024-03-20T10:00:00Z"
        }
    ],
    "score_trends": {
        "total_scores": [85, 87, 90],
        "logical_scores": [90, 92, 95],
        "communication_scores": [85, 88, 90],
        "technical_scores": [85, 88, 90],
        "problem_solving_scores": [80, 82, 85],
        "motivation_scores": [80, 82, 85],
        "culture_fit_scores": [85, 86, 90]
    }
}
```

## 5. AIプロンプト設計

### 5.1 評価生成プロンプト

```
あなたは面接評価のエキスパートとして、以下の基準で面接回答を評価してください：

# 評価対象情報
企業：{company_name}（指定がある場合）
求人：{job_posting_details}（指定がある場合）
面接フェーズ：{interview_phase}

# 評価項目
1. 論理的思考力（ストーリー構成、論理展開）
2. コミュニケーション力（説明の明確さ、対話力）
3. 技術力（専門知識、スキルの深さ）
4. 問題解決能力（課題分析、解決アプローチ）
5. 志望度・意欲（熱意、モチベーション）
6. カルチャーフィット（企業文化との適合性）

# 質問・回答履歴
[
  {
    "question": "質問内容",
    "answer": "回答内容"
  }
]

# 出力形式
{
    "total_rank": "A-E",
    "total_score": 0-100,
    "overall_comment": "総評",
    "scores": {
        "logical_score": 0-100,
        "communication_score": 0-100,
        "technical_score": 0-100,
        "problem_solving_score": 0-100,
        "motivation_score": 0-100,
        "culture_fit_score": 0-100
    },
    "answer_evaluations": [
        {
            "question_id": "uuid",
            "logical_score": 0-100,
            "communication_score": 0-100,
            "technical_score": 0-100,
            "problem_solving_score": 0-100,
            "motivation_score": 0-100,
            "culture_fit_score": 0-100,
            "question_comment": "回答に対する詳細なコメント",
            "strengths": ["良かった点の配列"],
            "improvements": ["改善点の配列"]
        }
    ]
}
```

## 6. エラー処理

### 6.1 評価生成エラー
- **AIサービス接続エラー**
  - エラーコード: 503
  - メッセージ: "評価生成サービスが一時的に利用できません"
  - リトライ処理: 最大3回まで自動リトライ

- **評価生成タイムアウト**
  - エラーコード: 504
  - メッセージ: "評価生成がタイムアウトしました"
  - 処理: 30秒でタイムアウト、ユーザーに再試行を促す

### 6.2 データ取得エラー
- **評価データ不存在**
  - エラーコード: 404
  - メッセージ: "指定された評価データが見つかりません"

- **不正なパラメータ**
  - エラーコード: 400
  - メッセージ: "無効なパラメータが指定されました"
