package engine

import (
	"errors"
	"fmt"
	"math"
	"math/rand"

	"gennady-neural-spikes-model/cmd/config"
)

type Engine struct {
	Neurons  []Neuron
	Synapses []Synapse
	Config   *config.Config

	Tick uint64

	pending [][]SpikeEvent
	fired   []uint32
}

func NewEngine(cfg *config.Config) *Engine {
	if err := validateConfig(cfg); err != nil {
		panic(err)
	}

	pendingSize := int(cfg.MaxDelay) + 1

	return &Engine{
		Config:  cfg,
		pending: make([][]SpikeEvent, pendingSize),
		fired:   make([]uint32, 0, 128),
	}
}

func validateConfig(cfg *config.Config) error {
	if cfg == nil {
		return errors.New("config is nil")
	}
	if cfg.NeuronLimit <= 0 {
		return errors.New("NeuronLimit must be > 0")
	}
	if cfg.SynapseLimit <= 0 {
		return errors.New("SynapseLimit must be > 0")
	}
	if cfg.SynapsesPerNeuron <= 0 {
		return errors.New("SynapsesPerNeuron must be > 0")
	}
	if cfg.WorldSizeX <= 0 || cfg.WorldSizeY <= 0 || cfg.WorldSizeZ <= 0 {
		return errors.New("world size must be > 0")
	}
	if cfg.MaxAxonLength <= 0 {
		return errors.New("MaxAxonLength must be > 0")
	}
	if cfg.LocalConnectRadius <= 0 {
		return errors.New("LocalConnectRadius must be > 0")
	}
	if cfg.LocalConnectRadius > cfg.MaxAxonLength {
		return errors.New("LocalConnectRadius cannot exceed MaxAxonLength")
	}
	if cfg.MinDelay == 0 {
		return errors.New("MinDelay must be >= 1")
	}
	if cfg.MaxDelay < cfg.MinDelay {
		return errors.New("MaxDelay must be >= MinDelay")
	}
	if cfg.ExcitatoryRatio < 0 || cfg.InhibitoryRatio < 0 {
		return errors.New("polarity ratios cannot be negative")
	}
	if math.Abs(float64(cfg.ExcitatoryRatio+cfg.InhibitoryRatio-1.0)) > 0.0001 {
		return fmt.Errorf("ExcitatoryRatio + InhibitoryRatio must be 1.0")
	}
	if cfg.InputRatio < 0 || cfg.OutputRatio < 0 || cfg.InputRatio+cfg.OutputRatio >= 1.0 {
		return errors.New("invalid input/output ratios")
	}
	if cfg.LongConnectionProb < 0 || cfg.LongConnectionProb > 1 {
		return errors.New("LongConnectionProb must be in [0,1]")
	}
	if cfg.LeakFactor < 0 || cfg.LeakFactor > 1 {
		return errors.New("LeakFactor must be in [0,1]")
	}
	if cfg.ExcitatoryWeightMin > cfg.ExcitatoryWeightMax {
		return errors.New("ExcitatoryWeightMin must be <= ExcitatoryWeightMax")
	}
	if cfg.InhibitoryWeightMin > cfg.InhibitoryWeightMax {
		return errors.New("InhibitoryWeightMin must be <= InhibitoryWeightMax")
	}
	if cfg.InhibitoryWeightMax > 0 {
		return errors.New("InhibitoryWeightMax must be <= 0")
	}
	if cfg.ExcitatoryWeightMin < 0 {
		return errors.New("ExcitatoryWeightMin must be >= 0")
	}
	if cfg.WeightMin > cfg.WeightMax {
		return errors.New("WeightMin must be <= WeightMax")
	}
	if cfg.STDPWindowTicks == 0 {
		return errors.New("STDPWindowTicks must be > 0")
	}
	if cfg.SynapseUsageDecayInterval == 0 {
		return errors.New("SynapseUsageDecayInterval must be > 0")
	}
	return nil
}

func (e *Engine) InitSpatial3D() {
	e.Neurons = nil
	e.Synapses = nil
	e.Tick = 0

	for i := range e.pending {
		e.pending[i] = e.pending[i][:0]
	}

	e.initNeurons3D()
	e.initSynapsesSpatial()
}

func (e *Engine) initNeurons3D() {
	rng := rand.New(rand.NewSource(e.Config.Seed))
	e.Neurons = make([]Neuron, e.Config.NeuronLimit)

	inputCount := int(float32(e.Config.NeuronLimit) * e.Config.InputRatio)
	outputCount := int(float32(e.Config.NeuronLimit) * e.Config.OutputRatio)

	for i := 0; i < e.Config.NeuronLimit; i++ {
		noiseAmp := e.Config.ThresholdNoise
		thresholdNoise := rng.Float32()*(noiseAmp*2) - noiseAmp

		role := RoleHidden
		switch {
		case i < inputCount:
			role = RoleInput
		case i >= e.Config.NeuronLimit-outputCount:
			role = RoleOutput
		}

		polarity := choosePolarity(rng, e.Config.ExcitatoryRatio)

		e.Neurons[i] = Neuron{
			ID: uint32(i),

			Polarity: polarity,
			Role:     role,

			Position: Vec3{
				X: rng.Float32() * e.Config.WorldSizeX,
				Y: rng.Float32() * e.Config.WorldSizeY,
				Z: rng.Float32() * e.Config.WorldSizeZ,
			},

			Charge:          e.Config.RestCharge,
			RestCharge:      e.Config.RestCharge,
			BaseThreshold:   e.Config.BaseThreshold + thresholdNoise,
			ResetCharge:     e.Config.ResetCharge,
			Adaptation:      0,
			AdaptationDecay: e.Config.AdaptationDecay,
			AdaptationStep:  e.Config.AdaptationStep,
			LastSpikeTick:   -1,
			FireCount:       0,

			CooldownTicks: e.Config.CooldownTicks,
			CooldownLeft:  0,

			Outgoing:      nil,
			Incoming:      nil,
			FiredLastTick: false,
		}
	}
}

func choosePolarity(rng *rand.Rand, excitatoryRatio float32) NeuronPolarity {
	if rng.Float32() < excitatoryRatio {
		return PolarityExcitatory
	}
	return PolarityInhibitory
}

func (e *Engine) initSynapsesSpatial() {
	rng := rand.New(rand.NewSource(e.Config.Seed + 1))
	n := len(e.Neurons)
	if n == 0 {
		return
	}

	used := make(map[uint64]struct{})
	maxAxonLenSq := e.Config.MaxAxonLength * e.Config.MaxAxonLength
	localRadiusSq := e.Config.LocalConnectRadius * e.Config.LocalConnectRadius

	for from := 0; from < n; from++ {
		if len(e.Synapses) >= e.Config.SynapseLimit {
			return
		}

		created := 0
		attempts := 0
		maxAttempts := e.Config.SynapsesPerNeuron * 60

		for created < e.Config.SynapsesPerNeuron && attempts < maxAttempts {
			attempts++

			to := rng.Intn(n)
			if to == from {
				continue
			}

			key := uint64(from)<<32 | uint64(to)
			if _, exists := used[key]; exists {
				continue
			}

			distSq := distanceSq(e.Neurons[from].Position, e.Neurons[to].Position)
			if distSq > maxAxonLenSq {
				continue
			}

			if !shouldConnect(rng, distSq, localRadiusSq, maxAxonLenSq, e.Config.LongConnectionProb) {
				continue
			}

			if len(e.Synapses) >= e.Config.SynapseLimit {
				return
			}

			used[key] = struct{}{}

			dist := float32(math.Sqrt(float64(distSq)))
			delay := computeDelay(dist, e.Config.MaxAxonLength, e.Config.MinDelay, e.Config.MaxDelay)
			weight := e.computeWeight(e.Neurons[from].Polarity, rng)

			s := Synapse{
				ID:      uint32(len(e.Synapses)),
				FromID:  uint32(from),
				ToID:    uint32(to),
				Weight:  weight,
				Delay:   delay,
				Enabled: true,
			}

			e.Synapses = append(e.Synapses, s)
			e.Neurons[from].Outgoing = append(e.Neurons[from].Outgoing, s.ID)
			e.Neurons[to].Incoming = append(e.Neurons[to].Incoming, s.ID)
			created++
		}
	}
}

func shouldConnect(
	rng *rand.Rand,
	distSq float32,
	localRadiusSq float32,
	maxAxonLenSq float32,
	longProb float32,
) bool {
	if distSq <= localRadiusSq {
		return true
	}

	// плавное падение вероятности по расстоянию
	distNorm := distSq / maxAxonLenSq
	p := longProb * (1.0 - distNorm)
	if p < 0 {
		p = 0
	}
	return rng.Float32() < p
}

func (e *Engine) computeWeight(polarity NeuronPolarity, rng *rand.Rand) float32 {
	if polarity == PolarityInhibitory {
		return e.Config.InhibitoryWeightMin + rng.Float32()*(e.Config.InhibitoryWeightMax-e.Config.InhibitoryWeightMin)
	}
	return e.Config.ExcitatoryWeightMin + rng.Float32()*(e.Config.ExcitatoryWeightMax-e.Config.ExcitatoryWeightMin)
}

func computeDelay(dist, maxDist float32, minDelay, maxDelay uint16) uint16 {
	if maxDelay <= minDelay || maxDist <= 0 {
		return minDelay
	}

	norm := dist / maxDist
	if norm < 0 {
		norm = 0
	}
	if norm > 1 {
		norm = 1
	}

	span := float32(maxDelay - minDelay)
	return minDelay + uint16(norm*span)
}

func distanceSq(a, b Vec3) float32 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return dx*dx + dy*dy + dz*dz
}

func (e *Engine) InjectNow(targetNeuronID uint32, delta float32) error {
	if int(targetNeuronID) >= len(e.Neurons) {
		return errors.New("targetNeuronID is out of range")
	}

	slot := e.currentSlot()
	e.pending[slot] = append(e.pending[slot], SpikeEvent{
		TargetNeuronID: targetNeuronID,
		Delta:          delta,
	})

	return nil
}

func (e *Engine) InjectAfter(delay uint16, targetNeuronID uint32, delta float32) error {
	if int(targetNeuronID) >= len(e.Neurons) {
		return errors.New("targetNeuronID is out of range")
	}
	slot := e.slotForDelay(delay)
	e.pending[slot] = append(e.pending[slot], SpikeEvent{
		TargetNeuronID: targetNeuronID,
		Delta:          delta,
	})
	return nil
}

func (e *Engine) Run(steps int) []TickStats {
	stats := make([]TickStats, 0, steps)
	for i := 0; i < steps; i++ {
		stats = append(stats, e.Step())
	}
	return stats
}

func (e *Engine) Step() TickStats {
	stats := TickStats{
		Tick: e.Tick,
	}

	for i := range e.Neurons {
		e.Neurons[i].FiredLastTick = false
	}

	// phase 1: deliver pending events for current tick
	slot := e.currentSlot()
	events := e.pending[slot]
	stats.DeliveredEvents = len(events)

	for _, event := range events {
		if int(event.TargetNeuronID) >= len(e.Neurons) {
			continue
		}

		n := &e.Neurons[event.TargetNeuronID]
		if n.CooldownLeft > 0 {
			continue
		}

		n.Charge += event.Delta
	}

	e.pending[slot] = e.pending[slot][:0]

	// phase 2: update neurons / leak / threshold
	e.fired = e.fired[:0]
	var totalCharge float32

	for i := range e.Neurons {
		n := &e.Neurons[i]

		n.Adaptation *= n.AdaptationDecay

		if n.CooldownLeft > 0 {
			n.CooldownLeft--
			totalCharge += n.Charge
			continue
		}

		n.Charge = n.RestCharge + (n.Charge-n.RestCharge)*e.Config.LeakFactor

		effectiveThreshold := n.BaseThreshold + n.Adaptation
		if n.Charge >= effectiveThreshold {
			n.FiredLastTick = true
			e.fired = append(e.fired, n.ID)
			stats.FiredCount++
		} else {
			gap := effectiveThreshold - n.Charge
			if gap > 0 && gap <= 1.0 {
				stats.NearThresholdCount++
			}
		}

		totalCharge += n.Charge
	}

	if len(e.Neurons) > 0 {
		stats.MeanCharge = totalCharge / float32(len(e.Neurons))
	}

	// phase 3: emit spikes from fired neurons
	for _, neuronID := range e.fired {
		n := &e.Neurons[neuronID]

		n.Adaptation += n.AdaptationStep
		n.LastSpikeTick = int64(e.Tick)
		n.FireCount++

		n.Charge = n.ResetCharge
		n.CooldownLeft = n.CooldownTicks

		for _, synapseID := range n.Outgoing {
			if int(synapseID) >= len(e.Synapses) {
				continue
			}

			s := &e.Synapses[synapseID]
			if !s.Enabled {
				continue
			}

			s.UseCount++
			s.LastUsedTick = e.Tick

			delay := s.Delay
			if delay == 0 {
				delay = 1
			}

			futureSlot := e.slotForDelay(delay)
			e.pending[futureSlot] = append(e.pending[futureSlot], SpikeEvent{
				TargetNeuronID: s.ToID,
				Delta:          s.Weight,
			})
			stats.CreatedEvents++
		}
	}

	stats.UpdatedWeights += e.applySTDPPlasticity()
	stats.UpdatedWeights += e.applySynapseDecay()
	e.fillWeightStats(&stats)

	e.Tick++
	return stats
}

func (e *Engine) currentSlot() int {
	return int(e.Tick % uint64(len(e.pending)))
}

func (e *Engine) slotForDelay(delay uint16) int {
	if int(delay) == 0 {
		delay = 1
	}
	return int((e.Tick + uint64(delay)) % uint64(len(e.pending)))
}

func (e *Engine) PrintDelayHistogram() {
	if len(e.Synapses) == 0 {
		fmt.Println("Delay histogram: no synapses")
		return
	}

	hist := make(map[uint16]int)
	for _, s := range e.Synapses {
		hist[s.Delay]++
	}

	fmt.Println("=== Delay Histogram ===")
	for d := e.Config.MinDelay; d <= e.Config.MaxDelay; d++ {
		fmt.Printf("delay=%d count=%d\n", d, hist[d])
	}
	fmt.Println("=======================")
}

func (e *Engine) PrintOutgoingHistogram() {
	if len(e.Neurons) == 0 {
		fmt.Println("Outgoing histogram: no neurons")
		return
	}

	hist := make(map[int]int)
	for _, n := range e.Neurons {
		hist[len(n.Outgoing)]++
	}

	fmt.Println("=== Outgoing Histogram ===")
	for i := 0; i <= e.Config.SynapsesPerNeuron; i++ {
		fmt.Printf("outgoing=%d count=%d\n", i, hist[i])
	}
	fmt.Println("==========================")
}

func (e *Engine) ValidateTopology() error {
	for i, n := range e.Neurons {
		if n.ID != uint32(i) {
			return fmt.Errorf("neuron id mismatch: index=%d id=%d", i, n.ID)
		}

		for _, synID := range n.Outgoing {
			if int(synID) >= len(e.Synapses) {
				return fmt.Errorf("neuron %d has invalid outgoing synapse id %d", n.ID, synID)
			}

			s := e.Synapses[synID]
			if s.FromID != n.ID {
				return fmt.Errorf(
					"outgoing mismatch: neuron=%d synapse=%d synapse.FromID=%d",
					n.ID, s.ID, s.FromID,
				)
			}
		}

		for _, synID := range n.Incoming {
			if int(synID) >= len(e.Synapses) {
				return fmt.Errorf("neuron %d has invalid incoming synapse id %d", n.ID, synID)
			}
			s := e.Synapses[synID]
			if s.ToID != n.ID {
				return fmt.Errorf(
					"incoming mismatch: neuron=%d synapse=%d synapse.ToID=%d",
					n.ID, s.ID, s.ToID,
				)
			}
		}
	}

	for i, s := range e.Synapses {
		if s.ID != uint32(i) {
			return fmt.Errorf("synapse id mismatch: index=%d id=%d", i, s.ID)
		}
		if int(s.FromID) >= len(e.Neurons) {
			return fmt.Errorf("synapse %d has invalid FromID=%d", s.ID, s.FromID)
		}
		if int(s.ToID) >= len(e.Neurons) {
			return fmt.Errorf("synapse %d has invalid ToID=%d", s.ID, s.ToID)
		}
	}

	return nil
}
