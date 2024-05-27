package main

import "fmt"

func main() {
	// Creating storage
	userStorage := NewInMemoryUserStorage()
	vehicleStorage := NewInMemoryVehicleStorage()
	rideStorage := NewInMemoryRideStorage()

	// Creating managers
	userMgr := NewUserManager(userStorage)
	vehicleMgr := NewVehicleManager(vehicleStorage, userMgr)
	rideMgr := NewRideManager(rideStorage, userMgr, vehicleMgr)

	// Adding users
	if err := userMgr.AddUser(User{ID: "1", Name: "Amar", Role: "Driver"}); err != nil {
		fmt.Println(err)
		return
	}
	if err := userMgr.AddUser(User{ID: "2", Name: "Chetan", Role: "Driver"}); err != nil {
		fmt.Println(err)
		return
	}
	if err := userMgr.AddUser(User{ID: "3", Name: "Bhuwan", Role: "Passenger"}); err != nil {
		fmt.Println(err)
		return
	}
	if err := userMgr.AddUser(User{ID: "4", Name: "Vijay", Role: "Passenger"}); err != nil {
		fmt.Println(err)
		return
	}

	// Adding vehicles
	if err := vehicleMgr.AddVehicle(Vehicle{ID: "1", OwnerID: "1", Model: "Toyota", Capacity: 4}); err != nil {
		fmt.Println(err)
		return
	}

	if err := vehicleMgr.AddVehicle(Vehicle{ID: "2", OwnerID: "2", Model: "XUV", Capacity: 7}); err != nil {
		fmt.Println(err)
		return
	}

	// Offering rides
	if err := rideMgr.OfferRide(Ride{ID: "101", DriverID: "1", VehicleID: "1", Source: "A", Destination: "B", AvailableSeats: 4}); err != nil {
		fmt.Println(err)
		return
	}
	if err := rideMgr.OfferRide(Ride{ID: "102", DriverID: "2", VehicleID: "2", Source: "B", Destination: "C", AvailableSeats: 4}); err != nil {
		fmt.Println(err)
		return
	}

	_, err := rideMgr.SelectRide("3", "A", "C", 3, string(MostVacantSeats))
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = rideMgr.SelectRide("4", "A", "B", 1, string(MostVacantSeats))
	if err != nil {
		fmt.Println(err)
		return
	}
	rideMgr.PrintRideStats()
}
