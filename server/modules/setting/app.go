package setting

var (
	AppName        string
	AppVer         string
	AppMode        string
	AppURL         string
	AppSubURL      string
	AppSubURLDepth int // Number of slashes
	AppPath        string
	AppConfPath    string
	AppDataPath    string
	AppWorkPath    string
)

const (
	APP_MODE_WEB = "web"
)

func parseApp() error {

	return nil
}
