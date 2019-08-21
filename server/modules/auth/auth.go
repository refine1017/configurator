package auth

import (
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"server/models"
)

func SignedInAdmin(ctx *macaron.Context, sess session.Store) *models.Admin {
	adminId := GetAdminId(sess)
	if adminId == "" {
		return nil
	}

	admin, err := models.GetAdminById(adminId)
	if err != nil {
		return nil
	}

	return admin
}
