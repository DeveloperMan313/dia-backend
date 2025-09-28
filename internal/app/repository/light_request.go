package repository

import (
	"fmt"
)

type LightRequestRepository struct {
}

func NewLightRequestRepository() (*LightRequestRepository, error) {
	return &LightRequestRepository{}, nil
}

type LightRequest struct {
	ID             int
	MaxTotalPowerW float32
}

type LightRequestViewEntry struct {
	Lamp   Lamp
	AreaM2 float32
	Number int
}

type LightRequestView struct {
	LightRequest LightRequest
	Entries      []LightRequestViewEntry
}

var lightRequestViewByID = map[int]LightRequestView{
	1: {
		LightRequest: LightRequest{
			ID:             1,
			MaxTotalPowerW: 200,
		},
		Entries: []LightRequestViewEntry{
			{
				Lamp:   lamps[1],
				AreaM2: 40,
				Number: 1,
			},
			{
				Lamp:   lamps[2],
				AreaM2: 25,
				Number: 10,
			},
		},
	},
}

func (*LightRequestRepository) GetLightRequestEntryCntByID(id int) (int, error) {
	if len(lightRequestViewByID) == 0 {
		return 0, fmt.Errorf("массив пустой")
	}

	lightRequestView, found := lightRequestViewByID[id]
	if !found {
		return 0, fmt.Errorf("не найдено")
	}

	return len(lightRequestView.Entries), nil
}

func (*LightRequestRepository) GetLightRequestViewByID(id int, lampRepo *LampRepository) (*LightRequestView, error) {
	if len(lightRequestViewByID) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	lightRequestView, found := lightRequestViewByID[id]
	if !found {
		return nil, fmt.Errorf("не найдено")
	}

	return &lightRequestView, nil
}
