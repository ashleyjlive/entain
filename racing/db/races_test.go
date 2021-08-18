package db_test

import (
	"database/sql"
	"strconv"
	"testing"
	"time"

	"github.com/ashleyjlive/entain/racing/db"
	"github.com/ashleyjlive/entain/racing/proto/racing"
	"github.com/golang/protobuf/ptypes"
	"syreclabs.com/go/faker"
)

func GetTestDB() (*sql.DB, error) {
	return sql.Open("sqlite3", "races_testdata/test.db")
}

func EnsureDB(t *testing.T) {
	_, err := GetTestDB()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
}

func TestRepo(t *testing.T) {
	racingDB, err := GetTestDB()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	racesRepo := db.NewRacesRepo(racingDB)
	if err := racesRepo.Init(); err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
}

func TestNewRepo(t *testing.T) {
	racingDB, err := GetTestDB()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	racesRepo := db.NewRacesRepo(racingDB)
	err1 := racesRepo.Clear()

	if err1 != nil {
		t.Fatalf("Unable to clear test database %v", err)
	}
}

func TestPopulateRepo(t *testing.T) {
	racingDB, err := GetTestDB()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	racesRepo := db.NewRacesRepo(racingDB)
	err = racesRepo.Clear()
	if err != nil {
		t.Fatalf("Unable to clear test database %v", err)
	}

	if err != nil {
		t.Fatalf("Unabale to convert time.")
	}
	races := getRaces()
	err = racesRepo.InsertRace(races[0])
	if err != nil {
		t.Fatalf("Unable to insert record into database.")
	}
}

func TestPopulateAndFetchRepo(t *testing.T) {
	racingDB, err := GetTestDB()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	racesRepo := db.NewRacesRepo(racingDB)
	err = racesRepo.Clear()
	if err != nil {
		t.Fatalf("Unable to clear test database %v", err)
	}

	races := getRaces()
	err = racesRepo.InsertRace(races[0])
	if err != nil {
		t.Fatalf("Unable to insert record into database.")
	}
	var a []int64
	filter := racing.ListRacesRequestFilter{MeetingIds: a}
	a = append(a, races[0].MeetingId)
	rsp, err := racesRepo.List(&filter)
	if err != nil {
		t.Fatalf("Unable to retrieve races list.")
	}
	if len(rsp) != 1 {
		t.Fatalf("Returned incorrect amount of races.")
	}
	if rsp[0].MeetingId != a[0] {
		t.Fatalf("Incorrect item returned.")
	}
}

func TestPopulateAndFilterVisible(t *testing.T) {
	racingDB, err := GetTestDB()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	racesRepo := db.NewRacesRepo(racingDB)
	err = racesRepo.Clear()
	if err != nil {
		t.Fatalf("Unable to clear test database %v", err)
	}

	races := getRaces()
	err = racesRepo.InsertRace(races[0])
	if err != nil {
		t.Fatalf("Unable to insert record into database.")
	}
	visible := true
	filter := racing.ListRacesRequestFilter{Visible: &visible}
	rsp, err := racesRepo.List(&filter)
	if err != nil {
		t.Fatalf("Unable to retrieve races list.")
	}
	if len(rsp) == 0 && races[0].Visible == true {
		t.Fatalf("Unable to retrieve visible races list.")
	}
}

func getRaces() []*racing.Race {
	var (
		races []*racing.Race
	)
	for i := 1; i <= 100; i++ {
		meetingId, _ := strconv.Atoi(faker.Number().Between(1, 10))
		name := faker.Team().Name()
		number, _ := strconv.Atoi(faker.Number().Between(1, 12))
		visible := randBool()
		tm, _ := ptypes.TimestampProto(
			faker.Time().Between(
				time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)))
		race :=
			racing.Race{Id: int64(i), MeetingId: int64(meetingId),
				Name: name, Number: int64(number),
				Visible: visible, AdvertisedStartTime: tm}
		races = append(races, &race)
	}
	return races
}

func randBool() bool {
	v := faker.Number().Between(0, 1)
	if v == "0" {
		return false
	}
	return true
}
