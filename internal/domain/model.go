package domain

import "ris/pkg/utills"

type RawLaureate struct {
	Id         string  `json:"id"`
	FirstName  string  `json:"firstname"`
	Surname    *string `json:"surname"`
	Motivation string  `json:"motivation"`
	Share      string  `json:"share"`
}

func (r *RawLaureate) ToLaureate() Laureate {
	return Laureate{
		Id:         int32(utills.ParseStringToInt(r.Id)),
		Firstname:  r.FirstName,
		Surname:    utills.GetVal(r.Surname),
		Motivation: r.Motivation,
		Share:      int32(utills.ParseStringToInt(r.Share)),
	}
}

type Laureate struct {
	Id         int32
	Firstname  string
	Surname    string
	Motivation string
	Share      int32
}

type RawPrize struct {
	Year              string        `json:"year"`
	Category          string        `json:"category"`
	Laureates         []RawLaureate `json:"laureates"`
	OverallMotivation *string       `json:"overallMotivation"`
}

func (r *RawPrize) ToPrize() Prize {
	laureates := make([]Laureate, 0, len(r.Laureates))
	for _, rawLaureate := range r.Laureates {
		laureates = append(laureates, rawLaureate.ToLaureate())
	}
	return Prize{
		Year:              r.Year,
		Category:          r.Category,
		Laureates:         laureates,
		OverallMotivation: utills.GetVal(r.OverallMotivation),
	}
}

type Prize struct {
	Year              string
	Category          string
	Laureates         []Laureate
	OverallMotivation string
}

type NobelResponse struct {
	Prizes []RawPrize `json:"prizes"`
}

func (r *NobelResponse) ToPrizes() []Prize {
	prizes := make([]Prize, 0, len(r.Prizes))
	for _, rawPrize := range r.Prizes {
		prizes = append(prizes, rawPrize.ToPrize())
	}
	return prizes
}
