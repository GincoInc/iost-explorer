package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/iost-official/explorer/backend/model"
	"github.com/iost-official/explorer/backend/model/db"
	"github.com/labstack/echo"
)

type AccountOutput struct {
	Name    string  `json:"address"`
	Balance float64 `json:"balance"`
	TxCount int     `json:"txCount"`
}

type AccountsOutput struct {
	AccountList []*AccountOutput `json:"accountList"`
	Page        int              `json:"page"`
	PagePrev    int              `json:"pagePrev"`
	PageNext    int              `json:"pageNext"`
	PageLast    int              `json:"pageLast"`
	TotalLen    int              `json:"totalLen"`
}

type AccountTxsOutput struct {
	Name string `json:"address"`
	// TxnList  []*db.JsonFlatTx `json:"txnList"`
	TxnLen   int `json:"txnLen"`
	PageLast int `json:"pageLast"`
}

func convertAccOutputs(acc []*db.Account) []*AccountOutput {
	var ret []*AccountOutput
	for _, a := range acc {
		txCount, err := db.GetAccountTxNumber(a.Name)
		if err != nil {
			log.Printf("get account tx number failed. account=%v, err=%v", a.Name, txCount)
		}
		ret = append(ret, &AccountOutput{
			Name:    a.Name,
			Balance: a.AccountInfo.Balance * 1e8,
			TxCount: txCount,
		})
	}
	return ret
}

func calLastPage(total int) int {
	var lastPage int
	if total%AccountEachPage == 0 {
		lastPage = total / AccountEachPage
	} else {
		lastPage = total/AccountEachPage + 1
	}

	if lastPage > AccountMaxPage { // ?
		lastPage = AccountMaxPage
	}
	return lastPage
}

func GetAccounts(c echo.Context) error {
	page := c.QueryParam("p")

	pageInt, _ := strconv.Atoi(page)
	if pageInt <= 0 {
		pageInt = 1
	}

	start := (pageInt - 1) * AccountEachPage
	accountList, err := db.GetAccounts(start, AccountEachPage)
	if err != nil {
		return err
	}

	accountTotalLen, err := db.GetAccountsTotalLen()
	if err != nil {
		return err
	}
	lastPage := calLastPage(accountTotalLen)

	output := &AccountsOutput{
		AccountList: convertAccOutputs(accountList),
		Page:        pageInt,
		PagePrev:    pageInt - 1,
		PageNext:    pageInt + 1,
		PageLast:    lastPage,
		TotalLen:    accountTotalLen,
	}

	return c.JSON(http.StatusOK, FormatResponse(output))
}

func GetAccountDetail(c echo.Context) error {
	// TODO 检查地址格式
	address := c.Param("id")
	if address == "" {
		return errors.New("Address required")
	}

	account, err := db.GetAccountByName(address)
	if err != nil {
		return err
	}
	accOutput := convertAccOutputs([]*db.Account{account})[0]

	marketInfo, err := model.GetMarketInfo()
	price, _ := strconv.ParseFloat(marketInfo.Price, 32)
	value := account.AccountInfo.Balance * price

	return c.JSON(http.StatusOK, FormatResponse(struct {
		*AccountOutput
		Value float64 `json:"value"`
		Code  string  `json:"code"`
	}{
		accOutput,
		value,
		"",
	}))
}

func GetAccountTxs(c echo.Context) error {
	address := c.Param("id")
	if address == "" {
		return errors.New("address requied")
	}

	page := c.QueryParam("p")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		pageInt = 1
	}

	start := (pageInt - 1) * AccountEachPage
	_, err = db.GetAccountTxByName(address, start, AccountEachPage)
	if err != nil {
		return err
	}

	totalLen, err := db.GetAccountTxNumber(address)
	if err != nil {
		log.Printf("get account tx number failed. account=%v, err=%v", address, err)
	}

	pageLast := calLastPage(totalLen)

	output := &AccountTxsOutput{
		address,
		// txnList,
		totalLen,
		pageLast,
	}

	return c.JSON(http.StatusOK, FormatResponse(output))
}
