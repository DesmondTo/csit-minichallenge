package flight

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/DesmondTo/minichallenge/internal/database/flight"
)

const dateFormat = "2006-01-02"

type Flight struct {
	City             string
	DepartureDate    string
	DepartureAirline string
	DeparturePrice   int32
	ReturnDate       string
	ReturnAirline    string
	ReturnPrice      int32
}

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

func GetCheapest(collection *mongo.Collection, departureDate time.Time, returnDate time.Time, destination string) ([]Flight, error) {
	sortOpts := bson.D{{Key: "price", Value: 1}}
	projOpts := bson.D{
		{Key: "airlinename", Value: 1},
		{Key: "price", Value: 1},
	}

	filters := bson.D{
		{Key: "srccity", Value: "Singapore"},
		{Key: "destcity", Value: destination},
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

	filters = bson.D{
		{Key: "srccity", Value: destination},
		{Key: "destcity", Value: "Singapore"},
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

	var flights []Flight
	dd := departureDate.Format(dateFormat)
	rd := returnDate.Format(dateFormat)
	for _, departFlight := range departFlights {
		for _, returnFlight := range returnFlights {
			flight := Flight{
				City:             destination,
				DepartureDate:    dd,
				DepartureAirline: departFlight.AirlineName,
				DeparturePrice:   departFlight.Price,
				ReturnDate:       rd,
				ReturnAirline:    returnFlight.AirlineName,
				ReturnPrice:      returnFlight.Price,
			}
			flights = append(flights, flight)
		}
	}

	return flights, nil
}
