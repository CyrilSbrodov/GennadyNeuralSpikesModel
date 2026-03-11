package engine

type NeuronPolarity uint8

const (
	PolarityExcitatory NeuronPolarity = iota
	PolarityInhibitory
)

type NeuronRole uint8

const (
	RoleInput NeuronRole = iota
	RoleHidden
	RoleOutput
)

type Vec3 struct {
	X float32
	Y float32
	Z float32
}

type Neuron struct {
	ID uint32

	Polarity NeuronPolarity
	Role     NeuronRole

	Position Vec3

	Charge          float32
	RestCharge      float32
	BaseThreshold   float32
	ResetCharge     float32
	Adaptation      float32
	AdaptationDecay float32
	AdaptationStep  float32
	LastSpikeTick   int64
	FireCount       uint64

	CooldownTicks uint16
	CooldownLeft  uint16

	Outgoing []uint32
	Incoming []uint32

	FiredLastTick bool
}
