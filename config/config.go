package config

import(
	"os"
	"github.com/joho/godotenv"
	"log"
	"strconv"
	"fmt"
)




type ConfigError struct{
	msg string
}
func (ce ConfigError)Error() string{
	return ce.msg
}


type Config struct{
	//main config
	Host string
	Port string
	//mailer config
	Email string
	EmailPassword string
	SMTPHost string
	SMTPPort string
	//redis config
	RedisAddr string
	RedisPassword string
	//database config
	DataBaseURL string
	//security config
	JWTSecret string
	HashSecret string

}

func Init() (*Config, error){
	// Load .env if it exists. Don't panic if it's missing — production
	// environments typically inject env vars directly.
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatal("config: failed to load .env: "+err.Error())
		e := ConfigError{msg:"config: failed to load .env: "+err.Error()}
		return nil , e
	}

	c := &Config{
		// main
		Host: getEnv("HOST", "0.0.0.0"),
		Port: getEnv("PORT", "8080"),

		// mailer
		Email:         getEnv("EMAIL", ""),
		EmailPassword: getEnv("EMAIL_PASSWORD", ""),
		SMTPHost:      getEnv("SMTP_HOST", ""),
		SMTPPort:      getEnv("SMTP_PORT", ""),

		// redis
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		// database
		DataBaseURL: getEnv("DB_URL", ""),

		// security
		JWTSecret:  getEnv("JWT_SECRET", ""),
		HashSecret: getEnv("HASH_SECRET", ""),
	}

	err := c.validate()
	if err != nil{
		return nil , err
	}
	
	return c , nil
}



func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}



func (c *Config) validate() error{
	required := map[string]string{
		"DB_URL":      c.DataBaseURL,
		"JWT_SECRET":  c.JWTSecret,
		"HASH_SECRET": c.HashSecret,
	}

	for key, val := range required {
		if val == "" {
			e := ConfigError{msg:fmt.Sprintf("config: required environment variable %q is not set", key)}
			return e
		}
	}

	if _, err := strconv.Atoi(c.Port); err != nil {
		e := ConfigError{msg:fmt.Sprintf("config: PORT must be a number, got %q", c.Port)}
		return e
	}
	return nil
}