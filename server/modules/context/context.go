package context

import (
	"fmt"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/sirupsen/logrus"
	"gopkg.in/macaron.v1"
	"html"
	"html/template"
	"net/http"
	"path"
	"server/models"
	"server/modules/auth"
	"server/modules/code"
	"server/modules/setting"
	"strings"
	"time"
)

type Context struct {
	*macaron.Context
	Cache   cache.Cache
	csrf    csrf.CSRF
	Flash   *session.Flash
	Session session.Store

	admin *models.Admin

	Link        string // current request URL
	EscapedLink string
}

// Contexter initializes a classic context for a request.
func Contexter() macaron.Handler {
	return func(c *macaron.Context, cache cache.Cache, sess session.Store, f *session.Flash, x csrf.CSRF) {
		ctx := &Context{
			Context: c,
			Cache:   cache,
			csrf:    x,
			Flash:   f,
			Session: sess,
			Link:    setting.AppSubURL + strings.TrimSuffix(c.Req.URL.EscapedPath(), "/"),
		}
		c.Data["Link"] = ctx.Link
		ctx.Data["PageStartTime"] = time.Now()

		ctx.admin = auth.SignedInAdmin(ctx.Context, ctx.Session)

		// If request sends files, parse them here otherwise the Query() can't be parsed and the CsrfToken will be invalid.
		if ctx.Req.Method == "POST" && strings.Contains(ctx.Req.Header.Get("Content-Type"), "multipart/form-data") {
			if err := ctx.Req.ParseMultipartForm(setting.AttachmentMaxSize << 20); err != nil && !strings.Contains(err.Error(), "EOF") { // 32MB max size
				ctx.Refuse(code.POST_SIZE_LIMITED, "ParseMultipartForm", err)
				return
			}
		}

		ctx.Data["CsrfToken"] = html.EscapeString(x.GetToken())
		ctx.Data["CsrfTokenHtml"] = template.HTML(`<input type="hidden" name="_csrf" value="` + ctx.Data["CsrfToken"].(string) + `">`)
		logrus.Debugf("Session ID: %s", sess.ID())
		logrus.Debugf("CSRF Token: %v", ctx.Data["CsrfToken"])

		c.Map(ctx)
	}
}

func (ctx *Context) GetAdmin() *models.Admin {
	return ctx.admin
}

func (ctx *Context) GetAdminName() string {
	if ctx.admin != nil {
		return ctx.admin.Username
	}
	return ""
}

func (ctx *Context) Refuse(code int, message string, err error) {
	var res struct {
		Code    int         `json:"code"`
		Message interface{} `json:"message"`
	}

	var fullMsg = message
	if err != nil {
		fullMsg = fmt.Sprintf("%s(err: %v)", message, err)
	}

	logrus.Errorf(fullMsg)
	if macaron.Env != macaron.PROD {
		res.Message = fullMsg
	} else {
		res.Message = message
	}

	res.Code = code
	res.Message = fullMsg

	ctx.JSON(http.StatusOK, res)
}

func (ctx *Context) Ack(data interface{}) {
	var res struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
	}

	res.Code = code.SUCCESS
	res.Data = data

	ctx.JSON(http.StatusOK, res)
}

func (ctx *Context) SetFileHeader(filename string) {
	ctx.Resp.Header().Set("Content-Description", "File Transfer")
	ctx.Resp.Header().Set("Content-Type", "application/octet-stream")
	ctx.Resp.Header().Set("Content-Disposition", "attachment; filename="+filename)
	ctx.Resp.Header().Set("Content-Transfer-Encoding", "binary")
	ctx.Resp.Header().Set("Expires", "0")
	ctx.Resp.Header().Set("Cache-Control", "must-revalidate")
	ctx.Resp.Header().Set("Pragma", "public")
}

// ServeFile serves given file to response.
func (ctx *Context) ServeFile(file string, names ...string) {
	var name string
	if len(names) > 0 {
		name = names[0]
	} else {
		name = path.Base(file)
	}
	ctx.Resp.Header().Set("Content-Description", "File Transfer")
	ctx.Resp.Header().Set("Content-Type", "application/octet-stream")
	ctx.Resp.Header().Set("Content-Disposition", "attachment; filename="+name)
	ctx.Resp.Header().Set("Content-Transfer-Encoding", "binary")
	ctx.Resp.Header().Set("Expires", "0")
	ctx.Resp.Header().Set("Cache-Control", "must-revalidate")
	ctx.Resp.Header().Set("Pragma", "public")
	http.ServeFile(ctx.Resp, ctx.Req.Request, file)
}
