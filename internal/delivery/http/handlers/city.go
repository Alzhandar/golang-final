package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"restaurant-management/internal/models"
	"restaurant-management/internal/usecase"
)

type CityHandler struct {
	cityUC usecase.CityUseCase
}

func NewCityHandler(cityUC usecase.CityUseCase) *CityHandler {
	return &CityHandler{
		cityUC: cityUC,
	}
}

func (h *CityHandler) Register(e *echo.Group) {
	cities := e.Group("/cities")
	cities.POST("", h.Create)
	cities.GET("/:id", h.GetByID)
	cities.PUT("/:id", h.Update)
	cities.DELETE("/:id", h.Delete)
	cities.GET("", h.List)
}

// Create godoc
// @Summary Создать новый город
// @Description Создает новый город в системе
// @Tags cities
// @Accept json
// @Produce json
// @Param city body models.City true "Данные города"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/cities [post]
func (h *CityHandler) Create(c echo.Context) error {
	var city models.City
	if err := c.Bind(&city); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректные данные города",
		})
	}

	id, err := h.cityUC.Create(c.Request().Context(), &city)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "город успешно создан",
	})
}

// GetByID godoc
// @Summary Получить город по ID
// @Description Возвращает город по его ID
// @Tags cities
// @Accept json
// @Produce json
// @Param id path int true "ID города"
// @Success 200 {object} models.City
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/cities/{id} [get]
func (h *CityHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID города",
		})
	}

	city, err := h.cityUC.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, city)
}

// Update godoc
// @Summary Обновить данные города
// @Description Обновляет данные существующего города
// @Tags cities
// @Accept json
// @Produce json
// @Param id path int true "ID города"
// @Param city body models.City true "Обновленные данные города"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/cities/{id} [put]
func (h *CityHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID города",
		})
	}

	var city models.City
	if err := c.Bind(&city); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректные данные города",
		})
	}

	city.ID = id
	if err := h.cityUC.Update(c.Request().Context(), &city); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "город успешно обновлен",
	})
}

// Delete godoc
// @Summary Удалить город
// @Description Удаляет город по его ID
// @Tags cities
// @Accept json
// @Produce json
// @Param id path int true "ID города"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/cities/{id} [delete]
func (h *CityHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID города",
		})
	}

	if err := h.cityUC.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "город успешно удален",
	})
}

// List godoc
// @Summary Получить список всех городов
// @Description Возвращает список всех городов
// @Tags cities
// @Accept json
// @Produce json
// @Success 200 {array} models.City
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/cities [get]
func (h *CityHandler) List(c echo.Context) error {
	cities, err := h.cityUC.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, cities)
}
