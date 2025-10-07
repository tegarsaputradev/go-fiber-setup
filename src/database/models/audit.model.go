package models

type AuditLog struct {
	BaseModel
	EntityModel string `json:"entity_model" gorm:"size:255"`
	EntityID    uint   `json:"entity_id" gorm:"index"`
	UserID      uint   `json:"user_id" gorm:"index"`
	Diff        string `json:"diff" gorm:"type:json"`
	Action      string `json:"action" gorm:"size:20"`
	OldData     string `json:"old_data" gorm:"type:json"`
	NewData     string `json:"new_data" gorm:"type:json"`
}
