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
	{
		ID:                 1,
		Title:              "Умная люстра",
		PowerW:             100,
		LuminousFluxLm:     1000,
		ScatteringAngleDeg: 100,
		ImageURL:           "https://placehold.co/300",
	},
	{
		ID:                 2,
		Title:              "Slim Magnetic",
		PowerW:             100,
		LuminousFluxLm:     1000,
		ScatteringAngleDeg: 100,
		ImageURL:           "https://placehold.co/300",
	},
	{
		ID:                 3,
		Title:              "Светильник потолочный",
		PowerW:             100,
		LuminousFluxLm:     1000,
		ScatteringAngleDeg: 100,
		ImageURL:           "https://placehold.co/300",
	},
	{
		ID:                 4,
		Title:              "Потолочный светильник",
		PowerW:             100,
		LuminousFluxLm:     1000,
		ScatteringAngleDeg: 100,
		ImageURL:           "https://placehold.co/300",
	},
	{
		ID:                 5,
		Title:              "Подвесной светильник",
		PowerW:             100,
		LuminousFluxLm:     1000,
		ScatteringAngleDeg: 100,
		ImageURL:           "https://placehold.co/300",
	},
	{
		ID:                 6,
		Title:              "Esthetic Magnetic",
		PowerW:             100,
		LuminousFluxLm:     1000,
		ScatteringAngleDeg: 100,
		ImageURL:           "https://placehold.co/300",
	},
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
