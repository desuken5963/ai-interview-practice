# AI面接評価フィードバック機能

## 1. 概要
AI面接練習後の評価・フィードバック機能の基本設計を記載する。

## 2. 画面設計

### 2.1 面接評価結果画面

#### 2.1.1 総合評価セクション
- **総合ランク表示**
  - A～Eの5段階評価を表示
  - ランクに応じたカラーリング（A: 青, B: 緑, C: 黄, D: オレンジ, E: 赤）

- **総合スコア**
  - 0-100点での数値評価
  - 各評価軸のスコアを重み付けして算出

- **総評コメント**
  - 面接全体を通しての評価コメント
  - 特に良かった点と改善点を箇条書きで表示

#### 2.1.2 詳細評価セクション
- **レーダーチャート**
  - 以下の5つの評価軸でのスコアを図示
    1. 論理性（ストーリー構成、説明の流れ）
    2. 具体性（具体例の提示、数値での説明）
    3. ポジティブ感（前向きな姿勢、熱意）
    4. 適切性（質問意図との整合性、企業文化との適合）
    5. 簡潔性（無駄のない回答、要点の明確さ）

- **各軸の詳細スコア**
  - 各評価軸の点数（0-100点）
  - 評価の根拠となるポイントの説明

#### 2.1.3 質問別評価セクション
- **質問と回答の履歴**
  - 面接での各質問と回答内容を時系列で表示
  - 各回答に対する個別評価を表示
    - スコア（0-100点）
    - 良かった点
    - 改善点
    - AIからのアドバイス

### 2.2 面接履歴一覧画面

#### 2.2.1 練習履歴リスト
- **練習セッション一覧**
  - 実施日時
  - 企業名（設定時）
  - 求人情報（設定時）
  - 総合ランク
  - 総合スコア

#### 2.2.2 成長推移グラフ
- **時系列での評価推移**
  - 総合スコアの推移
  - 各評価軸のスコア推移
  - 期間選択機能（週間/月間/年間）

## 3. データ構造

### 3.1 面接評価テーブル (interview_evaluations)
| カラム名 | 型 | 説明 | NULL |
|---------|-----|------|------|
| id | UUID | 主キー | NO |
| session_id | UUID | 面接セッションID (FK) | NO |
| total_rank | ENUM('A','B','C','D','E') | 総合ランク | NO |
| total_score | INTEGER | 総合スコア（0-100） | NO |
| overall_comment | TEXT | 総評コメント | NO |
| logical_score | INTEGER | 論理性スコア（0-100） | NO |
| concrete_score | INTEGER | 具体性スコア（0-100） | NO |
| positive_score | INTEGER | ポジティブ感スコア（0-100） | NO |
| relevance_score | INTEGER | 適切性スコア（0-100） | NO |
| concise_score | INTEGER | 簡潔性スコア（0-100） | NO |
| created_at | TIMESTAMP | 作成日時 | NO |

### 3.2 回答評価テーブル (answer_evaluations)
| カラム名 | 型 | 説明 | NULL |
|---------|-----|------|------|
| id | UUID | 主キー | NO |
| answer_id | UUID | 面接回答ID (FK) | NO |
| score | INTEGER | 回答スコア（0-100） | NO |
| strengths | TEXT[] | 良かった点（配列） | NO |
| improvements | TEXT[] | 改善点（配列） | NO |
| advice | TEXT | AIからのアドバイス | NO |
| created_at | TIMESTAMP | 作成日時 | NO |

## 4. API設計

### 4.1 評価生成API

#### POST /api/v1/interview-sessions/{session_id}/evaluate
面接セッション全体の評価を生成

- リクエストボディ
```json
{
    "session_id": "uuid",
    "questions": [
        {
            "id": "uuid",
            "content": "質問内容",
            "answer": "回答内容"
        }
    ]
}
```

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
            "concrete_score": 85,
            "positive_score": 80,
            "relevance_score": 85,
            "concise_score": 85
        },
        "answer_evaluations": [
            {
                "answer_id": "uuid",
                "score": 85,
                "strengths": ["良かった点1", "良かった点2"],
                "improvements": ["改善点1", "改善点2"],
                "advice": "アドバイス内容"
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
        "concrete_scores": [85, 88, 90],
        "positive_scores": [80, 82, 85],
        "relevance_scores": [85, 87, 90],
        "concise_scores": [85, 86, 90]
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
1. 論理性（ストーリー構成、説明の流れ）
2. 具体性（具体例の提示、数値での説明）
3. ポジティブ感（前向きな姿勢、熱意）
4. 適切性（質問意図との整合性、企業文化との適合）
5. 簡潔性（無駄のない回答、要点の明確さ）

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
        "concrete_score": 0-100,
        "positive_score": 0-100,
        "relevance_score": 0-100,
        "concise_score": 0-100
    },
    "answer_evaluations": [
        {
            "answer_id": "uuid",
            "score": 0-100,
            "strengths": ["良かった点"],
            "improvements": ["改善点"],
            "advice": "アドバイス"
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
