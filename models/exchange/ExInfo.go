package exchange

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeReq/BookReq"
	"enchainer/models/exchange/exchangeReq/TradeReq"
	"reflect"
)

type ExchangeInfo struct {
	Exchange models.Exchange
	ReqType  string
}

var ExInfo = map[reflect.Type]ExchangeInfo{
	reflect.TypeOf(BookReq.BinanceBookParams{}): {models.BINANCE, "Book"},
	reflect.TypeOf(BookReq.GateioBookParams{}):  {models.GATEIO, "Book"},
	reflect.TypeOf(BookReq.HuobiBookParams{}):   {models.HUOBI, "Book"},
	reflect.TypeOf(BookReq.OkxBookParams{}):     {models.OKX, "Book"},
	reflect.TypeOf(BookReq.MexcBookParams{}):    {models.MEXC, "Book"},
	reflect.TypeOf(BookReq.BybitBookParams{}):   {models.BYBIT, "Book"},
	reflect.TypeOf(BookReq.KucoinBookParams{}):  {models.KUCOIN, "Book"},

	reflect.TypeOf(TradeReq.BinanceTradeParams{}): {models.BINANCE, "Trade"},
	reflect.TypeOf(TradeReq.GateioTradeParams{}):  {models.GATEIO, "Trade"},
	reflect.TypeOf(TradeReq.MexcTradeParams{}):    {models.HUOBI, "Trade"},
	reflect.TypeOf(TradeReq.BybitTradeParams{}):   {models.OKX, "Trade"},
}
