package payments

import (
	"log/slog"

	"github.com/bjarke-xyz/uber-clone-backend/internal/core"
)

type PaymentsService struct {
	pubsub core.Pubsub
	logger *slog.Logger
}

func NewService(pubsub core.Pubsub, logger *slog.Logger) *PaymentsService {
	return &PaymentsService{
		pubsub: pubsub,
		logger: logger,
	}
}

func (s *PaymentsService) CalculatePrice(distanceInMeters int) int {
	return calculatePrice(distanceInMeters)
}

func (s *PaymentsService) GetCurrencies() []Currency {
	return []Currency{
		{
			Symbol: "EUR",
			Icon:   "€",
		},
	}
}
