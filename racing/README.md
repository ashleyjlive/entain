# Racing Service
The racing service provides the ability to access racing events over a gRPC server.

## Building
To build a executable simply call

    $ go build

This will place a racing executable in the root directory.

## Command Line

You may provide optional command line arguments to the executable.
Currently, you may configure:
- `grpc-endpoint` - This is the endpoint that the front facing API server will speak to.
- `db_path` - This is the path of the database that the service will utilise.

For example:

    $ ./racing --grpc-endpoint=localhost:8080 --db_path:/foo/bar/db.db

## API

Please [see](proto/README.md) the documentation for the protobuf definitions.

## Testing

To test individual packages:

    $ go test ./...