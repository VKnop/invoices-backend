package model

type ProductInvoice struct {
	ID_PRODUCT int `json:"product_id"`
	ID_INVOICE int `json:"invoice_id"`
	SALE       int `json:"sale"`
}
