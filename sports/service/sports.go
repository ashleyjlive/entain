package service

import (
	"github.com/ashleyjlive/entain/sports/db"
	"github.com/ashleyjlive/entain/sports/proto/sports"
	"golang.org/x/net/context"
)

type Sports interface {
	// ListEvents will return a collection of events.
	ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error)
}

// sportsService implements the events interface.
type sportsService struct {
	eventsRepo db.EventsRepo
}

// NewSportsService instantiates and returns a new sportsService.
func NewSportsService(eventsRepo db.EventsRepo) Sports {
	return &sportsService{eventsRepo}
}

// Entry point for a ListEvents request.
func (s *sportsService) ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error) {
	evts, err := s.eventsRepo.List(in)
	if err != nil {
		return nil, err
	}

	return &sports.ListEventsResponse{Events: evts}, nil
}
