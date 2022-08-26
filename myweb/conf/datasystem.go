package conf

type AppData struct {
	AppId       int    `db:"aid"`
	AppName     string `db:"aname"`
	AppType     string `db:"atype"`
	CreateTime  string `db:"create_time"`
	DevelopPath string `db:"develop_path"`
	Ip          string `db:"ip"`
}
type LogData struct {
	AppId      int    `db:"aid"`
	AppName    string `db:"aname"`
	LogId      int    `db:"lid"`
	CreateTime string `db:"create_time"`
	LogPath    string `db:"log_path"`
	Topic      string `db:"topic"`
}

type EtcdData struct {
	Path  string
	Topic string
}
