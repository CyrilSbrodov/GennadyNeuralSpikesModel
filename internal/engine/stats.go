package engine

type TickStats struct {
	Tick               uint64
	DeliveredEvents    int
	CreatedEvents      int
	FiredCount         int
	MeanCharge         float32
	NearThresholdCount int
	UpdatedWeights     int
	MeanWeight         float32
	MeanAbsWeight      float32
}
