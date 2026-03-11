package engine

type NeuronPolarity uint8

// Если нейрон возбуждающий — его синапсы обычно положительные. Если нейрон возбуждающий — его синапсы обычно положительные
const (
	PolarityExcitatory NeuronPolarity = iota // Возбуждающий
	PolarityInhibitory                       // Тормозящий
)

type NeuronRole uint8

const (
	RoleInput  NeuronRole = iota // Получает сенсорные спайки от encoder
	RoleHidden                   // Участвует во внутренней динамике
	RoleOutput                   // Потом может стать слоем ответа/действия
)

type Vec3 struct {
	X float32
	Y float32
	Z float32
}

type Neuron struct {
	ID uint32

	Polarity NeuronPolarity // Полярность нейрона
	Role     NeuronRole     // Роль нейрона

	Position Vec3 // Координаты нейрона в 3D-пространстве.

	Charge          float32 // Текущий заряд / мембранный потенциал нейрона
	RestCharge      float32 // Потенциал покоя. К этому потенциалу стремится нейрон
	BaseThreshold   float32 // Базовый порог срабатывания. Это исходный порог, не учитывающий текущую усталость/адаптацию
	ResetCharge     float32 // Заряд после выстрела
	Adaptation      float32 // Текущая адаптация / усталость нейрона. Если нейрон часто стреляет, его становится сложнее снова возбудить
	AdaptationDecay float32 // Коэффициент затухания адаптации. Показывает, с какой скоростью нейрон “забывает усталость”
	AdaptationStep  float32 // Насколько увеличивается адаптация после одного спайка. Это “цена одного выстрела” для нейрона
	LastSpikeTick   int64   // Тик последнего спайка этого нейрона
	FireCount       uint64  // Сколько раз нейрон стрелял вообще. Для диагностики.

	CooldownTicks uint16 // Длительность периода покоя после спайка
	CooldownLeft  uint16 // Сколько тиков покоя осталось прямо сейчас

	Outgoing []uint32 // Список ID/индексов исходящих синапсов
	Incoming []uint32 // Список входящих синапсов

	FiredLastTick bool // Стрелял ли нейрон на прошлом тике
}
