package config

type Config struct {
	NeuronLimit       int
	SynapseLimit      int
	SynapsesPerNeuron int
	Seed              int64

	WorldSizeX float32
	WorldSizeY float32
	WorldSizeZ float32

	MaxAxonLength      float32
	LocalConnectRadius float32

	LongConnectionProb float32

	ExcitatoryRatio float32
	InhibitoryRatio float32

	InputRatio  float32
	OutputRatio float32

	MinDelay uint16
	MaxDelay uint16

	NeuronCharge      float32
	NeuronRestCharge  float32
	NeuronThreshold   float32
	NeuronResetCharge float32
}

func DefaultConfig() *Config {
	return &Config{
		NeuronLimit:       1000,
		SynapseLimit:      8000,
		SynapsesPerNeuron: 8,
		Seed:              0,

		WorldSizeX: 25,
		WorldSizeY: 25,
		WorldSizeZ: 25,

		MaxAxonLength:      10,
		LocalConnectRadius: 4,

		LongConnectionProb: 0.10,

		ExcitatoryRatio: 0.8,
		InhibitoryRatio: 0.2,

		InputRatio:  0.03,
		OutputRatio: 0.03,

		MinDelay: 1,
		MaxDelay: 6,

		NeuronCharge:      -70,
		NeuronRestCharge:  -70,
		NeuronThreshold:   -60,
		NeuronResetCharge: -75,
	}
}
