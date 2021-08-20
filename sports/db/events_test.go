package db_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ashleyjlive/entain/sports/db"
	"github.com/ashleyjlive/entain/sports/proto/sports"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestRepo(t *testing.T) {
	eventDB, err := GetTestDB("events", "TestRepo")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	eventsRepo := db.NewEventsRepo(eventDB)
	if err := eventsRepo.Init(false); err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
}

func TestNewRepo(t *testing.T) {
	racingDB, err := GetTestDB("events", "TestNewRepo")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	eventsRepo := db.NewEventsRepo(racingDB)
	events, _ := eventsRepo.ListAll()

	if len(events) != 0 {
		t.Fatal("New repo contains elements")
	}
}

func TestList(t *testing.T) {
	eventDB, err := GetTestDB("events", "GetList")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	eventsRepo := db.NewEventsRepo(eventDB)
	if err := eventsRepo.Init(false); err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	tm1 := timestamppb.New(time.Now().AddDate(0, 0, 2))
	race1 :=
		sports.Event{Id: int64(1),
			Name:                "Test1",
			Category:            "Soccer",
			AdvertisedStartTime: tm1}
	eventsRepo.InsertRace(&race1)

	tm2 := timestamppb.New(time.Now().AddDate(0, 0, -2))
	race2 :=
		sports.Event{Id: int64(2),
			Name:                "Test2",
			Category:            "Basketball",
			AdvertisedStartTime: tm2}
	eventsRepo.InsertRace(&race2)

	events, err := eventsRepo.ListAll()
	if err != nil {
		t.Fatalf("Unable to return list of events %v", err)
	}

	if events[0].Id != race1.Id || events[1].Id != race2.Id {
		t.Fatalf("List returned unexpected order V1: %v, V2: %v",
			events[0].Id, events[1].Id)
	}
}

func TestListOrdered(t *testing.T) {
	eventDB, err := GetTestDB("events", "GetListOrdered")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	eventsRepo := db.NewEventsRepo(eventDB)
	if err := eventsRepo.Init(false); err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	tm1 := timestamppb.New(time.Now().AddDate(0, 0, 2))
	race1 :=
		sports.Event{Id: int64(1),
			Name:                "Test1",
			Category:            "Soccer",
			AdvertisedStartTime: tm1}
	eventsRepo.InsertRace(&race1)

	tm2 := timestamppb.New(time.Now().AddDate(0, 0, -2))
	race2 :=
		sports.Event{Id: int64(2),
			Name:                "Test2",
			Category:            "Basketball",
			AdvertisedStartTime: tm2}
	eventsRepo.InsertRace(&race2)

	rq := sports.ListEventsRequest{}
	// List events by default sorts by advertised start time.
	events, err := eventsRepo.List(&rq)
	if err != nil {
		t.Fatalf("Unable to return list of events %v", err)
	}

	// Expect Race2 to be first and Race1 last due to AdvertisedStartTime.
	if events[0].Id != race2.Id || events[1].Id != race1.Id {
		t.Fatalf("List returned unexpected order V1: %v, V2: %v",
			events[0].Id, events[1].Id)
	}
}

func TestListCategoryFilter(t *testing.T) {
	eventDB, err := GetTestDB("events", "GetListCategoryFilter")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	eventsRepo := db.NewEventsRepo(eventDB)
	if err := eventsRepo.Init(false); err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	tm1 := timestamppb.New(time.Now().AddDate(0, 0, 2))
	event1 :=
		sports.Event{Id: int64(1),
			Name:                "Test1",
			Category:            "Soccer",
			AdvertisedStartTime: tm1}
	eventsRepo.InsertRace(&event1)

	tm2 := timestamppb.New(time.Now().AddDate(0, 0, -2))
	event2 :=
		sports.Event{Id: int64(2),
			Name:                "Test2",
			Category:            "Basketball",
			AdvertisedStartTime: tm2}
	eventsRepo.InsertRace(&event2)

	cat := "Soccer"
	filter := sports.ListEventsRequestFilter{Category: &cat}
	rq := sports.ListEventsRequest{Filter: &filter}
	// List events by default sorts by advertised start time.
	events, err := eventsRepo.List(&rq)
	if err != nil {
		t.Fatalf("Unable to return list of events %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("Category filter did not return expected dataset length.")
	}

	if events[0].Id != event1.Id {
		t.Fatalf("List returned unexpected order V1: %v, V2: %v",
			events[0].Id, events[1].Id)
	}
}

func TestListIdFilter(t *testing.T) {
	eventDB, err := GetTestDB("events", "TestListIdFilter")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	eventsRepo := db.NewEventsRepo(eventDB)
	if err := eventsRepo.Init(false); err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	tm1 := timestamppb.New(time.Now().AddDate(0, 0, 20))
	event1 :=
		sports.Event{Id: int64(10),
			Name:                "Foo",
			Category:            "Soccer",
			AdvertisedStartTime: tm1}
	eventsRepo.InsertRace(&event1)

	tm2 := timestamppb.New(time.Now().AddDate(0, 0, -4))
	event2 :=
		sports.Event{Id: int64(212),
			Name:                "Baz",
			Category:            "Basketball",
			AdvertisedStartTime: tm2}
	eventsRepo.InsertRace(&event2)

	var ids []int64
	ids = append(ids, event1.Id)
	filter := sports.ListEventsRequestFilter{Ids: ids}
	rq := sports.ListEventsRequest{Filter: &filter}
	// List events by default sorts by advertised start time.
	events, err := eventsRepo.List(&rq)
	if err != nil {
		t.Fatalf("Unable to return list of events %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("Category filter did not return expected dataset length.")
	}

	if events[0].Id != event1.Id {
		t.Fatalf("List returned unexpected order V1: %v, V2: %v",
			events[0].Id, events[1].Id)
	}
}

// Helpers //

func GetTestDB(testType string, testName string) (*sql.DB, error) {
	dir := filepath.Join(testType+"_testdata", testName)
	err := os.RemoveAll(dir)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		return nil, err
	}
	return sql.Open("sqlite3", filepath.Join(dir, "test.db"))
}
