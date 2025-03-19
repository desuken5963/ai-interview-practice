# 企業求人管理登録機能

## 1. 概要
企業情報および求人情報の登録・編集・削除を行う機能の基本設計を記載する。

## 2. 画面設計

### 2.1 企業/求人管理ページ
- 企業一覧をカード形式で表示
- 企業情報の新規登録ボタン
- 各企業カードに求人情報の件数を表示
- 各企業カードに面接練習開始ボタンを配置
- ページネーション機能を実装

### 2.2 企業情報登録/編集モーダル
- 企業名（必須）入力フォーム
- 事業内容（任意）入力エリア
- カスタム項目の動的追加/削除機能
- 保存/キャンセルボタン
- バリデーションの即時フィードバック

### 2.3 求人一覧モーダル
- 企業ごとの求人情報一覧を表示
- 新規登録ボタン
- 各求人の編集/削除ボタン
- ページネーション機能を実装

### 2.4 求人情報登録/編集モーダル
- 求人タイトル（必須）入力フォーム
- 仕事内容（任意）入力エリア
- カスタム項目の動的追加/削除機能
- 保存/キャンセルボタン
- バリデーションの即時フィードバック

## 3. データ構造

### 3.1 企業情報テーブル (companies)
| カラム名 | 型 | 説明 | NULL | 制約 |
|---------|-----|------|------|------|
| id | INT | 主キー（自動採番） | NO | PRIMARY KEY |
| name | VARCHAR(100) | 企業名 | NO | UNIQUE |
| business_description | TEXT | 事業内容 | YES | - |
| created_at | TIMESTAMP | 作成日時 | NO | - |
| updated_at | TIMESTAMP | 更新日時 | NO | - |

### 3.2 企業追加情報テーブル (company_custom_fields)
| カラム名 | 型 | 説明 | NULL | 制約 |
|---------|-----|------|------|------|
| id | INT | 主キー（自動採番） | NO | PRIMARY KEY |
| company_id | INT | 企業ID (FK) | NO | FOREIGN KEY |
| field_name | VARCHAR(50) | 項目名 | NO | - |
| content | TEXT | 内容 | NO | - |
| created_at | TIMESTAMP | 作成日時 | NO | - |
| updated_at | TIMESTAMP | 更新日時 | NO | - |

### 3.3 求人情報テーブル (job_postings)
| カラム名 | 型 | 説明 | NULL | 制約 |
|---------|-----|------|------|------|
| id | INT | 主キー（自動採番） | NO | PRIMARY KEY |
| company_id | INT | 企業ID (FK) | NO | FOREIGN KEY |
| title | VARCHAR(100) | 求人タイトル | NO | UNIQUE |
| description | TEXT | 仕事内容 | YES | - |
| created_at | TIMESTAMP | 作成日時 | NO | - |
| updated_at | TIMESTAMP | 更新日時 | NO | - |

### 3.4 求人追加情報テーブル (job_custom_fields)
| カラム名 | 型 | 説明 | NULL | 制約 |
|---------|-----|------|------|------|
| id | INT | 主キー（自動採番） | NO | PRIMARY KEY |
| job_id | INT | 求人ID (FK) | NO | FOREIGN KEY |
| field_name | VARCHAR(50) | 項目名 | NO | - |
| content | TEXT | 内容 | NO | - |
| created_at | TIMESTAMP | 作成日時 | NO | - |
| updated_at | TIMESTAMP | 更新日時 | NO | - |

## 4. API設計

### 4.1 共通エラーレスポンス形式
```json
{
    "error": {
        "code": "ERROR_CODE",
        "message": "エラーメッセージ",
        "details": [
            {
                "field": "エラーが発生したフィールド名",
                "message": "詳細なエラーメッセージ"
            }
        ]
    }
}
```

### 4.2 企業求人情報API

#### GET /api/v1/companies
企業求人情報の取得
- 全ての企業情報とそれに紐付く求人情報を取得します

- レスポンスパターン
  1. 企業情報のみ
  ```json
  {
      "companies": [
          {
              "id": 1,
              "name": "企業名",
              "business_description": "事業内容",
              "custom_fields": [],
              "job_postings": [],
              "created_at": "2024-01-01T00:00:00Z",
              "updated_at": "2024-01-01T00:00:00Z"
          }
      ],
      "total": 100,
      "page": 1,
      "limit": 10
  }
  ```

  2. 企業情報とカスタムフィールド
  ```json
  {
      "companies": [
          {
              "id": 1,
              "name": "企業名",
              "business_description": "事業内容",
              "custom_fields": [
                  {
                      "id": 1,
                      "field_name": "企業理念",
                      "content": "企業理念の内容",
                      "created_at": "2024-01-01T00:00:00Z",
                      "updated_at": "2024-01-01T00:00:00Z"
                  }
              ],
              "job_postings": [],
              "created_at": "2024-01-01T00:00:00Z",
              "updated_at": "2024-01-01T00:00:00Z"
          }
      ],
      "total": 100,
      "page": 1,
      "limit": 10
  }
  ```

  3. 企業情報と求人情報
  ```json
  {
      "companies": [
          {
              "id": 1,
              "name": "企業名",
              "business_description": "事業内容",
              "custom_fields": [],
              "job_postings": [
                  {
                      "id": 1,
                      "company_id": 1,
                      "title": "求人タイトル",
                      "description": "仕事内容",
                      "custom_fields": [],
                      "created_at": "2024-01-01T00:00:00Z",
                      "updated_at": "2024-01-01T00:00:00Z"
                  }
              ],
              "created_at": "2024-01-01T00:00:00Z",
              "updated_at": "2024-01-01T00:00:00Z"
          }
      ],
      "total": 100,
      "page": 1,
      "limit": 10
  }
  ```

  4. 企業情報と求人情報（求人にカスタムフィールドあり）
  ```json
  {
      "companies": [
          {
              "id": 1,
              "name": "企業名",
              "business_description": "事業内容",
              "custom_fields": [],
              "job_postings": [
                  {
                      "id": 1,
                      "company_id": 1,
                      "title": "求人タイトル",
                      "description": "仕事内容",
                      "custom_fields": [
                          {
                              "id": 1,
                              "field_name": "必要なスキル",
                              "content": "スキルの内容",
                              "created_at": "2024-01-01T00:00:00Z",
                              "updated_at": "2024-01-01T00:00:00Z"
                          }
                      ],
                      "created_at": "2024-01-01T00:00:00Z",
                      "updated_at": "2024-01-01T00:00:00Z"
                  }
              ],
              "created_at": "2024-01-01T00:00:00Z",
              "updated_at": "2024-01-01T00:00:00Z"
          }
      ],
      "total": 100,
      "page": 1,
      "limit": 10
  }
  ```

- クエリパラメータ
  | フィールド | ルール | デフォルト値 |
  |------------|--------|--------------|
  | page | 任意, 1以上の整数 | 1 |
  | limit | 任意, 1-100の整数 | 6 |

- ステータスコード
  - 200: 取得成功
  - 400: 不正なクエリパラメータ
  - 500: サーバーエラー

#### POST /api/v1/companies
企業情報の作成

- バリデーションルール
  | フィールド | ルール |
  |------------|--------|
  | name | 必須, 1-100文字 |
  | business_description | 任意, 最大1000文字 |
  | custom_fields | 任意, 配列 |
  | custom_fields[].field_name | 配列内の場合必須, 1-50文字 |
  | custom_fields[].content | 配列内の場合必須, 最大500文字 |

- ステータスコード
  - 201: 作成成功
  - 400: バリデーションエラー
  - 500: サーバーエラー

- リクエストボディ
```json
{
    "name": "企業名",
    "business_description": "事業内容",
    "custom_fields": [
        {
            "field_name": "企業理念",
            "content": "企業理念の内容"
        }
    ]
}
```

#### PUT /api/v1/companies/{id}
企業情報の更新
- リクエストボディ: POST と同様

#### DELETE /api/v1/companies/{id}
企業情報と関連する求人情報を削除
- レスポンス: 204 No Content

### 4.3 求人情報API

#### POST /api/v1/job-postings
求人情報の作成

- バリデーションルール
  | フィールド | ルール |
  |------------|--------|
  | company_id | 必須, 存在する企業ID |
  | title | 必須, 1-100文字 |
  | description | 任意, 最大1000文字 |
  | custom_fields | 任意, 配列 |
  | custom_fields[].field_name | 配列内の場合必須, 1-50文字 |
  | custom_fields[].content | 配列内の場合必須, 最大500文字 |

- ステータスコード
  - 201: 作成成功
  - 400: バリデーションエラー
  - 500: サーバーエラー

- リクエストボディ
```json
{
    "company_id": 1,
    "title": "求人タイトル",
    "description": "仕事内容",
    "custom_fields": [
        {
            "field_name": "必要なスキル",
            "content": "スキルの内容"
        }
    ]
}
```

#### PUT /api/v1/job-postings/{id}
求人情報の更新
- リクエストボディ: POST と同様

#### DELETE /api/v1/job-postings/{id}
求人情報の削除
- レスポンス: 204 No Content