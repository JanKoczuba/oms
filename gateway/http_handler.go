package main

import (
	"errors"
	common "github.com/JanKoczuba/commons"
	pb "github.com/JanKoczuba/commons/api"
	"github.com/JanKoczuba/oms-gateway/gateway"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type handler struct {
	gateway gateway.OrdersGateway
}

func NewHandler(gateway gateway.OrdersGateway,
) *handler {
	return &handler{gateway}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.HandleCreateOrder)
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")

	var items []*pb.ItemsWithQuantity
	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())

		return
	}

	order, err := h.gateway.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})

	rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}
		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusOK, order)
}

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return common.ErrNoItems
	}

	for _, item := range items {
		if item.ID == "" {
			return errors.New("item id is required")
		}
		if item.Quantity <= 0 {
			return errors.New("item quantity is required")
		}
	}

	return nil
}
