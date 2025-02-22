package controller

import (
	"net/http"
	"strconv"

	"github.com/GincoInc/iost-explorer/backend/model/db"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type SearchOutput struct {
	Search string `json:"search"`
	Type   string `json:"type"`
	Text   string `json:"text,omitempty"`
}

func GetSearch(c echo.Context) error {
	search := c.Param("id")
	if search == "" {
		return errors.New("Nothing to search")
	}

	output := &SearchOutput{
		Search: search,
	}

	account, _ := db.GetAccountByName(search)
	if account != nil {
		output.Type = "account"
		return c.JSON(http.StatusOK, FormatResponse(output))
	}

	tx, _ := db.GetTxByHash(search)
	if tx != nil {
		output.Type = "tx"
		return c.JSON(http.StatusOK, FormatResponse(output))
	}

	blkHash, _, _ := db.GetBlockByHash(search)
	if blkHash != nil {
		output.Type = "block"
		output.Text = strconv.FormatInt(blkHash.Number, 10)
		return c.JSON(http.StatusOK, FormatResponse(output))
	}

	if searchInt64, _ := strconv.ParseInt(search, 10, 64); searchInt64 > 0 {
		block, _ := db.GetBlockByHeight(searchInt64)
		if block != nil {
			output.Type = "block"
		}
	}
	return c.JSON(http.StatusOK, FormatResponse(output))
}
