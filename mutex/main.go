package main

import (
	"fmt"
	"sync"
	"time"
)

type User struct {
	ID   int
	Name string
}

var userStorage = map[int]*User{}
var mutex sync.Mutex

func AddUser(user *User) {
	mutex.Lock()
	userStorage[user.ID] = user
	mutex.Unlock()
}

func DeleteUser(userID int) {
	mutex.Lock()
	delete(userStorage, userID)
	mutex.Unlock()
}

func LogStorage() {
	mutex.Lock()
	fmt.Println(userStorage)
	mutex.Unlock()
}

func main() {
	go AddUser(&User{1, "Weerasak"})
	// go DeleteUser(1)
	LogStorage()
	time.Sleep(1 * time.Second)
}
