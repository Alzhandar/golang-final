package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"restaurant-management/internal/models"
	"restaurant-management/internal/usecase"
)

type RestaurantHandler struct {
	restaurantUC usecase.RestaurantUseCase
}

func NewRestaurantHandler(restaurantUC usecase.RestaurantUseCase) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantUC: restaurantUC,
	}
}

func (h *RestaurantHandler) Register(e *echo.Group) {
	restaurants := e.Group("/restaurants")
	restaurants.POST("", h.Create)
	restaurants.GET("/:id", h.GetByID)
	restaurants.GET("/city/:cityID", h.GetByCity)
	restaurants.PUT("/:id", h.Update)
	restaurants.DELETE("/:id", h.Delete)
	restaurants.GET("", h.List)
}

// Create godoc
// @Summary Создать новый ресторан
// @Description Создает новый ресторан в системе
// @Tags restaurants
// @Accept json
// @Produce json
// @Param restaurant body models.Restaurant true "Данные ресторана"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /restaurants [post]
func (h *RestaurantHandler) Create(c echo.Context) error {
	var restaurant models.Restaurant
	if err := c.Bind(&restaurant); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректные данные ресторана",
		})
	}

	id, err := h.restaurantUC.Create(c.Request().Context(), &restaurant)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "ресторан успешно создан",
	})
}

// GetByID godoc
// @Summary Получить ресторан по ID
// @Description Возвращает ресторан по его ID
// @Tags restaurants
// @Accept json
// @Produce json
// @Param id path int true "ID ресторана"
// @Success 200 {object} models.Restaurant
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /restaurants/{id} [get]
func (h *RestaurantHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID ресторана",
		})
	}

	restaurant, err := h.restaurantUC.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, restaurant)
}

// GetByCity godoc
// @Summary Получить рестораны по ID города
// @Description Возвращает список ресторанов, находящихся в указанном городе
// @Tags restaurants
// @Accept json
// @Produce json
// @Param cityID path int true "ID города"
// @Success 200 {array} models.Restaurant
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /restaurants/city/{cityID} [get]
func (h *RestaurantHandler) GetByCity(c echo.Context) error {
	cityIDStr := c.Param("cityID")
	cityID, err := strconv.ParseInt(cityIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID города",
		})
	}

	restaurants, err := h.restaurantUC.GetByCity(c.Request().Context(), cityID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, restaurants)
}

// Update godoc
// @Summary Обновить данные ресторана
// @Description Обновляет данные существующего ресторана
// @Tags restaurants
// @Accept json
// @Produce json
// @Param id path int true "ID ресторана"
// @Param restaurant body models.Restaurant true "Обновленные данные ресторана"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /restaurants/{id} [put]
func (h *RestaurantHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID ресторана",
		})
	}

	var restaurant models.Restaurant
	if err := c.Bind(&restaurant); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректные данные ресторана",
		})
	}

	restaurant.ID = id
	if err := h.restaurantUC.Update(c.Request().Context(), &restaurant); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ресторан успешно обновлен",
	})
}

// Delete godoc
// @Summary Удалить ресторан
// @Description Удаляет ресторан по его ID
// @Tags restaurants
// @Accept json
// @Produce json
// @Param id path int true "ID ресторана"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /restaurants/{id} [delete]
func (h *RestaurantHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID ресторана",
		})
	}

	if err := h.restaurantUC.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ресторан успешно удален",
	})
}

// List godoc
// @Summary Получить список ресторанов
// @Description Возвращает список ресторанов с фильтрацией по активности
// @Tags restaurants
// @Accept json
// @Produce json
// @Param active query bool false "Фильтр по активности ресторанов" default(true)
// @Success 200 {array} models.Restaurant
// @Failure 500 {object} map[string]interface{}
// @Router /restaurants [get]
func (h *RestaurantHandler) List(c echo.Context) error {
	activeStr := c.QueryParam("active")

	active := true
	if activeStr != "" {
		a, err := strconv.ParseBool(activeStr)
		if err == nil {
			active = a
		}
	}

	restaurants, err := h.restaurantUC.List(c.Request().Context(), active)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, restaurants)
}
