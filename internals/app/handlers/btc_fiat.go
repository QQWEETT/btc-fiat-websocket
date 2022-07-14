package handlers

import (
	"bytes"
	"encoding/xml"
	"github.com/aglyzov/charmap"
	"io/ioutil"
	"log"
	"net/http"
	"quest/internals/app/models"
	"quest/internals/app/processors"
)

type BTCFiatHandler struct {
	processor *processors.BTCFiatProcessor
}

func NewBTCFiatHandler(processor *processors.BTCFiatProcessor) *BTCFiatHandler {
	handler := new(BTCFiatHandler)
	handler.processor = processor
	return handler
}

func (handler *BTCFiatHandler) LastBTCFiat(w http.ResponseWriter, r *http.Request) {
	list, err := handler.processor.LastBTCFiat()

	if err != nil {
		WrapError(w, err)
	}

	var m = map[string]interface{}{
		"data": list,
	}

	WrapOK(w, m)
}

func (handler *BTCFiatHandler) CreateBTCFiat() {
	for {
		resp, _ := http.Get("http://www.cbr.ru/scripts/XML_daily.asp")
		defer resp.Body.Close()
		byts, _ := ioutil.ReadAll(resp.Body)
		res1 := bytes.Trim(byts, "<? ?>")
		res2 := bytes.Replace(res1, []byte(","), []byte("."), 100)
		s := charmap.CP1251_to_UTF8(res2)
		var fiat models.Fiat
		xml.Unmarshal(s, &fiat)
		err := handler.processor.CreateBTCFiat(fiat)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (handler *BTCFiatHandler) SocketHandler(w http.ResponseWriter, r *http.Request) {

	list, err := handler.processor.BTCFiatForSocket()
	if err != nil {
		log.Println(err)
	}
	var m = map[string]interface{}{
		"data": list,
	}

	WrapOK(w, m)
}
