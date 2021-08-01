// Package database for all database assets
package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const imagesTableName = "images"

// ImageObject for images that will come from a database
type ImageObject struct {
	ID       int    `db:"id"`
	UserID   int    `db:"user_id"`
	DeviceID int    `db:"device_id"`
	Path     string `db:"path"`
	Active   int    `db:"active"`
}

// Populate populates the image object with the data from database row
func (image *ImageObject) Populate(rows *sqlx.Rows) error {
	if rows.Next() {
		err := rows.StructScan(image)
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
func (image *ImageObject) Load(database *sqlx.DB) error {
	query := fmt.Sprintf(specificItemLoad, imagesTableName, image.ID)
	results, err := database.Queryx(query)

	if err != nil {
		return err
	}

	err = image.Populate(results)
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
func (image *ImageObject) LoadByField(database *sqlx.DB, field string) error {
	// query := fmt.Sprintf("select * from %s where device_id=%s", imagesTableName, field)
	results, err := database.Queryx("select * from ? where device_id=?", imagesTableName, field)

	if err != nil {
		return err
	}

	err = image.Populate(results)
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
func (image *ImageObject) Create(database *sqlx.DB) error {
	values := fmt.Sprintf("%d,%d,'%s',1", image.UserID, image.DeviceID, image.Path)
	query := fmt.Sprintf(createItem, imagesTableName, "user_id,device_id,path,active", values)

	result, err := database.Exec(query)
	if err != nil {
		return err
	}

	nextID, err := result.LastInsertId()

	if err != nil {
		return err
	}

	image.ID = int(nextID)

	return nil
}

// Update the item in the database, returning an error if failure
func (image *ImageObject) Update(database *sqlx.DB) error {
	values := fmt.Sprintf("user_id=%d,device_id=%d,path='%s',active=%d", image.UserID, image.DeviceID, image.Path, image.Active)
	query := fmt.Sprintf(updateItem, imagesTableName, values)

	_, err := database.Exec(query)

	return err
}

// UpdateMany items in the database using specified criteria
func (image *ImageObject) UpdateMany(database *sqlx.DB, values, criteria map[string]string) error {
	valueUpdates := getValues(values)
	restrictions := getCriteria(criteria)

	query := fmt.Sprintf(updateManyItems, imagesTableName, valueUpdates, restrictions)

	_, err := database.Exec(query)

	return err
}

// Remove the item from the database, returning an error if failure
func (image *ImageObject) Remove(database *sqlx.DB) error {
	query := fmt.Sprintf(deleteItem, imagesTableName, image.ID)

	_, err := database.Exec(query)

	return err
}

// Query the items from the database, returning an nil if failure
func (image *ImageObject) Query(database *sqlx.DB, criteria map[string]string) *[]Access {
	objects := make([]Access, 0)
	restrictions := getCriteria(criteria)

	query := fmt.Sprintf(queryMany, imagesTableName, restrictions)

	results, err := database.Query(query)
	if err != nil {
		return nil
	}
	for results.Next() {
		var ID, userID, deviceID, active int
		var path string
		err = results.Scan(&ID, &userID, &deviceID, &path, &active)
		if err == nil {
			var image = ImageObject{ID: ID, UserID: userID, DeviceID: deviceID, Path: path, Active: active}
			objects = append(objects, &image)
		}
	}

	return &objects
}
