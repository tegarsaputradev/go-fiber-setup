package callbacks

import (
	"encoding/json"
	"go-rest-setup/src/database/models"
	"reflect"

	"gorm.io/gorm"
)

func LogUpdate(db *gorm.DB) {
	if db.Error != nil || db.Statement.Schema == nil {
		return
	}

	if db.Statement.Schema != nil && db.Statement.Schema.Table == "audit_logs" {
		return
	}

	entityModel := db.Statement.Schema.Table
	entityValue := db.Statement.Dest

	// ambil primary key
	pk := getPrimaryKeyValue(db)
	if pk == nil {
		return
	}

	oldData := map[string]interface{}{}
	db.Session(&gorm.Session{NewDB: true}).Table(entityModel).Where("id = ?", pk).Take(&oldData)

	oldJSON, _ := json.Marshal(oldData)
	newJSON, _ := json.Marshal(entityValue)

	diff := diffJSON(oldData, entityValue)

	var userID uint
	if val, ok := db.Statement.Context.Value("user_id").(uint); ok {
		userID = val
	}

	audit := models.AuditLog{
		EntityModel: entityModel,
		EntityID:    pk.(uint),
		UserID:      userID,
		Action:      "UPDATE",
		OldData:     string(oldJSON),
		NewData:     string(newJSON),
		Diff:        string(diff),
	}

	if audit.Diff == "" {
		audit.Diff = "{}"
	}

	db.Session(&gorm.Session{NewDB: true}).Create(&audit)
}

func getPrimaryKeyValue(db *gorm.DB) interface{} {
	if db.Statement == nil || db.Statement.Schema == nil {
		return nil
	}

	for _, field := range db.Statement.Schema.Fields {
		if field.PrimaryKey {
			value, _ := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
			return value
		}
	}
	return nil
}

func diffJSON(oldVal interface{}, newVal interface{}) []byte {
	oldMap := structToMap(oldVal)
	newMap := structToMap(newVal)

	before := map[string]interface{}{}
	after := map[string]interface{}{}

	for key, newValue := range newMap {
		oldValue, exists := oldMap[key]
		if !exists || !reflect.DeepEqual(oldValue, newValue) {
			before[key] = oldValue
			after[key] = newValue
		}
	}

	// tambahkan field yang dihapus
	for key, oldValue := range oldMap {
		if _, exists := newMap[key]; !exists {
			before[key] = oldValue
			after[key] = nil
		}
	}

	diff := map[string]interface{}{
		"before": before,
		"after":  after,
	}

	out, _ := json.Marshal(diff)
	return out
}

func structToMap(v interface{}) map[string]interface{} {
	b, _ := json.Marshal(v)
	var result map[string]interface{}
	_ = json.Unmarshal(b, &result)
	return result
}
