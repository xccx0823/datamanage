package datawatch

type Option func(*SourceDataWatcher)

// SetOptions 设置参数
func (sdw *SourceDataWatcher) SetOptions(serverId uint32, options ...Option) {
	sdw.ServerID = serverId
	for _, option := range options {
		option(sdw)
	}
}

func WithDB(host string, port uint16, user, password string) Option {
	return func(sdw *SourceDataWatcher) {
		sdw.Host = host
		sdw.Port = port
		sdw.User = user
		sdw.Password = password
	}
}

func WithCharset(charset string) Option {
	return func(sdw *SourceDataWatcher) {
		sdw.Charset = charset
	}
}

func WithMonitorSyncTime(syncTime int) Option {
	return func(sdw *SourceDataWatcher) {
		sdw.monitorSyncTime = syncTime
	}
}
