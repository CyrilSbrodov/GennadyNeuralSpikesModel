package engine

import "testing"

func TestSpikeDeliveryAndFire(t *testing.T) {
	e, err := NewEngine(2, []SynapseDef{{Source: 0, Target: 1, Weight: 80, Delay: 1, Plastic: true}}, Config{
		LeakFactor: 1, MinWeight: 0, MaxWeight: 100, LearnRate: 1, ForgetRate: 0.1, LearnWindowTicks: 2, ForgetAfterInactive: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	e.Neurons[0].Threshold = 10
	e.Neurons[1].Threshold = 50

	fired := e.Step([]InputSignal{{NeuronID: 0, Charge: 10}})
	if len(fired) != 1 || fired[0] != 0 {
		t.Fatalf("expected neuron 0 to fire, got %#v", fired)
	}
	fired = e.Step(nil)
	if len(fired) != 1 || fired[0] != 1 {
		t.Fatalf("expected neuron 1 to fire after delayed spike, got %#v", fired)
	}
}

func TestLearningIncreasesWeight(t *testing.T) {
	e, _ := NewEngine(2, []SynapseDef{{Source: 0, Target: 1, Weight: 10, Delay: 1, Plastic: true}}, Config{
		LeakFactor: 1, MinWeight: 0, MaxWeight: 100, LearnRate: 2, ForgetRate: 0.1, LearnWindowTicks: 3, ForgetAfterInactive: 100,
	})
	e.Neurons[0].Threshold = 10
	e.Neurons[1].Threshold = 5

	e.Step([]InputSignal{{NeuronID: 0, Charge: 10}})
	e.Step([]InputSignal{{NeuronID: 1, Charge: 5}})

	if e.SynWeight[0] <= 10 {
		t.Fatalf("expected weight to increase, got %f", e.SynWeight[0])
	}
}
