package repository

import (
	"dia-backend/internal/app/ds"
	"fmt"

	"gorm.io/gorm"
)

type CalcRequestRepository struct {
	db *gorm.DB
}

func NewCalcRequestRepository(db *gorm.DB) *CalcRequestRepository {
	return &CalcRequestRepository{
		db: db,
	}
}

type CalcRequestViewEntry struct {
	Lamp   ds.Lamp
	AreaM2 float32
	Number int
}

type CalcRequestView struct {
	CalcRequest ds.CalcRequest
	Entries     []CalcRequestViewEntry
}

var calcRequests = []ds.CalcRequest{
	{
		ID:             1,
		MaxTotalPowerW: 200,
		TotalPowerW:    158,
	},
}

var calcRequestToLamps = []ds.CalcRequestToLamp{
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

func (*CalcRequestRepository) GetCalcRequestEntryCntByID(id int) (int, error) {
	if len(calcRequests) == 0 {
		return 0, fmt.Errorf("массив пустой")
	}

	var calcRequest *ds.CalcRequest = nil
	for _, req := range calcRequests {
		if req.ID == id {
			calcRequest = &req
		}
	}
	if calcRequest == nil {
		return 0, fmt.Errorf("не найдено")
	}

	var calcRequestEntryCnt int = 0
	for _, reqToLamp := range calcRequestToLamps {
		if reqToLamp.RequestID == calcRequest.ID {
			calcRequestEntryCnt++
		}
	}

	return calcRequestEntryCnt, nil
}

func (*CalcRequestRepository) GetCalcRequestViewByID(id int, lampRepo *LampRepository) (*CalcRequestView, error) {
	if len(calcRequests) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	var calcRequest *ds.CalcRequest = nil
	for _, req := range calcRequests {
		if req.ID == id {
			calcRequest = &req
		}
	}
	if calcRequest == nil {
		return nil, fmt.Errorf("не найдено")
	}

	calcRequestView := CalcRequestView{
		CalcRequest: *calcRequest,
	}

	for _, reqToLamp := range calcRequestToLamps {
		if reqToLamp.RequestID == calcRequest.ID {
			lamp, err := lampRepo.GetLampByID(reqToLamp.LampID)
			if err != nil {
				return nil, fmt.Errorf("не найдено")
			}
			calcRequestView.Entries = append(calcRequestView.Entries, CalcRequestViewEntry{
				Lamp:   *lamp,
				AreaM2: reqToLamp.AreaM2,
				Number: reqToLamp.Number,
			})
		}
	}

	return &calcRequestView, nil
}
