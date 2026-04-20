package model

type InvoiceResponse struct {
	ID             int    `json:"id"`
	CURRENT_STATUS string `json:"current_status"`
}
