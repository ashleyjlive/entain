# Sports Service

The sports service provides the ability to retrieve sports events from a given
SQL database using its gRPC API.

# Building
To build a executable simply call

    $ go build

This will place a sports executable in the root directory.

## Command Line

You may provide optional command line arguments to the executable.
Currently, you may configure:
- `grpc-endpoint` - This is the endpoint that the front facing API server will speak to.
- `db_path` - This is the path of the database that the service will utilise.

For example:

    $ ./sports --grpc-endpoint=localhost:8080 --db_path:/foo/bar/db.db

## API

Please [see](proto/README.md) the documentation for the protobuf definitions.

## Testing

To test individual packages:

    $ go test ./...

## Directory Structure

### [`db`](db/README.md)

Defines the data storage implementation for the sports service.
### [`proto`](proto/README.md)

Declares the gRPC API to use when interacting with the front facing API server.
### [`service`](service/README.md)

Implements the service and its interactions with the data storage implementation.