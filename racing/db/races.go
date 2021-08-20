package db

import (
	"database/sql"
	"errors"
	"strings"
	"sync"
	"time"
	"unicode"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ashleyjlive/entain/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init(bool) error
	Clear() error
	InsertRace(*racing.Race) error
	// List will return a list of races.
	List(request *racing.ListRacesRequest) ([]*racing.Race, error)
	ListAll() ([]*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init(seed bool) error {
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

// Clears all data in the races repository.
func (r *racesRepo) Clear() error {
	return r.clear()
}

// Allows insertions of a race into the repository.
func (r *racesRepo) InsertRace(race *racing.Race) error {
	return r.insert(race)
}

func (r *racesRepo) List(request *racing.ListRacesRequest) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, request.Filter)
	query = r.applyOrdering(query, request.OrderBy)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

func (r *racesRepo) ListAll() ([]*racing.Race, error) {
	return r.listAll()
}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if filter.Visible != nil {
		// Visibility is optional. If nil, then allow non filtering of
		// visibility.
		clauses = append(clauses, "visible = ?")
		args = append(args, filter.Visible)
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (r *racesRepo) applyOrdering(query string, orderBy *string) string {
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

func (m *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts := timestamppb.New(advertisedStart)

		if advertisedStart.Before(time.Now()) {
			// All races that have an `advertised_start_time` in the past should
			// reflect `CLOSED`
			// Note that this depends on the system having the correct time.
			race.Status = racing.Race_CLOSED
		} else {
			race.Status = racing.Race_OPEN
		}
		race.AdvertisedStartTime = ts

		races = append(races, &race)
	}

	return races, nil
}
