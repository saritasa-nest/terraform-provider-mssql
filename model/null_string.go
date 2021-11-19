package model

import (
	"database/sql/driver"
	"errors"
)

type NullString string

// Scan implements sql.Scanner interface
func (s *NullString) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	strVal, ok := value.(string)
	if !ok {
		return errors.New("column is not a string")
	}
	*s = NullString(strVal)
	return nil
}

// Value implements driver.Valuer interface
func (s NullString) Value() (driver.Value, error) {
	if len(s) == 0 { // if nil or empty string
		return nil, nil
	}
	return string(s), nil
}

// ToString Get underlying string value
func (s NullString) ToString() string {
	return string(s)
}

// ValueOrSqlNull if value is empty, return string = "NULL" (for building queries)
func (s NullString) ValueOrSqlNull() string {
	if len(s) == 0 {
		return "NULL"
	}
	return string(s)
}
