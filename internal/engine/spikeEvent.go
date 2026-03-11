package engine

type SpikeEvent struct {
	TargetNeuronID uint32
	Delta          float32
}

func (se *SpikeEvent) buildBuffer() [][]SpikeEvent {
	return [][]SpikeEvent{}
}
