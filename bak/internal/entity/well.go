package entity

import "time"

type Well struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	GammaG    float64   `json:"gamma"`     // Относительная плотность газа (по воздуху)
	Temp      float64   `json:"temp"`      // Температура пласта, К
	TempUst   float64   `json:"tempust"`   // Температура на устье, К
	Depth     float64   `json:"depth"`     // Глубина скважины, м
	Pbuf      float64   `json:"pbuf"`      // Буферное давление, Па
	Ptb       float64   `json:"ptb"`       // Затрубное давление, Па
	Ppl       float64   `json:"ppl"`       // Пастовое давление, Па
	Pz        float64   `json:"pz"`        // Забойноеое давление, Па
	Q         float64   `json:"q"`         // Дебит газа, м³/сут
	Roughness float64   `json:"roughness"` // Шероховатость трубы, м
	Diameter  float64   `json:"diameter"`  // Диаметр НКТ, м
	A         float64   `json:"a"`         // пластовый коэф-т
	B         float64   `json:"b"`         // пластовый коэф-т
	Mu        float64   `json:"mu"`        // вязкость газа мПа*с
	WGF       float64   `json:"wgf"`       // водогазовый фактор
	Rog       float64   `json:"rog"`       //Плотность воды, кг/м3
	Hw        float64   `json:"hw"`        // Высота столба ГЖС, м
	Qmin      float64   `json:"qmin"`      // м3/сут
	Pmax      float64   `json:"pmax"`
	Status    string    `json:"status"` //
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}
