package location

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/lithammer/shortuuid/v4"

	"github.com/tbe-team/raybot/internal/model"
	"github.com/tbe-team/raybot/internal/pubsub"
	"github.com/tbe-team/raybot/internal/repository"
	"github.com/tbe-team/raybot/internal/service"
	"github.com/tbe-team/raybot/internal/storage/db"
)

type Service struct {
	locationRepo repository.LocationRepository
	publisher    message.Publisher
	dbProvider   db.Provider
}

func NewService(locationRepo repository.LocationRepository, publisher message.Publisher, dbProvider db.Provider) *Service {
	return &Service{locationRepo: locationRepo, publisher: publisher, dbProvider: dbProvider}
}

func (s Service) UpdateLocation(ctx context.Context, params service.UpdateLocationParams) error {
	if err := s.locationRepo.UpdateLocation(ctx, s.dbProvider.DB(), model.Location{
		CurrentLocation: params.CurrentLocation,
		UpdatedAt:       time.Now(),
	}); err != nil {
		return fmt.Errorf("update location: %w", err)
	}

	if err := s.publishLocationUpdatedEvent(params.CurrentLocation); err != nil {
		return fmt.Errorf("publish location updated event: %w", err)
	}

	return nil
}

func (s Service) publishLocationUpdatedEvent(location string) error {
	ev := pubsub.RobotLocationUpdatedEvent{
		Location: location,
	}
	payload, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("json marshal location updated event: %w", err)
	}

	msg := message.NewMessage(shortuuid.New(), payload)
	if err := s.publisher.Publish(pubsub.TopicRobotLocationUpdated, msg); err != nil {
		return fmt.Errorf("publisher publish location updated event: %w", err)
	}

	return nil
}
