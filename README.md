# Hotel reservation backend

## Project outline
- Users -> book room from a hotel
- Admins -> going to check reservation
- Authentication and authorization -> JWT tokens
- Hotels -> CRUD API -> JSON
- Rooms -> CRUD API -> JSON
- Scripts -> database management -> seeding, migration 

## Resources
### Mongodb driver
Documentation
```
https://gofiber.io
```
Install Mongodb
```
go get go.mongodb.org/mongo-driver/mongo
```
### gofiber
Documentation
```
https://gofiber.io
```
Install gofiber
```
go get github.com/gofiber/fiber/v2
```

## Docker
### Installing mongodb as a Docker container
```
docker run --name mongodb -d mongo:latest -p 27017:27017
```