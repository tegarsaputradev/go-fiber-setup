package audit

import (
	"go-rest-setup/src/core/audit/callbacks"

	"gorm.io/gorm"
)

func RegisterAuditLogCallbacks(db *gorm.DB) {
	db.Callback().Create().After("gorm:create").Register("audit_log_create", callbacks.LogCreate)
	db.Callback().Update().After("gorm:update").Register("audit_log_update", callbacks.LogUpdate)
	db.Callback().Delete().After("gorm:delete").Register("audit_log_delete", callbacks.LogDelete)
}
