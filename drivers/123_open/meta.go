package open123

import (
	"time"

	"github.com/OpenListTeam/OpenList/v4/internal/driver"
)

// var _ driver.IDriver = (*Open123)(nil)

// Open123 123pan open api 配置
// 参考： https://123yunpan.yuque.com/org-wiki-123yunpan-muaork/cr6ced/hpengmyg32blkbg8
type Open123 struct {
	// 个人开发者（client_id）/第三方授权应用(app_id),使用时需根据refresh_token判断具体是哪个，下同
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"` // 个人开发者（client_secret）/第三方授权应用(secret_id)
	AccessToken  string    `json:"access_token"`
	ExpiredAt    time.Time `json:"expired_at"`    // token过期时间
	RefreshToken string    `json:"refresh_token"` // 仅第三方授权应用
	Scope        string    `json:"scope"`         // 仅第三方授权应用
	TokenType    string    `json:"token_type"`    // 仅第三方授权应用
	// Username string `json:"username"`
	// Password string `json:"password"`
	qpsInstance map[string]*APIInfo // QPS限制是和账号绑定的
}

// Init  初始化
func (d *Open123) Init() error {
	d.qpsInstance = map[string]*APIInfo{
		accessTokenAPI:    initAPIInfo(baseURL+accessTokenAPI, 1),
		refreshTokenAPI:   initAPIInfo(baseURL+refreshTokenAPI, 0),
		userInfoAPI:       initAPIInfo(baseURL+userInfoAPI, 1),
		fileListAPI:       initAPIInfo(baseURL+fileListAPI, 3),
		downloadInfoAPI:   initAPIInfo(baseURL+downloadInfoAPI, 0),
		mkdirAPI:          initAPIInfo(baseURL+mkdirAPI, 2),
		moveAPI:           initAPIInfo(baseURL+moveAPI, 1),
		renameAPI:         initAPIInfo(baseURL+renameAPI, 0),
		trashAPI:          initAPIInfo(baseURL+trashAPI, 1),
		preupCreateAPI:    initAPIInfo(baseURL+preupCreateAPI, 0),
		sliceUploadAPI:    initAPIInfo(sliceUploadAPI, 0),
		uploadCompleteAPI: initAPIInfo(baseURL+uploadCompleteAPI, 0),
		uploadURLAPI:      initAPIInfo(baseURL+uploadURLAPI, 0),
		singleUploadAPI:   initAPIInfo(singleUploadAPI, 0),
	}
	return nil
}

// Name 获取驱动名称
func (d *Open123) Name() string {
	return driver.DriverName123Open
}

// SortSupported 是否支持排序
func (d *Open123) SortSupported() bool {
	return false
}

// IsAccessTokenExpired 检测 AccessToken 是否过期
func (d *Open123) IsAccessTokenExpired() bool {
	return time.Now().After(d.ExpiredAt)
}

// ConfigMeta 获取配置参数
func ConfigMeta() []driver.FieldMeta {
	return []driver.FieldMeta{
		driver.FieldMeta{
			Name:     "client_id",
			Type:     driver.TypeText,
			Required: true,
			Help:     "个人开发者（client_id）或第三方授权应用（app_id）",
		},
		driver.FieldMeta{
			Name:     "client_secret",
			Type:     driver.TypeText,
			Required: true,
			Help:     "个人开发者（client_secret）或第三方授权应用（secret_id）",
		},
		driver.FieldMeta{
			Name:     "refresh_token",
			Type:     driver.TypeText,
			Required: false,
			Help:     "第三方授权应用所需的refresh_token",
		},
	}
}
