//go:build !solution

package hotelbusiness

import (
	"sort"
)

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	if len(guests) == 0 {
		return []Load{}
	}

	events := make([]Load, 0, 2*len(guests))
	for _, guest := range guests {
		events = append(events, Load{StartDate: guest.CheckInDate, GuestCount: 1})
		events = append(events, Load{StartDate: guest.CheckOutDate, GuestCount: -1})
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].StartDate == events[j].StartDate {
			return events[i].GuestCount < events[j].GuestCount
		}
		return events[i].StartDate < events[j].StartDate
	})

	GuestNumPerDay := make([]Load, 0, len(guests))
	PrevCnt, PrevDate := events[0].GuestCount, events[0].StartDate
	CurrCnt, CurrDate := PrevCnt, 0
	PrevDataCnt := 0
	for i := 1; i < len(events); i++ {
		CurrDate = events[i].StartDate
		CurrCnt += events[i].GuestCount
		if PrevDate != CurrDate {
			if PrevDataCnt != PrevCnt {
				GuestNumPerDay = append(
					GuestNumPerDay,
					Load{StartDate: PrevDate, GuestCount: PrevCnt},
				)
				PrevDataCnt = PrevCnt
			}
			PrevDate = CurrDate
		}
		PrevCnt = CurrCnt
	}
	GuestNumPerDay = append(GuestNumPerDay, Load{StartDate: PrevDate, GuestCount: PrevCnt})
	return GuestNumPerDay
}
