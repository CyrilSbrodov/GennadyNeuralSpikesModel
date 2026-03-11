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

	RestCharge      float32
	BaseThreshold   float32
	ThresholdNoise  float32
	ResetCharge     float32
	LeakFactor      float32
	CooldownTicks   uint16
	AdaptationDecay float32
	AdaptationStep  float32

	ExcitatoryWeightMin float32
	ExcitatoryWeightMax float32
	InhibitoryWeightMin float32
	InhibitoryWeightMax float32
	WeightMin           float32
	WeightMax           float32

	HebbianEnable             bool
	HebbianLearningRate       float32
	HebbianDecay              float32
	STDPWindowTicks           uint64
	STDPPotentiation          float32
	STDPDepression            float32
	SynapseUsageDecayInterval uint64
	SynapseUsageWeightDecay   float32
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

		RestCharge:      -70,
		BaseThreshold:   -60,
		ThresholdNoise:  2,
		ResetCharge:     -75,
		LeakFactor:      0.98,
		CooldownTicks:   2,
		AdaptationDecay: 0.92,
		AdaptationStep:  0.8,

		ExcitatoryWeightMin: 1.5,
		ExcitatoryWeightMax: 4.0,
		InhibitoryWeightMin: -4.0,
		InhibitoryWeightMax: -1.5,
		WeightMin:           -5.0,
		WeightMax:           5.0,

		HebbianEnable:             true,
		HebbianLearningRate:       0.03,
		HebbianDecay:              0.005,
		STDPWindowTicks:           8,
		STDPPotentiation:          1.0,
		STDPDepression:            0.4,
		SynapseUsageDecayInterval: 16,
		SynapseUsageWeightDecay:   0.02,
	}
}
