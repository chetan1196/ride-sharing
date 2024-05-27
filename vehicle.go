package main

import "fmt"

type Vehicle struct {
	ID       string
	OwnerID  string
	Model    string
	Capacity int
}

type vehicleManager struct {
	storage VehicleStorage
	userMgr *userManager
}

func NewVehicleManager(storage VehicleStorage, userMgr *userManager) *vehicleManager {
	return &vehicleManager{storage: storage, userMgr: userMgr}
}

func (vm *vehicleManager) AddVehicle(vehicle Vehicle) error {
	// if _, err := vm.userMgr.GetUserByID(vehicle.OwnerID); err != nil {
	// 	return fmt.Errorf("owner %s not found: %v", vehicle.OwnerID, err)
	// }
	if err := vm.storage.AddVehicle(vehicle); err != nil {
		return fmt.Errorf("could not add vehicle: %v", err)
	}
	fmt.Printf("Vehicle added: %+v\n", vehicle)
	return nil
}

func (vm *vehicleManager) GetVehicleByID(vehicleID string) (Vehicle, error) {
	vehicle, err := vm.storage.GetVehicleByID(vehicleID)
	if err != nil {
		return Vehicle{}, fmt.Errorf("could not find vehicle %s: %v", vehicleID, err)
	}
	return vehicle, nil
}

func (vm *vehicleManager) ValidateVehicle(vehicleID string, availSeats int) error {
	vehicle, err := vm.storage.GetVehicleByID(vehicleID)
	if err != nil {
		return fmt.Errorf("failed to find the vehicle %s: %v", vehicleID, err)
	}
	if vehicle.Capacity < availSeats {
		return fmt.Errorf("available seats exceedes vehicle capacity")
	}
	return nil
}
