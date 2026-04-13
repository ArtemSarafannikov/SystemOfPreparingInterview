package model

import (
	"fmt"
	"reflect"
	"strings"
)

// TableSchema .
type TableSchema struct {
	tableName   string
	columnNames []string
	returning   string
}

// NewSchema .
func NewSchema[T any](tableName string) *TableSchema {
	columns := parseDBTags[T](tableName)
	return &TableSchema{
		tableName:   tableName,
		columnNames: columns,
		returning:   fmt.Sprintf("RETURNING %s", strings.Join(columns, ",")),
	}
}

// TableName .
func (s *TableSchema) TableName() string {
	return s.tableName
}

// FieldName .
func (s *TableSchema) FieldName(field string) string {
	return fmt.Sprintf("%s.%s", s.tableName, field)
}

// Columns .
func (s *TableSchema) Columns() []string {
	return s.columnNames
}

// Returning .
func (s *TableSchema) Returning() string {
	return s.returning
}

func parseDBTags[T any](tableName string) []string {
	var model T
	r := reflect.Indirect(reflect.ValueOf(model)).Type()

	cols := make([]string, 0, r.NumField())
	for i := 0; i < r.NumField(); i++ {
		colName := r.Field(i).Tag.Get("db")
		if colName == "" || colName == "-" {
			continue
		}

		cols = append(cols, fmt.Sprintf("%s.%s", tableName, colName))
	}

	return cols
}
