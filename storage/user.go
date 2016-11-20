package storage

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	Name      string
	Available bool
	LastSeen  time.Time
}

type Persister interface {
	AddUser(usr *User) error
	GetUser(usrName string) (*User, error)
	UpdateUser(usr *User) error
}

type DB struct {
	*sql.DB
}

func (db DB) AddUser(usr *User) error {
	_, err := db.Exec("INSERT INTO names(name,available,last_seen) VALUES (?,?,?) ", usr.Name, usr.Available, usr.LastSeen.Format(time.RFC3339Nano))
	return err
}

func (db DB) GetUser(usrName string) (*User, error) {
	var name, lastSeenString string
	var available bool
	row := db.QueryRow("SELECT * FROM names WHERE name=?", usrName)
	err := row.Scan(&name, &available, &lastSeenString)
	lastSeen, _ := time.Parse(time.RFC3339Nano, lastSeenString)
	return &User{Name: name, Available: available, LastSeen: lastSeen}, err
}

//noinspection GoPlaceholderCount
func (db DB) UpdateUser(usr *User) error {
	if usr.Name == "" {
		return fmt.Errorf("User name is required to update a user\n")
	}
	_, err := db.Exec("UPDATE names SET available=?, last_seen=? WHERE name=?", usr.Available, usr.LastSeen.Format(time.RFC3339Nano), usr.Name)
	return err
}

func (db DB) CreateSchema() error {
	sqlStmt := `create table if not exists names (name text not null primary key, available bool, last_seen text);
		delete from names;`
	_, err := db.Exec(sqlStmt)
	return err
}
