# remote-storage
This project is the backend part of a Remote Storage Application. The application allows authorized users to manage their remote directory, including uploading and downloading files, renaming, moving, and copying files within the remote directory.
The project is built using a microservice architecture with the Go kit (gokit) framework and includes two microservices:
1)Storage: This service provides a separate filesystem for every user, enabling them to control it.
2)Authentication Microservice: This microservice uses JWT authentication to check user credentials, validate, and refresh JWT tokens.
Microservices are discovered with Consul service. Each microservice contains a client package that returns a service, working with the Remote Procedure Call (RPC) approach using HTTP/JSON as the transport protocol. Microservices include logging middleware, and the storage service includes authentication middleware.
For storing user data, I use a PostgreSQL database. The database package provides a higher level of abstraction for working with the database.

