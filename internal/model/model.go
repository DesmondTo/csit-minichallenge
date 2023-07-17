package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Flight struct {
	ID              primitive.ObjectID `bson:"_id"`
	Airline         string             `bson:"airline"`
	AirlineID       int32              `bson:"airline_id"`
	AirlineName     string             `bson:"airlinename"`
	Codeshare       string             `bson:"codeshare"`
	Date            time.Time          `bson:"date"`
	DestAirport     string             `bson:"destairport"`
	DestAirportID   int32              `bson:"destairportid"`
	DestAirportName string             `bson:"destairportname"`
	DestCity        string             `bson:"destcity"`
	DestCountry     string             `bson:"destcountry"`
	Eq              string             `bson:"eq"`
	Price           int32              `bson:"price"`
	SrcAirport      string             `bson:"srcairport"`
	SrcAirportID    int32              `bson:"srcairportid"`
	SrcAirportName  string             `bson:"srcairportname"`
	SrcCity         string             `bson:"srccity"`
	SrcCountry      string             `bson:"srccountry"`
	Stop            int32              `bson:"stop"`
}

type Hotel struct {
	ID    primitive.ObjectID `bson:"_id"`
	City  string             `bson:"city"`
	Date  time.Time          `bson:"date"`
	Name  string             `bson:"hotelName"`
	Price int32              `bson:"price"`
}
