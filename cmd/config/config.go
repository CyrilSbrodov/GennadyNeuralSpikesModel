package config

type Config struct {
	NeuronLimit       int   // Лимит нейронов
	SynapseLimit      int   // Лимит синапсов
	SynapsesPerNeuron int   // Количество синапсов у нейрона
	Seed              int64 // Seed

	WorldSizeX float32 // Размер мира по оси X
	WorldSizeY float32 // Размер мира по оси Y
	WorldSizeZ float32 // Размер мира по оси Z

	MaxAxonLength      float32 // Максимальная длина аксона. То есть нейрон не может соединиться дальше этого расстояния
	LocalConnectRadius float32 // Радиус локального соединения. Если нейроны ближе этого расстояния — связь почти гарантирована. Это создаёт локальные кластеры.

	LongConnectionProb float32 // Вероятность дальних соединений. Это имитирует long-range connections мозга.

	ExcitatoryRatio float32 // Соотношение типов нейронов
	InhibitoryRatio float32 // Соотношение типов нейронов.

	InputRatio  float32 // Сколько нейронов являются InputRatio
	OutputRatio float32 // Сколько нейронов являются OutputRatio

	MinDelay uint16 // Минимальная задержка передачи сигнала.
	MaxDelay uint16 // Максимальная задержка передачи сигнала.

	RestCharge      float32 // Потенциал покоя. К этому значению заряд стремится со временем.
	BaseThreshold   float32 // Базовый порог возбуждения. Если заряд выше порога - нейрон стреляет.
	ThresholdNoise  float32 // Шум порога. Добавляется случайное отклонение: threshold = BaseThreshold ± noise. Это делает нейроны немного разными. Иначе они будут вести себя слишком синхронно.
	ResetCharge     float32 // Заряд после спайка.
	LeakFactor      float32 // Утечка заряда.
	CooldownTicks   uint16  // Рефрактерный период. После спайка нейрон не может стрелять несколько тиков.
	AdaptationDecay float32 // Скорость восстановления нейрона после активности. Каждый тик Adaptation *= AdaptationDecay
	AdaptationStep  float32 // Насколько увеличивается усталость после спайка. После спайка Adaptation += AdaptationStep

	ExcitatoryWeightMin float32 // Минимальный вес возбуждающих нейронов.
	ExcitatoryWeightMax float32 // Максимальный вес возбуждающих нейронов.
	InhibitoryWeightMin float32 // Минимальный вес тормозных синапсов.
	InhibitoryWeightMax float32 // Максимальный вес тормозных синапсов.
	WeightMin           float32 // Абсолютная граница максимального веса.
	WeightMax           float32 // Абсолютная граница минимального веса.

	HebbianEnable             bool    // Включает или выключает обучение
	HebbianLearningRate       float32 // Скорость обучения. Чем больше значение - тем быстрее меняются веса
	HebbianDecay              float32 // Фоновое ослабление веса. Используется как мягкий стабилизатор
	STDPWindowTicks           uint64  // Временное окно STDP. Если нейрон A стрелял за несколько тиков до B, связь может усиливаться.
	STDPPotentiation          float32 // Сила усиления связи.
	STDPDepression            float32 // Сила ослабления связи.
	SynapseUsageDecayInterval uint64  // Как часто проверять забытые связи.
	SynapseUsageWeightDecay   float32 // Насколько ослаблять неиспользуемые синапсы. Это имитирует забывание мозга.
}

func DefaultConfig() *Config {
	return &Config{
		NeuronLimit:       1000, // Лимит нейронов
		SynapseLimit:      8000, // Лимит синапсов
		SynapsesPerNeuron: 8,    // Количество синапсов у нейрона
		Seed:              0,    // Seed

		WorldSizeX: 25, // Размер мира по оси X
		WorldSizeY: 25, // Размер мира по оси Y
		WorldSizeZ: 25, // Размер мира по оси Z

		MaxAxonLength:      10, // Максимальная длина аксона. То есть нейрон не может соединиться дальше этого расстояния
		LocalConnectRadius: 4,  // Радиус локального соединения. Если нейроны ближе этого расстояния — связь почти гарантирована. Это создаёт локальные кластеры.

		LongConnectionProb: 0.10, // Вероятность дальних соединений. Это имитирует long-range connections мозга.

		ExcitatoryRatio: 0.8, // Соотношение типов нейронов.
		InhibitoryRatio: 0.2, // Соотношение типов нейронов.

		InputRatio:  0.03, // Сколько нейронов являются InputRatio
		OutputRatio: 0.03, // Сколько нейронов являются OutputRatio

		MinDelay: 1, // Минимальная задержка передачи сигнала.
		MaxDelay: 6, // Максимальная задержка передачи сигнала.

		RestCharge:      -70,  // Потенциал покоя. К этому значению заряд стремится со временем.
		BaseThreshold:   -60,  // Базовый порог возбуждения. Если заряд выше порога - нейрон стреляет.
		ThresholdNoise:  2,    // Шум порога. Добавляется случайное отклонение: threshold = BaseThreshold ± noise. Это делает нейроны немного разными. Иначе они будут вести себя слишком синхронно.
		ResetCharge:     -75,  // Заряд после спайка.
		LeakFactor:      0.98, // Утечка заряда.
		CooldownTicks:   2,    // Рефрактерный период. После спайка нейрон не может стрелять несколько тиков.
		AdaptationDecay: 0.96, // Скорость восстановления нейрона после активности. Каждый тик Adaptation *= AdaptationDecay
		AdaptationStep:  0.5,  // Насколько увеличивается усталость после спайка. После спайка Adaptation += AdaptationStep

		// Диапазон весов возбуждающих нейронов.
		ExcitatoryWeightMin: 1.5,  // Минимальный вес возбуждающих нейронов.
		ExcitatoryWeightMax: 4.0,  // Максимальный вес возбуждающих нейронов.
		InhibitoryWeightMin: -4.0, // Минимальный вес тормозных синапсов.
		InhibitoryWeightMax: -1.5, // Максимальный вес тормозных синапсов.
		WeightMin:           -5.0, // Абсолютная граница максимального веса.
		WeightMax:           5.0,  // Абсолютная граница минимального веса.

		// Пластичность (обучение)
		HebbianEnable:             true,  // Включает или выключает обучение
		HebbianLearningRate:       0.02,  // Скорость обучения. Чем больше значение - тем быстрее меняются веса
		HebbianDecay:              0.005, // Фоновое ослабление веса. Используется как мягкий стабилизатор
		STDPWindowTicks:           8,     // Временное окно STDP. Если нейрон A стрелял за несколько тиков до B, связь может усиливаться.
		STDPPotentiation:          1.0,   // Сила усиления связи.
		STDPDepression:            0.25,  // Сила ослабления связи.
		SynapseUsageDecayInterval: 32,    // Как часто проверять забытые связи.
		SynapseUsageWeightDecay:   0.01,  // Насколько ослаблять неиспользуемые синапсы. Это имитирует забывание мозга.
	}
}
