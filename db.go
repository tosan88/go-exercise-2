package main

import (
	"log"
	"time"
)

func (c *botClient) addUserToDB(usr *user) {
	_, err := c.namesDB.Exec("INSERT INTO names(name,available,last_seen) VALUES (?,?,?) ", usr.name, usr.available, usr.lastSeen.Format(time.RFC3339Nano))
	if err != nil {
		log.Printf("ERROR - insert into names: %v\n", err)
		return
	}
}

func (c *botClient) getUserFromDB(usrName string) (user, error) {
	var name, lastSeenString string
	var available bool
	row := c.namesDB.QueryRow("SELECT * FROM names WHERE name=?", usrName)
	err := row.Scan(&name, &available, &lastSeenString)
	lastSeen, _ := time.Parse(time.RFC3339Nano, lastSeenString)
	return user{name: name, available: available, lastSeen: lastSeen}, err
}

func (c *botClient) modifyUserInDB(usr *user) {
	storedUsr, err := c.getUserFromDB(usr.name)
	if err != nil {
		log.Printf("ERROR - getting user %v: %v\n", usr.name, err)
		return
	}
	if storedUsr.name == "" || storedUsr.name != usr.name {
		log.Printf("WARN - no stored user in DB or other user is returned: %v vs %v\n", usr.name, storedUsr.name)
		return
	}
	_, err = c.namesDB.Exec("UPDATE names SET available=?, last_seen=? WHERE name=?", usr.available, usr.lastSeen.Format(time.RFC3339Nano), usr.name)
	if err != nil {
		log.Printf("ERROR - update of names: %v\n", err)
		return
	}
}
