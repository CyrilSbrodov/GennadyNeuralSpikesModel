package engine

type TickStats struct {
	Tick            uint64
	DeliveredEvents int
	CreatedEvents   int
	FiredCount      int
	MeanCharge      float32
}
