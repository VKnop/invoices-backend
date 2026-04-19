package usecase

import (
	"korp/model"
	"korp/repository"
)

type InvoiceUsecase struct {
	repository repository.InvoiceRepository
}

func NewInvoiceUseCase(repo repository.InvoiceRepository) InvoiceUsecase {
	return InvoiceUsecase{
		repository: repo,
	}
}

func (iu *InvoiceUsecase) GetInvoices() ([]model.Invoice, error) {
	return iu.repository.GetInvoices()
}

func (iu *InvoiceUsecase) CreateInvoice(invoice model.Invoice) (model.Invoice, error) {

	invoiceId, err := iu.repository.CreateInvoice(invoice)
	if err != nil {
		return model.Invoice{}, err
	}

	invoice.ID = invoiceId

	return invoice, nil
}
