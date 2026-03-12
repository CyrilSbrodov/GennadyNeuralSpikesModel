package engine

import (
	"encoding/gob"
	"fmt"
	"os"

	"gennady-neural-spikes-model/cmd/config"
)

// BrainState — плоский снапшот, который gob любит
type BrainState struct {
	Neurons  []Neuron
	Synapses []Synapse
	Tick     uint64
}

func (e *Engine) Save(filename string) error {
	state := BrainState{
		Neurons:  e.Neurons,
		Synapses: e.Synapses,
		Tick:     e.Tick,
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create state file: %w", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(state); err != nil {
		return fmt.Errorf("encode brain: %w", err)
	}

	fmt.Printf("✅ Мозг сохранён: %s (нейронов: %d, синапсов: %d, тик: %d)\n",
		filename, len(e.Neurons), len(e.Synapses), e.Tick)
	return nil
}

// LoadOrCreate — главная магия
// Если файл есть — загружает, если нет — создаёт новый через NewEngine
func LoadOrCreate(cfg *config.Config) (*Engine, error) {
	filename := cfg.StateFile // берём из yaml

	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("🔄 Загружаем мозг из %s...\n", filename)

		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		state := BrainState{}
		if err := gob.NewDecoder(file).Decode(&state); err != nil {
			return nil, fmt.Errorf("decode brain: %w", err)
		}

		e := &Engine{
			Neurons:  state.Neurons,
			Synapses: state.Synapses,
			Config:   cfg,
			Tick:     state.Tick,
			pending:  make([][]SpikeEvent, cfg.MaxDelay+1), // пересоздаём runtime
			fired:    make([]uint32, 0, 1024),
		}
		fmt.Printf("✅ Мозг загружен! Нейронов: %d, синапсов: %d, тик: %d\n",
			len(e.Neurons), len(e.Synapses), e.Tick)
		return e, nil
	}

	// файла нет — новый мозг
	fmt.Println("🆕 Файл состояния не найден, создаём новый мозг...")
	return NewEngine(cfg), nil // твоя существующая функция создания
}
