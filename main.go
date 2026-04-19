package main

import (
	"korp/controller"
	"korp/db"
	"korp/repository"
	"korp/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	service := gin.Default()

	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	//Camada de repository
	invoiceRepository := repository.NewInvoiceRepository(dbConnection)
	//Camada usecase
	invoiceUseCase := usecase.NewInvoiceUseCase(invoiceRepository)
	//Camada de controllers
	invoiceController := controller.NewInvoiceController(invoiceUseCase)

	service.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	service.GET("/invoices", invoiceController.GetInvoices)
	service.POST("/invoice", invoiceController.CreateInvoice)

	service.Run(":8001")
}
