package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/yourusername/crypto-exchange/internal/matching"
	"github.com/yourusername/crypto-exchange/internal/models"
	"github.com/yourusername/crypto-exchange/internal/ws"
)

type Server struct {
	router     *mux.Router
	wsServer   *ws.Server
	orderBooks map[string]*matching.OrderBook
	mu         sync.RWMutex
}

func NewServer() *Server {
	s := &Server{
		router:     mux.NewRouter(),
		wsServer:   ws.NewServer(),
		orderBooks: make(map[string]*matching.OrderBook),
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// API endpoints
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Public endpoints
	api.HandleFunc("/markets", s.getMarkets).Methods("GET")
	api.HandleFunc("/markets/{symbol}/orderbook", s.getOrderBook).Methods("GET")
	api.HandleFunc("/markets/{symbol}/trades", s.getTrades).Methods("GET")

	// Private endpoints (require authentication)
	private := api.PathPrefix("/private").Subrouter()
	private.Use(s.authMiddleware)

	private.HandleFunc("/orders", s.createOrder).Methods("POST")
	private.HandleFunc("/orders", s.getOrders).Methods("GET")
	private.HandleFunc("/orders/{id}", s.cancelOrder).Methods("DELETE")
	private.HandleFunc("/wallets", s.getWallets).Methods("GET")

	// WebSocket endpoint
	s.router.HandleFunc("/ws", s.wsServer.HandleWebSocket)
}

func (s *Server) Run(addr string) error {
	go s.wsServer.Run()
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) getMarkets(w http.ResponseWriter, r *http.Request) {
	markets := []models.MarketData{
		{
			Symbol:    "BTC/USDT",
			LastPrice: decimal.NewFromFloat(50000),
			Volume:    decimal.NewFromFloat(1000),
			High:      decimal.NewFromFloat(51000),
			Low:       decimal.NewFromFloat(49000),
			Change:    decimal.NewFromFloat(2.5),
			UpdatedAt: time.Now(),
		},
		// Добавьте другие торговые пары
	}

	json.NewEncoder(w).Encode(markets)
}

func (s *Server) getOrderBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := vars["symbol"]

	s.mu.RLock()
	orderBook, exists := s.orderBooks[symbol]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Market not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(orderBook.GetOrderBook(20))
}

func (s *Server) getTrades(w http.ResponseWriter, r *http.Request) {
	// Реализация получения истории сделок
}

func (s *Server) createOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	orderBook, exists := s.orderBooks[order.Symbol]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Market not found", http.StatusNotFound)
		return
	}

	trades, err := orderBook.PlaceOrder(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Отправляем обновления через WebSocket
	for _, trade := range trades {
		s.wsServer.BroadcastToChannel(order.Symbol, trade)
	}

	json.NewEncoder(w).Encode(trades)
}

func (s *Server) getOrders(w http.ResponseWriter, r *http.Request) {
	// Реализация получения списка ордеров пользователя
}

func (s *Server) cancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	// Реализация отмены ордера
}

func (s *Server) getWallets(w http.ResponseWriter, r *http.Request) {
	// Реализация получения списка кошельков пользователя
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Реализация проверки аутентификации
		next.ServeHTTP(w, r)
	})
} 