package models

import "encoding/xml"

type ValCurs struct {
	Id      int     `json:"id"`
	Date    string  `json:"date"`
	Code    string  `xml:"CharCode" json:"code"`
	Nominal int     `xml:"Nominal"`
	Value   float64 `xml:"Value" json:"value"`
}
type Fiat struct {
	XMLName xml.Name  `xml:"ValCurs"`
	Valute  []ValCurs `xml:"Valute"`
}
