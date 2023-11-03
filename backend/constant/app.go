package constant

const (
	Running     = "Running"
	UnHealthy   = "UnHealthy"
	Error       = "Error"
	Stopped     = "Stopped"
	Installing  = "Installing"
	DownloadErr = "DownloadErr"
	Upgrading   = "Upgrading"
	UpgradeErr  = "UpgradeErr"
	Rebuilding  = "Rebuilding"
	Syncing     = "Syncing"
	SyncSuccess = "SyncSuccess"
	SyncErr     = "SyncErr"

	ContainerPrefix = "1Panel-"

	AppNormal   = "Normal"
	AppTakeDown = "TakeDown"

	AppOpenresty = "openresty"
	AppMysql     = "mysql"
	AppMariaDB   = "mariadb"
	AppRedis     = "redis"

	AppResourceLocal  = "local"
	AppResourceRemote = "remote"

	CPUS          = "CPUS"
	MemoryLimit   = "MEMORY_LIMIT"
	HostIP        = "HOST_IP"
	ContainerName = "CONTAINER_NAME"
)

type AppOperate string

var (
	Up      AppOperate = "up"
	Down    AppOperate = "down"
	Start   AppOperate = "start"
	Stop    AppOperate = "stop"
	Restart AppOperate = "restart"
	Delete  AppOperate = "delete"
	Sync    AppOperate = "sync"
	Backup  AppOperate = "backup"
	Restore AppOperate = "restore"
	Update  AppOperate = "update"
	Rebuild AppOperate = "rebuild"
	Upgrade AppOperate = "upgrade"
)
