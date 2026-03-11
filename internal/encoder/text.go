package encoder

import (
	"strings"

	"gennady-neural-spikes-model/internal/engine"
)

func EncodeTextToSpikes(text string, t0 uint64, cfg EncoderConfig) []engine.TimedSpikeEvent {
	if err := cfg.Validate(); err != nil {
		return nil
	}

	tokens := strings.Fields(text)
	if len(tokens) == 0 {
		return nil
	}

	events := make([]engine.TimedSpikeEvent, 0, len(tokens)*cfg.NeuronsPerInput)
	for tokenIdx, token := range tokens {
		windowStart := t0 + uint64(tokenIdx)*cfg.UnitWindowTicks
		channel := hashToken(token) % cfg.InputChannels
		intensity := tokenIntensity(token)
		events = append(events, buildTemporalSpikes(windowStart, channel, intensity, cfg)...)
	}

	return events
}

func buildTemporalSpikes(windowStart uint64, channel int, intensity float32, cfg EncoderConfig) []engine.TimedSpikeEvent {
	if intensity <= 0 {
		return nil
	}

	delta := cfg.IntensityScale * intensity
	out := make([]engine.TimedSpikeEvent, 0, cfg.SpikesPerWindow*cfg.NeuronsPerInput)

	for ni := 0; ni < cfg.NeuronsPerInput; ni++ {
		target := cfg.neuronForChannel(channel, ni)

		switch cfg.Coding {
		case CodingLatency:
			maxSteps := maxInt(1, int(cfg.UnitWindowTicks/cfg.TimeStepTicks)-1)
			latencyStep := int((1.0 - intensity) * float32(maxSteps))
			when := windowStart + uint64(latencyStep)*cfg.TimeStepTicks
			out = append(out, engine.TimedSpikeEvent{Tick: when, Event: engine.SpikeEvent{TargetNeuronID: target, Delta: delta}})
		default:
			count := int(float32(cfg.SpikesPerWindow)*intensity) + 1
			for si := 0; si < count; si++ {
				offset := uint64(si) * cfg.TimeStepTicks
				if offset >= cfg.UnitWindowTicks {
					break
				}
				out = append(out, engine.TimedSpikeEvent{Tick: windowStart + offset, Event: engine.SpikeEvent{TargetNeuronID: target, Delta: delta}})
			}
		}
	}

	return out
}

func hashToken(token string) int {
	h := 2166136261
	for _, r := range token {
		h ^= int(r)
		h *= 16777619
	}
	if h < 0 {
		return -h
	}
	return h
}

func tokenIntensity(token string) float32 {
	if token == "" {
		return 0
	}
	letters := 0
	for _, r := range token {
		if r >= '0' && r <= '9' {
			letters += 2
			continue
		}
		if r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r >= 'А' && r <= 'я' {
			letters++
		}
	}
	intensity := float32(letters) / float32(maxInt(1, len([]rune(token))))
	if intensity < 0.1 {
		return 0.1
	}
	if intensity > 1 {
		return 1
	}
	return intensity
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
