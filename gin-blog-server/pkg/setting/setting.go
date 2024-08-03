package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	JwtSecret       string
	PageSize        int
	RuntimeRootPath string

	ImagePrefixUrl string
	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

var Cfg *ini.File

func Setup() {
	log.Println("setting Setup...")
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	err = Cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo AppSetting err: %v", err)
	}
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	log.Printf("AppSetting: %+v\n", AppSetting)

	err = Cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo ServerSetting err: %v", err)
	}
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

	log.Printf("ServerSetting: %+v\n", ServerSetting)

	err = Cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo DatabaseSetting err: %v", err)
	}
	log.Printf("DatabaseSetting: %+v\n", DatabaseSetting)

	err = Cfg.Section("redis").MapTo(RedisSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
	log.Printf("RedisSetting: %+v\n", RedisSetting)
}

// func LoadBase() {
// 	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
// 	// RunMode = Cfg.Section("").Key("RUN_MODE").String()
// }
// func LoadServer() {

// 	sec, err := Cfg.GetSection("server")
// 	if err != nil {
// 		log.Fatalf("Fail to get section 'server': %v", err)
// 	}
// 	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
// 	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second // 60s
// 	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second

// }

// func LoadApp() {
// 	sec, err := Cfg.GetSection("app")
// 	if err != nil {
// 		log.Fatalf("Fail to get section 'app': %v", err)
// 	}
// 	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
// 	JwtSecret = sec.Key("JWT_SECRET").MustString("23347$040412")
// }
