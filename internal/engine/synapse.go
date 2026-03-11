package engine

type Synapse struct {
	ID      uint32
	FromID  uint32
	ToID    uint32
	Weight  float32
	Delay   uint16
	Enabled bool

	UseCount        uint64
	LastUsedTick    uint64
	PlasticityScore float32
}
