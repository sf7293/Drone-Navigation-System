# Drone-Navigation-System
A simple project which calculates location of space sectors by x,y,z coordinations and velocity of a drone.
To run the project, Do these steps:
- `go mod vendor` (and `go mod tidy` may be needed)
- `go run main.go`

This project has only one POST API, you could access it with this simple curl:
`curl --location --request POST 'http://127.0.0.1:8080/v1/dns/location' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'x=0.2' \
--data-urlencode 'y=2.3' \
--data-urlencode 'z=3.4' \
--data-urlencode 'vel=5'`

To build the binary of the project 
