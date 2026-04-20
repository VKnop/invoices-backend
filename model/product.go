package model

type Product struct {
	ID          int    `json:"id"`
	CODE        string `json:"code"`
	DESCRIPTION string `json:"description"`
	BALANCE     int    `json:"balance"`
}
