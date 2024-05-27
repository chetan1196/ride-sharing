package main

import (
	"testing"
)

// Test offering a ride
func TestOfferRide(t *testing.T) {
	userStorage := NewInMemoryUserStorage()
	vehicleStorage := NewInMemoryVehicleStorage()
	rideStorage := NewInMemoryRideStorage()

	userMgr := NewUserManager(userStorage)
	vehicleMgr := NewVehicleManager(vehicleStorage, userMgr)
	rideMgr := NewRideManager(rideStorage, userMgr, vehicleMgr)

	user := User{ID: "1", Name: "Amar", Role: "Driver"}
	vehicle := Vehicle{ID: "1", OwnerID: "1", Model: "Toyota", Capacity: 4}
	userMgr.AddUser(user)
	vehicleMgr.AddVehicle(vehicle)
	ride := Ride{ID: "1", DriverID: "1", VehicleID: "1", Source: "A", Destination: "B", AvailableSeats: 3}

	if err := rideMgr.OfferRide(ride); err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	retrievedRide, err := rideStorage.GetRideByID(ride.ID)
	if err != nil {
		t.Fatalf("Expected to retrieve ride, but got error %v", err)
	}
	if retrievedRide != ride {
		t.Fatalf("Expected ride to be %v, but got %v", ride, retrievedRide)
	}
}

// Test selecting a ride
func TestSelectRide(t *testing.T) {
	userStorage := NewInMemoryUserStorage()
	vehicleStorage := NewInMemoryVehicleStorage()
	rideStorage := NewInMemoryRideStorage()

	userMgr := NewUserManager(userStorage)
	vehicleMgr := NewVehicleManager(vehicleStorage, userMgr)
	rideMgr := NewRideManager(rideStorage, userMgr, vehicleMgr)

	user := User{ID: "1", Name: "Amar", Role: "Driver"}
	vehicle := Vehicle{ID: "1", OwnerID: "1", Model: "Toyota", Capacity: 4}
	userMgr.AddUser(user)
	vehicleMgr.AddVehicle(vehicle)
	ride := Ride{ID: "1", DriverID: "1", VehicleID: "1", Source: "A", Destination: "B", AvailableSeats: 3}
	rideMgr.OfferRide(ride)

	user2 := User{ID: "1", Name: "Chetan", Role: "Passenger"}
	selectedRoute, err := rideMgr.SelectRide(user2.ID, "A", "B", 1, "Preferred Vehicle=Toyota")
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if selectedRoute[0].ID != ride.ID {
		t.Fatalf("Expected selected ride ID to be %v, but got %v", ride.ID, selectedRoute[0].ID)
	}

	if selectedRoute[0].AvailableSeats != 2 {
		t.Fatalf("Expected available seats to be 2, but got %v", selectedRoute[0].AvailableSeats)
	}
}

// Test ending a ride
func TestEndRide(t *testing.T) {
	userStorage := NewInMemoryUserStorage()
	vehicleStorage := NewInMemoryVehicleStorage()
	rideStorage := NewInMemoryRideStorage()

	userMgr := NewUserManager(userStorage)
	vehicleMgr := NewVehicleManager(vehicleStorage, userMgr)
	rideMgr := NewRideManager(rideStorage, userMgr, vehicleMgr)

	user := User{ID: "1", Name: "Amar", Role: "Driver"}
	vehicle := Vehicle{ID: "1", OwnerID: "1", Model: "Toyota", Capacity: 4}
	userMgr.AddUser(user)
	vehicleMgr.AddVehicle(vehicle)
	ride := Ride{ID: "1", DriverID: "1", VehicleID: "1", Source: "A", Destination: "B", AvailableSeats: 3}
	rideMgr.OfferRide(ride)

	if err := rideMgr.EndRide(ride.ID); err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if _, err := rideStorage.GetRideByID(ride.ID); err == nil {
		t.Fatalf("Expected ride to be deleted, but it still exists")
	}
}

func TestFindMultipleRidesForMultipleSegments(t *testing.T) {
	// Create a storage
	rideStorage := NewInMemoryRideStorage()
	userStorage := NewInMemoryUserStorage()
	vehicleStorage := NewInMemoryVehicleStorage()

	userMgr := NewUserManager(userStorage)
	vehicleMgr := NewVehicleManager(vehicleStorage, userMgr)

	// Create a ride manager
	rideMgr := NewRideManager(rideStorage, userMgr, vehicleMgr)
	_ = userMgr.AddUser(User{ID: "1", Name: "Amar1", Role: "Driver"})
	_ = userMgr.AddUser(User{ID: "2", Name: "Amar2", Role: "Driver"})
	_ = userMgr.AddUser(User{ID: "3", Name: "Amar3", Role: "Driver"})
	_ = userMgr.AddUser(User{ID: "4", Name: "Amar4", Role: "Driver"})

	_ = vehicleMgr.AddVehicle(Vehicle{ID: "1", OwnerID: "1", Model: "XUV", Capacity: 7})
	_ = vehicleMgr.AddVehicle(Vehicle{ID: "2", OwnerID: "2", Model: "XUV", Capacity: 3})
	_ = vehicleMgr.AddVehicle(Vehicle{ID: "3", OwnerID: "3", Model: "XUV", Capacity: 5})
	_ = vehicleMgr.AddVehicle(Vehicle{ID: "4", OwnerID: "4", Model: "XUV", Capacity: 6})

	// Offer rides for different segments of the journey
	rides := []Ride{
		{ID: "1", DriverID: "1", VehicleID: "1", Source: "A", Destination: "B", AvailableSeats: 3},
		{ID: "2", DriverID: "2", VehicleID: "2", Source: "B", Destination: "C", AvailableSeats: 2},
		{ID: "3", DriverID: "3", VehicleID: "3", Source: "C", Destination: "D", AvailableSeats: 4},
		{ID: "4", DriverID: "4", VehicleID: "4", Source: "D", Destination: "E", AvailableSeats: 2},
	}
	for _, ride := range rides {
		err := rideMgr.OfferRide(ride)
		if err != nil {
			t.Fatalf("Error offering ride: %v", err)
		}
	}

	user := User{ID: "5", Name: "Amar", Role: "Passenger"}
	// Search for rides from A to E (no direct route)
	selectedRoutes, err := rideMgr.SelectRide(user.ID, "A", "E", 2, string(MostVacantSeats))
	if err != nil {
		t.Fatalf("Error finding rides: %v", err)
	}

	maxSeatAvail := 999
	for _, route := range selectedRoutes {
		if maxSeatAvail > route.AvailableSeats {
			maxSeatAvail = route.AvailableSeats
		}
	}

	// Check if multiple rides are available for the route
	expectedRideCount := 2 // A-B, B-C, C-D, D-E
	if maxSeatAvail != expectedRideCount {
		t.Fatalf("Expected %d rides available for the route, but got %d", expectedRideCount, maxSeatAvail)
	}
}
