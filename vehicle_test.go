package main

import (
	"testing"
)

// Test adding a vehicle
func TestAddVehicle(t *testing.T) {
	vehicleStorage := NewInMemoryVehicleStorage()
	vehicleMgr := NewVehicleManager(vehicleStorage, nil)

	vehicle := Vehicle{ID: "1", OwnerID: "1", Model: "Toyota", Capacity: 4}
	if err := vehicleMgr.AddVehicle(vehicle); err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	retrievedVehicle, err := vehicleStorage.GetVehicleByID(vehicle.ID)
	if err != nil {
		t.Fatalf("Expected to retrieve vehicle, but got error %v", err)
	}
	if retrievedVehicle != vehicle {
		t.Fatalf("Expected vehicle to be %v, but got %v", vehicle, retrievedVehicle)
	}
}
