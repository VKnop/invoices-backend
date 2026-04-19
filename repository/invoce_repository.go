package repository

import (
	"database/sql"
	"fmt"
	"korp/model"
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
	query, err := ir.connection.Prepare("INSERT INTO invoices" +
		"(code, description, balance)" +
		" VALUES ($1, $2, $3) RETURNING id")
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

	query.Close()
	return id, nil
}
