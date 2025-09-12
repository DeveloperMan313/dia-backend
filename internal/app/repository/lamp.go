package repository

import (
	"fmt"
	"strings"
)

type LampRepository struct {
}

func NewLampRepository() (*LampRepository, error) {
	return &LampRepository{}, nil
}

type Lamp struct {
	ID                 int
	Title              string
	PowerW             float32
	LuminousFluxLm     float32
	ScatteringAngleDeg float32
	ImageURL           string
}

var lamps = []Lamp{
	{ID: 1, Title: "Умная потолочная люстра с перламутром", PowerW: 40, LuminousFluxLm: 3200, ScatteringAngleDeg: 120, ImageURL: "http://localhost:9000/lamp-images/1.jpg"},
	{ID: 2, Title: "Slim Magnetic Трековый светильник 26W 4000K Most чёрный", PowerW: 26, LuminousFluxLm: 2210, ScatteringAngleDeg: 60, ImageURL: "http://localhost:9000/lamp-images/2.jpg"},
	{ID: 3, Title: "Светильник потолочный светодиодный Trio 8W 3000K белый", PowerW: 8, LuminousFluxLm: 680, ScatteringAngleDeg: 120, ImageURL: "http://localhost:9000/lamp-images/3.jpg"},
	{ID: 4, Title: "Потолочный светильник", PowerW: 40, LuminousFluxLm: 3200, ScatteringAngleDeg: 120, ImageURL: "http://localhost:9000/lamp-images/4.jpg"},
	{ID: 5, Title: "Подвесной светильник со стеклянными плафонами", PowerW: 40, LuminousFluxLm: 3200, ScatteringAngleDeg: 120, ImageURL: "http://localhost:9000/lamp-images/5.jpg"},
	{ID: 6, Title: "Esthetic Magnetic Трековый светильник 3W 3000K (чёрный)", PowerW: 3, LuminousFluxLm: 255, ScatteringAngleDeg: 60, ImageURL: "http://localhost:9000/lamp-images/6.jpg"},
	{ID: 7, Title: "Подвесной светильник", PowerW: 40, LuminousFluxLm: 3200, ScatteringAngleDeg: 120, ImageURL: "http://localhost:9000/lamp-images/7.jpg"},
	{ID: 8, Title: "Трековый светильник 100W 4200K Full Light N05 Slim Magnetic", PowerW: 100, LuminousFluxLm: 8500, ScatteringAngleDeg: 60, ImageURL: "http://localhost:9000/lamp-images/8.jpg"},
	{ID: 9, Title: "Светильник встраиваемый светодиодный Forte 15W 4000K титан", PowerW: 15, LuminousFluxLm: 1275, ScatteringAngleDeg: 120, ImageURL: "http://localhost:9000/lamp-images/9.jpg"},
	{ID: 10, Title: "Подвесной светильник со стеклянными плафонами", PowerW: 40, LuminousFluxLm: 3200, ScatteringAngleDeg: 120, ImageURL: "http://localhost:9000/lamp-images/10.jpg"},
	{ID: 11, Title: "Светильник потолочный светодиодный Tend 9W 4000K черный", PowerW: 9, LuminousFluxLm: 765, ScatteringAngleDeg: 120, ImageURL: "http://localhost:9000/lamp-images/11.jpg"},
	{ID: 12, Title: "Подвесной светодиодный светильник", PowerW: 40, LuminousFluxLm: 3200, ScatteringAngleDeg: 120, ImageURL: "http://localhost:9000/lamp-images/12.jpg"},
}

func (*LampRepository) GetLampByID(id int) (*Lamp, error) {
	if len(lamps) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	for _, lamp := range lamps {
		if lamp.ID == id {
			return &lamp, nil
		}
	}

	return nil, fmt.Errorf("не найдено")
}

func (*LampRepository) GetLamps() ([]Lamp, error) {
	if len(lamps) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return lamps, nil
}

func (r *LampRepository) GetLampsByTitle(title string) ([]Lamp, error) {
	lamps, err := r.GetLamps()
	if err != nil {
		return []Lamp{}, err
	}

	var result []Lamp
	for _, lamp := range lamps {
		if strings.Contains(strings.ToLower(lamp.Title), strings.ToLower(title)) {
			result = append(result, lamp)
		}
	}

	return result, nil
}
