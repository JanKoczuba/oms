package main

import "net/http"

type handler struct {
	// gateway
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customents/{customerID}/orders", func(writer http.ResponseWriter, request *http.Request) {})

}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {

}
