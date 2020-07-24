package config

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/rl404/point-system/internal/model"
	"github.com/rl404/point-system/internal/utils"
)

// tableList is list of all required tables.
var tableList = []interface{}{
	model.UserPoint{},
	model.Log{},
}

// validateDB to validate db structure.
func (c *Config) validateDB(db *gorm.DB) error {
	// Schema check.
	if !c.isSchemaExist(db) {
		err := c.createSchema(db)
		if err != nil {
			return err
		}
	}

	// Tables check.
	existingTables := c.getExistingTables(db)
	for _, model := range tableList {
		tableName := db.NewScope(model).TableName()
		if !utils.InArray(existingTables, tableName) {
			err := c.createTable(db, model)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// isSchemaExist to check if schema is exist.
func (c *Config) isSchemaExist(db *gorm.DB) (isExist bool) {
	db.Raw("SELECT EXISTS(SELECT 1 FROM pg_namespace WHERE nspname = ?)", c.Schema).Row().Scan(&isExist)
	return isExist
}

// createSchema to create new schema.
func (c *Config) createSchema(db *gorm.DB) error {
	query := fmt.Sprintf("CREATE SCHEMA \"%s\" AUTHORIZATION \"%s\"", c.Schema, c.User)
	return db.Exec(query).Error
}

// getExistingTables to get list of existing tables.
func (c *Config) getExistingTables(db *gorm.DB) (tables []string) {
	rows, _ := db.Raw("SELECT concat(table_schema, '.', table_name)  FROM information_schema.tables WHERE table_schema = ?", c.Schema).Rows()
	defer rows.Close()
	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)
		tables = append(tables, tableName)
	}
	return tables
}

// CreateTable to create table.
func (c *Config) createTable(db *gorm.DB, model interface{}) error {
	return db.CreateTable(model).Error
}
