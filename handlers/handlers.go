package handlers

import (
	internalCtx "UserService/internal/context"
	"UserService/internal/models"
	"UserService/internal/service"
	"UserService/pkg/logger"
	"encoding/json"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type TaskHandler struct {
	logger       logger.Logger
	serviceUser  service.UserService
	serviceAdmin service.AdminService
	serviceOrder service.OrderService
}

func NewHandler(logger logger.Logger, serviceUser service.UserService, serviceAdmin service.AdminService, serviceOrder service.OrderService) TaskHandler {
	return TaskHandler{
		serviceUser:  serviceUser,
		serviceAdmin: serviceAdmin,
		serviceOrder: serviceOrder,
		logger:       logger,
	}
}

type registerReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (t *TaskHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request registerReq

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	log := t.logger.With(zap.String("handler", "Register"))

	err = t.serviceUser.Register(r.Context(), models.RegisterRequest{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		handleError(w, log, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "Created",
	})
}

type verifyReq struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func (t *TaskHandler) Verify(w http.ResponseWriter, r *http.Request) {

	var request verifyReq

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "error from Decode", http.StatusBadRequest)
		return
	}

	log := t.logger.With(zap.String("handler", "Verify"))

	id, err := t.serviceUser.Verify(r.Context(), models.RegisterRequest{
		Email: request.Email,
		OTP:   request.OTP,
	})
	if err != nil {
		log.Error("verify", zap.Error(err))
		handleError(w, log, err)
		return
	}
	_ = json.NewEncoder(w).Encode(id)
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (t *TaskHandler) Login(w http.ResponseWriter, r *http.Request) {
	var login loginReq

	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	log := t.logger.With(zap.String("handler", "Login"))

	accessToken, refreshToken, err := t.serviceUser.Login(r.Context(), models.LoginRequest{
		Email:    login.Email,
		Password: login.Password,
	})
	if err != nil {
		handleError(w, log, err)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

type logOutReq struct {
	RefreshToken string `json:"refreshToken"`
}

func (t *TaskHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	var request logOutReq

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "error from Decode", http.StatusBadRequest)
		return
	}

	log := t.logger.With(zap.String("handler", "LogOut"))

	err = t.serviceUser.LogOut(r.Context(), models.RefreshAccessTokens{
		RefreshToken: request.RefreshToken,
	})
	if err != nil {
		handleError(w, log, err)
		return
	}

	w.WriteHeader(http.StatusOK)

}

type refreshTokenReq struct {
	RefreshToken string `json:"refreshToken"`
}

func (t *TaskHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var request refreshTokenReq

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "error from Decode", http.StatusBadRequest)
		return
	}

	log := t.logger.With(zap.String("handler", "RefreshToken"))

	tokens, err := t.serviceUser.RefreshToken(r.Context(), models.HashToken{
		TokenHash: request.RefreshToken,
	})
	if err != nil {
		handleError(w, log, err)
		return
	}

	_ = json.NewEncoder(w).Encode(tokens)

}

func (t *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {

	userIDStr := r.Context().Value(internalCtx.UserIDKey)
	userID, ok := userIDStr.(int)
	if !ok {
		http.Error(w, "user not exist in context", http.StatusNotFound)
		return
	}

	log := t.logger.With(zap.String("handler", "Get"))

	user, err := t.serviceUser.Get(r.Context(), userID)
	if err != nil {
		handleError(w, log, err)
	}

	_ = json.NewEncoder(w).Encode(user)
}

type updateReq struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func (t *TaskHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var update updateReq

	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userIDStr := r.Context().Value(internalCtx.UserIDKey)
	userID, ok := userIDStr.(int)
	if !ok {
		http.Error(w, "user not exist in context", http.StatusBadRequest)
		return
	}

	if userID < 1 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	log := t.logger.With(zap.String("handler", "UpdateProfile"))

	err = t.serviceUser.UpdateProfile(r.Context(), userID, models.User{
		Name:  update.Name,
		Phone: update.Phone,
	})
	if err != nil {
		handleError(w, log, err)
		return
	}
}

func (t *TaskHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {

	userIDStr := r.Context().Value(internalCtx.UserIDKey)
	userID, ok := userIDStr.(int)
	if !ok {
		http.Error(w, "user not exist in context", http.StatusBadRequest)
		return
	}

	if userID < 1 {
		http.Error(w, "user not exist in context", http.StatusBadRequest)
		return
	}

	log := t.logger.With(zap.String("handler", "DeleteProfile"))

	err := t.serviceUser.DeleteProfile(r.Context(), userID)
	if err != nil {
		handleError(w, log, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type changePassReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (t *TaskHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var request changePassReq

	userIDStr := r.Context().Value(internalCtx.UserIDKey)
	userID, ok := userIDStr.(int)
	if !ok {
		http.Error(w, "user not found in context", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "error from Decode", http.StatusBadRequest)
		return
	}

	log := t.logger.With(zap.String("handler", "ChangePassword"))

	err = t.serviceUser.ChangePassword(r.Context(), userID, models.Password{
		OldPassword: request.OldPassword,
		NewPassword: request.NewPassword,
	})
	if err != nil {
		handleError(w, log, err)
		return
	}
}

func (t *TaskHandler) GetOrdersProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Context().Value(internalCtx.UserIDKey)
	userID, ok := userIDStr.(int)
	if !ok {
		http.Error(w, "user not found in context", http.StatusBadRequest)
		return
	}

	if userID < 1 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	log := t.logger.With(zap.String("handler", "GetOrdersProfile"))

	userOrder, err := t.serviceUser.GetOrdersProfile(r.Context(), userID)
	if err != nil {
		handleError(w, log, err)
		return
	}

	_ = json.NewEncoder(w).Encode(userOrder)
}

type orderReq struct {
	Product string `json:"product"`
	Price   int    `json:"price"`
	Status  string `json:"status"`
}

func (t *TaskHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order orderReq

	userIDValue := r.Context().Value(internalCtx.UserIDKey)

	userID, ok := userIDValue.(int)
	if !ok {
		http.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	log := t.logger.With(zap.String("handler", "CreateOrder"))

	err = t.serviceOrder.CreateOrder(r.Context(), models.Order{
		Product: order.Product,
		Price:   order.Price,
		UserID:  userID,
		Status:  order.Status,
	})
	if err != nil {
		handleError(w, log, err)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "order successfully created",
	})
}

func (t *TaskHandler) GetMyOrders(w http.ResponseWriter, r *http.Request) {

	userIDStr := r.Context().Value(internalCtx.UserIDKey)
	userID := userIDStr.(int)

	if userID < 1 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	log := t.logger.With(zap.String("handler", "GetMyOrders"))

	myOrders, err := t.serviceOrder.GetMyOrders(r.Context(), userID)
	if err != nil {
		handleError(w, log, err)
		return
	}

	_ = json.NewEncoder(w).Encode(myOrders)
}

func (t *TaskHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if id < 1 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	log := t.logger.With(zap.String("handler", "GetOrderByID"))

	order, err := t.serviceOrder.GetOrderByID(r.Context(), id)
	if err != nil {
		handleError(w, log, err)
		return
	}

	_ = json.NewEncoder(w).Encode(order)

}

type updateOrderReq struct {
	Product string `json:"product"`
	Price   int    `json:"price"`
	Status  string `json:"status"`
}

func (t *TaskHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var order updateOrderReq

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if id < 1 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	log := t.logger.With(zap.String("handler", "UpdateOrder"))

	err = t.serviceOrder.UpdateOrder(r.Context(), id, models.Order{
		Product: order.Product,
		Price:   order.Price,
		Status:  order.Status,
	})
	if err != nil {
		handleError(w, log, err)
		return
	}
}

func (t *TaskHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if id < 1 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	log := t.logger.With(zap.String("handler", "DeleteOrder"))

	err = t.serviceOrder.DeleteOrders(r.Context(), id)
	if err != nil {
		handleError(w, log, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (t *TaskHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	log := t.logger.With(zap.String("handler", "GetAllUsers"))

	users, err := t.serviceAdmin.GetAllUsers(r.Context())
	if err != nil {
		handleError(w, log, err)
		return
	}

	_ = json.NewEncoder(w).Encode(users)
}

func (t *TaskHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {

	log := t.logger.With(zap.String("handler", "GetAllOrders"))

	orders, err := t.serviceAdmin.GetAllOrders(r.Context())
	if err != nil {
		handleError(w, log, err)
		return
	}

	_ = json.NewEncoder(w).Encode(orders)
}

type roleReq struct {
	Role string `json:"role"`
}

func (t *TaskHandler) ChangeRole(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if id < 1 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var request roleReq

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "error from Decode", http.StatusBadRequest)
		return
	}

	log := t.logger.With(zap.String("handler", "ChangeRole"))

	user, err := t.serviceAdmin.ChangeRole(r.Context(), id, models.User{
		Role: request.Role})
	if err != nil {
		handleError(w, log, err)
		return
	}

	_ = json.NewEncoder(w).Encode(user)

}

func (t *TaskHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {

	orderID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "error from strconv.Atoi %w", http.StatusBadRequest)
		return
	}

	if orderID < 1 {
		http.Error(w, "invalid orderID", http.StatusBadRequest)
		return
	}

	userIDStr := r.Context().Value(internalCtx.UserIDKey)
	userID, ok := userIDStr.(int)
	if !ok {
		http.Error(w, "user not found in context", http.StatusBadRequest)
		return
	}

	if userID < 1 {
		http.Error(w, "invalid userID", http.StatusBadRequest)
		return
	}

	log := t.logger.With(zap.String("handler", "CancelOrder"))

	err = t.serviceOrder.CancelOrder(r.Context(), userID, orderID)
	if err != nil {
		handleError(w, log, err)
		return
	}

}
