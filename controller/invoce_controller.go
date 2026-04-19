package controller

import (
	"korp/model"
	"korp/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type invoiceController struct {
	invoiceUseCase usecase.InvoiceUsecase
}

func NewInvoiceController(usecase usecase.InvoiceUsecase) invoiceController {
	return invoiceController{
		invoiceUseCase: usecase,
	}
}

func (i *invoiceController) GetInvoices(ctx *gin.Context) {
	invoices, err := i.invoiceUseCase.GetInvoices()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, invoices)
}

func (i *invoiceController) CreateInvoice(ctx *gin.Context) {

	var invoice model.Invoice
	err := ctx.BindJSON(&invoice)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	insertedInvoice, err := i.invoiceUseCase.CreateInvoice(invoice)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, insertedInvoice)
}
