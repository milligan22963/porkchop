// Package database for all database assets
package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const userTableName = "users"

// UserObject for users that will come from a database
type UserObject struct {
	ID             int       `db:"id"`
	FirstName      string    `db:"fname"`
	LastName       string    `db:"lname"`
	NickName       string    `db:"nname"`
	UserName       string    `db:"uname"`
	Password       string    `db:"password"`
	PasswordChange time.Time `db:"password_change"`
	EmailAddress   string    `db:"email"`
	Phone          string    `db:"phone"`
	Age            int       `db:"age"`
	AcceptsCookies int       `db:"accepts_cookies"`
	FilterContent  int       `db:"filter_content"`
	LastLogin      time.Time `db:"last_login"`
	Token          string    `db:"token"`
	Active         int       `db:"active"`
}

// Populate populates the settings object with the data from database row
func (user *UserObject) Populate(rows *sqlx.Rows) error {
	if rows.Next() {
		err := rows.StructScan(user)
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

// Load the user object from the database response
func (user *UserObject) Load(database *sqlx.DB) error {
	query := fmt.Sprintf(specificItemLoad, userTableName, user.ID)
	results, err := database.Queryx(query)

	if err != nil {
		return err
	}

	err = user.Populate(results)
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
func (user *UserObject) LoadByField(database *sqlx.DB, field string) error {
	// query := fmt.Sprintf("select * from %s where uname='%s'", userTableName, field)
	results, err := database.Queryx("select * from ? where uname='?'", userTableName, field)

	if err != nil {
		return err
	}

	err = user.Populate(results)
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
func (user *UserObject) Create(database *sqlx.DB) error {
	values := fmt.Sprintf("'%s,'%s','%s','%s','%s','%s',%d,%d,%d,1", user.FirstName, user.LastName, user.NickName, user.UserName, user.EmailAddress, user.Phone, user.Age, user.AcceptsCookies, user.FilterContent)
	query := fmt.Sprintf(createItem, userTableName, "fname,lname,nname,uname,email,phone,age,accepts_cookie,filter_content,active", values)

	result, err := database.Exec(query)
	if err != nil {
		return err
	}

	nextID, err := result.LastInsertId()

	if err != nil {
		return err
	}

	user.ID = int(nextID)

	return nil
}

// Update the item in the database, returning an error if failure
func (user *UserObject) Update(database *sqlx.DB) error {
	values := fmt.Sprintf("fname='%s',lname='%s',nname='%s',uname='%s',email='%s',phone='%s',age=%d,accepts_cookies=%d,filter_content=%d,active=%d", user.FirstName, user.LastName, user.NickName, user.UserName, user.EmailAddress, user.Phone, user.Age, user.AcceptsCookies, user.FilterContent, user.Active)
	query := fmt.Sprintf(updateItem, userTableName, values)

	_, err := database.Exec(query)

	return err
}

// UpdateMany items in the database using specified criteria
func (user *UserObject) UpdateMany(database *sqlx.DB, values, criteria map[string]string) error {
	valueUpdates := getValues(values)
	restrictions := getCriteria(criteria)

	query := fmt.Sprintf(updateManyItems, userTableName, valueUpdates, restrictions)

	_, err := database.Exec(query)

	return err
}

// Remove the item from the database, returning an error if failure
func (user *UserObject) Remove(database *sqlx.DB) error {
	query := fmt.Sprintf(deleteItem, userTableName, user.ID)

	_, err := database.Exec(query)

	return err
}

// Query the items from the database, returning an nil if failure
func (user *UserObject) Query(database *sqlx.DB, criteria map[string]string) *[]Access {
	objects := make([]Access, 0)
	restrictions := getCriteria(criteria)

	query := fmt.Sprintf(queryMany, userTableName, restrictions)

	results, err := database.Queryx(query)
	if err != nil {
		return nil
	}
	for results.Next() {
		var userObj = UserObject{}
		err = results.StructScan(&userObj)
		if err == nil {
			objects = append(objects, &userObj)
		}
	}

	return &objects
}
