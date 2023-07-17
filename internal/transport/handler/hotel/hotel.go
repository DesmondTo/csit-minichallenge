package hotel

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	fs "github.com/DesmondTo/minichallenge/internal/service/hotel"
)

type HotelHandler struct {
	hotelColl *mongo.Collection
}

func NewHandler(client *mongo.Client) *HotelHandler {
	return &HotelHandler{hotelColl: client.Database("minichallenge").Collection("hotels")}
}

func (fh *HotelHandler) GetCheapest(c *gin.Context) {
	checkInDate := c.Query("checkInDate")
	checkOutDate := c.Query("checkOutDate")
	destination := c.Query("destination")

	cid, cod, dest, err := checkAndParseQueryParams(checkInDate, checkOutDate, destination)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	coll := fh.hotelColl
	hotels, err := fs.GetCheapest(coll, cid, cod, dest)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, hotels)
}

func checkAndParseQueryParams(checkInDate, checkOutDate, destination string) (time.Time, time.Time, string, error) {
	checkInTime, err := parseDate(checkInDate)
	if err != nil || checkInDate == "" {
		return time.Time{}, time.Time{}, "", errors.New("Invalid or missing check-in date")
	}

	checkOutTime, err := parseDate(checkOutDate)
	if err != nil || checkOutDate == "" {
		return time.Time{}, time.Time{}, "", errors.New("Invalid or missing check-out date")
	}

	if destination == "" {
		return time.Time{}, time.Time{}, "", errors.New("Missing destination")
	}

  if checkInTime.After(checkOutTime) {
    return time.Time{}, time.Time{}, "", errors.New("Check-in date must be before or equal to check-out date")
  }

	return checkInTime, checkOutTime, destination, nil
}

func parseDate(dateString string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil || dateString == "" {
		return time.Time{}, err
	}

	return date, nil
}
