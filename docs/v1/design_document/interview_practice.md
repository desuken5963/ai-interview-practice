# 面接練習機能

## 1. 概要
テキストチャットおよび音声通話形式での、AI面接練習機能の基本設計を記載する。

## 2. 画面設計

### 2.1 面接練習開始画面

#### 2.1.1 企業・求人設定
- **企業情報の指定**
  - 登録済みの企業 (companies テーブル) から選択可能
  - 企業を指定しないことも可能

- **求人情報の指定**
  - 企業を選択している場合、その企業に紐づく求人 (job_postings テーブル) から選択可能
  - 企業を指定していない場合は選択不可
  - 指定しないことも可能

#### 2.1.2 面接方法の選択
- **テキスト**  
  - 面接官の質問を画面に表示
  - 画面下部のテキスト入力欄に回答を入力  
  - 入力内容をテキストとして API に送信

- **音声**  
  - 面接官の質問を画面に表示すると同時に、外部 API 等を利用して音声で読み上げ  
  - 画面下部の録音ボタン（REC On/Off）で回答を録音し、録音データをテキスト変換して API に送信  

> **注意点**  
> - ユーザーの音声回答は最終的にテキスト化し、プロンプトの一部として AI へ送信する  
> - 音声合成・音声認識は外部 API を利用する想定  

#### 2.1.3 面接設定
1. **質問数**  
   - 面接で行う質問の回数  
   - 選択肢: 5問 / 10問 / 15問  
   - 面接時間の目安はこの質問数に基づき想定
   - 必須選択事項

2. **想定面接官の役職**  
   - 例: 「人事担当」「現場責任者」「経営者」など  
   - 未選択のままでも可

3. **想定面接フェーズ (シチュエーション)**  
   - 以下のいずれかを選択  
     - 一次面接  
     - 二次面接  
     - 最終面接  
     - 圧迫面接  
     - インターン面接  
   - 未選択も可能

4. **アイスブレイクの有無**  
   - 面接の最初にアイスブレイクを実施するかどうかを選択  
   - 選択しない場合は挨拶の次に最初の本質問を直接開始  
   - 選択しないことも可能

5. **自己紹介の有無**  
   - 面接の最初に自己紹介を求めるかどうかを選択  
   - 選択しない場合はアイスブレイクまたは本質問へ直接開始  
   - 選択しないことも可能

### 2.2 AI面接練習画面
- 面接練習開始画面で設定した内容に基づき、AI面接を実施
- 画面はリロードやページ遷移を挟まずに進行
- 面接官の質問は、面接フェーズごとに分けて管理(フェーズごとに叩くapiが変わる)

#### 2.2.1 質問表示
- 画面中央に「AI面接官の質問テキスト」を表示
- 音声面接の場合は同時に読み上げを実行

#### 2.2.2 回答入力
- **テキスト面接**  
  - 画面下部のテキスト入力欄に回答を入力  
  - 送信ボタン押下で回答をAI APIに送信

- **音声面接**  
  - 画面下部の録音ボタンを押下して回答を録音  
  - 録音終了後、録音データをテキストに変換し、AI APIに送信

#### 2.2.3 面接継続
- 1問ごとに「AIの質問 → ユーザー回答」というサイクルを繰り返す
- 選択した「質問数」に達するか、ユーザーが中断した場合に面接を終了

### 2.3 AI面接練習完了画面
- 面接終了後に「AI面接練習が完了した」旨を表示
- 必要に応じて振り返り要素（回答のテキスト一覧など）を表示することを検討

## 3. データ構造

### 3.1 面接セッションテーブル (interview_sessions)
| カラム名 | 型 | 説明 | NULL |
|---------|-----|------|------|
| id | UUID | 主キー | NO |
| company_id | UUID | 企業ID (FK) | YES |
| job_posting_id | UUID | 求人ID (FK) | YES |
| interview_phase | TEXT | 面接フェーズ（例：一次面接、最終面接など） | YES |
| interviewer_role | TEXT | 面接官役職（例：人事担当、現場責任者など） | YES |
| question_count | INTEGER | 質問数（5, 10, 15のいずれか） | NO |
| include_self_introduction | BOOLEAN | 自己紹介実施有無 | NO |
| include_ice_break | BOOLEAN | アイスブレイク実施有無 | NO |
| status | ENUM('CREATED','GREETING','SELF_INTRODUCTION','ICE_BREAK','MAIN','PAUSED','COMPLETED','TERMINATED','CLOSING') | 実施状態 | NO |
| started_at | TIMESTAMP | 開始日時 | NO |
| ended_at | TIMESTAMP | 終了日時 | YES |

> **セッションステータスの定義**
> - CREATED: セッション作成直後の初期状態
> - GREETING: 挨拶フェーズ実施中
> - SELF_INTRODUCTION: 自己紹介フェーズ実施中
> - ICE_BREAK: アイスブレイクフェーズ実施中
> - MAIN: 主要質問フェーズ実施中
> - PAUSED: 一時中断中（再開可能）
> - COMPLETED: 全質問完了による正常終了
> - TERMINATED: 面接の途中終了（一時中断からの中止、またはエラーによる異常終了）
> - CLOSING: 面接セッションの最終状態（フィードバック完了後、または定期クリーンアップ対象）

> ※ステータスはENUM型で管理し、アプリケーション全体で一貫性のある値を使用します。

> **ステータス遷移の補足**
> - 正常終了：MAIN → COMPLETED → (フィードバック) → CLOSING
> - 中断からの終了：PAUSED → TERMINATED
> - エラー時：任意のステータス → TERMINATED
> - 定期クリーンアップ：COMPLETED/TERMINATED → CLOSING
> - フィードバック完了：COMPLETED → CLOSING

> **補足**
> - セッション作成時のステータスは必ず`GREETING`から開始
> - 面接の進行順序：GREETING → [SELF_INTRODUCTION] → [ICE_BREAK] → MAIN → COMPLETED → (フィードバック) → CLOSING
> - 一時中断した場合：PAUSED → (再開) → 元のステータス or (中止) → TERMINATED

### 3.2 面接質問テーブル (interview_questions)
| カラム名 | 型 | 説明 | NULL |
|---------|-----|------|------|
| id | UUID | 主キー | NO |
| session_id | UUID | セッションID (FK) | NO |
| content | TEXT | 質問内容 | NO |
| sequence | INTEGER | 質問順序（1から開始） | NO |
| created_at | TIMESTAMP | 作成日時 | NO |

### 3.3 面接回答テーブル (interview_answers)
| カラム名 | 型 | 説明 | NULL |
|---------|-----|------|------|
| id | UUID | 主キー | NO |
| question_id | UUID | 質問ID (FK) | NO |
| content | TEXT | 回答内容 | NO |
| created_at | TIMESTAMP | 作成日時 | NO |

> **テーブル設計の補足**
> - 質問（interview_questions）と回答（interview_answers）を分離することで、未回答の質問の管理が容易になります
> - 質問の順序は`sequence`で明示的に管理し、ユーザーへの提示や分析に使用します
> - 進行状況は質問テーブルの`sequence`の最大値と`interview_sessions.question_count`の差分で把握します
> - 回答テーブルは質問への回答が存在する場合のみレコードが作成されます
> - 音声回答は一時的にテキスト変換のみに使用し、変換後の音声データは保持しません

## 4. API設計

### 4.1 面接セッションAPI

#### POST /api/v1/interview-sessions
面接セッションを開始し、面接練習画面へ遷移するためのセッション情報を作成

- リクエストボディ
```json
{
    "company_id": "uuid（任意）",
    "job_posting_id": "uuid（任意）",
    "interview_phase": "string（面接フェーズ。例：一次面接、最終面接など）",
    "interviewer_role": "string（面接官役職。例：人事担当、現場責任者など）",
    "question_count": "integer（5, 10, 15のいずれか）",
    "include_self_introduction": "boolean",
    "include_ice_break": "boolean"
}
```

- ステータスコード
  - 201: 作成成功
  - 400: リクエストパラメータ不正
  - 401: 認証エラー
  - 500: サーバーエラー

- バリデーション
  - `company_id`: 存在する企業IDであること（指定時）
  - `job_posting_id`: 
    - 指定時は存在する求人IDであること
    - 指定時は`company_id`が必須で、その企業に紐づく求人であること
  - `interview_phase`: 
    - 任意（nullまたは空文字を許容）
    - 指定時は最大文字数は100文字
  - `interviewer_role`: 
    - 任意（nullまたは空文字を許容）
    - 指定時は最大文字数は100文字
  - `question_count`: 
    - 必須
    - 5, 10, 15のいずれかであること
  - `include_ice_break`: 
    - 必須
    - true/falseのいずれかであること
  - `include_self_introduction`: 
    - 必須
    - true/falseのいずれかであること

> **補足**
> - セッション作成時のステータスは必ず`GREETING`から開始
> - 面接の進行順序：GREETING → [SELF_INTRODUCTION] → [ICE_BREAK] → MAIN → COMPLETED → (フィードバック) → CLOSING
>   - 一時中断した場合：PAUSED → (再開) → 元のステータス or (中止) → TERMINATED
> - `include_self_introduction`と`include_ice_break`の値に応じてリクエスト先のAPIが変化
>   - `include_self_introduction: true`の場合：greeting API → self-introduction API → ...
>   - `include_ice_break: true`の場合：[self-introduction API →] ice-break API → main API
>   - 両方`false`の場合：greeting API → main API
> - `company_id`と`job_posting_id`は任意だが、`job_posting_id`を指定する場合は`company_id`も必須
> - 作成成功時は自動的に面接練習画面へ遷移
> - `interview_phase`と`interviewer_role`は自由入力可能（画面上ではプルダウンと直接入力の併用）

#### GET /api/v1/interview-sessions/{session_id}/greeting
面接開始時の挨拶を取得

- レスポンスボディ
```json
{
    "question": {
        "id": "uuid",
        "content": "はじめまして。本日は面接にお時間をいただき、ありがとうございます。私は面接官の○○と申します。よろしくお願いいたします。",
        "sequence": 1
    },
    "next_status": "SELF_INTRODUCTION or ICE_BREAK or MAIN",
    "audio_enabled": true
}
```

- ステータスコード
  - 200: 取得成功
  - 401: 認証エラー
  - 404: セッションが存在しない
  - 409: セッションのステータスが不正（CREATED以外）
  - 500: サーバーエラー

> **補足**
> - セッションのステータスが`CREATED`の場合のみ呼び出し可能
> - 挨拶文は面接官の役職や面接フェーズに応じて適切に生成
> - `next_status`は設定値に応じて以下のように変化
>   - `include_self_introduction: true`の場合：SELF_INTRODUCTION
>   - `include_self_introduction: false`かつ`include_ice_break: true`の場合：ICE_BREAK
>   - 両方`false`の場合：MAIN
> - `audio_enabled`は音声読み上げの要否を示す（将来の拡張用）

#### POST /api/v1/interview-sessions/{session_id}/self-introduction
挨拶フェーズの回答を送信し、自己紹介の質問を取得

- リクエストボディ
```json
{
    "previous_answer": "挨拶への回答内容のテキスト",
    "current_status": "GREETING"
}
```

- レスポンス
```json
{
    "status": "SELF_INTRODUCTION",
    "next_question": {
        "id": "uuid",
        "content": "自己紹介の質問内容",
        "sequence": 2
    },
    "audio_enabled": true
}
```

- ステータスコード
  - 200: 送信成功
  - 400: リクエストパラメータ不正
  - 401: 認証エラー
  - 404: セッションが存在しない
  - 409: セッションのステータスが不正（GREETING以外）
  - 500: サーバーエラー

> **補足**
> - セッションのステータスが`GREETING`の場合のみ呼び出し可能
> - 回答内容は最大10000文字まで許容
> - `audio_enabled`は音声読み上げの要否を示す（将来の拡張用）

#### POST /api/v1/interview-sessions/{session_id}/ice-break
挨拶フェーズまたは自己紹介フェーズの回答を送信し、アイスブレイクの質問を取得

- リクエストボディ
```json
{
    "previous_answer": "挨拶または自己紹介の回答内容のテキスト",
    "current_status": "GREETING or SELF_INTRODUCTION"
}
```

- レスポンス
```json
{
    "status": "ICE_BREAK",
    "next_question": {
        "id": "uuid",
        "content": "アイスブレイクの質問内容",
        "sequence": 2 or 3
    },
    "audio_enabled": true
}
```

- ステータスコード
  - 200: 送信成功
  - 400: リクエストパラメータ不正
  - 401: 認証エラー
  - 404: セッションが存在しない
  - 409: セッションのステータスが不正（GREETINGまたはSELF_INTRODUCTION以外）
  - 500: サーバーエラー

> **補足**
> - セッションのステータスが`GREETING`または`SELF_INTRODUCTION`の場合のみ呼び出し可能
> - `include_self_introduction`の値によって呼び出し元のステータスが変化
>   - `true`の場合：SELF_INTRODUCTIONから呼び出し
>   - `false`の場合：GREETINGから呼び出し
> - 回答内容は最大10000文字まで許容
> - `sequence`は`include_self_introduction`の値に応じて変化
>   - `true`の場合：3（挨拶→自己紹介→アイスブレイク）
>   - `false`の場合：2（挨拶→アイスブレイク）
> - `audio_enabled`は音声読み上げの要否を示す（将来の拡張用）

#### POST /api/v1/interview-sessions/{session_id}/question
任意の問いかけフェーズの回答を送信し、次の質問を取得。設定された質問数に達した場合は終了フェーズへ移行。

- リクエストボディ
```json
{
    "previous_answer": "前のフェーズでの回答内容のテキスト",
    "current_status": "GREETING or SELF_INTRODUCTION or ICE_BREAK or MAIN"
}
```

- レスポンス（質問が残っている場合）
```json
{
    "status": "MAIN",
    "next_question": {
        "id": "uuid",
        "content": "次の質問内容",
        "sequence": 2 to N
    },
    "audio_enabled": true,
    "remaining_questions": 4,
    "should_end_session": false
}
```

- レスポンス（最後の質問への回答を受け取った場合）
```json
{
    "status": "COMPLETED",
    "next_question": null,
    "remaining_questions": 0,
    "should_end_session": true
}
```

- ステータスコード
  - 200: 送信成功
  - 400: リクエストパラメータ不正
  - 401: 認証エラー
  - 404: セッションが存在しない
  - 409: セッションのステータスが不正（GREETING、SELF_INTRODUCTION、ICE_BREAK、MAIN以外）
  - 500: サーバーエラー

> **補足**
> - セッションのステータスが`GREETING`、`SELF_INTRODUCTION`、`ICE_BREAK`、`MAIN`のいずれかの場合に呼び出し可能
> - `include_self_introduction`と`include_ice_break`の値によって呼び出し元のステータスが変化
>   - 両方`false`の場合：GREETINGから呼び出し
>   - `include_self_introduction: true`かつ`include_ice_break: false`の場合：SELF_INTRODUCTIONから呼び出し
>   - `include_ice_break: true`の場合：ICE_BREAKから呼び出し
>   - それ以外の場合：MAINから呼び出し
> - 回答内容は最大10000文字まで許容
> - `sequence`は前のフェーズに応じて変化
>   - GREETINGからの場合：2
>   - SELF_INTRODUCTIONからの場合：3
>   - ICE_BREAKからの場合：3または4
>   - MAINからの場合：4以降
> - 質問数のカウント方法
>   - 挨拶（1問）、自己紹介（設定に応じて1問）、アイスブレイク（設定に応じて1問）を含めて`question_count`を計算
>   - 例：`question_count: 5`、`include_self_introduction: true`、`include_ice_break: true`の場合
>     - 挨拶(1) + 自己紹介(1) + アイスブレイク(1) + 主質問(2)で合計5問
> - 最後の質問への回答を受け取った時点で、セッションのステータスを`COMPLETED`に更新
> - `should_end_session`が`true`の場合は次の質問を生成せず、クライアントは直接面接練習完了画面へ遷移
> - `audio_enabled`は音声読み上げの要否を示す（将来の拡張用）

### 4.2 音声処理API

#### POST /api/v1/speech-to-text
音声データをテキストに変換

- リクエストボディ
```json
{
    "audio_data": "base64エンコードされた音声データ"
}
```

- レスポンス
```json
{
    "text": "変換されたテキスト",
    "confidence": 0.95
}
```

#### POST /api/v1/text-to-speech
テキストを音声に変換

- リクエストボディ
```json
{
    "text": "音声化するテキスト"
}
```

- レスポンス
```json
{
    "audio_data": "base64エンコードされた音声データ"
}
```

## 5. AIプロンプト設計

### 5.1 面接官プロンプト
```