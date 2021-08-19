package db

import (
	"database/sql"
	"errors"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/ashleyjlive/entain/sports/proto/sports"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"
)

// EventsRepo provides repository access to events.
type EventsRepo interface {
	// Init will initialise our events repository.
	Init(bool) error

	// Clears all entries from the repository.
	Clear() error

	// List will return a list of sports events.
	List(request *sports.ListEventsRequest) ([]*sports.Event, error)

	// Inserts a sport event entry into the repository.
	InsertRace(*sports.Event) error

	// Returns all sport events.
	ListAll() ([]*sports.Event, error)
}

type eventsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewEventsRepo creates a new events repository.
func NewEventsRepo(db *sql.DB) EventsRepo {
	return &eventsRepo{db: db}
}

// Init prepares the events repository dummy data.
func (r *eventsRepo) Init(seed bool) error {
	var err error

	r.init.Do(func() {
		if seed {
			if err = r.seed(); err != nil {
				err = r.init_tbl()
			}
		} else {
			err = r.init_tbl()
		}
	})

	return err
}

// Clears all data in the events repository.
func (r *eventsRepo) Clear() error {
	return r.clear()
}

// Allows insertions of a race into the repository.
func (r *eventsRepo) InsertRace(event *sports.Event) error {
	return r.insert(event)
}

func (r *eventsRepo) List(request *sports.ListEventsRequest) ([]*sports.Event, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getEventsQueries()[eventsList]

	query, args = r.applyFilter(query, request.Filter)
	query = r.applyOrdering(query, request.OrderBy)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanEvents(rows)
}

func (r *eventsRepo) ListAll() ([]*sports.Event, error) {
	return r.listAll()
}

func (r *eventsRepo) applyFilter(query string, filter *sports.ListEventsRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if filter.Category != nil {
		clauses = append(clauses, "LOWER(`category`) LIKE ?")
		args = append(args, strings.ToLower(*filter.Category))
	}

	if len(filter.Ids) > 0 {
		clauses = append(clauses, "id IN ("+strings.Repeat("?,", len(filter.Ids)-1)+"?)")

		for _, id := range filter.Ids {
			args = append(args, id)
		}
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (r *eventsRepo) applyOrdering(query string, orderBy *string) string {
	const defaultOrder = " ORDER BY advertised_start_time"
	if orderBy == nil {
		query += defaultOrder
	} else {
		// DB implementation doesn't allow prepared statements for ORDER BY
		// input variables.
		// To allow for dynamic and safe ordering, we must sanitize and rebuild
		// in SQL friendly format.
		orderStr, err := toOrderBySql(*orderBy)
		if err != nil {
			query += defaultOrder
		} else if orderStr != nil {
			query += " ORDER BY " + *orderStr
		} else {
			query += defaultOrder
		}
	}
	return query
}

// Accepts `order_by` definition from Google API standards [1] and returns an
// output in SQL friendly format.
// [1] - https://cloud.google.com/apis/design/design_patterns#sorting_order
func toOrderBySql(input string) (*string, error) {
	var (
		terms []string
	)
	// 1. Splits the input string by comma.
	// 2. Then for each element, determine the words (max of 2, min of 1).
	// 3. Ensure that the first (column name) only contains valid chars.
	// 4. Ensure that if sort parameter is provided, it is "asc" or "desc" only.
	// 5. Rebuilds the string in CSV (SQL ORDER BY) format.
	for _, str := range strings.Split(input, ",") {
		words := strings.Fields(str)
		wordCount := len(words)
		if wordCount > 2 || wordCount < 1 {
			return nil, errors.New("invalid order by term count")
		}
		sortField := words[0]
		if strings.IndexFunc(sortField, isUnsafeColumnChar) != -1 {
			return nil, errors.New("invalid column name")
		}
		if wordCount == 2 {
			sort := words[1]
			if !(strings.EqualFold(sort, "asc") || strings.EqualFold(sort, "desc")) {
				return nil, errors.New("invalid order by dir parameter")
			}
			sortField += " " + sort
		}
		terms = append(terms, sortField)
	}
	output := strings.Join(terms, ",")
	return &output, nil
}

// Determines if the supplied rune is safe to be used in the column name.
// This err's on the side of caution and makes no attempt for numerics or
// escaped special chars.
func isUnsafeColumnChar(c rune) bool {
	switch c {
	case '_':
		// Edge case for column names with underscores.
		return false
	default:
		return !unicode.IsLetter(c)
	}
}

func (m *eventsRepo) scanEvents(
	rows *sql.Rows,
) ([]*sports.Event, error) {
	var events []*sports.Event

	for rows.Next() {
		var event sports.Event
		var advertisedStart time.Time

		if err := rows.Scan(&event.Id, &event.Name, &event.Category, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		event.AdvertisedStartTime = ts

		events = append(events, &event)
	}

	return events, nil
}
