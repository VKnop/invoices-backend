package repository

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"korp/model"
	"net/http"
)

type InvoiceRepository struct {
	connection *sql.DB
}

func NewInvoiceRepository(connection *sql.DB) InvoiceRepository {
	return InvoiceRepository{
		connection: connection,
	}
}

func (ir *InvoiceRepository) GetInvoices() ([]model.Invoice, error) {

	query := "SELECT id, current_status FROM invoices"
	rows, err := ir.connection.Query(query)

	if err != nil {
		fmt.Println(err)
		return []model.Invoice{}, err
	}

	var invoiceList []model.Invoice
	var invoice model.Invoice

	for rows.Next() {
		err = rows.Scan(
			&invoice.ID,
			&invoice.CURRENT_STATUS)

		if err != nil {
			fmt.Println(err)
			return []model.Invoice{}, err
		}

		invoiceList = append(invoiceList, invoice)
	}

	rows.Close()

	return invoiceList, nil
}

func (ir *InvoiceRepository) CreateInvoice(invoice model.Invoice) (int, error) {
	var id int
	query, err := ir.connection.Prepare("INSERT INTO invoices(current_status)" +
		" VALUES ($1) RETURNING id")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	err = query.QueryRow(
		invoice.CURRENT_STATUS).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	var productList []model.ProductInvoice = make([]model.ProductInvoice, len(invoice.PRODUCTS))
	copy(productList, invoice.PRODUCTS)

	_, err = ir.incrementPivotTable(productList, id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	query.Close()
	return id, nil
}

func (ir *InvoiceRepository) EditInvoice(id int, invoice model.Invoice) (model.Invoice, error) {

	query, err := ir.connection.Prepare("UPDATE invoices SET current_status = $1 WHERE id = $2 RETURNING id")
	if err != nil {
		fmt.Println(err)
		return model.Invoice{}, err
	}

	err = query.QueryRow(invoice.CURRENT_STATUS, id).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.Invoice{}, nil
		}
		return model.Invoice{}, err
	}

	newInvoice := model.Invoice{
		ID:             id,
		CURRENT_STATUS: invoice.CURRENT_STATUS}

	query.Close()

	return newInvoice, nil
}

func (ir *InvoiceRepository) incrementPivotTable(products []model.ProductInvoice, invoice_id int) (int, error) {

	var id int
	var product_invoice model.ProductInvoice

	for i := 0; i < len(products); i++ {
		product_invoice = products[i]

		product_invoice.ID_INVOICE = invoice_id
		query, err := ir.connection.Prepare("INSERT INTO product_invoice" +
			"(invoice_id, product_id, sale)" +
			"VALUES ($1,$2,$3) RETURNING id")
		if err != nil {
			fmt.Println("Esse produto não existe", err)
			return 0, err
		}

		updateProducts(product_invoice.ID_PRODUCT, product_invoice.SALE)
		err = query.QueryRow(product_invoice.ID_INVOICE, product_invoice.ID_PRODUCT, product_invoice.SALE).Scan(&id)

		if err != nil {
			fmt.Println(err)
			return 0, err
		}

		query.Close()
	}
	return 0, nil

}

func updateProducts(id int, sale int) {

	url := fmt.Sprintf("http://localhost:8000/product/%v", id)
	req, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer req.Body.Close()

	if req.StatusCode != http.StatusOK {
		fmt.Printf("Retornou um erro: %d\n", req.StatusCode)
		return
	}

	respbody, _ := io.ReadAll(req.Body)
	var jsonStr = string(respbody)
	var product model.Product
	err = json.Unmarshal([]byte(jsonStr), &product)
	if err != nil {
		fmt.Println("Error:", err)
	}

	editProduct := model.Product{
		ID:          product.ID,
		CODE:        product.CODE,
		DESCRIPTION: product.DESCRIPTION,
		BALANCE:     (product.BALANCE - sale)}

	body, err := json.Marshal(editProduct)

	put, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		fmt.Println(err)
		return
	}
	put.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(put)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Retornou um erro: %d\n", resp.StatusCode)
		return
	}

}
