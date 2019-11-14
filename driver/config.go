package driver

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// Driver configuration struct
type Config struct {
	// base64 encoded namespace
	Namespace string
	// path to the entity file for this driver
	EntityFile string
	// local site router address
	SiteRouter string
	// default report rate
	ReportRate time.Duration
	Params     map[string]interface{}
}

func (cfg Config) GetString(key string) string {
	if cfg.Params == nil {
		log.Fatalf("No params in %v", cfg)
	}
	if v, found := cfg.Params[key]; !found {
		log.Fatalf("No key %s found in cfg %v", key, cfg)
	} else if val, ok := v.(string); !ok {
		log.Fatalf("Key %s was not a string. Value was %v, type was %T", key, v, v)
	} else {
		return val
	}
	return ""
}

func (cfg Config) GetInt(key string) int {
	if cfg.Params == nil {
		log.Fatalf("No params in %v", cfg)
	}
	if v, found := cfg.Params[key]; !found {
		log.Fatalf("No key %s found in cfg %v", key, cfg)
	} else if val, ok := v.(int); !ok {
		log.Fatalf("Key %s was not a int. Value was %v, type was %T", key, v, v)
	} else {
		return val
	}
	return -1
}

func (cfg Config) GetStringSlice(key string) []string {
	if cfg.Params == nil {
		log.Fatalf("No params in %v", cfg)
	}
	if v, found := cfg.Params[key]; !found {
		log.Fatalf("No key %s found in cfg %v", key, cfg)
	} else if val, ok := v.([]interface{}); !ok {
		log.Fatalf("Key %s was not a slice. Value was %v, type was %T", key, v, v)
	} else {
		var ret []string
		for _, newv := range val {
			if s, ok := newv.(string); ok {
				ret = append(ret, s)
			} else {
				ret = append(ret, fmt.Sprintf("%v", newv))
			}
		}
		return ret
	}
	return []string{}
}

func ReadConfigFromEnviron() (Config, error) {
	cfg := Config{}
	//read configuration from environment; fill in with defaults
	cfg.Namespace = os.Getenv("XBOS_DEFAULT_NAMESPACE")
	if cfg.Namespace == "" {
		return cfg, fmt.Errorf("Could not get XBOS_DEFAULT_NAMESPACE")
	}

	cfg.EntityFile = os.Getenv("XBOS_ENTITY_FILE")
	if cfg.EntityFile == "" {
		return cfg, fmt.Errorf("Could not get WAVE_DEFAULT_ENTITY")
	}

	cfg.SiteRouter = os.Getenv("XBOS_SITE_ROUTER")
	if cfg.SiteRouter == "" {
		cfg.SiteRouter = "127.0.0.1:4516"
	}

	var err error
	cfg.ReportRate, err = time.ParseDuration(os.Getenv("XBOS_REPORT_RATE"))
	return cfg, err
}

func ReadConfigFromFile(filename string) (Config, error) {
	var cfg Config
	var s = struct {
		Namespace  string
		EntityFile string
		SiteRouter string
		ReportRate string
		Params     map[string]interface{}
	}{}
	if _, err := toml.DecodeFile(filename, &s); err != nil {
		return cfg, err
	}
	if s.SiteRouter == "" {
		cfg.SiteRouter = "127.0.0.1:4516"
	} else {
		cfg.SiteRouter = s.SiteRouter
	}
	cfg.Namespace = s.Namespace
	cfg.EntityFile = s.EntityFile
	cfg.Params = s.Params
	var err error
	cfg.ReportRate, err = time.ParseDuration(s.ReportRate)
	return cfg, err
}
