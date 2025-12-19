package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	AppName     string
	Env         string
	Port        int
	Debug       bool
	LogLevel    string
	LogFilePath string
	Database    DatabaseConfig
	Redis       RedisConfig
	JWT         JWTConfig
	CORS        CORSConfig
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string
	ExpireHours int
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins []string
}

// AppConfig 全局配置实例
var AppConfig *Config

// InitConfig 初始化配置
func InitConfig() error {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")

	// 设置默认值
	setDefaults()

	// 自动读取环境变量
	viper.AutomaticEnv()

	// 读取配置文件（如果不存在也不报错）
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	AppConfig = &Config{
		AppName:     viper.GetString("APP_NAME"),
		Env:         viper.GetString("APP_ENV"),
		Port:        viper.GetInt("APP_PORT"),
		Debug:       viper.GetBool("APP_DEBUG"),
		LogLevel:    viper.GetString("LOG_LEVEL"),
		LogFilePath: viper.GetString("LOG_FILE_PATH"),
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
			TimeZone: viper.GetString("DB_TIMEZONE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetInt("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:      viper.GetString("JWT_SECRET"),
			ExpireHours: viper.GetInt("JWT_EXPIRE_HOURS"),
		},
		CORS: CORSConfig{
			AllowedOrigins: viper.GetStringSlice("CORS_ALLOWED_ORIGINS"),
		},
	}

	return nil
}

// setDefaults 设置默认配置
func setDefaults() {
	viper.SetDefault("APP_NAME", "Shoppee")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", 8080)
	viper.SetDefault("APP_DEBUG", true)
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_FILE_PATH", "./logs/app.log")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "shoppee")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_TIMEZONE", "Asia/Shanghai")

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)

	viper.SetDefault("JWT_SECRET", "your-super-secret-key")
	viper.SetDefault("JWT_EXPIRE_HOURS", 24)

	viper.SetDefault("CORS_ALLOWED_ORIGINS", []string{"*"})
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
		c.Database.TimeZone,
	)
}

// GetRedisAddr 获取Redis地址
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}
