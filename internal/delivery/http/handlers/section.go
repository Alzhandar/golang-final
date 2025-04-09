package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"restaurant-management/internal/models"
	"restaurant-management/internal/usecase"
)

type SectionHandler struct {
	sectionUC usecase.SectionUseCase
}

func NewSectionHandler(sectionUC usecase.SectionUseCase) *SectionHandler {
	return &SectionHandler{
		sectionUC: sectionUC,
	}
}

func (h *SectionHandler) Register(e *echo.Group) {
	sections := e.Group("/sections")
	sections.POST("", h.Create)
	sections.GET("/:id", h.GetByID)
	sections.GET("/restaurant/:restaurantID", h.GetByRestaurant)
	sections.PUT("/:id", h.Update)
	sections.DELETE("/:id", h.Delete)
}

// Create godoc
// @Summary Создать новую секцию ресторана
// @Description Создает новую секцию для указанного ресторана
// @Tags sections
// @Accept json
// @Produce json
// @Param section body models.Section true "Данные секции"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /sections [post]
func (h *SectionHandler) Create(c echo.Context) error {
	var section models.Section
	if err := c.Bind(&section); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректные данные секции",
		})
	}

	id, err := h.sectionUC.Create(c.Request().Context(), &section)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "секция успешно создана",
	})
}

// GetByID godoc
// @Summary Получить секцию по ID
// @Description Возвращает секцию по её ID
// @Tags sections
// @Accept json
// @Produce json
// @Param id path int true "ID секции"
// @Success 200 {object} models.Section
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /sections/{id} [get]
func (h *SectionHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID секции",
		})
	}

	section, err := h.sectionUC.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, section)
}

// GetByRestaurant godoc
// @Summary Получить секции по ID ресторана
// @Description Возвращает список секций для указанного ресторана
// @Tags sections
// @Accept json
// @Produce json
// @Param restaurantID path int true "ID ресторана"
// @Success 200 {array} models.Section
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /sections/restaurant/{restaurantID} [get]
func (h *SectionHandler) GetByRestaurant(c echo.Context) error {
	restaurantIDStr := c.Param("restaurantID")
	restaurantID, err := strconv.ParseInt(restaurantIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID ресторана",
		})
	}

	sections, err := h.sectionUC.GetByRestaurant(c.Request().Context(), restaurantID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, sections)
}

// Update godoc
// @Summary Обновить данные секции
// @Description Обновляет данные существующей секции
// @Tags sections
// @Accept json
// @Produce json
// @Param id path int true "ID секции"
// @Param section body models.Section true "Обновленные данные секции"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /sections/{id} [put]
func (h *SectionHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID секции",
		})
	}

	var section models.Section
	if err := c.Bind(&section); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректные данные секции",
		})
	}

	section.ID = id
	if err := h.sectionUC.Update(c.Request().Context(), &section); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "секция успешно обновлена",
	})
}

// Delete godoc
// @Summary Удалить секцию
// @Description Удаляет секцию по её ID
// @Tags sections
// @Accept json
// @Produce json
// @Param id path int true "ID секции"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /sections/{id} [delete]
func (h *SectionHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID секции",
		})
	}

	if err := h.sectionUC.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "секция успешно удалена",
	})
}
