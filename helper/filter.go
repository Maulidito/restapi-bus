package helper

import (
	"fmt"
	"restapi-bus/models/request"
)

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

	if request.Limit != 0 {
		result += fmt.Sprintf(`LIMIT %d`, request.Limit)
	}
	return
}
