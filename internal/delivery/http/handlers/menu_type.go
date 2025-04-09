package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"restaurant-management/internal/models"
	"restaurant-management/internal/usecase"
)

type MenuTypeHandler struct {
	menuTypeUC usecase.MenuTypeUseCase
}

func NewMenuTypeHandler(menuTypeUC usecase.MenuTypeUseCase) *MenuTypeHandler {
	return &MenuTypeHandler{
		menuTypeUC: menuTypeUC,
	}
}

func (h *MenuTypeHandler) Register(e *echo.Group) {
	menuTypes := e.Group("/menu-types")
	menuTypes.POST("", h.Create)
	menuTypes.GET("/:id", h.GetByID)
	menuTypes.PUT("/:id", h.Update)
	menuTypes.DELETE("/:id", h.Delete)
	menuTypes.GET("", h.List)
}

// Create godoc
// @Summary Создать новый тип меню
// @Description Создает новый тип меню в системе
// @Tags menu-types
// @Accept json
// @Produce json
// @Param menuType body models.MenuType true "Данные типа меню"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/menu-types [post]
func (h *MenuTypeHandler) Create(c echo.Context) error {
	var menuType models.MenuType
	if err := c.Bind(&menuType); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректные данные типа меню",
		})
	}

	id, err := h.menuTypeUC.Create(c.Request().Context(), &menuType)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "тип меню успешно создан",
	})
}

// GetByID godoc
// @Summary Получить тип меню по ID
// @Description Возвращает тип меню по его ID
// @Tags menu-types
// @Accept json
// @Produce json
// @Param id path int true "ID типа меню"
// @Success 200 {object} models.MenuType
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/menu-types/{id} [get]
func (h *MenuTypeHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID типа меню",
		})
	}

	menuType, err := h.menuTypeUC.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, menuType)
}

// Update godoc
// @Summary Обновить данные типа меню
// @Description Обновляет данные существующего типа меню
// @Tags menu-types
// @Accept json
// @Produce json
// @Param id path int true "ID типа меню"
// @Param menuType body models.MenuType true "Обновленные данные типа меню"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/menu-types/{id} [put]
func (h *MenuTypeHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID типа меню",
		})
	}

	var menuType models.MenuType
	if err := c.Bind(&menuType); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректные данные типа меню",
		})
	}

	menuType.ID = id
	if err := h.menuTypeUC.Update(c.Request().Context(), &menuType); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "тип меню успешно обновлен",
	})
}

// Delete godoc
// @Summary Удалить тип меню
// @Description Удаляет тип меню по его ID
// @Tags menu-types
// @Accept json
// @Produce json
// @Param id path int true "ID типа меню"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/menu-types/{id} [delete]
func (h *MenuTypeHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID типа меню",
		})
	}

	if err := h.menuTypeUC.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "тип меню успешно удален",
	})
}

// List godoc
// @Summary Получить список всех типов меню
// @Description Возвращает список всех типов меню
// @Tags menu-types
// @Accept json
// @Produce json
// @Success 200 {array} models.MenuType
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/menu-types [get]
func (h *MenuTypeHandler) List(c echo.Context) error {
	menuTypes, err := h.menuTypeUC.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, menuTypes)
}
