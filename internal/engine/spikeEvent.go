package engine

type SpikeEvent struct {
	TargetNeuronID uint32
	Delta          float32
}

type TimedSpikeEvent struct {
	Tick  uint64
	Event SpikeEvent
}
