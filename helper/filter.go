package helper

import (
	"fmt"
	"restapi-bus/models/request"
)

func addFilterLimit(limit int, result string) string {

	if limit != 0 {
		result += fmt.Sprintf(` LIMIT %d`, limit)
	}
	return result
}

func RequestFilterBusToString(request *request.BusFilter) (result string) {

	if request.FrontNumberPlate != "" {
		result += " WHERE "
		result += fmt.Sprintf(`LEFT(number_plate,1) = "%s" `, request.FrontNumberPlate)
	}
	return
}

func RequestFilterAgencyToString(request *request.AgencyFilter) (result string) {
	result = "WHERE 1=1 "
	if request.Place != "" {
		result += fmt.Sprintf(`AND agency.place="%s"`, request.Place)
	}

	if request.Name != "" {

		result += fmt.Sprintf(`AND LEFT(agency.name,%d) = "%s"`, len(request.Name), request.Name)
	}

	if request.AboveBusCount+request.BelowBusCount != 0 {
		result = fmt.Sprintf(`LEFT JOIN bus on agency.agency_id = bus.agency_id %s GROUP BY agency.agency_id `, result)
		result += " HAVING 1=1 "

		if request.AboveBusCount != 0 {
			result += fmt.Sprintf(" AND COUNT(bus.bus_id) > %d", request.AboveBusCount)
		}
		if request.BelowBusCount != 0 {
			result += fmt.Sprintf(" AND COUNT(bus.bus_id) < %d", request.BelowBusCount)
		}

	}

	addFilterLimit(request.Limit, result)
	return
}

func RequestFilterCustomerToString(request *request.CustomerFilter) (result string) {
	result += "WHERE 1=1 "
	if request.Name != "" {
		result += fmt.Sprintf(` AND LEFT(name,%d) = "%s"`, len(request.Name), request.Name)
	}

	if request.FrontNumber != "" {
		result += fmt.Sprintf(` AND LEFT(phone_number,4) = "%s"`, request.FrontNumber)
	}

	addFilterLimit(request.Limit, result)
	return
}

func RequestFilterDriverToString(request *request.DriverFilter) (result string) {
	result += "WHERE 1 = 1 "
	if request.Name != "" {
		result += fmt.Sprintf(` AND LEFT(name,%d) = "%s"`, len(request.Name), request.Name)
	}

	addFilterLimit(request.Limit, result)
	return
}

func RequestFilterTicketToString(request *request.TicketFilter) (result string) {
	result += "WHERE 1 = 1 "

	if request.FromDate != "" && request.ToDate != "" {
		result += fmt.Sprintf(` AND date BETWEEN "%s" AND "%s" `, request.FromDate, request.ToDate)
	} else {

		if request.FromDate != "" {
			result += fmt.Sprintf(` AND date >  "%s" `, request.FromDate)
		}

		if request.ToDate != "" {
			result += fmt.Sprintf(` AND date >  "%s" `, request.FromDate)
		}
	}

	if request.OnDate != "" {
		result += fmt.Sprintf(` AND LEFT(date,%d) = "%s" `, len(request.OnDate), request.OnDate)
	}

	if request.DeparturePlace != "" {
		result += fmt.Sprintf(` AND LEFT(departure_place,%d) = "%s" `, len(request.DeparturePlace), request.DeparturePlace)
	}

	if request.ArrivalPlace != "" {
		result += fmt.Sprintf(` AND LEFT(arrival_place,%d) ="%s" `, len(request.ArrivalPlace), request.ArrivalPlace)
	}
	if request.Arrived != nil {

		result += fmt.Sprintf(` AND arrived = %v`, *request.Arrived)
	}

	if request.PriceBelow != 0 {
		result += fmt.Sprintf(` AND price <  %v `, request.PriceBelow)
	}

	if request.PriceAbove != 0 {
		result += fmt.Sprintf(` AND price >  "%v" `, request.PriceAbove)
	}

	addFilterLimit(request.Limit, result)
	return
}

func RequestFilterScheduleToString(request *request.ScheduleFilter) (result string) {
	result += " WHERE 1=1 "

	if request.Arrived != nil {
		result += fmt.Sprintf(" AND arrived = %v", request.Arrived)
	}

	if request.FromAgency != 0 {
		result += fmt.Sprintf(" AND to_agency = %d ", request.FromAgency)
	}

	if request.ToAgency != 0 {
		result += fmt.Sprintf(" AND from_agency = %d ", request.ToAgency)
	}

	if request.OnDate != "" {
		result += fmt.Sprintf(" AND date = %s ", request.OnDate)
	}

	if request.PriceBelow != 0 {
		result += fmt.Sprintf(` AND price <  %v `, request.PriceBelow)
	}

	if request.PriceAbove != 0 {
		result += fmt.Sprintf(` AND price >  "%v" `, request.PriceAbove)
	}

	addFilterLimit(request.Limit, result)
	return
}
