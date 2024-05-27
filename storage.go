package main

// UserStorage defines methods for user storage
type UserStorage interface {
	AddUser(user User) error
	GetUserByID(userID string) (User, error)
	GetAllUsers() map[string]User
}

// VehicleStorage defines methods for vehicle storage
type VehicleStorage interface {
	AddVehicle(vehicle Vehicle) error
	GetVehicleByID(vehicleID string) (Vehicle, error)
	GetAllVehicles() map[string]Vehicle
}

// RideStorage defines methods for ride storage
type RideStorage interface {
	AddRide(ride Ride) error
	GetRideByID(rideID string) (Ride, error)
	UpdateRide(ride Ride) error
	DeleteRide(rideID string) error
	GetAllRides() map[string]Ride
}
