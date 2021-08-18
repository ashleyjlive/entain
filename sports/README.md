# Sports Service

The sports service provides the ability to retrieve sports events from a given
SQL database using its gRPC API.

## Directory Structure

### [`db`](db/README.md)

Defines the data storage implementation for the sports service.
### [`proto`](proto/README.md)

Declares the gRPC API to use when interacting with the front facing API server.
### [`service`](service/README.md)

Implements the service and its interactions with the data storage implementation.
