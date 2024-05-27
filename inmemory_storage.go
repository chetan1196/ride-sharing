package main

import "fmt"

// InMemoryUserStorage implements UserStorage using a map
type InMemoryUserStorage struct {
	users map[string]User
}

func NewInMemoryUserStorage() UserStorage {
	return &InMemoryUserStorage{users: make(map[string]User)}
}

func (s *InMemoryUserStorage) AddUser(user User) error {
	if _, exists := s.users[user.ID]; exists {
		return fmt.Errorf("user already exists")
	}
	s.users[user.ID] = user
	return nil
}

func (s *InMemoryUserStorage) GetUserByID(userID string) (User, error) {
	user, exists := s.users[userID]
	if !exists {
		return User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *InMemoryUserStorage) GetAllUsers() map[string]User {
	return s.users
}

//////

// InMemoryVehicleStorage implements VehicleStorage using a map
type InMemoryVehicleStorage struct {
	vehicles map[string]Vehicle
}

func NewInMemoryVehicleStorage() VehicleStorage {
	return &InMemoryVehicleStorage{vehicles: make(map[string]Vehicle)}
}

func (s *InMemoryVehicleStorage) AddVehicle(vehicle Vehicle) error {
	if _, exists := s.vehicles[vehicle.ID]; exists {
		return fmt.Errorf("vehicle already exists")
	}
	s.vehicles[vehicle.ID] = vehicle
	return nil
}

func (s *InMemoryVehicleStorage) GetVehicleByID(vehicleID string) (Vehicle, error) {
	vehicle, exists := s.vehicles[vehicleID]
	if !exists {
		return Vehicle{}, fmt.Errorf("vehicle not found")
	}
	return vehicle, nil
}

func (s *InMemoryVehicleStorage) GetAllVehicles() map[string]Vehicle {
	return s.vehicles
}

//////

// InMemoryRideStorage implements RideStorage using a map
type InMemoryRideStorage struct {
	rides map[string]Ride
}

func NewInMemoryRideStorage() RideStorage {
	return &InMemoryRideStorage{rides: make(map[string]Ride)}
}

func (s *InMemoryRideStorage) AddRide(ride Ride) error {
	if _, exists := s.rides[ride.ID]; exists {
		return fmt.Errorf("ride already exists")
	}
	s.rides[ride.ID] = ride
	return nil
}

func (s *InMemoryRideStorage) GetRideByID(rideID string) (Ride, error) {
	ride, exists := s.rides[rideID]
	if !exists {
		return Ride{}, fmt.Errorf("ride not found")
	}
	return ride, nil
}

func (s *InMemoryRideStorage) UpdateRide(ride Ride) error {
	if _, exists := s.rides[ride.ID]; !exists {
		return fmt.Errorf("ride not found")
	}
	s.rides[ride.ID] = ride
	return nil
}

func (s *InMemoryRideStorage) DeleteRide(rideID string) error {
	if _, exists := s.rides[rideID]; !exists {
		return fmt.Errorf("ride not found")
	}
	delete(s.rides, rideID)
	return nil
}

func (s *InMemoryRideStorage) GetAllRides() map[string]Ride {
	return s.rides
}
