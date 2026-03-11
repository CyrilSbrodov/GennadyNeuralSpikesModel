package engine

import (
	"fmt"
	"sort"
)

type neuronDistance struct {
	ID   uint32
	Dist float64
}

func (e *Engine) FindNearestNeurons(seedID uint32, count int) []uint32 {
	if int(seedID) >= len(e.Neurons) || count <= 0 {
		return nil
	}

	seed := e.Neurons[seedID]
	items := make([]neuronDistance, 0, len(e.Neurons)-1)

	for _, n := range e.Neurons {
		if n.ID == seedID {
			continue
		}

		items = append(items, neuronDistance{
			ID:   n.ID,
			Dist: distance(seed.Position, n.Position),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Dist < items[j].Dist
	})

	if count > len(items) {
		count = len(items)
	}

	out := make([]uint32, 0, count+1)
	out = append(out, seedID)

	for i := 0; i < count; i++ {
		out = append(out, items[i].ID)
	}

	return out
}

func (e *Engine) InjectCluster(seedID uint32, neighborsCount int, delta float32) error {
	ids := e.FindNearestNeurons(seedID, neighborsCount)
	if len(ids) == 0 {
		return fmt.Errorf("unable to build cluster for seedID=%d", seedID)
	}

	for _, id := range ids {
		if err := e.InjectNow(id, delta); err != nil {
			return err
		}
	}

	return nil
}

type chargedNeuron struct {
	ID        uint32
	Charge    float32
	Threshold float32
	Gap       float32
}

func (e *Engine) PrintTopCharged(limit int) {
	if limit <= 0 || len(e.Neurons) == 0 {
		return
	}

	items := make([]chargedNeuron, 0, len(e.Neurons))

	for _, n := range e.Neurons {
		effectiveThreshold := n.BaseThreshold + n.Adaptation
		items = append(items, chargedNeuron{
			ID:        n.ID,
			Charge:    n.Charge,
			Threshold: effectiveThreshold,
			Gap:       effectiveThreshold - n.Charge,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Charge > items[j].Charge
	})

	if limit > len(items) {
		limit = len(items)
	}

	fmt.Println("=== Top Charged Neurons ===")
	for i := 0; i < limit; i++ {
		x := items[i]
		fmt.Printf(
			"#%d neuron=%d charge=%.3f threshold=%.3f gap=%.3f\n",
			i+1,
			x.ID,
			x.Charge,
			x.Threshold,
			x.Gap,
		)
	}
	fmt.Println("===========================")
}

func (e *Engine) PrintCurrentSlotTargets() {
	if len(e.pending) == 0 {
		return
	}

	slot := e.currentSlot()
	events := e.pending[slot]

	if len(events) == 0 {
		fmt.Println("Current slot targets: no events")
		return
	}

	targetCount := make(map[uint32]int)
	targetDelta := make(map[uint32]float32)

	for _, ev := range events {
		targetCount[ev.TargetNeuronID]++
		targetDelta[ev.TargetNeuronID] += ev.Delta
	}

	type item struct {
		ID    uint32
		Count int
		Delta float32
	}

	items := make([]item, 0, len(targetCount))
	for id, count := range targetCount {
		items = append(items, item{
			ID:    id,
			Count: count,
			Delta: targetDelta[id],
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Count == items[j].Count {
			return items[i].Delta > items[j].Delta
		}
		return items[i].Count > items[j].Count
	})

	fmt.Printf("=== Current Slot %d Targets ===\n", slot)
	maxPrint := 15
	if maxPrint > len(items) {
		maxPrint = len(items)
	}
	for i := 0; i < maxPrint; i++ {
		x := items[i]
		fmt.Printf(
			"#%d neuron=%d inputs=%d total_delta=%.3f\n",
			i+1,
			x.ID,
			x.Count,
			x.Delta,
		)
	}
	fmt.Println("==============================")
}
