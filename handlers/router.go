package handlers

import "net/http"

func New(handler TaskHandler, mw *Middleware) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/register", handler.Register)
	mux.HandleFunc("POST /auth/login", handler.Login)

	mux.Handle("POST /users/me", Auth(http.HandlerFunc(handler.Get)))
	mux.Handle("PUT /users/me", Auth(http.HandlerFunc(handler.UpdateProfile)))
	mux.Handle("DELETE /users/me", Auth(http.HandlerFunc(handler.DeleteProfile)))
	mux.Handle("POST /orders", Auth(http.HandlerFunc(handler.CreateOrder)))
	mux.Handle("GET /orders", Auth(http.HandlerFunc(handler.GetMyOrders)))
	mux.Handle("GET /orders/{id}", Auth(http.HandlerFunc(handler.GetOrderByID)))
	mux.Handle("PUT /orders/{id}", Auth(http.HandlerFunc(handler.UpdateOrder)))
	mux.Handle("DELETE /orders/{id}", Auth(http.HandlerFunc(handler.DeleteOrder)))

	mux.Handle("GET /admin/users", mw.AuthAdmin(http.HandlerFunc(handler.GetAllUsers)))
	mux.Handle("GET /admin/orders", mw.AuthAdmin(http.HandlerFunc(handler.GetAllOrders)))
	mux.Handle("PATCH /admin/users/{id}/role", mw.AuthAdmin(http.HandlerFunc(handler.ChangeRole)))

	return mux
}
