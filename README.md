# Chat Application

This is a simple chat application built using Golang and MySQL.

## Database Setup

To set up the database, please run the following command:

```
docker run -d 
--name mysql_db 
-e MYSQL_ROOT_PASSWORD=root 
-e MYSQL_DATABASE=chat 
-e MYSQL_USER=user 
-e MYSQL_PASSWORD=password 
-v my-volume:/docker-entrypoint-initdb.d 
-v "$(pwd)/schema.sql:/docker-entrypoint-initdb.d/schema.sql" 
-p 3306:3306 
mysql:5.7
```

To check if the database is set up correctly, run the following commands:

```docker ps```

```docker exec -it mysql_db bash```

```mysql -u user -ppassword chat```

```show tables;```

## Running the Application

To run the application, navigate to the `cmd` folder and run the following command:

`go run main.go`


The application has the following endpoints:

- `POST /signup`: creates a new user account
- `POST /login`: logs in a user and returns a JWT token
- `GET /logout`: logs out the currently authenticated user
- `POST /ws/createRoom`: creates a new chat room
- `GET /ws/joinRoom/:roomId`: joins a chat room with the specified `roomId`
- `GET /ws/getRooms`: gets a list of all available chat rooms
- `GET /ws/getClients/:roomId`: gets a list of all clients in the specified chat room

## Dependencies

- Golang 1.16 or higher
- MySQL 5.7 or higher
- Docker (optional)

## Contributing

If you would like to contribute to this project, please feel free to submit a pull request.


