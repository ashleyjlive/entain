package db

import (
	"database/sql"
	"time"

	"github.com/ashleyjlive/entain/sports/proto/sports"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"syreclabs.com/go/faker"
)

func (r *eventsRepo) init_tbl() error {
	statement, err := r.db.Prepare(`CREATE TABLE IF NOT EXISTS events (id INTEGER PRIMARY KEY, name TEXT, category TEXT, advertised_start_time DATETIME)`)
	if err == nil {
		_, err = statement.Exec()
		return err
	} else {
		return err
	}
}

func (r *eventsRepo) seed() error {
	statement, err := r.db.Prepare(`CREATE TABLE IF NOT EXISTS events (id INTEGER PRIMARY KEY, name TEXT, category TEXT, advertised_start_time DATETIME)`)
	if err == nil {
		_, err = statement.Exec()
	}

	for i := 1; i <= 100; i++ {
		statement, err = r.db.Prepare(`INSERT OR IGNORE INTO events(id, name, category, advertised_start_time) VALUES (?,?,?,?)`)
		if err == nil {
			_, err = statement.Exec(
				i,
				faker.Team().Name(),
				faker.RandomChoice(categories()),
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339),
			)
		}
	}

	return err
}

func (r *eventsRepo) listAll() ([]*sports.Event, error) {
	var events []*sports.Event
	rows, err :=
		r.db.Query(
			"SELECT id, name, category, advertised_start_time FROM events")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var event sports.Event
		var advertisedStart time.Time

		if err := rows.Scan(&event.Id, &event.Name, &event.Category, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}
		ts := timestamppb.New(advertisedStart)

		event.AdvertisedStartTime = ts

		events = append(events, &event)
	}
	return events, nil
}

func (r *eventsRepo) insert(event *sports.Event) error {
	var statement *sql.Stmt
	ts, err := ptypes.Timestamp(event.AdvertisedStartTime)
	if err != nil {
		return err
	}
	statement, err = r.db.Prepare(`INSERT INTO events(id, name, category, advertised_start_time) VALUES (?,?,?,?)`)
	if err == nil {
		_, err = statement.Exec(
			&event.Id,
			&event.Name,
			&event.Category,
			ts,
		)
	}
	return err
}

func (r *eventsRepo) clear() error {
	_, err := r.db.Exec("DELETE FROM events")
	return err
}

func categories() []string {
	return []string{
		"AFL",
		"Rugby League",
		"Soccer",
		"Basketball",
	}
}
