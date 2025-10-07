package callbacks

import (
	"encoding/json"
	"go-rest-setup/src/database/models"

	"gorm.io/gorm"
)

func LogCreate(db *gorm.DB) {
	if db.Error != nil || db.Statement.Schema == nil {
		return
	}
	if db.Statement.Schema != nil && db.Statement.Schema.Table == "audit_logs" {
		return
	}

	entityModel := db.Statement.Schema.Table
	entityValue := db.Statement.Dest

	newData, _ := json.Marshal(entityValue)

	var userID uint
	if val, ok := db.Statement.Context.Value("user_id").(uint); ok {
		userID = val
	}

	audit := models.AuditLog{
		EntityModel: entityModel,
		UserID:      userID,
		Action:      "CREATE",
		OldData:     "{}",
		NewData:     string(newData),
		Diff:        "{}",
	}

	db.Session(&gorm.Session{NewDB: true}).Create(&audit)
}
