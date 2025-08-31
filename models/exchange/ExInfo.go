package exchange

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeReq/BookReq"
	"enchainer/models/exchange/exchangeReq/ContractReq"
	"enchainer/models/exchange/exchangeReq/OtherReq"
	"enchainer/models/exchange/exchangeReq/TradeReq"
	"reflect"
)

type ExchangeInfo struct {
	Exchange   models.Exchange
	ReqType    string // TODO
	Commission float32
}

var ExInfo = map[reflect.Type]ExchangeInfo{
	reflect.TypeOf(BookReq.BinanceBookParams{}): {Exchange: models.BINANCE, ReqType: "Book"},
	reflect.TypeOf(BookReq.GateioBookParams{}):  {Exchange: models.GATEIO, ReqType: "Book"},
	reflect.TypeOf(BookReq.HuobiBookParams{}):   {Exchange: models.HUOBI, ReqType: "Book"},
	reflect.TypeOf(BookReq.OkxBookParams{}):     {Exchange: models.OKX, ReqType: "Book"},
	reflect.TypeOf(BookReq.MexcBookParams{}):    {Exchange: models.MEXC, ReqType: "Book"},
	reflect.TypeOf(BookReq.BybitBookParams{}):   {Exchange: models.BYBIT, ReqType: "Book"},
	reflect.TypeOf(BookReq.KucoinBookParams{}):  {Exchange: models.KUCOIN, ReqType: "Book"},
	reflect.TypeOf(BookReq.CoinexBookParams{}):  {Exchange: models.COINEX, ReqType: "Book"},

	reflect.TypeOf(TradeReq.BinanceTradeParams{}): {Exchange: models.BINANCE, ReqType: "Trade", Commission: 0.1},
	reflect.TypeOf(TradeReq.GateioTradeParams{}):  {Exchange: models.GATEIO, ReqType: "Trade", Commission: 0.1},
	reflect.TypeOf(TradeReq.OkxTradeParams{}):     {Exchange: models.OKX, ReqType: "Trade", Commission: 0.1},
	reflect.TypeOf(TradeReq.MexcTradeParams{}):    {Exchange: models.MEXC, ReqType: "Trade", Commission: 0.1},
	reflect.TypeOf(TradeReq.BybitTradeParams{}):   {Exchange: models.BYBIT, ReqType: "Trade", Commission: 0.1},
	reflect.TypeOf(TradeReq.CoinexTradeParams{}):  {Exchange: models.COINEX, ReqType: "Trade", Commission: 0.1},
	reflect.TypeOf(TradeReq.HuobiTradeParams{}):   {Exchange: models.HUOBI, ReqType: "Trade", Commission: 0.1},
	reflect.TypeOf(TradeReq.KucoinTradeParams{}):  {Exchange: models.KUCOIN, ReqType: "Trade", Commission: 0.1},

	reflect.TypeOf(OtherReq.CoinexTransferParams{}): {Exchange: models.COINEX, ReqType: "Transfer"},

	reflect.TypeOf(ContractReq.GateioContractParams{}): {Exchange: models.GATEIO, ReqType: "Contract"},
	reflect.TypeOf(ContractReq.HuobiContractParams{}):  {Exchange: models.HUOBI, ReqType: "Contract"},
}
