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
	ReqType    models.RqType
	Commission float32
}

var ExInfo = map[reflect.Type]ExchangeInfo{
	reflect.TypeOf(BookReq.BinanceBookParams{}): {Exchange: models.BINANCE, ReqType: models.ReqType.Book},
	reflect.TypeOf(BookReq.GateioBookParams{}):  {Exchange: models.GATEIO, ReqType: models.ReqType.Book},
	reflect.TypeOf(BookReq.HuobiBookParams{}):   {Exchange: models.HUOBI, ReqType: models.ReqType.Book},
	reflect.TypeOf(BookReq.OkxBookParams{}):     {Exchange: models.OKX, ReqType: models.ReqType.Book},
	reflect.TypeOf(BookReq.MexcBookParams{}):    {Exchange: models.MEXC, ReqType: models.ReqType.Book},
	reflect.TypeOf(BookReq.BybitBookParams{}):   {Exchange: models.BYBIT, ReqType: models.ReqType.Book},
	reflect.TypeOf(BookReq.KucoinBookParams{}):  {Exchange: models.KUCOIN, ReqType: models.ReqType.Book},
	reflect.TypeOf(BookReq.CoinexBookParams{}):  {Exchange: models.COINEX, ReqType: models.ReqType.Book},

	reflect.TypeOf(TradeReq.BinanceTradeParams{}): {Exchange: models.BINANCE, ReqType: models.ReqType.Trade, Commission: 0.1},
	reflect.TypeOf(TradeReq.GateioTradeParams{}):  {Exchange: models.GATEIO, ReqType: models.ReqType.Trade, Commission: 0.1},
	reflect.TypeOf(TradeReq.OkxTradeParams{}):     {Exchange: models.OKX, ReqType: models.ReqType.Trade, Commission: 0.1},
	reflect.TypeOf(TradeReq.MexcTradeParams{}):    {Exchange: models.MEXC, ReqType: models.ReqType.Trade, Commission: 0.1},
	reflect.TypeOf(TradeReq.BybitTradeParams{}):   {Exchange: models.BYBIT, ReqType: models.ReqType.Trade, Commission: 0.1},
	reflect.TypeOf(TradeReq.CoinexTradeParams{}):  {Exchange: models.COINEX, ReqType: models.ReqType.Trade, Commission: 0.1},
	reflect.TypeOf(TradeReq.HuobiTradeParams{}):   {Exchange: models.HUOBI, ReqType: models.ReqType.Trade, Commission: 0.1},
	reflect.TypeOf(TradeReq.KucoinTradeParams{}):  {Exchange: models.KUCOIN, ReqType: models.ReqType.Trade, Commission: 0.1},

	reflect.TypeOf(OtherReq.CoinexTransferParams{}): {Exchange: models.COINEX, ReqType: models.ReqType.Transfer},

	reflect.TypeOf(ContractReq.GateioContractParams{}): {Exchange: models.GATEIO, ReqType: models.ReqType.Contract},
	reflect.TypeOf(ContractReq.OkxContractParams{}):    {Exchange: models.OKX, ReqType: models.ReqType.Contract},
	reflect.TypeOf(ContractReq.HuobiContractParams{}):  {Exchange: models.HUOBI, ReqType: models.ReqType.Contract},
}
