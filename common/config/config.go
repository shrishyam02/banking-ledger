package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

//constants
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
	Brokers           string `env:"KAFKA_BROKERS" default:"kafka:9092"`
	SchemaRegistryURL string `env:"SCHEMA_REGISTRY_URL" default:"http://schema-registry:8081"`
}

// LoadConfig loads the overall application configuration.
func LoadConfig() (*Config, error) {
	services := make(map[string]*ServiceConfig)
	dbConfig, err := LoadDatabaseConfig()
	if err != nil {
		return nil, err
	}
	kafkaConfig, err := LoadKafkaConfig()
	if err != nil {
		return nil, err
	}

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
		Database: dbConfig,
		Kafka:    kafkaConfig,
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
func LoadDatabaseConfig() (*DatabaseConfig, error) {
	postgresConnStr := os.Getenv("POSTGRES_CONNECTION_STRING")
	if postgresConnStr == "" {
		return nil, errors.New("missing env")
	}

	mongoConnStr := os.Getenv("MONGODB_CONNECTION_STRING")
	if mongoConnStr == "" {
		return nil, errors.New("missing env")
	}

	return &DatabaseConfig{
		PostgresConnectionString: postgresConnStr,
		MongoDBConnectionString:  mongoConnStr,
	}, nil
}

// LoadKafkaConfig loads the kafka configuration.
func LoadKafkaConfig() (*KafkaConfig, error) {
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "kafka:9092" // Default if not set
	}

	schemaRegistryURL := os.Getenv("SCHEMA_REGISTRY_URL")
	if schemaRegistryURL == "" {
		schemaRegistryURL = "http://schema-registry:8081" // Default if not set
	}

	return &KafkaConfig{
		Brokers:           kafkaBrokers,
		SchemaRegistryURL: schemaRegistryURL,
	}, nil
}
