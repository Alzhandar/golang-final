package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"restaurant-management/internal/models"
	"restaurant-management/internal/usecase"
)

type RestaurantEventHandler struct {
	eventUC usecase.RestaurantEventUseCase
}

func NewRestaurantEventHandler(eventUC usecase.RestaurantEventUseCase) *RestaurantEventHandler {
	return &RestaurantEventHandler{
		eventUC: eventUC,
	}
}

func (h *RestaurantEventHandler) Register(e *echo.Group) {
	events := e.Group("/events")
	events.POST("", h.Create)
	events.GET("/:id", h.GetByID)
	events.GET("/type/:type", h.GetByType)
	events.PUT("/:id", h.Update)
	events.DELETE("/:id", h.Delete)
	events.GET("", h.List)
}

// Create godoc
// @Summary Создать новое событие ресторана
// @Description Создает новое событие ресторана в системе
// @Tags events
// @Accept json
// @Produce json
// @Param event body models.RestaurantEvent true "Данные события ресторана"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events [post]
func (h *RestaurantEventHandler) Create(c echo.Context) error {
	var event models.RestaurantEvent
	if err := c.Bind(&event); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректные данные события ресторана",
		})
	}

	id, err := h.eventUC.Create(c.Request().Context(), &event)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "событие ресторана успешно создано",
	})
}

// GetByID godoc
// @Summary Получить событие ресторана по ID
// @Description Возвращает событие ресторана по его ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "ID события ресторана"
// @Success 200 {object} models.RestaurantEvent
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{id} [get]
func (h *RestaurantEventHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID события ресторана",
		})
	}

	event, err := h.eventUC.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, event)
}

// GetByType godoc
// @Summary Получить события ресторана по типу
// @Description Возвращает список событий ресторана указанного типа
// @Tags events
// @Accept json
// @Produce json
// @Param type path string true "Тип события (wedding, birthday, corporate)"
// @Success 200 {array} models.RestaurantEvent
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/type/{type} [get]
func (h *RestaurantEventHandler) GetByType(c echo.Context) error {
	typeStr := c.Param("type")

	var eventType models.EventType
	switch typeStr {
	case "wedding":
		eventType = models.EventTypeWedding
	case "birthday":
		eventType = models.EventTypeBirthday
	case "corporate":
		eventType = models.EventTypeCorporate
	default:
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "неизвестный тип события, используйте: wedding, birthday, corporate",
		})
	}

	events, err := h.eventUC.GetByType(c.Request().Context(), eventType)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, events)
}

// Update godoc
// @Summary Обновить данные события ресторана
// @Description Обновляет данные существующего события ресторана
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "ID события ресторана"
// @Param event body models.RestaurantEvent true "Обновленные данные события ресторана"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{id} [put]
func (h *RestaurantEventHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID события ресторана",
		})
	}

	var event models.RestaurantEvent
	if err := c.Bind(&event); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректные данные события ресторана",
		})
	}

	event.ID = id
	if err := h.eventUC.Update(c.Request().Context(), &event); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "событие ресторана успешно обновлено",
	})
}

// Delete godoc
// @Summary Удалить событие ресторана
// @Description Удаляет событие ресторана по его ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "ID события ресторана"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /events/{id} [delete]
func (h *RestaurantEventHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID события ресторана",
		})
	}

	if err := h.eventUC.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "событие ресторана успешно удалено",
	})
}

// List godoc
// @Summary Получить список всех событий ресторана
// @Description Возвращает список всех событий ресторана
// @Tags events
// @Accept json
// @Produce json
// @Success 200 {array} models.RestaurantEvent
// @Failure 500 {object} map[string]interface{}
// @Router /events [get]
func (h *RestaurantEventHandler) List(c echo.Context) error {
	events, err := h.eventUC.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, events)
}
