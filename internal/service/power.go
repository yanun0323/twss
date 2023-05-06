package service

import (
	"stocker/internal/model"
	"time"

	"github.com/shopspring/decimal"
)

var (
	_upperRatio  = decimal.NewFromFloat(0.8)
	_candleRatio = decimal.NewFromFloat(1.0)
	_lowerRatio  = decimal.NewFromFloat(0.8)
)

func (svc Service) CalculatePower(in model.PowerInput) (model.PowerOutput, error) {
	trades, err := svc.Repo.ListTrade(svc.Ctx, in.ID, in.From, in.To)
	if err != nil {
		return model.PowerOutput{}, err
	}

	temp := decimal.Zero
	if len(trades) != 0 {
		temp = trades[0].PriceOpen
	}

	sum := decimal.Zero
	power := make(map[time.Time]decimal.Decimal, len(trades))
	for _, trade := range trades {
		sum = sum.Add(svc.calculateTrade(trade, temp))
		temp = trade.PriceClose
		power[trade.Date] = sum
	}

	return model.PowerOutput{
		PowerInput: in,
		Power:      power,
	}, nil
}

// 計算交易
func (svc Service) calculateTrade(trade model.Trade, lastClose decimal.Decimal) decimal.Decimal {
	open := trade.PriceOpen
	close := trade.PriceClose
	candlePower := close.Sub(open).Mul(_candleRatio)

	high, low := open, close
	if high.LessThan(low) {
		high, low = low, high
	}

	upperPower := decimal.Zero
	if trade.PriceHighest.GreaterThan(high) {
		upperPower = trade.PriceHighest.Sub(high).Mul(_upperRatio)
	}

	lowerPower := decimal.Zero
	if trade.PriceLowest.LessThan(low) {
		lowerPower = low.Sub(trade.PriceLowest).Mul(_lowerRatio)
	}

	// switch candlePower.Sign() {
	// case -1:
	// 	return candlePower.Add(lowerPower).Mul(trade.TradeShare)
	// case 1:
	// 	return candlePower.Sub(upperPower).Mul(trade.TradeShare)
	// }
	return candlePower.Sub(upperPower).Add(lowerPower).Mul(trade.TradeShare)
}
