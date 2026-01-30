package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"testovoe/internal/domain"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UseCase interface {
	CreateSub(ctx context.Context, userSub domain.UserSub) error
	UpdateSub(ctx context.Context, userSub domain.UserSub) error
	DeleteSub(ctx context.Context, subID, userID uuid.UUID) error
	GetSubs(ctx context.Context) ([]*domain.UserSub, error)
	GetUserSub(ctx context.Context, subID uuid.UUID) (*domain.UserSub, error)
	GetUserSubs(ctx context.Context, userID uuid.UUID) ([]*domain.UserSub, error)
	GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName, fromStr, toStr string) (int, error)
}

type HttpHandler struct {
	log     *slog.Logger
	useCase UseCase
}

func New(log *slog.Logger, useCase UseCase) *HttpHandler {
	return &HttpHandler{log: log, useCase: useCase}
}

// CreateSub
// @Summary Создать новую подписку
// @Description Создает запись об онлайн-подписке для конкретного пользователя
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param   input  body      domain.UserSub  true  "Данные подписки"
// @Success 201    {object}  map[string]string "Успешное создание"
// @Failure 400    {object}  map[string]string "Ошибка валидации или некорректный JSON"
// @Failure 500    {object}  map[string]string "Внутренняя ошибка сервера"
// @Router /api/v1/subscriptions [post]
func (h *HttpHandler) CreateSub(w http.ResponseWriter, r *http.Request) {
	const op = "httpHandlers.CreateSub"
	ctx := r.Context()

	log := h.log.With(
		slog.String("op", op),
		slog.String("method", r.Method),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req domain.UserSub

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		log.Error("validation error", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	req.StartedAt = time.Now()

	err = h.useCase.CreateSub(ctx, req)
	if err != nil {
		log.Error("create sub failed", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal server error"})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{"status": "sub created successfully"})
}

// UpdateSub
// @Summary Обновить запись о подписке
// @Description Обновляет запись об онлайн-подписке для конкретного пользователя
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param   input  body      domain.UserSub  true  "Данные подписки"
// @Success 201    {object}  map[string]string "Успешное обновление"
// @Failure 400    {object}  map[string]string "Ошибка валидации или некорректный JSON"
// @Failure 500    {object}  map[string]string "Внутренняя ошибка сервера"
// @Router /api/v1/subscriptions [put]
func (h *HttpHandler) UpdateSub(w http.ResponseWriter, r *http.Request) {
	const op = "httpHandlers.UpdateSub"
	ctx := r.Context()

	log := h.log.With(
		slog.String("op", op),
		slog.String("method", r.Method),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req domain.UserSub

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		log.Error("validation error", "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	subIDStr := chi.URLParam(r, "id")
	subID, err := uuid.Parse(subIDStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	req.ID = subID

	err = h.useCase.UpdateSub(ctx, req)
	if err != nil {
		log.Error("update sub failed", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal server error"})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{"status": "sub updated successfully"})
}

// DeleteSub
// @Summary Удаляет запись о подписке
// @Description Удаляет запись по ID подписки (path) и ID пользователя (query)
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param   id       path      string  true  "ID подписки (UUID)"
// @Param   user_id  query     string  true  "ID пользователя (UUID)"
// @Success 204    "No Content"
// @Failure 400    {object}  map[string]string "Ошибка валидации ID"
// @Failure 500    {object}  map[string]string "Внутренняя ошибка сервера"
// @Router /api/v1/subscriptions/{id} [delete]
func (h *HttpHandler) DeleteSub(w http.ResponseWriter, r *http.Request) {
	const op = "httpHandlers.DeleteSub"
	ctx := r.Context()

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	subIDStr := chi.URLParam(r, "id")
	subID, err := uuid.Parse(subIDStr)
	if err != nil {
		log.Warn("invalid sub id", "id", subIDStr)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid subscription id"})
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Warn("invalid user id", "id", userIDStr)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid user id"})
		return
	}

	if err := h.useCase.DeleteSub(ctx, subID, userID); err != nil {
		log.Error("failed to delete sub", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal server error"})
		return
	}

	render.Status(r, http.StatusNoContent)
}

// ListSubs
// @Summary Получить список подписок
// @Description Возвращает все подписки или подписки конкретного пользователя (если передан user_id)
// @Tags subscriptions
// @Produce  json
// @Param   user_id  query     string  false  "ID пользователя (UUID)"
// @Success 200      {array}   domain.UserSub "Список подписок"
// @Failure 500      {object}  map[string]string "Внутренняя ошибка сервера"
// @Router /api/v1/subscriptions [get]
func (h *HttpHandler) ListSubs(w http.ResponseWriter, r *http.Request) {
	const op = "httpHandlers.ListSubs"
	ctx := r.Context()

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	userIDStr := r.URL.Query().Get("user_id")

	var subs []*domain.UserSub
	var err error

	if userIDStr != "" {
		uid, parseErr := uuid.Parse(userIDStr)
		if parseErr != nil {
			log.Warn("invalid user_id", "id", userIDStr)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "invalid user_id format"})
			return
		}
		subs, err = h.useCase.GetUserSubs(ctx, uid)
	} else {
		subs, err = h.useCase.GetSubs(ctx)
	}

	if err != nil {
		log.Error("failed to fetch subs", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal server error"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, subs)
}

// GetTotalCost
// @Summary Рассчитать итоговую стоимость
// @Description Считает сумму трат за период. Формат дат: MM-YYYY
// @Tags subscriptions
// @Produce  json
// @Param   user_id      query     string  true  "ID пользователя (UUID)"
// @Param   service_name query     string  true  "Название сервиса"
// @Param   from         query     string  true  "Дата начала (01-2025)"
// @Param   to           query     string  true  "Дата окончания (03-2025)"
// @Success 200          {object}  map[string]int "Результат"
// @Failure 400          {object}  map[string]string "Ошибка валидации параметров"
// @Failure 500          {object}  map[string]string "Внутренняя ошибка сервера"
// @Router /api/v1/subscriptions/total [get]
func (h *HttpHandler) GetTotalCost(w http.ResponseWriter, r *http.Request) {
	const op = "httpHandler.GetTotalCost"
	ctx := r.Context()

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	userIDStr := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if userIDStr == "" || serviceName == "" || from == "" || to == "" {
		log.Warn("missing query params")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "missing query params"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Warn("invalid user id", "id", userIDStr)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid user id"})
		return
	}

	totalCost, err := h.useCase.GetTotalCost(ctx, userID, serviceName, from, to)
	if err != nil {
		log.Error("failed to fetch total cost", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal server error"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{"totalCost": totalCost})
}

// GetUserSub
// @Summary Получить одну подписку
// @Description Возвращает подписку по её ID (передается в пути)
// @Tags subscriptions
// @Produce  json
// @Param   id   path      string  true  "ID подписки (UUID)"
// @Success 200  {object}  domain.UserSub "Данные подписки"
// @Failure 400  {object}  map[string]string "Некорректный ID"
// @Failure 404  {object}  map[string]string "Подписка не найдена"
// @Failure 500  {object}  map[string]string "Внутренняя ошибка сервера"
// @Router /api/v1/subscriptions/{id} [get]
func (h *HttpHandler) GetUserSub(w http.ResponseWriter, r *http.Request) {
	const op = "httpHandler.GetUserSub"
	ctx := r.Context()

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Warn("invalid sub id", "id", userIDStr)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid subscription id"})
		return
	}

	sub, err := h.useCase.GetUserSub(ctx, userID)
	if err != nil {
		log.Error("failed to fetch sub", "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal server error"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, sub)
}
