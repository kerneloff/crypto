package matching

import (
	"container/heap"
	"errors"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/yourusername/crypto-exchange/internal/models"
)

var (
	ErrInvalidOrder     = errors.New("invalid order")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidPrice     = errors.New("invalid price")
	ErrInvalidQuantity  = errors.New("invalid quantity")
	ErrOrderNotFound    = errors.New("order not found")
)

type OrderBook struct {
	symbol    string
	bids      *PriceHeap
	asks      *PriceHeap
	mu        sync.RWMutex
	tradeChan chan models.Trade
	lastPrice decimal.Decimal
	lastUpdate time.Time
}

type PriceHeap struct {
	orders []*Order
	isBid  bool
}

type Order struct {
	ID        uint
	UserID    uint
	Price     decimal.Decimal
	Quantity  decimal.Decimal
	Timestamp time.Time
	index     int
}

func NewOrderBook(symbol string) *OrderBook {
	return &OrderBook{
		symbol:    symbol,
		bids:      &PriceHeap{isBid: true},
		asks:      &PriceHeap{isBid: false},
		tradeChan: make(chan models.Trade, 1000),
		lastUpdate: time.Now(),
	}
}

func (ob *OrderBook) validateOrder(order models.Order) error {
	if order.Quantity.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidQuantity
	}

	if order.Type == "limit" && order.Price.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidPrice
	}

	if order.Side != "buy" && order.Side != "sell" {
		return ErrInvalidOrder
	}

	if order.Type != "limit" && order.Type != "market" {
		return ErrInvalidOrder
	}

	return nil
}

func (ob *OrderBook) PlaceOrder(order models.Order) ([]models.Trade, error) {
	if err := ob.validateOrder(order); err != nil {
		return nil, err
	}

	ob.mu.Lock()
	defer ob.mu.Unlock()

	var trades []models.Trade
	remainingQuantity := order.Quantity

	if order.Side == "buy" {
		for remainingQuantity.GreaterThan(decimal.Zero) && ob.asks.Len() > 0 {
			bestAsk := ob.asks.orders[0]
			if order.Type == "limit" && bestAsk.Price.GreaterThan(order.Price) {
				break
			}

			tradeQuantity := decimal.Min(remainingQuantity, bestAsk.Quantity)
			trade := models.Trade{
				OrderID:  order.ID,
				Symbol:   ob.symbol,
				Price:    bestAsk.Price,
				Quantity: tradeQuantity,
				Side:     "buy",
			}

			trades = append(trades, trade)
			remainingQuantity = remainingQuantity.Sub(tradeQuantity)
			bestAsk.Quantity = bestAsk.Quantity.Sub(tradeQuantity)

			if bestAsk.Quantity.Equal(decimal.Zero) {
				heap.Pop(ob.asks)
			}

			ob.lastPrice = bestAsk.Price
			ob.lastUpdate = time.Now()
			ob.tradeChan <- trade
		}

		if remainingQuantity.GreaterThan(decimal.Zero) && order.Type == "limit" {
			heap.Push(ob.bids, &Order{
				ID:        order.ID,
				UserID:    order.UserID,
				Price:     order.Price,
				Quantity:  remainingQuantity,
				Timestamp: time.Now(),
			})
		}
	} else {
		for remainingQuantity.GreaterThan(decimal.Zero) && ob.bids.Len() > 0 {
			bestBid := ob.bids.orders[0]
			if order.Type == "limit" && bestBid.Price.LessThan(order.Price) {
				break
			}

			tradeQuantity := decimal.Min(remainingQuantity, bestBid.Quantity)
			trade := models.Trade{
				OrderID:  order.ID,
				Symbol:   ob.symbol,
				Price:    bestBid.Price,
				Quantity: tradeQuantity,
				Side:     "sell",
			}

			trades = append(trades, trade)
			remainingQuantity = remainingQuantity.Sub(tradeQuantity)
			bestBid.Quantity = bestBid.Quantity.Sub(tradeQuantity)

			if bestBid.Quantity.Equal(decimal.Zero) {
				heap.Pop(ob.bids)
			}

			ob.lastPrice = bestBid.Price
			ob.lastUpdate = time.Now()
			ob.tradeChan <- trade
		}

		if remainingQuantity.GreaterThan(decimal.Zero) && order.Type == "limit" {
			heap.Push(ob.asks, &Order{
				ID:        order.ID,
				UserID:    order.UserID,
				Price:     order.Price,
				Quantity:  remainingQuantity,
				Timestamp: time.Now(),
			})
		}
	}

	return trades, nil
}

func (ob *OrderBook) CancelOrder(orderID uint) error {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	// Поиск и удаление ордера из стакана
	for i, order := range ob.bids.orders {
		if order.ID == orderID {
			heap.Remove(ob.bids, i)
			ob.lastUpdate = time.Now()
			return nil
		}
	}

	for i, order := range ob.asks.orders {
		if order.ID == orderID {
			heap.Remove(ob.asks, i)
			ob.lastUpdate = time.Now()
			return nil
		}
	}

	return ErrOrderNotFound
}

func (ob *OrderBook) GetOrderBook(depth int) models.OrderBook {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	orderBook := models.OrderBook{
		Symbol: ob.symbol,
		Bids:   make([]models.OrderBookEntry, 0, depth),
		Asks:   make([]models.OrderBookEntry, 0, depth),
	}

	// Копируем биды
	for i := 0; i < depth && i < ob.bids.Len(); i++ {
		order := ob.bids.orders[i]
		orderBook.Bids = append(orderBook.Bids, models.OrderBookEntry{
			Price:    order.Price,
			Quantity: order.Quantity,
		})
	}

	// Копируем аски
	for i := 0; i < depth && i < ob.asks.Len(); i++ {
		order := ob.asks.orders[i]
		orderBook.Asks = append(orderBook.Asks, models.OrderBookEntry{
			Price:    order.Price,
			Quantity: order.Quantity,
		})
	}

	return orderBook
}

func (ob *OrderBook) GetLastPrice() decimal.Decimal {
	ob.mu.RLock()
	defer ob.mu.RUnlock()
	return ob.lastPrice
}

func (ob *OrderBook) GetLastUpdate() time.Time {
	ob.mu.RLock()
	defer ob.mu.RUnlock()
	return ob.lastUpdate
}

// Реализация heap.Interface для PriceHeap
func (h PriceHeap) Len() int { return len(h.orders) }

func (h PriceHeap) Less(i, j int) bool {
	if h.isBid {
		return h.orders[i].Price.GreaterThan(h.orders[j].Price)
	}
	return h.orders[i].Price.LessThan(h.orders[j].Price)
}

func (h PriceHeap) Swap(i, j int) {
	h.orders[i], h.orders[j] = h.orders[j], h.orders[i]
	h.orders[i].index = i
	h.orders[j].index = j
}

func (h *PriceHeap) Push(x interface{}) {
	n := len(h.orders)
	order := x.(*Order)
	order.index = n
	h.orders = append(h.orders, order)
}

func (h *PriceHeap) Pop() interface{} {
	old := h.orders
	n := len(old)
	order := old[n-1]
	old[n-1] = nil
	order.index = -1
	h.orders = old[0 : n-1]
	return order
} 