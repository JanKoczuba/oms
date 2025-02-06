package main

import (
	common "github.com/JanKoczuba/commons"
	pb "github.com/JanKoczuba/commons/api"
	"net/http"
)

type handler struct {
	client pb.OrderServiceClient
}

func NewHandler(client pb.OrderServiceClient,
) *handler {
	return &handler{client}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customents/{customerID}/orders", func(writer http.ResponseWriter, request *http.Request) {})
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")

	var items []*pb.ItemsWithQuantity
	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerId: customerID,
		Items:      items,
	})
}
