package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// constants
const (
	AccountService     = "ACCOUNT_SERVICE"
	TransactionService = "TRANSACTION_SERVICE"
	ProcessorService   = "PROCESSOR_SERVICE"
	LedgerService      = "LEDGER_SERVICE"
)

// Config holds the overall application configuration.
type Config struct {
	Services map[string]*ServiceConfig
	Database *DatabaseConfig
	Kafka    *KafkaConfig
	ApiAuth  *ApiAuth
}

// ServiceConfig holds the configuration for a specific service.
type ServiceConfig struct {
	Port     string `env:"PORT"`
	LogLevel string `env:"LOG_LEVEL" default:"info"`
}

// DatabaseConfig holds the database configuration.
type DatabaseConfig struct {
	PostgresConnectionString string `env:"POSTGRES_CONNECTION_STRING"`
	MongoDBConnectionString  string `env:"MONGODB_CONNECTION_STRING"`
}

// KafkaConfig holds the kafka configuration.
type KafkaConfig struct {
	Brokers string `env:"KAFKA_BROKERS" default:""`
}

type ApiAuth struct {
	UserName string `env:"API_AUTH_USERNAME" default:""`
	Password string `env:"API_AUTH_PASSWORD" default:""`
}

// LoadConfig loads the overall application configuration.
func LoadConfig() (*Config, error) {
	if os.Getenv("IS_LOCAL") == "true" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	services := make(map[string]*ServiceConfig)

	// Load configuration for each service
	serviceNames := []string{"ACCOUNT_SERVICE", "TRANSACTION_SERVICE", "PROCESSOR_SERVICE", "LEDGER_SERVICE"}
	for _, serviceName := range serviceNames {
		serviceConfig, err := LoadServiceConfig(serviceName)
		if err != nil {
			return nil, err
		}
		services[serviceName] = serviceConfig
	}

	return &Config{
		Services: services,
		Database: LoadDatabaseConfig(),
		Kafka:    LoadKafkaConfig(),
		ApiAuth:  LoadApiAuthConfig(),
	}, nil
}

// LoadServiceConfig loads the configuration for a specific service.
func LoadServiceConfig(serviceName string) (*ServiceConfig, error) {
	var config ServiceConfig

	v := reflect.ValueOf(&config).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		envTag := field.Tag.Get("env")
		defaultValue := field.Tag.Get("default")

		envVar := os.Getenv(serviceName + "_" + strings.ToUpper(envTag)) // e.g., ACCOUNT_SERVICE_PORT
		if envVar == "" && defaultValue != "" {
			envVar = defaultValue
		}

		if envVar != "" {
			switch field.Type.Kind() {
			case reflect.String:
				v.Field(i).SetString(envVar)
			case reflect.Int:
				intValue, err := strconv.Atoi(envVar)
				if err != nil {
					return nil, err
				}
				v.Field(i).SetInt(int64(intValue))
			case reflect.Bool:
				boolValue, err := strconv.ParseBool(envVar)
				if err != nil {
					return nil, err
				}
				v.Field(i).SetBool(boolValue)
			default:
				return nil, fmt.Errorf("unknown env format %v", v)
			}
		}
	}

	return &config, nil
}

// LoadDatabaseConfig loads the database configuration.
func LoadDatabaseConfig() *DatabaseConfig {
	postgresConnStr := os.Getenv("POSTGRES_CONNECTION_STRING")

	mongoConnStr := os.Getenv("MONGODB_CONNECTION_STRING")

	return &DatabaseConfig{
		PostgresConnectionString: postgresConnStr,
		MongoDBConnectionString:  mongoConnStr,
	}
}

// LoadKafkaConfig loads the kafka configuration.
func LoadKafkaConfig() *KafkaConfig {
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")

	return &KafkaConfig{
		Brokers: kafkaBrokers,
	}
}

func LoadApiAuthConfig() *ApiAuth {
	uname := os.Getenv("API_AUTH_USERNAME")
	pass := os.Getenv("API_AUTH_PASSWORD")

	return &ApiAuth{
		UserName: uname,
		Password: pass,
	}
}
