package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"math/rand"
	"os"

	"gennady-neural-spikes-model/cmd/config"
	"gennady-neural-spikes-model/internal/dataset"
	"gennady-neural-spikes-model/internal/encoder"
	"gennady-neural-spikes-model/internal/engine"
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
	fmt.Println("Loaded samples:", len(samples))
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

	//seedID := uint32(100)
	//
	//nearest := brain.FindNearestNeurons(seedID, 11)
	//fmt.Println("cluster ids:", nearest)
	//
	//if err := brain.InjectCluster(seedID, 11, 25); err != nil {
	//	panic(err)
	//}
	//
	//for i := 0; i < 10; i++ {
	//	fmt.Printf("\n--- BEFORE STEP tick=%d ---\n", brain.Tick)
	//	brain.PrintCurrentSlotTargets()
	//
	//	stats := brain.Step()
	//
	//	fmt.Printf(
	//		"tick=%d delivered=%d created=%d fired=%d mean_charge=%.2f\n",
	//		stats.Tick,
	//		stats.DeliveredEvents,
	//		stats.CreatedEvents,
	//		stats.FiredCount,
	//		stats.MeanCharge,
	//	)
	//
	//	brain.PrintTopCharged(10)
	//}

	encCfg := encoder.EncoderConfig{
		TimeStepTicks:   1,
		UnitWindowTicks: 1,
		IntensityScale:  12,
		InputChannels:   4,
		NeuronsPerInput: 2,
		SpikesPerWindow: 2,
		Coding:          encoder.CodingRate,
		InputNeuronIDs:  brain.InputNeuronIDs(),
		ImageWidth:      2,
		ImageHeight:     2,
		UseRGB:          false,
	}

	if err := encCfg.Validate(); err != nil {
		panic(err)
	}

	textSpikes := encoder.EncodeTextToSpikes(samples[0].Label, brain.Tick+1, encCfg)
	if err := brain.InjectTimed(textSpikes); err != nil {
		panic(err)
	}

	img, err := decodeImage(samples[0].ImagePath)
	if err != nil {
		panic(err)
	}
	imgSpikes := encoder.EncodeImageToSpikes(img, brain.Tick+3, encCfg)
	if err := brain.InjectTimed(imgSpikes); err != nil {
		panic(err)
	}

	stats := brain.Run(10)
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

func decodeImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}
