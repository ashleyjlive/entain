# racing/proto
This defines the protobuf message types for the racing API.

## Racing API

**RPCS**

- `ListRaces(ListRacesRequest) ListRacesResponse`

### ListRacesRequest
- Supports a filter parameter of type `ListRacesRequestFilter`.
- Supports a order_by parameter of the form defined by [Google API Design](https://cloud.google.com/apis/design/design_patterns#sorting_order).

### ListRacesRequestFilter
- A list of integer IDs can be supplied to perform a bulk lookup request. (optional)
- A category filter can be supplied which finds any races which match the category (optional, case-insensitive).

### ListRacesResponse
- Contains a list of all matching races for the given lookup.

### Race
- The ID of the race (int64).
- The unique identifier for the races meeting (int64).
- The name of the race (string).
- The visibility of the race (bool).
- The advertised start time of the race (Timestamp).
- A status flag indicating if the event is `OPEN` or `CLOSED` - this is based off the advertised start time with the systems current time.