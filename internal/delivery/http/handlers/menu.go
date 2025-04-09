package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"restaurant-management/internal/models"
	"restaurant-management/internal/usecase"
)

type MenuHandler struct {
	menuUC usecase.MenuUseCase
}

func NewMenuHandler(menuUC usecase.MenuUseCase) *MenuHandler {
	return &MenuHandler{
		menuUC: menuUC,
	}
}

func (h *MenuHandler) Register(e *echo.Group) {
	menus := e.Group("/menus")
	menus.POST("", h.Create)
	menus.GET("/:id", h.GetByID)
	menus.GET("/restaurant/:restaurantID", h.GetByRestaurant)
	menus.PUT("/:id", h.Update)
	menus.DELETE("/:id", h.Delete)

}

// Create godoc
// @Summary Создать новое меню
// @Description Создает новое меню для указанного ресторана
// @Tags menus
// @Accept json
// @Produce json
// @Param menu body models.Menu true "Данные меню"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /menus [post]
func (h *MenuHandler) Create(c echo.Context) error {
	var menu models.Menu
	if err := c.Bind(&menu); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректные данные меню",
		})
	}

	id, err := h.menuUC.Create(c.Request().Context(), &menu)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "меню успешно создано",
	})
}

// GetByID godoc
// @Summary Получить меню по ID
// @Description Возвращает меню по его ID
// @Tags menus
// @Accept json
// @Produce json
// @Param id path int true "ID меню"
// @Success 200 {object} models.Menu
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /menus/{id} [get]
func (h *MenuHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID меню",
		})
	}

	menu, err := h.menuUC.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, menu)
}

// GetByRestaurant godoc
// @Summary Получить меню по ID ресторана
// @Description Возвращает список меню для указанного ресторана
// @Tags menus
// @Accept json
// @Produce json
// @Param restaurantID path int true "ID ресторана"
// @Success 200 {array} models.Menu
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /menus/restaurant/{restaurantID} [get]
func (h *MenuHandler) GetByRestaurant(c echo.Context) error {
	restaurantIDStr := c.Param("restaurantID")
	restaurantID, err := strconv.ParseInt(restaurantIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID ресторана",
		})
	}

	menus, err := h.menuUC.GetByRestaurant(c.Request().Context(), restaurantID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, menus)
}

// Update godoc
// @Summary Обновить данные меню
// @Description Обновляет данные существующего меню
// @Tags menus
// @Accept json
// @Produce json
// @Param id path int true "ID меню"
// @Param menu body models.Menu true "Обновленные данные меню"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /menus/{id} [put]
func (h *MenuHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID меню",
		})
	}

	var menu models.Menu
	if err := c.Bind(&menu); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректные данные меню",
		})
	}

	menu.ID = id
	if err := h.menuUC.Update(c.Request().Context(), &menu); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "меню успешно обновлено",
	})
}

// Delete godoc
// @Summary Удалить меню
// @Description Удаляет меню по его ID
// @Tags menus
// @Accept json
// @Produce json
// @Param id path int true "ID меню"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /menus/{id} [delete]
func (h *MenuHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID меню",
		})
	}

	if err := h.menuUC.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "меню успешно удалено",
	})
}
