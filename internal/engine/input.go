package engine

import "fmt"

func (e *Engine) InputNeuronIDs() []uint32 {
	ids := make([]uint32, 0)
	for _, n := range e.Neurons {
		if n.Role == RoleInput {
			ids = append(ids, n.ID)
		}
	}
	return ids
}

func (e *Engine) InjectTimed(events []TimedSpikeEvent) error {
	for _, te := range events {
		if te.Tick < e.Tick {
			return fmt.Errorf("cannot inject into the past: tick=%d now=%d", te.Tick, e.Tick)
		}

		delayRaw := te.Tick - e.Tick
		delay := uint16(delayRaw)
		if delayRaw > uint64(^uint16(0)) {
			delay = ^uint16(0)
		}

		if err := e.InjectAfter(delay, te.Event.TargetNeuronID, te.Event.Delta); err != nil {
			return err
		}
	}
	return nil
}
