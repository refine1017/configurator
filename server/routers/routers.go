package routers

import (
	"github.com/go-macaron/bindata"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"server/modules/context"
	"server/modules/setting"
	"server/public"
	"server/routers/app"
	"server/routers/config"
	"server/routers/manage"
	"server/routers/system"
	"server/routers/user"
)

func routers(m *macaron.Macaron) {
	bindIgnErr := binding.BindIgnErr

	reqAdminLogin := context.ReqAdminLogin()
	reqAdminLog := context.ReqAdminLog()

	m.Group("v1", func() {

		m.Group("/user", func() {
			m.Post("/login", bindIgnErr(user.SignInForm{}), user.Login)
			m.Post("/logout", user.Logout)
			m.Get("/info", reqAdminLogin, user.Info)
			m.Get("/projects", reqAdminLogin, user.Projects)
		})

		m.Group("/app", func() {
			m.Get("/setting", app.Setting)
		}, reqAdminLogin)

		m.Group("/admin", func() {
			m.Get("/list", bindIgnErr(system.AdminQueryParams{}), system.FetchAdminPageList)
			m.Post("/create", bindIgnErr(system.AdminAddInfo{}), system.CreateAdminRow)
			m.Post("/:id", bindIgnErr(system.AdminEditInfo{}), system.UpdateAdminRow)
			m.Post("/:id/delete", system.DeleteAdminRow)
			m.Get("/logs", bindIgnErr(system.AdminLogQueryParams{}), system.FetchAdminLogPageList)
		}, reqAdminLogin)

		m.Group("/project", func() {
			m.Get("/list", bindIgnErr(manage.ProjectQueryParams{}), manage.FetchProjectPageList)
			m.Post("/create", bindIgnErr(manage.ProjectAddInfo{}), manage.CreateProjectRow)
			m.Post("/:id", bindIgnErr(manage.ProjectEditInfo{}), manage.UpdateProjectRow)
			m.Post("/:id/delete", manage.DeleteProjectRow)
		}, reqAdminLogin)

		m.Group("/environment", func() {
			m.Get("/list", bindIgnErr(manage.EnvironmentQueryParams{}), manage.FetchEnvironmentPageList)
			m.Post("/:projectId/create", bindIgnErr(manage.EnvironmentAddInfo{}), manage.CreateEnvironmentRow)
			m.Post("/:id", bindIgnErr(manage.EnvironmentEditInfo{}), manage.UpdateEnvironmentRow)
			m.Post("/:id/delete", manage.DeleteEnvironmentRow)
			m.Post("/:id/copy", bindIgnErr(manage.EnvironmentCopyInfo{}), manage.CopyEnvironmentRow)
		}, reqAdminLogin)

		m.Group("/config", func() {
			m.Get("/list", bindIgnErr(manage.ConfigQueryParams{}), manage.FetchConfigPageList)
			m.Post("/:envId/create", bindIgnErr(manage.ConfigAddInfo{}), manage.CreateConfigRow)
			m.Post("/:envId/:index", bindIgnErr(manage.ConfigEditInfo{}), manage.UpdateConfigRow)
			m.Post("/:envId/:index/delete", manage.DeleteConfigRow)
			m.Post("/:envId/push/:config/:server", manage.PushConfig)
			m.Post("/:envId/merge/:config/info", manage.MergeConfigInfo)
			m.Post("/:envId/merge/:config/", manage.MergeConfig)
		}, reqAdminLogin)

		m.Group("/config", func() {
			m.Get("/:envId/json/:config", manage.ExportConfigJson)
			m.Get("/:envId/lua/:config", manage.ExportConfigLua)
			m.Get("/:envId/excel", manage.ExportConfigExcel)
			m.Get("/:envId/proto", manage.ExportConfigProto)
			m.Get("/:envId/entitas", manage.ExportConfigEntitas)
			m.Get("/:envId/cs", manage.ExportConfigCS)
		})

		m.Group("/server", func() {
			m.Get("/list", bindIgnErr(manage.ServerQueryParams{}), manage.FetchServerPageList)
			m.Post("/:projectId/create", bindIgnErr(manage.ServerAddInfo{}), manage.CreateServerRow)
			m.Post("/:id", bindIgnErr(manage.ServerEditInfo{}), manage.UpdateServerRow)
			m.Post("/:id/delete", manage.DeleteServerRow)
		}, reqAdminLogin)

		m.Group("/fields", func() {
			m.Get("/list", bindIgnErr(manage.FieldQueryParams{}), manage.FetchFieldPageList)
			m.Post("/:envId/:configIndex/create", bindIgnErr(manage.FieldAddInfo{}), manage.CreateFieldRow)
			m.Post("/:envId/:configIndex/:index", bindIgnErr(manage.FieldEditInfo{}), manage.UpdateFieldRow)
			m.Post("/:envId/:configIndex/:index/delete", manage.DeleteFieldRow)
		}, reqAdminLogin)

		m.Group("/table", func() {
			m.Get("/list", bindIgnErr(config.TableQueryParams{}), config.FetchTablePageList)
			m.Post("/:envId/:collect", config.CreateConfigTableRow)
			m.Post("/:envId/:collect/upload", config.UploadConfigTable)
			m.Post("/:envId/:collect/:id", config.UpdateConfigTableRow)
			m.Post("/:envId/:collect/:id/delete", config.DeleteConfigTableRow)
		}, reqAdminLogin)

		m.Group("/table", func() {
			m.Get("/:envId/:collect/json", config.ExportConfigTableJsonFile)
			m.Get("/:envId/:collect/lua", config.ExportConfigTableLuaFile)
			m.Get("/:envId/:collect/excel", config.ExportConfigTableExcelFile)
		})

		m.Group("/kv", func() {
			m.Get("/list", bindIgnErr(config.KVQueryParams{}), config.FetchKVPageList)
			m.Post("/:envId/:collect", config.CreateConfigKVRow)
			m.Post("/:envId/:collect/upload", config.UploadConfigKV)
			m.Post("/:envId/:collect/:id", config.UpdateConfigKVRow)
			m.Post("/:envId/:collect/:id/delete", config.DeleteConfigKVRow)
		}, reqAdminLogin)

		m.Group("/kv", func() {
			m.Get("/:envId/:collect/json", config.ExportConfigKVJsonFile)
			m.Get("/:envId/:collect/lua", config.ExportConfigKVLuaFile)
			m.Get("/:envId/:collect/excel", config.ExportConfigKVExcelFile)
		})

		m.Group("/json", func() {
			m.Get("/list", bindIgnErr(config.JsonQueryParams{}), config.FetchJsonPageList)
			m.Post("/:envId/:collect", config.CreateConfigJsonRow)
			m.Post("/:envId/:collect/upload", config.UploadConfigJson)
			m.Post("/:envId/:collect/:id", config.UpdateConfigJsonRow)
			m.Post("/:envId/:collect/:id/delete", config.DeleteConfigJsonRow)
			m.Get("/:envId/:collect/:id/data", config.GetConfigJsonData)
			m.Post("/:envId/:collect/:id/data", config.SetConfigJsonData)
		}, reqAdminLogin)

		m.Group("/json", func() {
			m.Get("/:envId/:collect/export", config.ExportConfigJsonZipFile)
			m.Get("/:envId/:collect/:id/export", config.ExportConfigJsonData)
		})
	}, reqAdminLog)
}

func Startup() {
	m := macaron.Classic()

	routers(m)

	m.Use(cache.Cacher(cache.Options{
		Adapter:       setting.CacheService.Adapter,
		AdapterConfig: setting.CacheService.Conn,
		Interval:      setting.CacheService.Interval,
	}))
	m.Use(session.Sessioner(setting.SessionConfig))
	m.Use(csrf.Csrfer(csrf.Options{
		Secret:         setting.SecretKey,
		Cookie:         setting.CSRFCookieName,
		SetCookie:      true,
		Secure:         setting.SessionConfig.Secure,
		CookieHttpOnly: true,
		Header:         "X-Csrf-Token",
		CookiePath:     setting.AppSubURL,
	}))

	m.Use(macaron.Static("public",
		macaron.StaticOptions{
			FileSystem: bindata.Static(bindata.Options{
				Asset:      public.Asset,
				AssetDir:   public.AssetDir,
				AssetNames: public.AssetNames,
				Prefix:     "",
			}),
		},
	))

	m.Use(macaron.Renderer())
	m.Use(context.Contexter())
	m.SetAutoHead(true)

	m.Run(setting.Router.Host, setting.Router.Port)
}
