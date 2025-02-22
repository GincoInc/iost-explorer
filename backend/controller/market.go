package controller

import (
	"net/http"

	"github.com/GincoInc/iost-explorer/backend/model"
	"github.com/labstack/echo"
)

func GetMarket(c echo.Context) error {
	marketInfo, err := model.GetMarketInfo()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, FormatResponse(marketInfo))
}
