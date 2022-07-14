package models

type BTC struct {
	Id    int64   `json:"id"`
	Time  string  `json:"time"`
	Value float64 `json:"value"`
}
