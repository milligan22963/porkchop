// Package database for all database assets
package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const deviceUserMappingTableName = "device_user_mapping"

// DeviceUserMappingObject for mappings between devices and users that will come from a database
type DeviceUserMappingObject struct {
	ID       int `db:"id"`
	UserID   int `db:"user_id"`
	DeviceID int `db:"device_id"`
	Active   int `db:"active"`
}

// Populate populates the settings object with the data from database row
func (deviceUserMapping *DeviceUserMappingObject) Populate(rows *sqlx.Rows) error {
	if rows.Next() {
		err := rows.StructScan(deviceUserMapping)
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
func (deviceUserMapping *DeviceUserMappingObject) Load(database *sqlx.DB) error {
	query := fmt.Sprintf(specificItemLoad, deviceUserMappingTableName, deviceUserMapping.ID)
	results, err := database.Queryx(query)

	if err != nil {
		return err
	}

	err = deviceUserMapping.Populate(results)
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
func (deviceUserMapping *DeviceUserMappingObject) LoadByField(database *sqlx.DB, field string) error {
	// query := fmt.Sprintf("select * from %s where device_id=%s", deviceUserMappingTableName, field)
	results, err := database.Queryx("select * from ? where device_id=?", deviceUserMappingTableName, field)

	if err != nil {
		return err
	}

	err = deviceUserMapping.Populate(results)
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
func (deviceUserMapping *DeviceUserMappingObject) Create(database *sqlx.DB) error {
	values := fmt.Sprintf("%d,%d,1", deviceUserMapping.UserID, deviceUserMapping.DeviceID)
	query := fmt.Sprintf(createItem, deviceUserMappingTableName, "user_id,device_id,active", values)

	result, err := database.Exec(query)
	if err != nil {
		return err
	}

	nextID, err := result.LastInsertId()

	if err != nil {
		return err
	}

	deviceUserMapping.ID = int(nextID)

	return nil
}

// Update the item in the database, returning an error if failure
func (deviceUserMapping *DeviceUserMappingObject) Update(database *sqlx.DB) error {
	values := fmt.Sprintf("user_id=%d,device_id=%d,active=%d", deviceUserMapping.UserID, deviceUserMapping.DeviceID, deviceUserMapping.Active)
	query := fmt.Sprintf(updateItem, deviceUserMappingTableName, values)

	_, err := database.Exec(query)

	return err
}

// UpdateMany items in the database using specified criteria
func (deviceUserMapping *DeviceUserMappingObject) UpdateMany(database *sqlx.DB, values, criteria map[string]string) error {
	valueUpdates := getValues(values)
	restrictions := getCriteria(criteria)

	query := fmt.Sprintf(updateManyItems, deviceUserMappingTableName, valueUpdates, restrictions)

	_, err := database.Exec(query)

	return err
}

// Remove the item from the database, returning an error if failure
func (deviceUserMapping *DeviceUserMappingObject) Remove(database *sqlx.DB) error {
	query := fmt.Sprintf(deleteItem, deviceUserMappingTableName, deviceUserMapping.ID)

	_, err := database.Exec(query)

	return err
}

// Query the items from the database, returning an nil if failure
func (deviceUserMapping *DeviceUserMappingObject) Query(database *sqlx.DB, criteria map[string]string) *[]Access {
	objects := make([]Access, 0)
	restrictions := getCriteria(criteria)

	query := fmt.Sprintf(queryMany, deviceUserMappingTableName, restrictions)

	results, err := database.Queryx(query)
	if err != nil {
		return nil
	}
	for results.Next() {
		var deviceUserMapping = DeviceUserMappingObject{}
		err = results.StructScan(&deviceUserMapping)
		if err == nil {
			objects = append(objects, &deviceUserMapping)
		}
	}

	return &objects
}
