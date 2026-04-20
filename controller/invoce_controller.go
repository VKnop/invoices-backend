package controller

import (
	"korp/model"
	"korp/usecase"
	"net/http"
	"strconv"

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

	var invoiceResponse []model.InvoiceResponse
	for i := 0; i < len(invoices); i++ {
		newInvoice := model.InvoiceResponse{
			ID:             invoices[i].ID,
			CURRENT_STATUS: invoices[i].CURRENT_STATUS}
		invoiceResponse = append(invoiceResponse, newInvoice)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, invoiceResponse)
}

func (i *invoiceController) CreateInvoice(ctx *gin.Context) {

	var invoice model.Invoice
	err := ctx.BindJSON(&invoice)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	insertedInvoice, err := i.invoiceUseCase.CreateInvoice(invoice)

	invoiceResponse := model.InvoiceResponse{
		ID:             insertedInvoice.ID,
		CURRENT_STATUS: insertedInvoice.CURRENT_STATUS}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, invoiceResponse)
}

func (p *invoiceController) EditInvoice(ctx *gin.Context) {

	id, er := strconv.Atoi(ctx.Param("id"))
	if er != nil {
		response := model.Response{
			Message: "Id da nota precisa ser um numero",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var invoice model.Invoice
	err := ctx.BindJSON(&invoice)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	changedInvoice, errr := p.invoiceUseCase.EditInvoice(id, invoice)

	invoiceResponse := model.InvoiceResponse{
		ID:             changedInvoice.ID,
		CURRENT_STATUS: changedInvoice.CURRENT_STATUS}

	if errr != nil {
		ctx.JSON(http.StatusInternalServerError, errr)
		return
	}

	if changedInvoice.ID == (0) {
		ctx.JSON(http.StatusBadRequest, "Nota não existe")
		return
	}

	ctx.JSON(http.StatusOK, invoiceResponse)
}
