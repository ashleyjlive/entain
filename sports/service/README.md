# sports/service

Defines the sports service and associated interface functions.

## Creating service

To create a new sports service:
```
sportsService := sports.NewSportsService(eventsRepo)
```
This service can then be registered under the sports server when initialising:
```
sports.RegisterSportsServer(grpcServer,sportsService)
```

## Interface

`ListEvents(context.Context, *sports.ListEventsRequest) (*sports.ListEventsResponse, error)`
