package controller2

import (
	"github.com/iost-official/explorer/backend/util/mail"
	"github.com/labstack/echo"
	"net/http"
	"log"
)

func SendMail(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	to := c.FormValue("email")
	content := c.FormValue("content")

	err := mail.SendMail(to, content)

	var errMsg string
	if err != nil {
		log.Println("SendMail error:", err)
		errMsg = err.Error()
	}

	return c.JSON(http.StatusOK, &CommOutput{0, errMsg})
}
