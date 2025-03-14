package entity

import "time"

// JobPosting は求人情報の作成リクエストを表すエンティティです
type JobPosting struct {
	ID           int              `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID    int              `json:"company_id"`
	Title        string           `json:"title"`
	Description  *string          `json:"description,omitempty"`
	CustomFields []JobCustomField `json:"custom_fields" gorm:"foreignKey:JobID"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// JobCustomField は求人のカスタムフィールドを表すエンティティです
type JobCustomField struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	JobID     int       `json:"job_id"`
	FieldName string    `json:"field_name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// JobResponse は求人一覧のレスポンスを表すエンティティです
type JobResponse struct {
	JobPostings []JobPosting `json:"jobPostings"`
	Total       int          `json:"total"`
	Page        int          `json:"page"`
	Limit       int          `json:"limit"`
}
