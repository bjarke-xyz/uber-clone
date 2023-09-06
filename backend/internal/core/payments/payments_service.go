package payments

type PaymentsService struct {
}

func NewService() *PaymentsService {
	return &PaymentsService{}
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
