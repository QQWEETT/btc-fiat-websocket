package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"quest/internals/app/models"
	"quest/internals/app/processors"
	"strconv"
	"time"
)

type BTCHandler struct {
	processor *processors.BTCProcessor
}

func NewBTCHandler(processor *processors.BTCProcessor) *BTCHandler {
	handler := new(BTCHandler)
	handler.processor = processor
	return handler
}

func (handler *BTCHandler) CreateBTC() {
	for {
		time.Sleep(10 * time.Second)
		var newBTC models.BTC
		req, err := http.Get("https://api.kucoin.com/api/v1/market/stats?symbol=BTC-USDT")
		if err != nil {
			fmt.Println(err)
		}
		defer req.Body.Close()
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
		}

		var objmap map[string]json.RawMessage
		err = json.Unmarshal(b, &objmap)
		if err != nil {
			fmt.Println(err)
		}
		var internalMap map[string]string
		json.Unmarshal(objmap["data"], &internalMap)
		newBTC.Value, err = strconv.ParseFloat(internalMap["buy"], 64)
		if err != nil {
			fmt.Println(err)
		}
		date := time.Now()
		newBTC.Time = date.Format("2006-01-02 15:04:05")

		err = handler.processor.CreateBTC(newBTC)
		if err != nil {
			log.Println(err)
			return
		}

	}

}

func (handler *BTCHandler) LastBTC(w http.ResponseWriter, r *http.Request) {
	list, err := handler.processor.LastBTC()

	if err != nil {
		WrapError(w, err)
	}
	var m = map[string]interface{}{
		"history": list,
	}
	WrapOK(w, m)
}

func (handler *BTCHandler) FindBTC(w http.ResponseWriter, r *http.Request) {
	var filter models.Filter
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		WrapError(w, err)
		return
	}

	s, err := handler.processor.ListBTC(filter)

	if err != nil {
		WrapError(w, err)
		return
	}
	var m = map[string]interface{}{
		"result": "OK",
		"data":   s,
	}
	WrapOK(w, m)
}
