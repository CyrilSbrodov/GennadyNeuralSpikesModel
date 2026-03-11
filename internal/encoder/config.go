package encoder

import (
	"fmt"
)

type TemporalCoding string

const (
	CodingRate    TemporalCoding = "rate"
	CodingLatency TemporalCoding = "latency"
)

type EncoderConfig struct {
	TimeStepTicks   uint64
	UnitWindowTicks uint64
	IntensityScale  float32
	InputChannels   int
	NeuronsPerInput int
	SpikesPerWindow int
	Coding          TemporalCoding

	InputNeuronIDs []uint32

	ImageWidth  int
	ImageHeight int
	UseRGB      bool
}

func (cfg EncoderConfig) Validate() error {
	if cfg.TimeStepTicks == 0 {
		return fmt.Errorf("TimeStepTicks must be > 0")
	}
	if cfg.UnitWindowTicks == 0 {
		return fmt.Errorf("UnitWindowTicks must be > 0")
	}
	if cfg.IntensityScale <= 0 {
		return fmt.Errorf("IntensityScale must be > 0")
	}
	if cfg.InputChannels <= 0 {
		return fmt.Errorf("InputChannels must be > 0")
	}
	if cfg.NeuronsPerInput <= 0 {
		return fmt.Errorf("NeuronsPerInput must be > 0")
	}
	if cfg.SpikesPerWindow <= 0 {
		return fmt.Errorf("SpikesPerWindow must be > 0")
	}
	if len(cfg.InputNeuronIDs) == 0 {
		return fmt.Errorf("InputNeuronIDs cannot be empty")
	}
	if cfg.Coding != CodingRate && cfg.Coding != CodingLatency {
		return fmt.Errorf("unsupported coding mode: %s", cfg.Coding)
	}

	neededNeurons := cfg.InputChannels * cfg.NeuronsPerInput
	if neededNeurons > len(cfg.InputNeuronIDs) {
		return fmt.Errorf(
			"input mapping needs %d neurons, but only %d RoleInput neurons are available",
			neededNeurons,
			len(cfg.InputNeuronIDs),
		)
	}
	if cfg.ImageWidth <= 0 || cfg.ImageHeight <= 0 {
		return fmt.Errorf("ImageWidth and ImageHeight must be > 0")
	}

	return nil
}

func (cfg EncoderConfig) neuronForChannel(channel, offset int) uint32 {
	idx := channel*cfg.NeuronsPerInput + offset
	return cfg.InputNeuronIDs[idx]
}
