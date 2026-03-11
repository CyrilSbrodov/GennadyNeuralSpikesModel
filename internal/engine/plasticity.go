package engine

import "math"

func (e *Engine) applySTDPPlasticity() int {
	if !e.Config.HebbianEnable {
		return 0
	}

	updated := 0
	window := int64(e.Config.STDPWindowTicks)

	for _, postID := range e.fired {
		post := &e.Neurons[postID]
		for _, synID := range post.Incoming {
			if int(synID) >= len(e.Synapses) {
				continue
			}

			s := &e.Synapses[synID]
			if !s.Enabled {
				continue
			}

			pre := e.Neurons[s.FromID]
			delta := e.Config.HebbianLearningRate * e.Config.HebbianDecay
			if pre.LastSpikeTick >= 0 {
				dt := int64(e.Tick) - pre.LastSpikeTick
				if dt >= 0 && dt <= window {
					factor := 1 - float32(dt)/float32(e.Config.STDPWindowTicks)
					delta = e.Config.HebbianLearningRate * e.Config.STDPPotentiation * factor
				} else {
					delta = -e.Config.HebbianLearningRate * e.Config.STDPDepression
				}
			}

			if delta == 0 {
				continue
			}

			s.Weight = e.clampSynapseWeight(s.Weight+delta, pre.Polarity)
			s.PlasticityScore += delta
			updated++
		}
	}

	return updated
}

func (e *Engine) applySynapseDecay() int {
	if e.Config.SynapseUsageDecayInterval == 0 || e.Tick == 0 || e.Tick%e.Config.SynapseUsageDecayInterval != 0 {
		return 0
	}

	updated := 0
	decay := e.Config.SynapseUsageWeightDecay

	for i := range e.Synapses {
		s := &e.Synapses[i]
		if !s.Enabled {
			continue
		}

		if s.UseCount > 0 && e.Tick-s.LastUsedTick < e.Config.SynapseUsageDecayInterval {
			continue
		}

		newWeight := s.Weight * (1.0 - decay)
		if s.Weight > 0 && newWeight < 0 {
			newWeight = 0
		}
		if s.Weight < 0 && newWeight > 0 {
			newWeight = 0
		}

		polarity := e.Neurons[s.FromID].Polarity
		s.Weight = e.clampSynapseWeight(newWeight, polarity)
		updated++
	}

	return updated
}

func (e *Engine) clampSynapseWeight(weight float32, polarity NeuronPolarity) float32 {
	if weight < e.Config.WeightMin {
		weight = e.Config.WeightMin
	}
	if weight > e.Config.WeightMax {
		weight = e.Config.WeightMax
	}

	if polarity == PolarityExcitatory {
		if weight < e.Config.ExcitatoryWeightMin {
			weight = e.Config.ExcitatoryWeightMin
		}
		if weight > e.Config.ExcitatoryWeightMax {
			weight = e.Config.ExcitatoryWeightMax
		}
		return weight
	}

	if weight < e.Config.InhibitoryWeightMin {
		weight = e.Config.InhibitoryWeightMin
	}
	if weight > e.Config.InhibitoryWeightMax {
		weight = e.Config.InhibitoryWeightMax
	}

	return float32(math.Min(float64(weight), 0))
}

func (e *Engine) fillWeightStats(stats *TickStats) {
	if len(e.Synapses) == 0 {
		return
	}

	var total float32
	var totalAbs float32
	for _, s := range e.Synapses {
		total += s.Weight
		if s.Weight < 0 {
			totalAbs -= s.Weight
		} else {
			totalAbs += s.Weight
		}
	}

	stats.MeanWeight = total / float32(len(e.Synapses))
	stats.MeanAbsWeight = totalAbs / float32(len(e.Synapses))
}
