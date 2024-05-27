package main

import "fmt"

type User struct {
	ID   string
	Name string
	Role Role // Driver or Passenger
}

type Role string

const (
	Passenger Role = "Passenger"
	Driver    Role = "Driver"
)

type userManager struct {
	storage UserStorage
}

func NewUserManager(storage UserStorage) *userManager {
	return &userManager{storage: storage}
}

func (um *userManager) AddUser(user User) error {
	if err := um.storage.AddUser(user); err != nil {
		return fmt.Errorf("could not add user: %v", err)
	}
	fmt.Printf("User added: %v\n", user)
	return nil
}

func (um *userManager) GetUserByID(userID string) (User, error) {
	user, err := um.storage.GetUserByID(userID)
	if err != nil {
		return User{}, fmt.Errorf("could not find user: %v", err)
	}
	return user, nil
}

func (um *userManager) IsDriver(userID string) error {
	user, err := um.storage.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("could not find user %s: %v", userID, err)
	}
	if user.Role != Driver {
		return fmt.Errorf("user %v is not a driver", userID)
	}
	return nil
}
