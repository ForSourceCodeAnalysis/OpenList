package driver

// 网盘驱动名称
const (
	//DriverName115 115pan
	DriverName115 = "115"
	// DriverName115Open 115pan 使用open api
	DriverName115Open = "115Open"
	// DriverName115Share 115pan 分享链接
	DriverName115Share = "115Share"
	// DriverName123 123pan
	DriverName123 = "123"
	// DriverName123Link 123pan 直链
	DriverName123Link = "123Link"
	// DriverName123Open 123pan 使用open api
	DriverName123Open = "123Open"
	// DriverName123Share 123pan 分享链接
	DriverName123Share = "123Share"
	// DriverName139 中国移动网盘
	DriverName139 = "139"
	// DriverName189 天翼云盘
	DriverName189 = "189"
	// DriverName189PC 天翼云盘个人盘
	DriverName189PC = "189PC"
	// DriverNameAlias 别名存储
	DriverNameAlias = "Alias"
	// DriverNameAliyundrive 阿里云盘
	DriverNameAliyundrive = "Aliyundrive"
	// DriverNameAliyundriveOpen 阿里云盘使用open api
	DriverNameAliyundriveOpen = "AliyundriveOpen"
	// DriverNameAliyundriveShare 阿里云盘分享链接
	DriverNameAliyundriveShare = "AliyundriveShare"
	// DriverNameAzureBlob azure blob
	DriverNameAzureBlob = "AzureBlob"
	// DriverNameBaiduNetdisk 百度网盘
	DriverNameBaiduNetdisk = "BaiduNetdisk"
	// DriverNameBaiduPhoto 百度一刻相册
	DriverNameBaiduPhoto = "BaiduPhoto"
	// DriverNameBaiduShare 百度网盘分享链接
	DriverNameBaiduShare = "BaiduShare"
	// DriverNameChaoxing 超星云盘
	DriverNameChaoxing = "Chaoxing"
	// DriverNameCloudreve cloudreve
	DriverNameCloudreve = "Cloudreve"
	// DriverNameCloudreveV4 cloudreve v4
	DriverNameCloudreveV4 = "CloudreveV4"
	// DriverNameCrypt 加密存储
	DriverNameCrypt = "Crypt"
	// DriverNameDoubao 豆包AI网盘
	DriverNameDoubao = "Doubao"
	// DriverNameDoubaoShare 豆包分享链接
	DriverNameDoubaoShare = "DoubaoShare"
	// DriverNameDropbox dropbox
	DriverNameDropbox = "Dropbox"
	// DriverNameFebbox febbox
	DriverNameFebbox = "Febbox"
	// DriverNameFtp ftp
	DriverNameFtp = "Ftp"
	// DriverNameGithub github存储
	DriverNameGithub = "Github"
	// DriverNameGithubReleases github release
	DriverNameGithubReleases = "GithubReleases"
	// DriverNameGoogleDrive 谷歌云盘
	DriverNameGoogleDrive = "GoogleDrive"
	// DriverNameGooglePhoto 谷歌相册
	DriverNameGooglePhoto = "GooglePhoto"
	// DriverNameHalalcloud halalcloud
	DriverNameHalalcloud = "Halalcloud"
	// DriverNameIlanzou 蓝奏云优享版
	DriverNameIlanzou = "Ilanzou"
	// DriverNameIpfsAPI ipfs
	DriverNameIpfsAPI = "IpfsAPI"
	// DriverNameKodbox kodbox 可道云
	DriverNameKodbox = "Kodbox"
	// DriverNameLanzou 蓝奏云
	DriverNameLanzou = "Lanzou"
	// DriverNameLenovoNasShare 联想nas分享
	DriverNameLenovoNasShare = "LenovonasShare"
	// DriverNameLocal 本地存储
	DriverNameLocal = "Local"
	// DriverNameMediatrack 分秒帧
	DriverNameMediatrack = "Mediatrack"
	// DriverNameMega  mega
	DriverNameMega = "Mega"
	// DriverNameMisskey misskey
	DriverNameMisskey = "Misskey"
	// DriverNameMopan 知行魔盘
	DriverNameMopan = "Mopan"
	// DriverNameNeteaseMusic 网易云音乐
	DriverNameNeteaseMusic = "NeteaseMusic"
	// DriverNameOnedrive onedrive
	DriverNameOnedrive = "Onedrive"
	// DriverNameOnedriveApp onedrive app
	DriverNameOnedriveApp = "OnedriveApp"
	// DriverNameOnedriveSharelink onedrive分享链接
	DriverNameOnedriveSharelink = "OnedriveSharelink"
	// DriverNameOpenlist openlist挂载
	DriverNameOpenlist = "Openlist"
	// DriverNamePikpak pikpak
	DriverNamePikpak = "Pikpak"
	// DriverNamePikpakShare pikpak分享链接
	DriverNamePikpakShare = "PikpakShare"
	// DriverNameQuarkOpen 夸克网盘 open api
	DriverNameQuarkOpen = "QuarkOpen"
	// DriverNameQuarkUc uc网盘
	DriverNameQuarkUc = "QuarkUc"
	// DriverNameQuarkUcTv uc 网盘 tv版
	DriverNameQuarkUcTv = "QuarkUcTv"
	// DriverNameS3 amazon s3
	DriverNameS3 = "S3"
	// DriverNameSeafile seafile
	DriverNameSeafile = "Seafile"
	// DriverNameSftp sftp
	DriverNameSftp = "Sftp"
	// DriverNameSmb smb
	DriverNameSmb = "Smb"
	// DriverNameTeambition 钉盘
	DriverNameTeambition = "Teambition"
	// DriverNameTerabox terabox
	DriverNameTerabox = "Terabox"
	// DriverNameThunder 迅雷云盘
	DriverNameThunder = "Thunder"
	// DriverNameThunderBrowser 迅雷浏览器
	DriverNameThunderBrowser = "ThunderBrowser"
	// DriverNameThunderX 迅雷X
	DriverNameThunderX = "ThunderX"
	// DriverNameTrainbit trainbit
	DriverNameTrainbit = "Trainbit"
	// DriverNameURLTree url tree
	DriverNameURLTree = "URLTree"
	// DriverNameVirtual 虚拟存储
	DriverNameVirtual = "Virtual"
	// DriverNameWebdav webdav
	DriverNameWebdav = "Webdav"
	// DriverNameWeiyun 微云
	DriverNameWeiyun = "Weiyun"
	// DriverNameWopan 沃盘
	DriverNameWopan = "Wopan"
	// DriverNameYandexDisk yandex disk
	DriverNameYandexDisk = "YandexDisk"
)

// 字段类型定义，方便前端识别
const (
	// TypeString string
	TypeString = "string"
	// TypeSelect select
	TypeSelect = "select"
	// TypeBool bool
	TypeBool = "bool"
	// TypeText text
	TypeText = "text"
	// TypeNumber number
	TypeNumber = "number"
	// TypeArray array
	TypeArray = "array"
)

// hash 类型
const (
	HashTypeMD5   = "md5"
	HashTypeSHA1  = "sha1"
	HashTypeSH256 = "sha256"
)
