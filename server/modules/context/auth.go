package context

import (
	"gopkg.in/macaron.v1"
	"net/http"
	"server/models"
)

func ReqAdminLogin() macaron.Handler {
	return func(ctx *Context) {
		if ctx.admin == nil {
			ctx.Refuse(http.StatusForbidden*100, "Not login", nil)
			return
		}
	}
}

func ReqAdminLog() macaron.Handler {
	return func(ctx *Context) {
		if ctx.admin == nil {
			return
		}

		models.RecordAdminLog(ctx.GetAdminName(), "router", ctx.Req.RequestURI)
	}
}
