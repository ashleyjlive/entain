package service

import (
	"github.com/ashleyjlive/entain/racing/db"
	"github.com/ashleyjlive/entain/racing/proto/racing"
	"golang.org/x/net/context"
)

type Racing interface {
	// ListRaces will return a collection of races.
	ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error)
	GetRace(ctx context.Context, in *racing.GetRaceRequest) (*racing.Race, error)
}

// racingService implements the Racing interface.
type racingService struct {
	racesRepo db.RacesRepo
}

// NewRacingService instantiates and returns a new racingService.
func NewRacingService(racesRepo db.RacesRepo) Racing {
	return &racingService{racesRepo}
}

func (s *racingService) ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error) {
	races, err := s.racesRepo.List(in)
	if err != nil {
		return nil, err
	}

	return &racing.ListRacesResponse{Races: races}, nil
}

// Retrieves a single race by its identifier, or, returns an error.
func (s *racingService) GetRace(ctx context.Context, in *racing.GetRaceRequest) (*racing.Race, error) {
	race, err := s.racesRepo.Get(in)
	if err != nil {
		return nil, err
	}
	return race, nil
}
