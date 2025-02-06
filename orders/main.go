package main

import "context"

func main() {

	store := NewStore()
	svc := newService(store)

	svc.CreateOrder(context.Background())
}
