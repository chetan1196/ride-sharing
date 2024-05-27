package main

import (
	"testing"
)

// Test adding a user
func TestAddUser(t *testing.T) {
	userStorage := NewInMemoryUserStorage()
	userMgr := NewUserManager(userStorage)

	user := User{ID: "1", Name: "Amar", Role: "Driver"}
	if err := userMgr.AddUser(user); err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	retrievedUser, err := userStorage.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("Expected to retrieve user, but got error %v", err)
	}
	if retrievedUser != user {
		t.Fatalf("Expected user to be %v, but got %v", user, retrievedUser)
	}
}
