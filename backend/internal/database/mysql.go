package database

import (
	"agent_learning/internal/config"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	sqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitMySQL(cfg *config.Config) (*gorm.DB, error) {
	gormLogger := logger.New(
		log.New(os.Stdout, "[gorm] ", log.LstdFlags),
		logger.Config{
			SlowThreshold:             600 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unknown database") {
			if createErr := ensureDatabaseExists(cfg.MySQLDSN); createErr != nil {
				return nil, fmt.Errorf("auto create database failed: %w", createErr)
			}
			db, err = gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{
				Logger: gormLogger,
			})
		}
		if err != nil {
			return nil, err
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(2 * time.Hour)

	return db, nil
}

func ensureDatabaseExists(dsn string) error {
	cfg, err := sqldriver.ParseDSN(dsn)
	if err != nil {
		return err
	}
	dbName := strings.TrimSpace(cfg.DBName)
	if dbName == "" {
		return fmt.Errorf("dsn has empty database name")
	}

	cfg.DBName = ""
	rootConn, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}
	defer rootConn.Close()

	safeDBName := strings.ReplaceAll(dbName, "`", "``")
	createSQL := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci",
		safeDBName,
	)
	_, err = rootConn.Exec(createSQL)
	return err
}
