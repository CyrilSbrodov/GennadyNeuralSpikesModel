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

	Charge      float32
	RestCharge  float32
	Threshold   float32
	ResetCharge float32

	CooldownTicks uint16
	CooldownLeft  uint16

	Outgoing []uint32

	FiredLastTick bool
}
