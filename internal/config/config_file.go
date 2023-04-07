package config

import (
	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

type API struct {
	ListenAddr     string
	PlannerFileDir string
}

type Redis struct {
	Addr     string
	Password string
	DB       int
}

type MySQL struct {
	Debug        bool
	Username     string
	Password     string
	Ip           string
	Port         int
	DatabaseName string
}

type InfluxDB struct {
	ServerURL string
	AuthToken string
	Org       string
	Bucket    string
}

type Log struct {
	Path  string
	Level string
}

type App struct {
	AppID string
	Token string
}

type Ftp struct {
	FtpUrl   string
	Username string
	Password string
	Timeout  int32 // ftpclient连接超时时间, 单位是秒
}

type Upgrade struct {
	UploadPath string
}
type Tecs struct {
	ListenAddr string
	Url        string
}

type CnfFile struct {
	API          API
	MySQL        MySQL
	IOTInfuxDB   InfluxDB
	IOTRedis     Redis
	SessionRedis Redis
	Log          Log
	Apps         []App
	Ftp          Ftp
	Upgrade      Upgrade
	Tecs         Tecs
}

var DefaultConfigFile string

func LoadConfigFile(configFile string) *CnfFile {
	var conf CnfFile
	if _, err := toml.DecodeFile(pathJoin(configFile), &conf); err != nil {
		logrus.Fatal(err, "decode config failed")
	}
	return &conf
}

// WatchConfigFile 监听文件改变
func WatchConfigFile(configFile string, fnc func()) {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Fatal(err)
	}
	defer watcher.Close()
	// Add a path.
	err = watcher.Add(pathJoin(configFile))
	if err != nil {
		logrus.Fatal(err)
	}
	// Start listening for events.
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			logrus.Println("event:", event)
			if event.Op == fsnotify.Write {
				fnc()
				return
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logrus.Println("error:", err)
		}
	}
}

func pathJoin(configFile string) string {
	if configFile == "" {
		runDir, err := os.Getwd()
		if err != nil {
			logrus.Fatal(err)
		}
		configFile = path.Join(runDir, "config", "tmios.conf")
	}
	return configFile
}
