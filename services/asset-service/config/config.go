package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

var (
	confPath string
	Conf     *Config
)

func init() {
	flag.StringVar(&confPath, "c", "", "config file path")
}

func Init() (err error) {
	file, err := ioutil.ReadFile(confPath)
	if err != nil {
		currentPath, _ := os.Getwd()
		fmt.Println("current path: ", currentPath, " given config path: ", confPath)
		return
	}
	err = yaml.Unmarshal(file, &Conf)
	return
}

type Duration time.Duration

// UnmarshalText unmarshal text to duration.
func (d *Duration) UnmarshalText(text []byte) error {
	tmp, err := time.ParseDuration(string(text))
	if err == nil {
		*d = Duration(tmp)
	}
	return err
}

// Config config
type Config struct {
	Log    *Log              `yaml:"log"`
	SqlMap map[string]*Mysql `yaml:"mysql"`
	Redis  *Redis            `yaml:"redis"`
}

type Log struct {
	Info        string `yaml:"info"`
	Error       string `yaml:"error"`
	Stdout      bool   `yaml:"stdout"`
	MysqlLog    bool   `yaml:"mysql"`
	Performance string `yaml:"performance"`
	Debug       bool   `yaml:"debug"`
}

type RPCServer struct {
	Network           string   `yaml:"network"`
	Addr              string   `yaml:"addr"`
	Timeout           Duration `yaml:"timeout"`
	IdleTimeout       Duration `yaml:"idleTimeout"`
	MaxLifeTime       Duration `yaml:"maxLifeTime"`
	ForceCloseWait    Duration `yaml:"ForceCloseWait"`
	KeepAliveInterval Duration `yaml:"KeepAliveInterval"`
	KeepAliveTimeout  Duration `yaml:"KeepAliveTimeout"`
	JaegerAddr        string   `yaml:"jaegerAddr"`
}

type Mysql struct {
	DSN         string   `yaml:"dsn"`
	MaxConn     int      `yaml:"max_connection"`
	MaxIdle     int      `yaml:"max_idle_connection"`
	MaxLifeTime Duration `yaml:"max_life_time"`
}

type Redis struct {
	Network      string   `yaml:"network"`
	Addr         string   `yaml:"addr"`
	Password     string   `yaml:"password"`
	DB           int      `yaml:"db"`
	DialTimeout  Duration `yaml:"dialTimeout"`
	ReadTimeout  Duration `yaml:"readTimeout"`
	WriteTimeout Duration `yaml:"writeTimeout"`
	PoolSize     int      `yaml:"poolsize"`
	MinIdleConns int      `yaml:"minIdleConns"`
	IdleTimeout  Duration `yaml:"idleTimeout"`
}

type GrpcClient struct {
}

func GetConfPath() string {
	return confPath
}
