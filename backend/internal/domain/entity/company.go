package entity

import "time"

// Company は企業情報を表すエンティティです
type Company struct {
	ID                  int                  `json:"id" gorm:"primaryKey"`
	Name                string               `json:"name" gorm:"not null"`
	BusinessDescription *string              `json:"business_description" gorm:"type:text"`
	CustomFields        []CompanyCustomField `json:"custom_fields" gorm:"foreignKey:CompanyID"`
	JobPostings         []JobPosting         `json:"job_postings" gorm:"foreignKey:CompanyID"`
	CreatedAt           time.Time            `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt           time.Time            `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// CompanyCustomField は企業の追加情報を表すエンティティです
type CompanyCustomField struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	CompanyID int       `json:"company_id" gorm:"not null"`
	FieldName string    `json:"field_name" gorm:"not null;type:varchar(50)"`
	Content   string    `json:"content" gorm:"not null;type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// CompanyResponse は企業情報のレスポンス形式を表します
type CompanyResponse struct {
	Companies []Company `json:"companies"`
	Total     int       `json:"total"`
	Page      int       `json:"page"`
	Limit     int       `json:"limit"`
}
