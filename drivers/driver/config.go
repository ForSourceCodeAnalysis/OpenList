package driver

// Config 针对网盘类型的配置，是在代码中写死的
// 总的来说，这个配置项完全没有必要
// type Config struct {
// 	Name      string `json:"name"`
// 	LocalSort bool   `json:"local_sort"` // 是否本地排序，如果api本身没有提供排序功能，需要将这个值设置为true
// 	// 感觉这两个配置没有设置的必要，本身Storage里面已经有Proxy的配置了
// 	// OnlyLocal         bool   `json:"only_local"` // 是否仅使用本地代理
// 	// OnlyProxy         bool   `json:"only_proxy"` // 仅使用指定代理
// 	// 这个配置没有必要，因为Storage里面已经有缓存的配置了，每个存储都配置Config，不就相当于每个存储都又搞了一个Additional了嘛
// 	// NoCache           bool   `json:"no_cache"`
// 	// NoUpload          bool   `json:"no_upload"` 不启用上传功能
// 	// 这个配置没有使用
// 	// NeedMs            bool   `json:"need_ms"` // if need get message from user, such as validate code
// 	DefaultRoot string `json:"default_root"`
// 	// 莫名其妙
// 	CheckStatus bool `json:"-"`
// 	// 创建存储时，如果有这个值，会弹出提示信息，
// 	// 这种提示我觉得更应该放到文档里面，而不是代码里面写死，或者如果确实需要提示，应该做成可以配置的
// 	// Alert string `json:"alert"` //info,success,warning,danger

// 	//
// 	NoOverwriteUpload bool `json:"-"` // whether to support overwrite upload
// 	// cao 又是一个不明所以的项
// 	ProxyRangeOption bool `json:"-"`
// }

/// web代理
/// 网页预览、下载和直接链接是否通过中转。如果你打开此项，建议你设置site_url，以帮助alist更好的工作。
/// Web代理：是使用网页时候的策略，默认为本地代理，如果填写了代理URL并且启用了Web代理使用的是代理URL
/// 怎么理解官网的说明呢？
/// 经过验证，这里说的链接并不是说页面上复制链接时，直接复制成对应网盘链接，而是指预览，下载时是否要经过代理中转
/// 比如，预览文件时，首先请求 fs/get接口获取文件信息，文件信息中会返回文件直链
/// 1. 如果未开启web代理，那么会直接访问文件直链
/// 2. 如果开启了web代理，那么会访问 p/path，对应server/handles/down.go.Proxy方法
/// 下载同理，虽然复制的链接是openlist的地址，但是如果没有开启web代理，会302跳转到直链，如果开启了web代理，则会走代理传输流量
