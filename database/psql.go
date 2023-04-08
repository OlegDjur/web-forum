package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	DBname   string
	SSLMode  string
}

func initConfig() *Config {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	config := &Config{
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		DBname:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}

	return config
}

func NewDB() *sqlx.DB {
	c := initConfig()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.Username, c.Password, c.Host, c.Port, c.DBname)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("error to connect database: %v", err)
	}

	m, err := migrate.New("file://schema", connStr)
	if err != nil {
		log.Printf("error migreation: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("error migration up: %v", err)
	}

	return db
}
