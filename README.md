# Ride-Sharing Application

This is a simple ride-sharing application written in Go. It allows users to offer and consume shared rides. The application supports user management, vehicle management, ride offering, ride selection, and ride statistics.

## Features

- **User Roles**: Users can either offer a shared ride (Driver) or consume a shared ride (Passenger).
- **Ride Selection**: Users can search and select from multiple available rides on a route with the same source and destination.
- **Statistics**: Retrieve and display total rides offered/taken by all users.

## Requirements

- **Go 1.22  and above**: 

## Build and Run
3. Build and run the application using the following command:
   *go build -o ride-sharing && ./ride-sharing*

## Sample Output
```User added: {1 Amar Driver}
User added: {2 Chetan Driver}
User added: {3 Bhuwan Passenger}
User added: {4 Vijay Passenger}
Vehicle added: {ID:1 OwnerID:1 Model:Toyota Capacity:4}
Vehicle added: {ID:2 OwnerID:2 Model:XUV Capacity:7}
Ride offered: {ID:101 DriverID:1 VehicleID:1 Source:A Destination:B AvailableSeats:4}
Ride offered: {ID:102 DriverID:2 VehicleID:2 Source:B Destination:C AvailableSeats:4}
No rides available directly: searching for rides through indirect routes.
Indirect Rides selected: [{ID:101 DriverID:1 VehicleID:1 Source:A Destination:B AvailableSeats:4} {ID:102 DriverID:2 VehicleID:2 Source:B Destination:C AvailableSeats:4}]
Ride selected: {ID:101 DriverID:1 VehicleID:1 Source:A Destination:B AvailableSeats:0}
Ride statistics:
User Amar: Offered:1: Taken: 0
User Chetan: Offered:1: Taken: 0
User Bhuwan: Offered:0: Taken: 2
User Vijay: Offered:0: Taken: 1
