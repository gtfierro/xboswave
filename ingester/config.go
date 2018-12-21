package main

import (
	"encoding/base64"
	"log"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	//"github.com/op/go-logging"
	logrus "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type BTrDBConfig struct {
	Address string
}
type InfluxDBConfig struct {
	Address  string
	Database string
	Username string
	Password string
}

type Database struct {
	BTrDB    *BTrDBConfig
	InfluxDB *InfluxDBConfig
}

type Config struct {
	Database Database

	WAVEMQ struct {
		EntityFile         string
		SiteRouter         string
		SubscriptionExpiry int64
		ClientID           string
	}

	Output struct {
		LogLevel string
	}
}

func init() {
	prefix := os.Getenv("GOPATH")
	// switch prefix to default GOPATH /home/{user}/go
	if prefix == "" {
		u, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		prefix = filepath.Join(u.HomeDir, "go")
	}
	// set defaults for config
	viper.SetDefault("Output.LogLevel", "debug")
	// Database
	// no btrdb default
	viper.SetDefault("Database.InfluxDB.Address", "http://127.0.0.1:8086")
	viper.SetDefault("Database.InfluxDB.Database", "xbos")
	viper.SetDefault("Database.InfluxDB.Username", "")
	viper.SetDefault("Database.InfluxDB.Password", "")

	viper.SetDefault("WAVEMQ.EntityFile", "")
	viper.SetDefault("WAVEMQ.SiteRouter", "127.0.0.1:4516")
	viper.SetDefault("WAVEMQ.SubscriptionExpiry", 60)

	randombytes := make([]byte, 8)
	rand.Read(randombytes)
	randomid := base64.URLEncoding.EncodeToString(randombytes)
	viper.SetDefault("WAVEMQ.ClientID", randomid)

}

func getCfg() *Config {
	//level, err := logging.LogLevel(viper.GetString("LogLevel"))
	//if err != nil {
	//	level = logging.DEBUG
	//}

	level, err := logrus.ParseLevel(viper.GetString("Output.LogLevel"))
	if err != nil {
		level = logrus.DebugLevel
	}
	logrus.SetLevel(level)

	cfg := &Config{
		Database: Database{},
	}

	if viper.IsSet("Database.BTrDB") {
		cfg.Database.BTrDB = &BTrDBConfig{
			Address: viper.GetString("Database.BTrDB.Address"),
		}
	} else if viper.IsSet("Database.InfluxDB") {
		cfg.Database.InfluxDB = &InfluxDBConfig{
			Address:  viper.GetString("Database.InfluxDB.Address"),
			Database: viper.GetString("Database.InfluxDB.Database"),
			Username: viper.GetString("Database.InfluxDB.Username"),
			Password: viper.GetString("Database.InfluxDB.Password"),
		}
	}

	cfg.WAVEMQ.EntityFile = viper.GetString("WAVEMQ.EntityFile")
	cfg.WAVEMQ.SiteRouter = viper.GetString("WAVEMQ.SiteRouter")
	cfg.WAVEMQ.SubscriptionExpiry = int64(viper.GetInt("WAVEMQ.SubscriptionExpiry"))
	cfg.WAVEMQ.ClientID = viper.GetString("WAVEMQ.ClientID")

	return cfg
}

func ReadConfig(file string) (*Config, error) {
	if len(file) > 0 {
		viper.SetConfigFile(file)
	}
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	viper.AutomaticEnv()

	return getCfg(), nil
}

func ReadConfigFromString(configString string) (*Config, error) {
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(strings.NewReader(configString)); err != nil {
		return nil, err
	}
	viper.AutomaticEnv()

	return getCfg(), nil
}
