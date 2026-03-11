package main

import (
	"fmt"
	"gennady-neural-spikes-model/cmd/config"
	"gennady-neural-spikes-model/internal/dataset"
	"gennady-neural-spikes-model/internal/engine"
	"log/slog"
	"math/rand"
)

func main() {
	cfg := config.DefaultConfig()
	if cfg.Seed == 0 {
		cfg.Seed = int64(rand.Uint64())
	}
	slog.Info("Start with:", "seed", cfg.Seed)
	brain := engine.NewEngine(cfg)
	brain.InitSpatial3D()

	if err := brain.ValidateTopology(); err != nil {
		panic(err)
	}

	brain.PrintTopologyStats()
	brain.PrintDelayHistogram()
	brain.PrintOutgoingHistogram()

	samples, err := dataset.Load("data")
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded samples:", samples)
	//_ = brain.InjectNow(0, 30)
	//_ = brain.InjectNow(1, 18)
	//_ = brain.InjectNow(2, 22)
	//_ = brain.InjectNow(3, 20)
	//_ = brain.InjectNow(4, 18)
	//_ = brain.InjectNow(5, 30)
	//_ = brain.InjectNow(6, 20)
	//_ = brain.InjectNow(7, 18)
	//_ = brain.InjectNow(8, 25)
	//_ = brain.InjectNow(9, 20)
	//_ = brain.InjectNow(10, 30)
	//_ = brain.InjectNow(11, 22)

	seedID := uint32(100)

	nearest := brain.FindNearestNeurons(seedID, 11)
	fmt.Println("cluster ids:", nearest)

	if err := brain.InjectCluster(seedID, 11, 25); err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		fmt.Printf("\n--- BEFORE STEP tick=%d ---\n", brain.Tick)
		brain.PrintCurrentSlotTargets()

		stats := brain.Step()

		fmt.Printf(
			"tick=%d delivered=%d created=%d fired=%d mean_charge=%.2f\n",
			stats.Tick,
			stats.DeliveredEvents,
			stats.CreatedEvents,
			stats.FiredCount,
			stats.MeanCharge,
		)

		brain.PrintTopCharged(10)
	}

	if err := brain.InjectCluster(seedID, 11, 10); err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		fmt.Printf("\n--- BEFORE STEP tick=%d ---\n", brain.Tick)
		brain.PrintCurrentSlotTargets()

		stats := brain.Step()

		fmt.Printf(
			"tick=%d delivered=%d created=%d fired=%d mean_charge=%.2f\n",
			stats.Tick,
			stats.DeliveredEvents,
			stats.CreatedEvents,
			stats.FiredCount,
			stats.MeanCharge,
		)

		brain.PrintTopCharged(10)
	}

	if err := brain.InjectCluster(seedID, 11, 5); err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		fmt.Printf("\n--- BEFORE STEP tick=%d ---\n", brain.Tick)
		brain.PrintCurrentSlotTargets()

		stats := brain.Step()

		fmt.Printf(
			"tick=%d delivered=%d created=%d fired=%d mean_charge=%.2f\n",
			stats.Tick,
			stats.DeliveredEvents,
			stats.CreatedEvents,
			stats.FiredCount,
			stats.MeanCharge,
		)

		brain.PrintTopCharged(10)
	}

	stats := brain.Run(20)

	for _, s := range stats {
		fmt.Printf(
			"tick=%d delivered=%d created=%d fired=%d mean_charge=%.2f\n",
			s.Tick,
			s.DeliveredEvents,
			s.CreatedEvents,
			s.FiredCount,
			s.MeanCharge,
		)
	}
}
