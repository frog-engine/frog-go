package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

// AppConfig 定义配置结构
type Config struct {
	Server struct {
		Addr         string `yaml:"addr"`
    Port         string `yaml:"port"`
		ReadTimeout  string `yaml:"read_timeout"`
		WriteTimeout string `yaml:"write_timeout"`
		Env          string `yaml:"env"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`

	Cache struct {
		Type  string `yaml:"type"`
		Redis struct {
			Addr     string `yaml:"addr"`
			Password string `yaml:"password"`
			DB       int    `yaml:"db"`
		} `yaml:"redis"`
		MaxSize    int `yaml:"max_size"`
		Expiration int `yaml:"expiration"`
	} `yaml:"cache"`

	Image struct {
		MaxWidth        int      `yaml:"max_width"`
		MaxHeight       int      `yaml:"max_height"`
		AllowedFormats  []string `yaml:"allowed_formats"`
		Quality         int      `yaml:"quality"`
		DownloadTimeout string   `yaml:"download_timeout"`
		ProcessTimeout  string   `yaml:"process_timeout"`
	} `yaml:"image"`

	Log struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
		Output string `yaml:"output"`
	} `yaml:"log"`
}

var (
	configInstance *Config
	configLock     sync.Mutex
)

// LoadConfig 从指定文件加载配置
func LoadConfig[T any](filePath string) (*T, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file: %w", err)
	}
	defer file.Close()

	var cfg T
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config file: %w", err)
	}
	return &cfg, nil
}

// GetConfig 获取全局配置实例
func GetConfig() (*Config) {
	// 如果配置实例为空，则加载配置
	if configInstance == nil {
    env := "test"
    Init(&env)
	}
	return configInstance
}

func Init(env *string) {
	configLock.Lock()
	defer configLock.Unlock()

	// 如果配置实例为空，则加载配置
	if configInstance == nil {
		var filePath string
		if *env == "production" {
			filePath = "config/config_prod.yaml"
		} else if *env == "test" {
			filePath = "config/config_test.yaml"
		} else {
			filePath = "config/config_test.yaml"
		}

		cfg, err := LoadConfig[Config](filePath); 
    if err == nil {
      configInstance = cfg
      log.Println("Configuration initialized successfully.")
    } else {
      log.Fatalf("Failed to initialize config: %v", err)
    }
  }
}