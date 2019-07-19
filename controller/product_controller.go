package controller

import (
	"fmt"
	"ggfly/dao"
	"ggfly/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

func initProductController(group *gin.RouterGroup) {
	controller := group.Group("product")
	controller.GET("listOfProducts", listOfProducts)
	controller.POST("saveProduct", saveProduct)
	controller.POST("saveProducts", saveProducts)
	controller.POST("removeProduct", removeProduct)
	controller.POST("updateProduct", updateProduct)
}

func listOfProducts(c *gin.Context) {
	var request dao.ProductListRequest
	c.JSON(http.StatusOK, request)
}

func saveProducts(c *gin.Context) {
	var products []*dao.Product
	var err error

	if bind := tryBind(c, &products); !bind {
		return
	}

	count := len(products)
	if count == 0 {
		writeFailedResponse(c, dao.StatusParameterError, errors.New("产品列表为空"))
		return
	}

	for i, product := range products {
		err = product.SaveValid()
		if err != nil {
			msg := fmt.Sprintf("第%d个 product 参数错误", i+1)
			writeFailedResponse(c, dao.StatusParameterError, errors.WithMessage(err, msg))
			return
		}
	}

	err = service.Share().ProductService.SaveProducts(getRequestId(c), products)
	if err != nil {
		writeFailedResponse(c, dao.StatusMySqlError, err)
		return
	}

	writeSuccessResponse(c, nil)
}

func saveProduct(c *gin.Context) {
	var product dao.Product
	var err error

	if bind := tryBind(c, &product); !bind {
		return
	}

	if err = product.SaveValid(); err != nil {
		writeFailedResponse(c, dao.StatusParameterError, err)
		return
	}

	err = service.Share().ProductService.SaveProduct(getRequestId(c), &product)
	if err != nil {
		writeFailedResponse(c, dao.StatusMySqlError, err)
		return
	}

	writeSuccessResponse(c, product)
}

func updateProduct(c *gin.Context) {
	var product dao.Product
	var err error

	if bind := tryBind(c, &product); !bind {
		return
	}

	if err = product.UpdateValid(); err != nil {
		writeFailedResponse(c, dao.StatusParameterError, err)
		return
	}

	err = service.Share().ProductService.UpdateProduct(getRequestId(c), &product)
	if err != nil {
		writeFailedResponse(c, dao.StatusMySqlError, err)
		return
	}

	writeSuccessResponse(c, product)
}

func removeProduct(c *gin.Context) {
	var product dao.Product
	var err error

	if bind := tryBind(c, &product); !bind {
		return
	}

	if product.ProductID.IsZero() {
		writeFailedResponse(c, dao.StatusParameterError, errors.New("productId不能为空"))
		return
	}

	err = service.Share().ProductService.RemoveProduct(getRequestId(c), &product)
	if err != nil {
		writeFailedResponse(c, dao.StatusMySqlError, err)
		return
	}

	writeSuccessResponse(c, nil)
}
