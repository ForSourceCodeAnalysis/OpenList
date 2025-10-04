package _123_open

// 参考： https://123yunpan.yuque.com/org-wiki-123yunpan-muaork/cr6ced/ppsuasz6rpioqbyt
// 官方对接口有QPS限制，但是实现起来很麻烦，所以这里不做限制，接口如果返回了429（请求太频繁）错误，则暂停1秒再请求
const (
	baseURL = "https://open-api.123pan.com"

	// API： POST 域名 +/api/v1/access_token
	accessTokenAPI = "/api/v1/access_token"
	// API： POST   域名 + /api/v1/oauth2/access_token
	refreshTokenAPI = "/api/v1/oauth2/access_token"
	// API： POST   域名 + /upload/v1/file/mkdir
	mkdirAPI = "/upload/v1/file/mkdir"

	// 分片上传 v1 ------------------------------------------
	// API： POST   域名 + /upload/v1/file/create
	uploadCreateV1API = "/upload/v1/file/create"
	//API： POST 域名 + /upload/v1/file/get_upload_url
	getUploadURLAPI = "/upload/v1/file/get_upload_url"
	// API： POST 域名 + /upload/v1/file/list_upload_parts
	listUploadPartsAPI = "/upload/v1/file/list_upload_parts"
	//API： POST   域名 + /upload/v1/file/upload_complete
	uploadCompleteV1API = "/upload/v1/file/upload_complete"
	//API： POST   域名 + /upload/v1/file/upload_async_result
	uploadAsyncAPI = "/upload/v1/file/upload_async_result"

	// 分片上传 v2 ------------------------------------------
	//API： POST   域名 + /upload/v2/file/create
	uploadCreateV2API = "/upload/v2/file/create"
	//API： POST   上传域名 + /upload/v2/file/slice
	uploadSliceAPI = "/upload/v2/file/slice"
	// API： POST   域名 + /upload/v2/file/upload_complete
	uploadCompleteV2API = "/upload/v2/file/upload_complete"

	// 单步上传 ------------------------------------------
	//API： GET 域名 + /upload/v2/file/domain
	uploadDomainAPI = "/upload/v2/file/domain"
	//API： POST   上传域名 + /upload/v2/file/single/create
	uploadSingleCreateAPI = "/upload/v2/file/single/create"

	// 重命名---------------------------------------------
	// API：PUT 域名 + /api/v1/file/name
	renameAPI = "/api/v1/file/name"
	// API： POST 域名 + /api/v1/file/rename
	batchRenameAPI = "/api/v1/file/rename"

	// 删除---------------------------------------------
	// API： POST 域名 + /api/v1/file/trash
	trashAPI = "/api/v1/file/trash"
	//API： POST 域名 + /api/v1/file/recover
	recoverAPI = "/api/v1/file/recover"
	// API： POST 域名 + /api/v1/file/delete
	deleteAPI = "/api/v1/file/delete"

	// 文件详情---------------------------------------------
	// API： GET 域名 + /api/v1/file/detail
	detailAPI = "/api/v1/file/detail"
	// API：POST 域名 + /api/v1/file/infos
	infosAPI = "/api/v1/file/infos"

	// 文件列表---------------------------------------------
	// API： GET 域名 + /api/v2/file/list  注意：此接口查询结果包含回收站的文件，需自行根据字段trashed判断处理
	fileListAPI = "/api/v2/file/list"
	// API： GET 域名 + /api/v1/file/list
	fileListV1API = "/api/v1/file/list"

	// 移动---------------------------------------------
	// API： POST 域名 + /api/v1/file/move
	moveAPI = "/api/v1/file/move"

	// 下载---------------------------------------------
	// API：GET 域名 + /api/v1/file/download_info
	downloadInfoAPI = "/api/v1/file/download_info"

	// 分享管理---------------------------------------------
	// API： POST 域名 + /api/v1/share/content-payment/create
	createContentPaymentShare = "/api/v1/share/content-payment/create"
	// API： POST 域名 + /api/v1/share/create
	createShareAPI = "/api/v1/share/create"
	// API： PUT 域名 + /api/v1/share/list/info
	updateShareInfo = "/api/v1/share/list/info"
	// API： GET 域名 + /api/v1/share/list
	listShareAPI = "/api/v1/share/list"

	// 离线下载---------------------------------------------
	// API： POST 域名 + /api/v1/offline/download
	offlineDownloadAPI = "/api/v1/offline/download"
	//API： GET 域名 + /api/v1/offline/download/process
	offlineDownloadProcessAPI = "/api/v1/offline/download/process"

	// 获取用户信息---------------------------------------------
	// API： GET 域名 + /api/v1/user/info
	userInfoAPI = "/api/v1/user/info"

	// 直链---------------------------------------------
	// IP黑名单管理--------------------------------------
	// API： POST 域名 + /api/v1/developer/config/forbide-ip/switch
	switchForbideIPAPI = "/api/v1/developer/config/forbide-ip/switch"
	// API： POST 域名 + /api/v1/developer/config/forbide-ip/update
	updateForbideIPAPI = "/api/v1/developer/config/forbide-ip/update"
	// API： GET 域名 + /api/v1/developer/config/forbide-ip/list
	listForbideIPAPI = "/api/v1/developer/config/forbide-ip/list"
	// ---------------------------------------------------
	// API：GET 域名 + /api/v1/direct-link/offline/logs
	dlOfflineLogsAPI = "/api/v1/direct-link/offline/logs"
	// API：GET 域名 + /api/v1/direct-link/log
	dlLogAPI = "/api/v1/direct-link/log"
	// API： POST 域名 + /api/v1/direct-link/enable
	enableDlAPI = "/api/v1/direct-link/enable"
	//API： GET 域名 + /api/v1/direct-link/url
	dlURLAPI = "/api/v1/direct-link/url"
	//API： POST 域名 + /api/v1/direct-link/disable
	disableDlAPI = "/api/v1/direct-link/disable"

	// 图床---------------------------------------------
	//TODO
	// 视频转码---------------------------------------------
	//TODO

)
