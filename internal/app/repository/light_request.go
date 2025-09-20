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
	TotalPowerW    float32
}

type LightRequestToLamp struct {
	RequestID int
	LampID    int
	AreaM2    float32
	Number    int
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

var lightRequests = []LightRequest{
	{
		ID:             1,
		MaxTotalPowerW: 200,
		TotalPowerW:    158,
	},
}

var lightRequestToLamps = []LightRequestToLamp{
	{
		RequestID: 1,
		LampID:    2,
		AreaM2:    40,
		Number:    1,
	},
	{
		RequestID: 1,
		LampID:    3,
		AreaM2:    25,
		Number:    10,
	},
}

func (*LightRequestRepository) GetLightRequestEntryCntByID(id int) (int, error) {
	if len(lightRequests) == 0 {
		return 0, fmt.Errorf("массив пустой")
	}

	var lightRequest *LightRequest = nil
	for _, req := range lightRequests {
		if req.ID == id {
			lightRequest = &req
		}
	}
	if lightRequest == nil {
		return 0, fmt.Errorf("не найдено")
	}

	var lightRequestEntryCnt int = 0
	for _, reqToLamp := range lightRequestToLamps {
		if reqToLamp.RequestID == lightRequest.ID {
			lightRequestEntryCnt++
		}
	}

	return lightRequestEntryCnt, nil
}

func (*LightRequestRepository) GetLightRequestViewByID(id int, lampRepo *LampRepository) (*LightRequestView, error) {
	if len(lightRequests) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	var lightRequest *LightRequest = nil
	for _, req := range lightRequests {
		if req.ID == id {
			lightRequest = &req
		}
	}
	if lightRequest == nil {
		return nil, fmt.Errorf("не найдено")
	}

	lightRequestView := LightRequestView{
		LightRequest: *lightRequest,
	}

	for _, reqToLamp := range lightRequestToLamps {
		if reqToLamp.RequestID == lightRequest.ID {
			lamp, err := lampRepo.GetLampByID(reqToLamp.LampID)
			if err != nil {
				return nil, fmt.Errorf("не найдено")
			}
			lightRequestView.Entries = append(lightRequestView.Entries, LightRequestViewEntry{
				Lamp:   *lamp,
				AreaM2: reqToLamp.AreaM2,
				Number: reqToLamp.Number,
			})
		}
	}

	return &lightRequestView, nil
}
