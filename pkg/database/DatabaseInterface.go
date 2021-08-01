// Package database for all database assets
package database

import (
	"github.com/jmoiron/sqlx"
)

const (
	specificItemLoad = "select * from %s where id=%d"
	createItem       = "insert into %s (%s) values (%s)"
	updateItem       = "update %s set %s"
	updateManyItems  = "update %s set %s where %s"
	deleteItem       = "delete from %s where id=%d"
	queryMany        = "select * from %s where %s"
)

// Access interface for working with database objects
type Access interface {
	Populate(rows *sqlx.Rows) error
	Load(database *sqlx.DB) error
	LoadByField(database *sqlx.DB, field string) error
	Create(database *sqlx.DB) error
	Update(database *sqlx.DB) error
	UpdateMany(database *sqlx.DB, values, criteria map[string]string) error
	Remove(database *sqlx.DB) error
	Query(database *sqlx.DB, criteria map[string]string) *[]Access
}

func getValues(values map[string]string) string {
	first := true
	valueUpdates := ""
	for k, v := range values {
		if first {
			first = false
		} else {
			valueUpdates += ","
		}
		valueUpdates += k + "='" + v + "'"
	}
	return valueUpdates
}

func getCriteria(criteria map[string]string) string {
	first := true
	restrictions := ""
	for k, v := range criteria {
		if first {
			first = false
		} else {
			restrictions += " and "
		}
		restrictions += k + "='" + v + "'"
	}
	return restrictions
}
