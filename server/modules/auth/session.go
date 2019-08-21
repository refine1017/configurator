package auth

import (
	"github.com/go-macaron/session"
	"server/models"
)

func GetAdminId(sess session.Store) string {
	adminId := sess.Get("adminId")
	if adminId != nil {
		if v, ok := adminId.(string); ok {
			return v
		}
	}

	return ""
}

func SetAdmin(sess session.Store, admin *models.Admin) {
	_ = sess.Set("adminId", admin.Id.Hex())
	_ = sess.Set("adminName", admin.Username)
}

func DelAdmin(sess session.Store) {
	_ = sess.Delete("adminId")
	_ = sess.Delete("adminName")
}
