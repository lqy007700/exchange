package config

import (
	"github.com/go-micro/plugins/v4/config/encoder/yaml"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/reader"
	"go-micro.dev/v4/config/reader/json"
	"go-micro.dev/v4/config/source/file"
	"go-micro.dev/v4/logger"
	"time"
)

// Config config
type Config struct {
	Micro     *Micro            `yaml:"micro"`
	Log       *Log              `yaml:"log"`
	SqlMap    map[string]*Mysql `yaml:"mysql"`
	Redis     *Redis            `yaml:"redis"`
	RPCServer *RPCServer        `yaml:"rpcServer"`
	Kafka     *Kafka            `yaml:"kafka"`
}

type Micro struct {
	Name string `json:"name"`
}

type Log struct {
	Info   string `yaml:"info"`
	Error  string `yaml:"error"`
	Stdout bool   `yaml:"stdout"`
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

type Kafka struct {
	Brokers []string `yaml:"brokers"`
}

var (
	confPath = "/Users/lqy007700/Data/code/go-application/exchange/services/asset-service/config/config.yaml"
	Conf     *Config
)

func Init() (err error) {
	enc := yaml.NewEncoder()
	c, _ := config.NewConfig(config.WithReader(
		json.NewReader( // json reader for internal config merge
			reader.WithEncoder(enc),
		),
	))

	err = c.Load(file.NewSource(file.WithPath(confPath)))
	if err != nil {
		logger.Errorf("load config error: %v", err)
		return err
	}

	// read a database host
	if err := c.Scan(&Conf); err != nil {
		logger.Errorf("scan config error: %v", err)
		return err
	}

	d := map[string]*Mysql{}
	err = c.Get("mysql").Scan(&d)
	if err != nil {
		return err
	}
	Conf.SqlMap = d

	InitLogger()
	return nil
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
