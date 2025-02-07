package app

import (
	"github.com/gin-gonic/gin"
	"hac/internal/app/ds"
	"net/http"
	"strconv"
)

// Ping godoc
// @Summary      Show hello text
// @Description  very very friendly response
// @Tags         Tests
// @Produce      json
// @Success      200  {object}  string
// @Router       /ping/{name} [get]
func (a *Application) Ping(gCtx *gin.Context) {
	name := gCtx.Param("name")
	gCtx.String(http.StatusOK, "Hello %s", name)
}

//type GoodsListResponse = []ds.Goods

// GetAll godoc
// @Summary      Show all rows in db
// @Description  Return all product and info about rows
// @Tags         Tests
// @Produce      json
// @Success      200  {object} []ds.Goods
// @Router       /goods [get]

func (a *Application) GetObjectsByFloor(gCtx *gin.Context) {
	floor := gCtx.Param("floor")
	floor_int, err := strconv.Atoi(floor)
	if err != nil {
		answer := AnswerJSON{Status: "error", Description: "cant convert id to int"}
		gCtx.IndentedJSON(http.StatusInternalServerError, answer)
		return
	}
	objects, err := a.repo.GetObjectsByFloor(floor_int)
	if err != nil {
		answer := AnswerJSON{Status: "error", Description: "cant get all rows"}
		gCtx.IndentedJSON(http.StatusInternalServerError, answer)
		return
	}
	gCtx.IndentedJSON(http.StatusOK, objects)

}

// GetProduct godoc
// @Summary      Show product info by id
// @Description  Return all info of one product by id
//@Parameters	id
// @Tags         Tests
// @Produce      json
// @Success      200  {object}  ds.Goods
// @Router       /goods/{id} [get]
func (a *Application) GetObjectById(gCtx *gin.Context) {
	id_product := gCtx.Param("id")
	id_product_int, err := strconv.Atoi(id_product)
	if err != nil {
		answer := AnswerJSON{Status: "error", Description: "cant convert id to int"}
		gCtx.IndentedJSON(http.StatusInternalServerError, answer)
		return
	}
	product, err := a.repo.GetProductByID(uint(id_product_int))
	if err != nil {
		answer := AnswerJSON{Status: "error", Description: "cant get product by id"}
		gCtx.IndentedJSON(http.StatusInternalServerError, answer)
		return
	}
	gCtx.IndentedJSON(http.StatusOK, &product)
}

// ChangePrice godoc
// @Summary      Change price of product by id
// @Description  Change price of product by id. Price can't be 0
// @Tags         Tests
// @Produce      json
// @Success      200  {object}  ds.Goods
// @Router       /goods/{id} [put]
func (a *Application) ChangePrice(gCtx *gin.Context) {
	var params ds.Goods
	err := gCtx.BindJSON(&params)
	if err != nil {
		answer := AnswerJSON{Status: "error", Description: "cant parse json, "}
		gCtx.IndentedJSON(http.StatusRequestedRangeNotSatisfiable, answer)
		return
	}

	id_product := gCtx.Param("id_good")
	id_product_int, err := strconv.Atoi(id_product)
	if err != nil {
		answer := AnswerJSON{Status: "error", Description: "id must be integer"}
		gCtx.IndentedJSON(http.StatusRequestedRangeNotSatisfiable, answer)
		return
	}

	if params.Price == 0 {
		answer := AnswerJSON{Status: "error", Description: "product cant cost 0"}
		gCtx.IndentedJSON(http.StatusRequestedRangeNotSatisfiable, answer)
		return
	}

	err = a.repo.ChangeProduct(uint(id_product_int), params.Price)
	if err != nil {
		answer := AnswerJSON{Status: "error", Description: "cant change price"}
		gCtx.IndentedJSON(http.StatusInternalServerError, answer)
		return
	}
	product, err := a.repo.GetProductByID(uint(id_product_int))
	if err != nil {
		answer := AnswerJSON{Status: "error", Description: "cant get product by id"}
		gCtx.IndentedJSON(http.StatusInternalServerError, answer)
		return
	}
	gCtx.IndentedJSON(http.StatusOK, &product)
}

// PostProduct godoc
// @Summary      Add new row
// @Description  add new row with parameters in json
// @Tags         Tests
// @Produce      json
// @Success      200  {object}  ds.Goods
// @Router       /goods [post]
func (a *Application) PostProduct(gCtx *gin.Context) {
	var params ds.Goods
	err := gCtx.BindJSON(&params)
	if err != nil {
		answer := AnswerJSON{Status: "error", Description: "cant parse json"}
		gCtx.IndentedJSON(http.StatusRequestedRangeNotSatisfiable, answer)
		return
	}
	if params.Price <= 0 {
		answer := AnswerJSON{Status: "error", Description: "price cant be <= 0"}
		gCtx.IndentedJSON(http.StatusRequestedRangeNotSatisfiable, answer)
	}
	if params.Quantity <= 0 {
		answer := AnswerJSON{Status: "error", Description: "quantity cant be <= 0"}
		gCtx.IndentedJSON(http.StatusRequestedRangeNotSatisfiable, answer)
	}
	err = a.repo.CreateProduct(&params)
	if err != nil {
		answer := AnswerJSON{Status: "error", Description: "cant create product row"}
		gCtx.IndentedJSON(http.StatusInternalServerError, answer)
		return
	}
	gCtx.IndentedJSON(http.StatusOK, params)
}

// DeleteProduct godoc
// @Summary      Delete row by id
// @Description  Delete row by id. If there is not this id return error
// @Tags         Tests
// @Produce      json
// @Success      200  {object}  AnswerJSON
// @Router       /goods/{id} [delete]
func (a *Application) DeleteProduct(gCtx *gin.Context) {
	id_product := gCtx.Param("id")
	id_product_int, err := strconv.Atoi(id_product)
	if err != nil {
		answer := AnswerJSON{Status: "error", Description: "id must be integer"}
		gCtx.IndentedJSON(http.StatusRequestedRangeNotSatisfiable, answer)
		return
	}
	err = a.repo.DeleteProduct(uint(id_product_int))
	if err.Error() == "record not found" {
		answer := AnswerJSON{Status: "error", Description: "id not found"}
		gCtx.IndentedJSON(http.StatusRequestedRangeNotSatisfiable, answer)
		return
	}
	if err != nil {
		answer := AnswerJSON{Status: "error", Description: "cant delete row"}
		gCtx.IndentedJSON(http.StatusInternalServerError, answer)
		return
	}
	answer := AnswerJSON{Status: "successful", Description: "row was deleted"}
	gCtx.IndentedJSON(http.StatusOK, answer)
}

func (a *Application) AddBasketRow(gCtx *gin.Context) {
	var params ds.Basket
	err := gCtx.BindJSON(&params)
	if err != nil {
		answer := AnswerJSON{Status: "error", Description: "cant parse json"}
		gCtx.IndentedJSON(http.StatusRequestedRangeNotSatisfiable, answer)
		return
	}
	err = a.repo.CreateBasketRow(&params)
	if err != nil {
		answer := AnswerJSON{Status: "error", Description: "good cant be added to basket"}
		gCtx.IndentedJSON(http.StatusInternalServerError, answer)
		return
	}
	answer := AnswerJSON{Status: "successful", Description: "good was added to basket"}
	gCtx.IndentedJSON(http.StatusOK, answer)
}

func (a *Application) GetBasket(gCtx *gin.Context) {
	var params ds.Users
	err := gCtx.BindJSON(&params)
	if err != nil {
		answer := AnswerJSON{Status: "error", Description: "cant parse json"}
		gCtx.IndentedJSON(http.StatusRequestedRangeNotSatisfiable, answer)
		return
	}
	login_user := params.Login
	id_user, err := a.repo.GetIdByLogin(login_user)
	if err != nil {
		answer := AnswerJSON{Status: "error", Description: "cant find user with this login"}
		gCtx.IndentedJSON(http.StatusForbidden, answer)
		return
	}
	basket, err := a.repo.GetBasket(id_user)
	if err != nil {
		answer := AnswerJSON{Status: "error", Description: "cant get rows in basket"}
		gCtx.IndentedJSON(http.StatusInternalServerError, answer)
		return
	}
	gCtx.IndentedJSON(http.StatusOK, basket)
}
