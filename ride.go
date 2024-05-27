package main

import (
	"fmt"
	"strings"
	"sync"
)

type Strategy string

const (
	PreferredVehicle Strategy = "Preferred Vehicle"
	MostVacantSeats  Strategy = "Most Vacant"
)

type Ride struct {
	ID             string
	DriverID       string
	VehicleID      string
	Source         string
	Destination    string
	AvailableSeats int
}

type rideManager struct {
	mu          sync.Mutex
	storage     RideStorage
	rideStats   map[string]stats // total rides offered/taken by user
	userMgr     *userManager
	vehicleMgr  *vehicleManager
	activeRides map[string]bool // Mapping of ride ID to active status
}

type stats struct {
	offered int
	taken   int
}

func NewRideManager(storage RideStorage, usersMgr *userManager, vehicleMgr *vehicleManager) *rideManager {
	return &rideManager{
		mu:          sync.Mutex{},
		storage:     storage,
		rideStats:   make(map[string]stats),
		activeRides: make(map[string]bool),
		userMgr:     usersMgr,
		vehicleMgr:  vehicleMgr,
	}
}

func (rm *rideManager) GetDirectRides(source, destination string) []Ride {
	var result []Ride
	for _, ride := range rm.storage.GetAllRides() {
		if ride.Source == source && ride.Destination == destination && ride.AvailableSeats > 0 {
			result = append(result, ride)
		}
	}
	return result
}

// GetRidesByVehicle retrieves all rides associated with a vehicle.
func (rm *rideManager) GetRidesByVehicle(vehicleID string) []Ride {
	var vehicleRides []Ride
	for _, ride := range rm.storage.GetAllRides() {
		if ride.VehicleID == vehicleID {
			vehicleRides = append(vehicleRides, ride)
		}
	}
	return vehicleRides
}

func (rm *rideManager) GetRidesBySource(source string) []Ride {
	var rides []Ride
	for _, ride := range rm.storage.GetAllRides() {
		if ride.Source == source {
			rides = append(rides, ride)
		}
	}
	return rides
}

// GetRidesByDriver retrieves all rides associated with a driver.
func (rm *rideManager) GetRidesByDriver(driverID string) []Ride {
	var driverRides []Ride
	for _, ride := range rm.storage.GetAllRides() {
		if ride.DriverID == driverID {
			driverRides = append(driverRides, ride)
		}
	}
	return driverRides
}

func (rm *rideManager) OfferRide(ride Ride) error {
	// validate the driver
	if err := rm.userMgr.IsDriver(ride.DriverID); err != nil {
		return err
	}
	if err := rm.vehicleMgr.ValidateVehicle(ride.VehicleID, ride.AvailableSeats); err != nil {
		return err
	}

	// Check if the driver is already offering a ride
	for _, existingRide := range rm.GetRidesByDriver(ride.DriverID) {
		if rm.isActive(existingRide.ID) {
			return fmt.Errorf("driver %s is already offering a ride", ride.DriverID)
		}
	}

	// Check if the vehicle is already in use for a ride
	for _, existingRide := range rm.GetRidesByVehicle(ride.VehicleID) {
		if rm.isActive(existingRide.ID) {
			return fmt.Errorf("vehicle %s is already in use for a ride", ride.VehicleID)
		}
	}

	// If no conflicts, add the ride
	if err := rm.storage.AddRide(ride); err != nil {
		return fmt.Errorf("could not offer ride: %v", err)
	}
	rm.activeRides[ride.ID] = true
	fmt.Printf("Ride offered: %+v\n", ride)

	rm.updateOfferedStats(ride.DriverID)
	return nil
}

// isActive checks if a ride is active.
func (rm *rideManager) isActive(rideID string) bool {
	_, isActive := rm.activeRides[rideID]
	return isActive
}

func (rm *rideManager) EndRide(rideID string) error {
	if err := rm.storage.DeleteRide(rideID); err != nil {
		return fmt.Errorf("could not end ride: %v", err)
	}
	delete(rm.activeRides, rideID)
	fmt.Printf("Ride ended: %v\n", rideID)
	return nil
}

func (rm *rideManager) PrintRideStats() {
	fmt.Println("Ride statistics:")
	rm.mu.Lock()
	for _, user := range rm.userMgr.storage.GetAllUsers() {
		fmt.Printf("User %s: Offered:%d: Taken: %d\n", user.Name, rm.rideStats[user.ID].offered, rm.rideStats[user.ID].taken)
	}
	rm.mu.Unlock()
}

func (rm *rideManager) isPreferredVehicle(vehicleID, preferredVehicle string) bool {
	vehicle, err := rm.vehicleMgr.GetVehicleByID(vehicleID)
	if err != nil {
		return false
	}
	return len(preferredVehicle) == 0 || vehicle.Model == preferredVehicle
}

func (rm *rideManager) incrementTakenStats(userID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	stats := rm.rideStats[userID]
	stats.taken++
	rm.rideStats[userID] = stats
}
func (rm *rideManager) decrementTakenStats(userID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	stats := rm.rideStats[userID]
	stats.taken--
	rm.rideStats[userID] = stats
}

func (rm *rideManager) updateOfferedStats(driverID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	stats := rm.rideStats[driverID]
	stats.offered++
	rm.rideStats[driverID] = stats
}

// FindRides finds rides for the given source, destination, and required seats
func (rm *rideManager) FindInDirectRoute(userID, source, destination string, seats int, preferredVehicle string) ([]Ride, error) {
	var selectedRides []Ride
	visited := make(map[string]bool)

	var dfs func(current, dest string) bool
	dfs = func(current, dest string) bool {
		if current == dest {
			return true
		}

		visited[current] = true

		// Find rides from the current source
		rides := rm.GetRidesBySource(current)

		// Iterate through rides to find possible paths
		for _, ride := range rides {
			if !visited[ride.Destination] && ride.AvailableSeats >= seats && rm.isPreferredVehicle(ride.VehicleID, preferredVehicle) {
				selectedRides = append(selectedRides, ride)
				rm.incrementTakenStats(userID)
				ride.AvailableSeats -= seats
				rm.storage.UpdateRide(ride)
				if dfs(ride.Destination, dest) {
					return true
				}
				selectedRides = selectedRides[:len(selectedRides)-1] // Backtrack
				rm.decrementTakenStats(userID)
				ride.AvailableSeats += seats
				rm.storage.UpdateRide(ride)
			}
		}
		return false
	}

	// Perform DFS from the source to find rides to the destination
	if !dfs(source, destination) {
		return nil, fmt.Errorf("no rides available for the route")
	}

	return selectedRides, nil
}

func (rm *rideManager) SelectRide(userID, source, destination string, seats int, preference string) ([]Ride, error) {
	strategy := preference
	preferedVehicle := ""
	if strategy != string(MostVacantSeats) {
		strategy, preferedVehicle = strings.Split(strategy, "=")[0], strings.Split(strategy, "=")[1]
	}

	rides := rm.GetDirectRides(source, destination)
	if len(rides) == 0 {
		fmt.Println("No rides available directly: searching for rides through indirect routes.")
		indirectRoute, err := rm.FindInDirectRoute(userID, source, destination, seats, preferedVehicle)
		if err != nil {
			return nil, fmt.Errorf("failed to find indirect routes: %v", err)
		}
		fmt.Printf("Indirect Rides selected: %+v\n", indirectRoute)
		return indirectRoute, nil
	}

	var selectedRide Ride
	switch strategy {
	case string(PreferredVehicle):
		for _, ride := range rides {
			vehicle, _ := rm.vehicleMgr.GetVehicleByID(ride.VehicleID)
			if vehicle.Model == preferedVehicle && ride.AvailableSeats >= seats {
				selectedRide = ride
				break
			}
		}
	case string(MostVacantSeats):
		maxSeats := -1
		for _, ride := range rides {
			if ride.AvailableSeats >= seats && ride.AvailableSeats > maxSeats {
				maxSeats = ride.AvailableSeats
				selectedRide = ride
			}
		}
	default:
		return nil, fmt.Errorf("unknown selection strategy")
	}

	if selectedRide.ID == "" {
		return nil, fmt.Errorf("no suitable ride found")
	}

	selectedRide.AvailableSeats -= seats
	if err := rm.storage.UpdateRide(selectedRide); err != nil {
		return nil, fmt.Errorf("could not update ride: %v", err)
	}

	fmt.Printf("Ride selected: %+v\n", selectedRide)
	rm.incrementTakenStats(userID)
	return []Ride{selectedRide}, nil
}
