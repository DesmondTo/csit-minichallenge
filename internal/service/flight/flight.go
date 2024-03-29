package flight

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/DesmondTo/minichallenge/internal/database/flight"
)

const dateFormat = "2006-01-02"

type FlightDetail map[string]interface{}

func getCheapestPrice(collection *mongo.Collection, filters bson.D) (int32, error) {
	sortOpts := bson.D{{Key: "price", Value: 1}}
	projOpts := bson.D{{Key: "price", Value: 1}}

	flights, err := flight.Get(collection, filters, sortOpts, projOpts)

	if err != nil {
		return -1, err
	}

	if len(flights) == 0 {
		return -1, nil
	}

	return flights[0].Price, nil
}

func GetCheapest(collection *mongo.Collection, departureDate time.Time, returnDate time.Time, destination string) ([]FlightDetail, error) {
	sortOpts := bson.D{{Key: "price", Value: 1}}
	projOpts := bson.D{
		{Key: "destcity", Value: 1},
		{Key: "airlinename", Value: 1},
		{Key: "price", Value: 1},
	}
	pattern := fmt.Sprintf("^%s$", destination)

	// Search for cheapest flight from Singapore to destination
	filters := bson.D{
		{Key: "srccity", Value: primitive.Regex{Pattern: "^Singapore$", Options: "i"}},
		{Key: "destcity", Value: primitive.Regex{Pattern: pattern, Options: "i"}},
		{Key: "date", Value: departureDate},
	}
	price, err := getCheapestPrice(collection, filters)
	if err != nil {
		return nil, err
	}
	filters = append(filters, bson.E{Key: "price", Value: bson.D{{Key: "$eq", Value: price}}})
	departFlights, err := flight.Get(collection, filters, sortOpts, projOpts)
	if err != nil {
		return nil, err
	}

	// Search for cheapest flight from destination to Singapore
	filters = bson.D{
		{Key: "srccity", Value: primitive.Regex{Pattern: pattern, Options: "i"}},
		{Key: "destcity", Value: primitive.Regex{Pattern: "^Singapore$", Options: "i"}},
		{Key: "date", Value: returnDate},
	}
	price, err = getCheapestPrice(collection, filters)
	if err != nil {
		return nil, err
	}
	filters = append(filters, bson.E{Key: "price", Value: bson.D{{Key: "$eq", Value: price}}})
	returnFlights, err := flight.Get(collection, filters, sortOpts, projOpts)
	if err != nil {
		return nil, err
	}

	var flightDetails []FlightDetail = make([]FlightDetail, 0)
	dd := departureDate.Format(dateFormat)
	rd := returnDate.Format(dateFormat)
	for _, departFlight := range departFlights {
		for _, returnFlight := range returnFlights {
			flightDetail := FlightDetail{
				"City":              departFlight.DestCity,
				"Departure Date":    dd,
				"Departure Airline": departFlight.AirlineName,
				"Departure Price":   departFlight.Price,
				"Return Date":       rd,
				"Return Airline":    returnFlight.AirlineName,
				"Return Price":      returnFlight.Price,
			}
			flightDetails = append(flightDetails, flightDetail)
		}
	}

	return flightDetails, nil
}
