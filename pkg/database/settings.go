// Package database for all database assets
package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	settingsTableName          = "settings"
	settingsFieldSpecificQuery = "select * from %s where name='%s'"
)

// SettingsObject for settings that will come from a database
type SettingsObject struct {
	ID                  int    `db:"id"`
	UserDeviceMappingID int    `db:"user_device_mapping_id"`
	Name                string `db:"name"`
	Value               string `db:"value"`
	Active              int    `db:"active"`
}

// Populate populates the settings object with the data from database row
func (settings *SettingsObject) Populate(rows *sqlx.Rows) error {
	if rows.Next() {
		err := rows.StructScan(settings)
		if err != nil {
			logrus.Warnf("failed scanning results: %v", err)
		}
	} else {
		err := rows.Err()
		if err != nil {
			return fmt.Errorf("failed to find any results - error: %v", err)
		}
	}
	return nil
}

// Load the settings object from the database response
func (settings *SettingsObject) Load(database *sqlx.DB) error {
	query := fmt.Sprintf(specificItemLoad, settingsTableName, settings.ID)
	results, err := database.Queryx(query)

	if err != nil {
		return err
	}

	err = settings.Populate(results)
	if err != nil {
		logrus.Warnf("failed populating data: %v", err)
	}

	err = results.Close()
	if err != nil {
		logrus.Warnf("failed closing results: %v", err)
	}
	return nil
}

// LoadByField loads an object by a specific field know to said object
func (settings *SettingsObject) LoadByField(database *sqlx.DB, field string) error {
	query := fmt.Sprintf(settingsFieldSpecificQuery, settingsTableName, field)
	results, err := database.Queryx(query)

	if err != nil {
		return err
	}

	err = settings.Populate(results)
	if err != nil {
		logrus.Warnf("failed populating data: %v", err)
	}

	err = results.Close()
	if err != nil {
		logrus.Warnf("failed closing results: %v", err)
	}
	return nil
}

// Create adds the item to the database, returning an error if failure
func (settings *SettingsObject) Create(database *sqlx.DB) error {
	values := fmt.Sprintf("%d,'%s','%s',1", settings.UserDeviceMappingID, settings.Name, settings.Value)
	query := fmt.Sprintf(createItem, settingsTableName, "user_device_mapping_id,name,value,active", values)

	result, err := database.Exec(query)
	if err != nil {
		return err
	}

	nextID, err := result.LastInsertId()

	if err != nil {
		return err
	}

	settings.ID = int(nextID)

	return nil
}

// Update the item in the database, returning an error if failure
func (settings *SettingsObject) Update(database *sqlx.DB) error {
	values := fmt.Sprintf("user_device_mapping_id=%d,name='%s',value='%s',active=%d", settings.UserDeviceMappingID, settings.Name, settings.Value, settings.Active)
	query := fmt.Sprintf(updateItem, settingsTableName, values)

	_, err := database.Exec(query)

	return err
}

// UpdateMany items in the database using specified criteria
func (settings *SettingsObject) UpdateMany(database *sqlx.DB, values, criteria map[string]string) error {
	valueUpdates := getValues(values)
	restrictions := getCriteria(criteria)

	query := fmt.Sprintf(updateManyItems, settingsTableName, valueUpdates, restrictions)

	_, err := database.Exec(query)

	return err
}

// Remove the item from the database, returning an error if failure
func (settings *SettingsObject) Remove(database *sqlx.DB) error {
	query := fmt.Sprintf(deleteItem, settingsTableName, settings.ID)

	_, err := database.Exec(query)

	return err
}

// Query the items from the database, returning an nil if failure
func (settings *SettingsObject) Query(database *sqlx.DB, criteria map[string]string) *[]Access {
	objects := make([]Access, 0)
	restrictions := getCriteria(criteria)

	query := fmt.Sprintf(queryMany, settingsTableName, restrictions)

	results, err := database.Queryx(query)
	if err != nil {
		return nil
	}
	for results.Next() {
		var settings = SettingsObject{}
		err = results.StructScan(&settings)
		if err == nil {
			objects = append(objects, &settings)
		}
	}

	return &objects
}
