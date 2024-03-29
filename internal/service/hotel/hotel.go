package hotel

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/DesmondTo/minichallenge/internal/database/hotel"
	"github.com/DesmondTo/minichallenge/internal/model"
)

const dateFormat = "2006-01-02"

type HotelKey struct {
	City      string
	HotelName string
}
type HotelPrice map[HotelKey]int32
type HotelDetail map[string]interface{}

func sumHotelPrice(hotels []model.Hotel) HotelPrice {
	hotelPrice := make(HotelPrice)
	for _, hotel := range hotels {
		hotelKey := HotelKey{City: hotel.City, HotelName: hotel.Name}
		hotelPrice[hotelKey] += hotel.Price
	}

	return hotelPrice
}

func findCheapesPrice(hotelPrices HotelPrice) int32 {
	var price int32 = -1
	for _, hotelPrice := range hotelPrices {
		if price == -1 {
			price = hotelPrice
		} else if hotelPrice < price {
			price = hotelPrice
		}
	}

	return price
}

func GetCheapest(collection *mongo.Collection, checkInDate time.Time, checkOutDate time.Time, destination string) ([]HotelDetail, error) {
	sortOpts := bson.D{
		{Key: "hotelName", Value: 1},
		{Key: "date", Value: 1},
	}
	projOpts := bson.D{
		{Key: "city", Value: 1},
		{Key: "hotelName", Value: 1},
		{Key: "price", Value: 1},
		{Key: "date", Value: 1},
	}
	pattern := fmt.Sprintf("^%s$", destination)

  // Search for cheapest hotel in destination
	filters := bson.D{
		{Key: "city", Value: primitive.Regex{Pattern: pattern, Options: "i"}},
		{Key: "date", Value: bson.D{{Key: "$gte", Value: checkInDate}, {Key: "$lte", Value: checkOutDate}}},
	}
	hotels, err := hotel.Get(collection, filters, sortOpts, projOpts)
	if err != nil {
		return nil, err
	}

	hotelPrices := sumHotelPrice(hotels)
	cheapestPrice := findCheapesPrice(hotelPrices)

	var hotelDetails []HotelDetail = make([]HotelDetail, 0)
	for hotelKey, hotelPrice := range hotelPrices {
		if hotelPrice == cheapestPrice {
			hotelDetail := HotelDetail{
				"City":           hotelKey.City,
				"Check In Date":  checkInDate.Format(dateFormat),
				"Check Out Date": checkOutDate.Format(dateFormat),
				"Hotel":          hotelKey.HotelName,
				"Price":          hotelPrice,
			}
			hotelDetails = append(hotelDetails, hotelDetail)
		}
	}

	return hotelDetails, nil
}
