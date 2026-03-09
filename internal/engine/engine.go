package engine

import (
	"encoding/json"
	"fmt"
	"os"
)

type Neuron struct {
	Charge       float64 `json:"charge"`
	Threshold    float64 `json:"threshold"`
	Cooldown     int     `json:"cooldown"`
	LastFireTick int     `json:"lastFireTick"`
}

type Config struct {
	LeakFactor          float64 `json:"leakFactor"`
	MinWeight           float64 `json:"minWeight"`
	MaxWeight           float64 `json:"maxWeight"`
	LearnRate           float64 `json:"learnRate"`
	ForgetRate          float64 `json:"forgetRate"`
	LearnWindowTicks    int     `json:"learnWindowTicks"`
	ForgetAfterInactive int     `json:"forgetAfterInactive"`
}

type InputSignal struct {
	NeuronID int     `json:"neuronId"`
	Charge   float64 `json:"charge"`
}

type Snapshot struct {
	Tick       int       `json:"tick"`
	Neurons    []Neuron  `json:"neurons"`
	SynTarget  []int     `json:"synTarget"`
	SynSource  []int     `json:"synSource"`
	SynWeight  []float64 `json:"synWeight"`
	SynDelay   []int     `json:"synDelay"`
	SynPlastic []bool    `json:"synPlastic"`
	SynOffset  []int     `json:"synOffset"`
	Config     Config    `json:"config"`
}

type spikeEvent struct {
	SynapseID int
	TargetID  int
	Weight    float64
}

type Engine struct {
	Tick int

	Neurons []Neuron

	SynTarget   []int
	SynSource   []int
	SynWeight   []float64
	SynDelay    []int
	SynPlastic  []bool
	SynLastUsed []int
	SynOffset   []int

	InSynOffset []int
	InSynIndex  []int

	SpikeQueue [][]spikeEvent

	Config Config
}

type SynapseDef struct {
	Source  int     `json:"source"`
	Target  int     `json:"target"`
	Weight  float64 `json:"weight"`
	Delay   int     `json:"delay"`
	Plastic bool    `json:"plastic"`
}

func NewEngine(neuronCount int, synapses []SynapseDef, cfg Config) (*Engine, error) {
	if neuronCount <= 0 {
		return nil, fmt.Errorf("neuronCount must be > 0")
	}
	maxDelay := 1
	for _, s := range synapses {
		if s.Source < 0 || s.Source >= neuronCount || s.Target < 0 || s.Target >= neuronCount {
			return nil, fmt.Errorf("invalid synapse (%d -> %d)", s.Source, s.Target)
		}
		if s.Delay <= 0 {
			return nil, fmt.Errorf("delay must be > 0")
		}
		if s.Delay > maxDelay {
			maxDelay = s.Delay
		}
	}

	counts := make([]int, neuronCount)
	for _, s := range synapses {
		counts[s.Source]++
	}

	synOffset := make([]int, neuronCount+1)
	for i := 0; i < neuronCount; i++ {
		synOffset[i+1] = synOffset[i] + counts[i]
	}

	target := make([]int, len(synapses))
	source := make([]int, len(synapses))
	weight := make([]float64, len(synapses))
	delay := make([]int, len(synapses))
	plastic := make([]bool, len(synapses))
	lastUsed := make([]int, len(synapses))
	for i := range lastUsed {
		lastUsed[i] = -1
	}

	writePos := append([]int(nil), synOffset[:neuronCount]...)
	for _, s := range synapses {
		idx := writePos[s.Source]
		writePos[s.Source]++
		target[idx] = s.Target
		source[idx] = s.Source
		weight[idx] = s.Weight
		delay[idx] = s.Delay
		plastic[idx] = s.Plastic
	}

	inCounts := make([]int, neuronCount)
	for _, t := range target {
		inCounts[t]++
	}
	inOffset := make([]int, neuronCount+1)
	for i := 0; i < neuronCount; i++ {
		inOffset[i+1] = inOffset[i] + inCounts[i]
	}
	inIndex := make([]int, len(target))
	inWrite := append([]int(nil), inOffset[:neuronCount]...)
	for synID := range target {
		t := target[synID]
		idx := inWrite[t]
		inWrite[t]++
		inIndex[idx] = synID
	}

	neurons := make([]Neuron, neuronCount)
	for i := range neurons {
		neurons[i] = Neuron{Threshold: 75, Cooldown: 2, LastFireTick: -1}
	}

	return &Engine{
		Neurons:     neurons,
		SynTarget:   target,
		SynSource:   source,
		SynWeight:   weight,
		SynDelay:    delay,
		SynPlastic:  plastic,
		SynLastUsed: lastUsed,
		SynOffset:   synOffset,
		InSynOffset: inOffset,
		InSynIndex:  inIndex,
		SpikeQueue:  make([][]spikeEvent, maxDelay+1),
		Config:      cfg,
	}, nil
}

func (e *Engine) Step(inputs []InputSignal) []int {
	for _, in := range inputs {
		if in.NeuronID < 0 || in.NeuronID >= len(e.Neurons) {
			continue
		}
		e.Neurons[in.NeuronID].Charge += in.Charge
	}

	bucket := e.Tick % len(e.SpikeQueue)
	for _, ev := range e.SpikeQueue[bucket] {
		e.Neurons[ev.TargetID].Charge += ev.Weight
		e.SynLastUsed[ev.SynapseID] = e.Tick
	}
	e.SpikeQueue[bucket] = e.SpikeQueue[bucket][:0]

	fired := make([]int, 0)
	for i := range e.Neurons {
		n := &e.Neurons[i]
		n.Charge *= e.Config.LeakFactor
		if n.LastFireTick >= 0 && e.Tick-n.LastFireTick < n.Cooldown {
			continue
		}
		if n.Charge >= n.Threshold {
			n.Charge = 0
			n.LastFireTick = e.Tick
			fired = append(fired, i)
		}
	}

	for _, src := range fired {
		start, end := e.SynOffset[src], e.SynOffset[src+1]
		for synID := start; synID < end; synID++ {
			deliverTick := (e.Tick + e.SynDelay[synID]) % len(e.SpikeQueue)
			e.SpikeQueue[deliverTick] = append(e.SpikeQueue[deliverTick], spikeEvent{
				SynapseID: synID,
				TargetID:  e.SynTarget[synID],
				Weight:    e.SynWeight[synID],
			})
			e.SynLastUsed[synID] = e.Tick
		}
	}

	for _, dst := range fired {
		start, end := e.InSynOffset[dst], e.InSynOffset[dst+1]
		for i := start; i < end; i++ {
			synID := e.InSynIndex[i]
			if !e.SynPlastic[synID] {
				continue
			}
			src := e.SynSource[synID]
			srcFire := e.Neurons[src].LastFireTick
			if srcFire >= 0 && dst != src && e.Tick-srcFire <= e.Config.LearnWindowTicks {
				e.SynWeight[synID] += e.Config.LearnRate
				if e.SynWeight[synID] > e.Config.MaxWeight {
					e.SynWeight[synID] = e.Config.MaxWeight
				}
			}
		}
	}

	for synID := range e.SynWeight {
		if e.SynLastUsed[synID] >= 0 && e.Tick-e.SynLastUsed[synID] >= e.Config.ForgetAfterInactive {
			e.SynWeight[synID] -= e.Config.ForgetRate
			if e.SynWeight[synID] < e.Config.MinWeight {
				e.SynWeight[synID] = e.Config.MinWeight
			}
		}
	}

	e.Tick++
	return fired
}

func (e *Engine) Snapshot() Snapshot {
	return Snapshot{
		Tick:       e.Tick,
		Neurons:    append([]Neuron(nil), e.Neurons...),
		SynTarget:  append([]int(nil), e.SynTarget...),
		SynSource:  append([]int(nil), e.SynSource...),
		SynWeight:  append([]float64(nil), e.SynWeight...),
		SynDelay:   append([]int(nil), e.SynDelay...),
		SynPlastic: append([]bool(nil), e.SynPlastic...),
		SynOffset:  append([]int(nil), e.SynOffset...),
		Config:     e.Config,
	}
}

func (e *Engine) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(e)
}

func Load(path string) (*Engine, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var e Engine
	if err := json.NewDecoder(f).Decode(&e); err != nil {
		return nil, err
	}
	if len(e.SpikeQueue) == 0 {
		e.SpikeQueue = make([][]spikeEvent, 8)
	}
	if len(e.InSynOffset) == 0 || len(e.InSynIndex) == 0 {
		e.rebuildIncoming()
	}
	return &e, nil
}

func (e *Engine) rebuildIncoming() {
	neuronCount := len(e.Neurons)
	counts := make([]int, neuronCount)
	for _, t := range e.SynTarget {
		counts[t]++
	}
	e.InSynOffset = make([]int, neuronCount+1)
	for i := 0; i < neuronCount; i++ {
		e.InSynOffset[i+1] = e.InSynOffset[i] + counts[i]
	}
	e.InSynIndex = make([]int, len(e.SynTarget))
	write := append([]int(nil), e.InSynOffset[:neuronCount]...)
	for synID, t := range e.SynTarget {
		idx := write[t]
		write[t]++
		e.InSynIndex[idx] = synID
	}
}
