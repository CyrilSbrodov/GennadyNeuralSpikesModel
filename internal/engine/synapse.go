package engine

type Synapse struct {
	ID      uint32
	FromID  uint32
	ToID    uint32
	Weight  float32
	Delay   uint16
	Enabled bool
}
