package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type configType struct {
	Name        string `envconfig:"NAME" required:"false" default:"icon-go-api"`
	NetworkName string `envconfig:"NETWORK_NAME" required:"false" default:"mainnet"`

	// Ports
	APIPort     string `envconfig:"API_PORT" required:"false" default:"8001"`
	HealthPort  string `envconfig:"HEALTH_PORT" required:"false" default:"8180"`
	MetricsPort string `envconfig:"METRICS_PORT" required:"false" default:"9400"`

	// Prefix
	RestPrefix      string `envconfig:"REST_PREFIX" required:"false" default:"/api/v1"`
	WebsocketPrefix string `envconfig:"WEBSOCKET_PREFIX" required:"false" default:"/ws/v1"`
	HealthPrefix    string `envconfig:"HEALTH_PREFIX" required:"false" default:"/health"`
	MetricsPrefix   string `envconfig:"METRICS_PREFIX" required:"false" default:"/metrics"`

	// Endpoints
	MaxPageSize int `envconfig:"MAX_PAGE_SIZE" required:"false" default:"100"`
	MaxPageSkip int `envconfig:"MAX_PAGE_SKIP" required:"false" default:"1500000"`

	// CORS
	CORSAllowOrigins  string `envconfig:"CORS_ALLOW_ORIGINS" required:"false" default:"*"`
	CORSAllowHeaders  string `envconfig:"CORS_ALLOW_HEADERS" required:"false" default:"*"`
	CORSAllowMethods  string `envconfig:"CORS_ALLOW_METHODS" required:"false" default:"GET,POST,HEAD,PUT,DELETE,PATCH"`
	CORSExposeHeaders string `envconfig:"CORS_EXPOSE_HEADERS" required:"false" default:"*"`

	// Compress
	RestCompressLevel int `envconfig:"REST_COMPRESS_LEVEL" required:"false" default:"2"`

	// Monitoring
	HealthPollingInterval int `envconfig:"HEALTH_POLLING_INTERVAL" required:"false" default:"10"`

	// Logging
	LogLevel         string `envconfig:"LOG_LEVEL" required:"false" default:"INFO"`
	LogToFile        bool   `envconfig:"LOG_TO_FILE" required:"false" default:"false"`
	LogFileName      string `envconfig:"LOG_FILE_NAME" required:"false" default:"blocks-service.log"`
	LogFormat        string `envconfig:"LOG_FORMAT" required:"false" default:"json"`
	LogIsDevelopment bool   `envconfig:"LOG_IS_DEVELOPMENT" required:"false" default:"true"`

	// DB
	DbDriver             string `envconfig:"DB_DRIVER" required:"false" default:"postgres"`
	DbHost               string `envconfig:"DB_HOST" required:"false" default:"localhost"`
	DbPort               string `envconfig:"DB_PORT" required:"false" default:"5432"`
	DbUser               string `envconfig:"DB_USER" required:"false" default:"postgres"`
	DbPassword           string `envconfig:"DB_PASSWORD" required:"false" default:"changeme"`
	DbName               string `envconfig:"DB_DBNAME" required:"false" default:"postgres"`
	DbSslmode            string `envconfig:"DB_SSL_MODE" required:"false" default:"disable"`
	DbTimezone           string `envconfig:"DB_TIMEZONE" required:"false" default:"UTC"`
	DbMaxIdleConnections int    `envconfig:"DB_MAX_IDLE_CONNECTIONS" required:"false" default:"2"`
	DbMaxOpenConnections int    `envconfig:"DB_MAX_OPEN_CONNECTIONS" required:"false" default:"10"`

	// Redis
	RedisHost                     string `envconfig:"REDIS_HOST" required:"false" default:"localhost"`
	RedisPort                     string `envconfig:"REDIS_PORT" required:"false" default:"6379"`
	RedisPassword                 string `envconfig:"REDIS_PASSWORD" required:"false" default:""`
	RedisSentinelClientMode       bool   `envconfig:"REDIS_SENTINEL_CLIENT_MODE" required:"false" default:"false"`
	RedisSentinelClientMasterName string `envconfig:"REDIS_SENTINEL_CLIENT_MASTER_NAME" required:"false" default:"master"`
	RedisKeyPrefix                string `envconfig:"REDIS_KEY_PREFIX" required:"false" default:"icon_"`

	// Redis Channels
	// NOTE must add to redis client manually
	// src/redis/client.go:63
	RedisBlocksChannel         string `envconfig:"REDIS_BLOCKS_CHANNEL" required:"false" default:"blocks"`
	RedisTransactionsChannel   string `envconfig:"REDIS_TRANSACTIONS_CHANNEL" required:"false" default:"transactions"`
	RedisLogsChannel           string `envconfig:"REDIS_LOGS_CHANNEL" required:"false" default:"logs"`
	RedisTokenTransfersChannel string `envconfig:"REDIS_TOKEN_TRANSFERS_CHANNEL" required:"false" default:"token_transfers"`

	// GORM
	GormLoggingThresholdMilli int `envconfig:"GORM_LOGGING_THRESHOLD_MILLI" required:"false" default:"250"`

	// For testing only
	ServiceUrl string `envconfig:"SERVICE_URL" required:"false" default:"http://localhost:8000"`

	// RPC
	IconNodeServiceURL           []string      `envconfig:"ICON_NODE_SERVICE_URL" required:"false" default:"https://api.icon.community/api/v3"`
	IconNodeRpcRetrySleepSeconds time.Duration `envconfig:"ICON_NODE_RPC_SLEEP_SECONDS" required:"false" default:"1s"`
	IconNodeRpcRetryAttempts     int           `envconfig:"ICON_NODE_RPC_RETRY_ATTEMPTS" required:"false" default:"20"`

	// Stats endpoints
	StatsMarketCapUpdateTime         time.Duration `envconfig:"STATS_MARKET_CAP_UPDATE_TIME" required:"false" default:"5m"`
	StatsCirculatingSupplyUpdateTime time.Duration `envconfig:"STATS_CIRCULATING_SUPPLY_UPDATE_TIME" required:"false" default:"5m"`
}

// Config - runtime config struct
var Config configType

func loadEnv() {
	// Only for local tests - Env vars for actual runtime
	// Source: https://github.com/joho/godotenv/issues/43#issuecomment-503183127
	re := regexp.MustCompile(`^(.*src)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	_ = godotenv.Load(string(rootPath) + `/../.env.test`)
}

// ReadEnvironment - Read and store runtime config
func ReadEnvironment() {
	// For local tests
	loadEnv()

	err := envconfig.Process("", &Config)
	if err != nil {
		log.Fatalf("ERROR: envconfig - %s\n", err.Error())
	}
}
