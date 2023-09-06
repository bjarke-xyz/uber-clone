package payments

type Currency struct {
	Symbol string `json:"symbol"`
	Icon   string `json:"icon"`
}

func calculatePrice(distanceInMeters int) int {
	kilometers := distanceInMeters / 1000
	return (700 + (140 * kilometers))
}
