# 面接練習機能

## 1. 概要
テキストチャットおよび音声通話形式での面接練習、AIによる評価・フィードバック機能の基本設計を記載する。

## 2. 画面設計

### 2.1 面接練習開始画面
- 練習モード選択（テキスト/音声）
- 企業・求人情報の選択（任意）
  - 指定無し
  - 企業指定
  - 求人指定
- 面接設定項目(必須)
  - 面接の種類（一次面接、最終面接など）
  - 面接時間目安（15分、30分など）
  - 面接官の性格（フレンドリー、厳格など）

### 2.2 面接画面
- 面接官キャラクター表示エリア
  - アニメーションステート（待機/話す/聞く/考える）
  - 表情変化
- 対話エリア
  - 面接官（AI）の発言表示
  - ユーザーの回答表示/入力
    - テキストモード: 入力フォーム
    - 音声モード: 音声入力UI
- ステータス表示
  - 経過時間
  - 残り想定質問数
  - 音声認識状態（音声モード時）
- コントロール
  - 一時停止/再開
  - 音声入力ON/OFF
  - 終了

### 2.3 音声面接画面
- 音声入力/出力状態表示
- 音声認識テキストのリアルタイム表示
- マイク/スピーカー設定
- 経過時間表示
- 面接終了ボタン
- 一時停止/再開ボタン

### 2.4 評価・フィードバック画面
- 総合評価（レーダーチャート）
- 質問別スコア一覧
- AIからの改善アドバイス
- 練習履歴への保存ボタン
- 新規練習開始ボタン

## 3. データ構造

### 3.1 面接セッションテーブル (interview_sessions)
| カラム名 | 型 | 説明 | NULL |
|---------|-----|------|------|
| id | UUID | 主キー | NO |
| user_id | UUID | ユーザーID (FK) | NO |
| company_id | UUID | 企業ID (FK) | YES |
| job_posting_id | UUID | 求人ID (FK) | YES |
| interview_type | VARCHAR(20) | 面接種類 | NO |
| mode | VARCHAR(10) | テキスト/音声 | NO |
| duration_minutes | INTEGER | 面接時間（分） | NO |
| interviewer_personality | VARCHAR(20) | 面接官タイプ | NO |
| status | VARCHAR(20) | 実施状態 | NO |
| started_at | TIMESTAMP | 開始日時 | NO |
| ended_at | TIMESTAMP | 終了日時 | YES |

### 3.2 面接QAテーブル (interview_qa_logs)
| カラム名 | 型 | 説明 | NULL |
|---------|-----|------|------|
| id | UUID | 主キー | NO |
| session_id | UUID | セッションID (FK) | NO |
| question | TEXT | 質問内容 | NO |
| answer | TEXT | 回答内容 | NO |
| answer_audio_url | VARCHAR(255) | 音声回答URL | YES |
| sequence | INTEGER | 質問順序 | NO |
| created_at | TIMESTAMP | 作成日時 | NO |

### 3.3 評価スコアテーブル (interview_evaluations)
| カラム名 | 型 | 説明 | NULL |
|---------|-----|------|------|
| id | UUID | 主キー | NO |
| qa_log_id | UUID | QAログID (FK) | NO |
| category | VARCHAR(50) | 評価カテゴリ | NO |
| score | INTEGER | スコア（0-100） | NO |
| feedback | TEXT | フィードバック内容 | NO |
| created_at | TIMESTAMP | 作成日時 | NO |

## 4. API設計

### 4.1 面接セッションAPI

#### POST /api/v1/interview-sessions
面接セッションを開始

- リクエストボディ
```json
{
    "company_id": "uuid（任意）",
    "job_posting_id": "uuid（任意）",
    "interview_type": "FIRST",
    "mode": "TEXT",
    "duration_minutes": 30,
    "interviewer_personality": "FRIENDLY"
}
```

- レスポンス
```json
{
    "session_id": "uuid",
    "initial_question": "最初の質問内容",
    "started_at": "2024-01-01T00:00:00Z"
}
```

#### POST /api/v1/interview-sessions/{session_id}/qa
質問への回答を送信

- リクエストボディ（テキストモード）
```json
{
    "answer": "回答テキスト"
}
```

- リクエストボディ（音声モード）
```json
{
    "audio_data": "base64エンコードされた音声データ"
}
```

- レスポンス
```json
{
    "evaluation": {
        "scores": {
            "logic": 85,
            "clarity": 90,
            "relevance": 75
        },
        "feedback": "フィードバックコメント"
    },
    "next_question": "次の質問内容"
}
```

#### PUT /api/v1/interview-sessions/{session_id}/end
面接セッションを終了

- レスポンス
```json
{
    "session_summary": {
        "duration": "25分",
        "question_count": 8,
        "overall_score": 85,
        "category_scores": {
            "logic": 80,
            "clarity": 90,
            "relevance": 85
        },
        "improvement_points": [
            "改善ポイント1",
            "改善ポイント2"
        ]
    }
}
```

### 4.2 音声認識Websocket API

#### WS /api/v1/interview-sessions/{session_id}/voice
音声データのストリーミング処理用Websocket接続

- クライアントからサーバーへのメッセージ
```json
{
    "type": "audio_data",
    "data": "base64エンコードされた音声データ"
}
```

- サーバーからクライアントへのメッセージ
```json
{
    "type": "recognition_result",
    "text": "認識されたテキスト",
    "is_final": true
}
```

### 4.3 評価履歴API

#### GET /api/v1/interview-sessions/{session_id}
セッションの詳細情報を取得

- レスポンス
```json
{
    "session_details": {
        "id": "uuid",
        "interview_type": "FIRST",
        "mode": "TEXT",
        "started_at": "2024-01-01T00:00:00Z",
        "ended_at": "2024-01-01T00:30:00Z",
        "qa_logs": [
            {
                "question": "質問内容",
                "answer": "回答内容",
                "evaluation": {
                    "scores": {
                        "logic": 85,
                        "clarity": 90,
                        "relevance": 75
                    },
                    "feedback": "フィードバック内容"
                }
            }
        ],
        "overall_evaluation": {
            "total_score": 85,
            "category_scores": {
                "logic": 80,
                "clarity": 90,
                "relevance": 85
            },
            "improvement_suggestions": [
                "改善提案1",
                "改善提案2"
            ]
        }
    }
}
```

## 5. AIプロンプト設計

### 5.1 面接官プロンプト
```text
あなたは面接官として以下の設定で面接を行います：
- 企業: {company_name}
- 求人: {job_title}
- 面接タイプ: {interview_type}
- 面接官の性格: {interviewer_personality}

以下の情報を参考に、適切な質問を生成してください：
- 企業の事業内容: {business_description}
- 求人の詳細: {job_description}
- 面接の進行状況: {current_progress}
```

### 5.2 回答評価プロンプト
```text
以下の回答を評価し、各項目でスコアとフィードバックを提供してください：

質問: {question}
回答: {answer}

評価項目：
1. 論理性（構成、説得力）
2. 明確性（簡潔さ、わかりやすさ）
3. 関連性（質問との適合性）
4. ポジティブ度（前向きさ、意欲）
``` 