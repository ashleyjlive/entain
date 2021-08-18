# sports/proto
This defines the protobuf message types for the sports API.

## sports API

**RPCS**

- `ListEvents(ListEventsRequest) ListEventsResponse`

### ListEventsRequest
- Supports a filter parameter of type `ListEventsRequestFilter`.
- Supports a order_by parameter of the form defined by [Google API Design](https://cloud.google.com/apis/design/design_patterns#sorting_order).

### ListEventsRequestFilter
- A list of integer IDs can be supplied to perform a bulk lookup request. (optional)
- A category filter can be supplied which finds any races which match the category (optional, case-insensitive).

### ListEventsResponse
- Contains a list of all matching events for the given lookup.

### Event
- The ID of the event (int64).
- The name of the event (string).
- The category of the event (string: e.g. "AFL").
- The advertised start time of the event (Timestamp).