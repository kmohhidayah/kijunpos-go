package config

import (
	"fmt"
	"github/kijunpos/config/db"

	"github.com/spf13/viper"
)

type (
	App struct {
		Port    int
		Env     string
		Name    string
		Version string
	}
	Otel struct {
		ApiKey      string
		Env         string
		ServiceName string
		URL         string
		Insecure    bool
		IsEnabled   bool
	}

	Config struct {
		App       App
		Otel      Otel
		Databases []db.Config
	}
)

var configData *Config

func LoadConfig() error {
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error when read config")
		return err
	}

	configData = &Config{
		App: App{
			Port:    getRequiredInt("APP_PORT"),
			Env:     getRequiredString("APP_ENV"),
			Name:    getRequiredString("APP_NAME"),
			Version: getRequiredString("APP_VERSION"),
		},
		Otel: Otel{
			ApiKey:      getRequiredString("OTEL_API_KEY"),
			Env:         getRequiredString("APP_ENV"),
			ServiceName: fmt.Sprintf("%s-%s", getRequiredString("APP_NAME"), getRequiredString("APP_ENV")),
			URL:         getRequiredString("OTEL_URL"),
			Insecure:    getRequiredBool("OTEL_INSECURE"),
			IsEnabled:   getRequiredBool("OTEL_IS_ENABLED"),
		},
		Databases: []db.Config{
			{
				Name:        db.KIJUNDB,
				HostURL:     getRequiredString("KIJUNDB_URL"),
				MaxIdleConn: getRequiredInt("KIJUNDB_MAX_IDLE_CONNECTIONS"),
				MaxOpenConn: getRequiredInt("KIJUNDB_MAX_OPEN_CONNECTIONS"),
			},
		},
	}

	return nil
}

func GetConfig() *Config {
	return configData
}

// ConfigValue adalah interface constraint untuk tipe nilai yang didukung
type ConfigValue interface {
	string | int | bool
}

// getRequired adalah fungsi generic untuk mengambil nilai konfigurasi yang required
func getRequired[T ConfigValue](key string, getter func(string) T) T {
	if viper.IsSet(key) {
		return getter(key)
	}
	panic(fmt.Errorf("KEY %s IS MISSING", key))
}

// Helper functions untuk memudahkan penggunaan
func getRequiredString(key string) string {
	return getRequired(key, viper.GetString)
}

func getRequiredInt(key string) int {
	return getRequired(key, viper.GetInt)
}

func getRequiredBool(key string) bool {
	return getRequired(key, viper.GetBool)
}
