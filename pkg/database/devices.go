// Package database for all database assets
package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	devicesTableName          = "devices"
	devicesFieldSpecificQuery = "select * from %s where serial='%s'"
)

// DeviceObject for devices that will come from a database
type DeviceObject struct {
	ID       int    `db:"id"`
	Model    string `db:"model"`
	Serial   string `db:"serial"`
	Firmware string `db:"firmware"`
	Active   int    `db:"active"`
}

// Populate populates the device object with the data from database row
func (device *DeviceObject) Populate(rows *sqlx.Rows) error {
	if rows.Next() {
		err := rows.StructScan(device)
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

// Load the device object from the database response
func (device *DeviceObject) Load(database *sqlx.DB) error {
	query := fmt.Sprintf(specificItemLoad, devicesTableName, device.ID)
	results, err := database.Queryx(query)

	if err != nil {
		return err
	}

	err = device.Populate(results)
	if err != nil {
		logrus.Warnf("failed populating data: %v", err)
	}

	err = results.Close()
	if err != nil {
		logrus.Warnf("failed closing results: %v", err)
	}
	return nil
}

// LoadByField loads an object by a specific device field know to said object
func (device *DeviceObject) LoadByField(database *sqlx.DB, field string) error {
	query := fmt.Sprintf(devicesFieldSpecificQuery, devicesTableName, field)
	results, err := database.Queryx(query)

	if err != nil {
		return err
	}

	err = device.Populate(results)
	if err != nil {
		logrus.Warnf("failed populating data: %v", err)
	}

	err = results.Close()
	if err != nil {
		logrus.Warnf("failed closing results: %v", err)
	}
	return nil
}

// Create adds the device item to the database, returning an error if failure
func (device *DeviceObject) Create(database *sqlx.DB) error {
	values := fmt.Sprintf("'%s','%s','%s',1", device.Model, device.Serial, device.Firmware)
	query := fmt.Sprintf(createItem, devicesTableName, "model,serial,firmware,active", values)

	result, err := database.Exec(query)
	if err != nil {
		return err
	}

	nextID, err := result.LastInsertId()

	if err != nil {
		return err
	}

	device.ID = int(nextID)

	return nil
}

// Update the device item in the database, returning an error if failure
func (device *DeviceObject) Update(database *sqlx.DB) error {
	values := fmt.Sprintf("model='%s',serial='%s',firmware='%s',active=%d", device.Model, device.Serial, device.Firmware, device.Active)
	query := fmt.Sprintf(updateItem, devicesTableName, values)

	_, err := database.Exec(query)

	return err
}

// UpdateMany device items in the database using specified criteria
func (device *DeviceObject) UpdateMany(database *sqlx.DB, values, criteria map[string]string) error {
	valueUpdates := getValues(values)
	restrictions := getCriteria(criteria)

	query := fmt.Sprintf(updateManyItems, devicesTableName, valueUpdates, restrictions)

	_, err := database.Exec(query)

	return err
}

// Remove the device item from the database, returning an error if failure
func (device *DeviceObject) Remove(database *sqlx.DB) error {
	query := fmt.Sprintf(deleteItem, devicesTableName, device.ID)

	_, err := database.Exec(query)

	return err
}

// Query the items from the database, returning an nil if failure
func (device *DeviceObject) Query(database *sqlx.DB, criteria map[string]string) *[]Access {
	objects := make([]Access, 0)
	restrictions := getCriteria(criteria)

	query := fmt.Sprintf(queryMany, devicesTableName, restrictions)

	results, err := database.Queryx(query)
	if err != nil {
		return nil
	}
	for results.Next() {
		var device = DeviceObject{}
		err = results.StructScan(&device)
		if err == nil {
			objects = append(objects, &device)
		}
	}

	return &objects
}
