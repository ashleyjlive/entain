package db_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/ashleyjlive/entain/racing/db"
	"github.com/ashleyjlive/entain/racing/proto/racing"
	"github.com/golang/protobuf/ptypes"
	"syreclabs.com/go/faker"
)

func TestRepo(t *testing.T) {
	racingDB, err := GetTestDB("races", "TestRepo")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	racesRepo := db.NewRacesRepo(racingDB)
	if err := racesRepo.Init(false); err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
}

func TestNewRepo(t *testing.T) {
	racingDB, err := GetTestDB("races", "TestNewRepo")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	racesRepo := db.NewRacesRepo(racingDB)
	races, err := racesRepo.ListAll()

	if len(races) != 0 {
		t.Fatal("New repo contains elements")
	}
}

func TestPopulateRepo(t *testing.T) {
	racingDB, err := GetTestDB("races", "TestPopulateRepo")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	racesRepo := db.NewRacesRepo(racingDB)
	_ = racesRepo.Init(false)
	err = racesRepo.Clear()
	if err != nil {
		t.Fatalf("Unable to clear test database %v", err)
	}

	if err != nil {
		t.Fatalf("Unabale to convert time.")
	}
	races := GetRaces()
	err = racesRepo.InsertRace(races[0])
	if err != nil {
		t.Fatalf("Unable to insert record into database.")
	}
}

func TestPopulateAndFetchRepo(t *testing.T) {
	racingDB, err := GetTestDB("races", "TestPopulateAndFetchRepo")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	racesRepo := db.NewRacesRepo(racingDB)
	_ = racesRepo.Init(false)
	err = racesRepo.Clear()
	if err != nil {
		t.Fatalf("Unable to clear test database %v", err)
	}

	races := GetRaces()
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
	racingDB, err := GetTestDB("races", "TestPopulateAndFilterVisible")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	racesRepo := db.NewRacesRepo(racingDB)
	_ = racesRepo.Init(false)
	err = racesRepo.Clear()
	if err != nil {
		t.Fatalf("Unable to clear test database %v", err)
	}

	races := GetRaces()
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

func TestFetchAllEmpty(t *testing.T) {
	racingDB, err := GetTestDB("races", "TestFetchAllEmpty")
	if err != nil {
		t.Fatalf("Failed to open testdb %v", err)
	}
	racesRepo := db.NewRacesRepo(racingDB)
	_ = racesRepo.Init(false)

	races, err := racesRepo.ListAll()
	if err != nil {
		t.Fatalf("Failed to fetch all races %v.", err)
	}
	if len(races) > 0 {
		t.Fatalf("List all request returned invalid dataset.")
	}
}

func TestFetchAll(t *testing.T) {
	racingDB, err := GetTestDB("races", "TestPopulateAndFetchRepo")
	if err != nil {
		t.Fatalf("Failed to open testdb %v", err)
	}
	racesRepo := db.NewRacesRepo(racingDB)
	_ = racesRepo.Init(false)

	races, _ := racesRepo.ListAll()
	if len(races) > 0 {
		t.Fatalf("List all request returned invalid dataset.")
	}
	races = GetRaces()
	err = racesRepo.InsertRace(races[0])
	if err != nil {
		t.Fatalf("Failed to insert race record.")
	}

	outRaces, err := racesRepo.ListAll()
	if err != nil {
		t.Fatalf("Failed to list all list races %v", err)
	}
	if len(races) > 0 {
		if outRaces[0].Id != races[0].Id {
			t.Fatalf("Invalid race ID returned. Got %v, expected %v", outRaces[0].Id, races[0].Id)
		}
	} else {
		t.Fatalf("Failed to fetch inserted race.")
	}
}

// Helpers //

func GetRaces() []*racing.Race {
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
