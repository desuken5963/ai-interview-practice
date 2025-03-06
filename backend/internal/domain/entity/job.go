package entity

import "time"

// JobPosting は求人情報を表すエンティティです
type JobPosting struct {
	ID           int              `json:"id" gorm:"primaryKey"`
	CompanyID    int              `json:"company_id" gorm:"not null"`
	Title        string           `json:"title" gorm:"not null;type:varchar(100)"`
	Description  *string          `json:"description" gorm:"type:text"`
	CustomFields []JobCustomField `json:"custom_fields" gorm:"foreignKey:JobID"`
	CreatedAt    time.Time        `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time        `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// JobCustomField は求人の追加情報を表すエンティティです
type JobCustomField struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	JobID     int       `json:"job_id" gorm:"not null"`
	FieldName string    `json:"field_name" gorm:"not null;type:varchar(50)"`
	Content   string    `json:"content" gorm:"not null;type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// JobResponse は求人情報のレスポンス形式を表します
type JobResponse struct {
	Jobs  []JobPosting `json:"jobs"`
	Total int64        `json:"total"`
	Page  int          `json:"page"`
	Limit int          `json:"limit"`
}
