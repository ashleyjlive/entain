package db

const (
	eventsList = "list"
)

// Defines all the queries for the events table.
func getEventsQueries() map[string]string {
	return map[string]string{
		eventsList: `
			SELECT 
				id, 
				name,
				category, 
				advertised_start_time 
			FROM events
		`,
	}
}
