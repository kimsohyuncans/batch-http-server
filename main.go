package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	mutex sync.Mutex
	timer *time.Timer
	value map[string]interface{} = map[string]interface{}{}
)

var queue = []int64{}
var requestTotalCount int64
var batchReport []string
var notifier = make(chan bool)

func main() {
	e := echo.New()
	e.GET("/summary", handlerfunc)
	e.GET("/stats", getStatsHandler)
	e.Logger.Fatal(e.Start(":1323"))
}

func doExpensiveMemoryCPUWork() {
	mutex.Lock()
	// push result per request to value
	for _, v := range queue {
		value[fmt.Sprint(v)] = "pending"
	}

	// do the summary
	time.Sleep(1 * time.Second)

	// push result per request to value
	for _, v := range queue {
		value[fmt.Sprint(v)] = "done"
	}

	batchReport = append(batchReport, fmt.Sprintf("batch %d done proccesed %d request", len(batchReport)+1, len(queue)))

	// reset the queue cause all the request is proccesed
	queue = nil
	timer = nil
	mutex.Unlock()
	notifier <- true

	log.Println("done doing expensive work")

}

func handlerfunc(c echo.Context) error {

	requestTotalCount++
	requestID := time.Now().UnixMicro()
	mutex.Lock()
	queue = append(queue, requestID)
	mutex.Unlock()

	// the first request will schedule the doExpensiveMemoryCPUWork for next 1 second
	if timer == nil {
		timer = time.AfterFunc(1*time.Second, doExpensiveMemoryCPUWork)
	}

	// wail until the result is reserved in value
	for {
		time.Sleep(300 * time.Millisecond)
		mutex.Lock()
		if status := value[fmt.Sprint(requestID)]; status == "done" {
			mutex.Unlock()
			break
		}
		mutex.Unlock()
	}

	delete(value, fmt.Sprint(requestID))
	return c.JSON(200, map[string]interface{}{
		"requestID": requestID,
		"result":    value[fmt.Sprint(requestID)],
	})
}

func getStatsHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"queue":             queue,
		"value":             value,
		"requestTotalCount": requestTotalCount,
		"batchReport":       batchReport,
	})
}
