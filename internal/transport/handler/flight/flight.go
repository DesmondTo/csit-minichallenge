package flight

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	fs "github.com/DesmondTo/minichallenge/internal/service/flight"
)

type FlightHandler struct {
	flightColl *mongo.Collection
}

func NewHandler(client *mongo.Client) *FlightHandler {
	return &FlightHandler{flightColl: client.Database("minichallenge").Collection("flights")}
}

func (fh *FlightHandler) GetCheapest(c *gin.Context) {
	departureDate := c.Query("departureDate")
	returnDate := c.Query("returnDate")
	destination := c.Query("destination")

	dd, rd, dest, err := checkAndParseQueryParams(departureDate, returnDate, destination)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	coll := fh.flightColl
	flights, err := fs.GetCheapest(coll, dd, rd, dest)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, flights)
}

func checkAndParseQueryParams(departureDate, returnDate, destination string) (time.Time, time.Time, string, error) {
	departureTime, err := parseDate(departureDate)
	if err != nil || departureDate == "" {
		return time.Time{}, time.Time{}, "", errors.New("Invalid or missing departure date")
	}

	returnTime, err := parseDate(returnDate)
	if err != nil || returnDate == "" {
		return time.Time{}, time.Time{}, "", errors.New("Invalid or missing return date")
	}

	if destination == "" {
		return time.Time{}, time.Time{}, "", errors.New("Missing destination")
	}

	return departureTime, returnTime, destination, nil
}

func parseDate(dateString string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil || dateString == "" {
		return time.Time{}, err
	}

	return date, nil
}
