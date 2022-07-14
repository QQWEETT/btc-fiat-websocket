package handlers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/aglyzov/charmap"
	"io/ioutil"
	"log"
	"net/http"
	"quest/internals/app/models"
	"quest/internals/app/processors"
	"time"
)

type FiatHandler struct {
	processor *processors.FiatProcessor
}

func NewFiatHandler(processor *processors.FiatProcessor) *FiatHandler {
	handler := new(FiatHandler)
	handler.processor = processor
	return handler
}

func (handler *FiatHandler) CreateFiat() {
	for {

		resp, _ := http.Get("http://www.cbr.ru/scripts/XML_daily.asp")
		defer resp.Body.Close()
		byts, _ := ioutil.ReadAll(resp.Body)
		res1 := bytes.Trim(byts, "<? ?>")
		res2 := bytes.Replace(res1, []byte(","), []byte("."), 100)
		s := charmap.CP1251_to_UTF8(res2)
		var fiat models.Fiat
		xml.Unmarshal(s, &fiat)
		err := handler.processor.CreateFiat(fiat)
		if err != nil {
			log.Println(err)
			return
		}
		time.Sleep(24 * time.Hour)

	}

}

func (handler *FiatHandler) LastFiat(w http.ResponseWriter, r *http.Request) {
	list, err := handler.processor.LastFiat()

	if err != nil {
		WrapError(w, err)
	}

	var m = map[string]interface{}{
		"data": list,
	}

	WrapOK(w, m)
}

func (handler *FiatHandler) FindFiat(w http.ResponseWriter, r *http.Request) {
	var filter models.Filter
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		WrapError(w, err)
		return
	}

	s, err := handler.processor.ListFiat(filter)
	if err != nil {
		WrapError(w, err)
		return
	}
	var m = map[string]interface{}{
		"data": s,
	}
	WrapOK(w, m)
}
