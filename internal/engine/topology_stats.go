package engine

import (
	"fmt"
	"math"
)

type TopologyStats struct {
	NeuronCount  int
	SynapseCount int

	InputCount  int
	HiddenCount int
	OutputCount int

	ExcitatoryCount int
	InhibitoryCount int

	MinOutgoing  int
	MaxOutgoing  int
	MeanOutgoing float64

	PositiveSynapses int
	NegativeSynapses int
	MeanWeight       float64
	MeanAbsWeight    float64

	MinDelay  uint16
	MaxDelay  uint16
	MeanDelay float64

	MinDistance  float64
	MaxDistance  float64
	MeanDistance float64
}

func (e *Engine) CollectTopologyStats() TopologyStats {
	stats := TopologyStats{
		NeuronCount:  len(e.Neurons),
		SynapseCount: len(e.Synapses),
	}

	if len(e.Neurons) == 0 {
		return stats
	}

	// ---- neuron role / polarity stats ----
	for _, n := range e.Neurons {
		switch n.Role {
		case RoleInput:
			stats.InputCount++
		case RoleHidden:
			stats.HiddenCount++
		case RoleOutput:
			stats.OutputCount++
		}

		switch n.Polarity {
		case PolarityExcitatory:
			stats.ExcitatoryCount++
		case PolarityInhibitory:
			stats.InhibitoryCount++
		}
	}

	// ---- outgoing stats ----
	stats.MinOutgoing = math.MaxInt

	var outgoingSum int
	for _, n := range e.Neurons {
		c := len(n.Outgoing)
		outgoingSum += c

		if c < stats.MinOutgoing {
			stats.MinOutgoing = c
		}
		if c > stats.MaxOutgoing {
			stats.MaxOutgoing = c
		}
	}

	if len(e.Neurons) > 0 {
		stats.MeanOutgoing = float64(outgoingSum) / float64(len(e.Neurons))
	}

	if stats.MinOutgoing == math.MaxInt {
		stats.MinOutgoing = 0
	}

	// ---- synapse stats ----
	if len(e.Synapses) == 0 {
		return stats
	}

	stats.MinDelay = math.MaxUint16
	stats.MinDistance = math.MaxFloat64

	var delaySum uint64
	var weightSum float64
	var absWeightSum float64
	var distanceSum float64

	for _, s := range e.Synapses {
		if s.Weight >= 0 {
			stats.PositiveSynapses++
		} else {
			stats.NegativeSynapses++
		}

		weightSum += float64(s.Weight)
		absWeightSum += math.Abs(float64(s.Weight))

		if s.Delay < stats.MinDelay {
			stats.MinDelay = s.Delay
		}
		if s.Delay > stats.MaxDelay {
			stats.MaxDelay = s.Delay
		}
		delaySum += uint64(s.Delay)

		if int(s.FromID) >= len(e.Neurons) || int(s.ToID) >= len(e.Neurons) {
			continue
		}

		from := e.Neurons[s.FromID]
		to := e.Neurons[s.ToID]

		dist := distance(from.Position, to.Position)
		if dist < stats.MinDistance {
			stats.MinDistance = dist
		}
		if dist > stats.MaxDistance {
			stats.MaxDistance = dist
		}
		distanceSum += dist
	}

	stats.MeanWeight = weightSum / float64(len(e.Synapses))
	stats.MeanAbsWeight = absWeightSum / float64(len(e.Synapses))
	stats.MeanDelay = float64(delaySum) / float64(len(e.Synapses))
	stats.MeanDistance = distanceSum / float64(len(e.Synapses))

	if stats.MinDistance == math.MaxFloat64 {
		stats.MinDistance = 0
	}

	return stats
}

func (e *Engine) PrintTopologyStats() {
	s := e.CollectTopologyStats()

	fmt.Println("=== Topology Stats ===")
	fmt.Printf("Neurons: %d\n", s.NeuronCount)
	fmt.Printf("Synapses: %d\n", s.SynapseCount)
	fmt.Println()

	fmt.Println("Roles:")
	fmt.Printf("  Input:  %d\n", s.InputCount)
	fmt.Printf("  Hidden: %d\n", s.HiddenCount)
	fmt.Printf("  Output: %d\n", s.OutputCount)
	fmt.Println()

	fmt.Println("Polarities:")
	fmt.Printf("  Excitatory: %d\n", s.ExcitatoryCount)
	fmt.Printf("  Inhibitory: %d\n", s.InhibitoryCount)
	fmt.Println()

	fmt.Println("Outgoing per neuron:")
	fmt.Printf("  Min:  %d\n", s.MinOutgoing)
	fmt.Printf("  Max:  %d\n", s.MaxOutgoing)
	fmt.Printf("  Mean: %.2f\n", s.MeanOutgoing)
	fmt.Println()

	fmt.Println("Synapse weights:")
	fmt.Printf("  Positive:     %d\n", s.PositiveSynapses)
	fmt.Printf("  Negative:     %d\n", s.NegativeSynapses)
	fmt.Printf("  Mean weight:  %.4f\n", s.MeanWeight)
	fmt.Printf("  Mean |weight| %.4f\n", s.MeanAbsWeight)
	fmt.Println()

	fmt.Println("Delays:")
	fmt.Printf("  Min:  %d\n", s.MinDelay)
	fmt.Printf("  Max:  %d\n", s.MaxDelay)
	fmt.Printf("  Mean: %.2f\n", s.MeanDelay)
	fmt.Println()

	fmt.Println("Distances:")
	fmt.Printf("  Min:  %.4f\n", s.MinDistance)
	fmt.Printf("  Max:  %.4f\n", s.MaxDistance)
	fmt.Printf("  Mean: %.4f\n", s.MeanDistance)
	fmt.Println("======================")
}

func distance(a, b Vec3) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	dz := float64(a.Z - b.Z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
