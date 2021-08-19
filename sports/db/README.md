# sports/db

### [db.go](db.go)
This provides functionality in order to test the sports subsystem by populating
a given SQL database with randomly generated data.
<br>
To populate the DB with test data simply call:
```
err := repo.seed()
```
where `repo` is of type `eventsRepo`.

### [queries.go](queries.go)
`queries.go` provides a set of functions which have predefined queries for their 
associated tables.
<br>
For example, `getEventsQueries()` will return a map of queries that can be used
to retrieve specific data for the events dataset (such as a list all request). 

### [events.go](events.go)
`events.go` provides access to the storage layer for sports events.

**Usage**

To obtain accss to the repository, simply call:
```
eventsRepo := db.NewEventsRepo(eventsDB)
```
where `eventsDB` is of type `sql.DB`.

**Initialising**

To initialise and populate the table if the table does not exist:
```
if err := eventsRepo.Init(); err != nil {
    return err
}
```

**Interface**

The `EventsRepo` interface currently implements the following functions:
- `Init(bool) error` - The input boolean determines if data seeding is required.
- `Clear() error` - Clears all entries from the data storage. 
- `List(request *sports.ListEventsRequest) ([]*sports.Event, error)` - Lists all sporting events given the inbound request.
- `InsertRace(*sports.Event) error` - Inserts a race into the data storage.
- `ListAll()  ([]*sports.Event, error)` - Returns all entries in the data storage.