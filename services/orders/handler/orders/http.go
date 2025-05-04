package handler

import (
	"encoding/json"
	"fmt"
	"grpc-go/services/common/genproto/orders"
	"grpc-go/services/orders/types"
	"net/http"
)

type OrdersHttpHandler struct {
	orderService types.OrderService
}

func NewOrdersHttpHandler(orderService types.OrderService) *OrdersHttpHandler {
	return &OrdersHttpHandler{
		orderService: orderService,
	}
}

func (h *OrdersHttpHandler) RegisterRouter(router *http.ServeMux) {
	router.HandleFunc("POST /orders", h.CreateOrder)
}

func (h *OrdersHttpHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req orders.CreateOrderRequest

	if err := ParseJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	order := &orders.Order{
		OrderID:    42,
		CustomerID: req.CustomerId,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
	}

	if err := h.orderService.CreateOrder(r.Context(), order); err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}

	res := &orders.CreateOrderResponse{
		Status: "success",
	}
	WriteJSON(w, http.StatusOK, res)

}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(v)
}
