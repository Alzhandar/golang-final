package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"restaurant-management/internal/models"
	"restaurant-management/internal/usecase"
)

type TableHandler struct {
	tableUC usecase.TableUseCase
}

func NewTableHandler(tableUC usecase.TableUseCase) *TableHandler {
	return &TableHandler{
		tableUC: tableUC,
	}
}

func (h *TableHandler) Register(e *echo.Group) {
	tables := e.Group("/tables")
	tables.POST("", h.Create)
	tables.GET("/:id", h.GetByID)
	tables.GET("/section/:sectionID", h.GetBySection)
	tables.PUT("/:id", h.Update)
	tables.DELETE("/:id", h.Delete)
	tables.POST("/:id/qr", h.GenerateQR)
}

// Create godoc
// @Summary Создать новый столик
// @Description Создает новый столик для указанной секции
// @Tags tables
// @Accept json
// @Produce json
// @Param table body models.Table true "Данные столика"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/tables [post]
func (h *TableHandler) Create(c echo.Context) error {
	var table models.Table
	if err := c.Bind(&table); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректные данные столика",
		})
	}

	id, err := h.tableUC.Create(c.Request().Context(), &table)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "столик успешно создан",
		"qr":      table.QR,
	})
}

// GetByID godoc
// @Summary Получить столик по ID
// @Description Возвращает столик по его ID
// @Tags tables
// @Accept json
// @Produce json
// @Param id path int true "ID столика"
// @Success 200 {object} models.Table
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/tables/{id} [get]
func (h *TableHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID столика",
		})
	}

	table, err := h.tableUC.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, table)
}

// GetBySection godoc
// @Summary Получить столики по ID секции
// @Description Возвращает список столиков для указанной секции
// @Tags tables
// @Accept json
// @Produce json
// @Param sectionID path int true "ID секции"
// @Success 200 {array} models.Table
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/tables/section/{sectionID} [get]
func (h *TableHandler) GetBySection(c echo.Context) error {
	sectionIDStr := c.Param("sectionID")
	sectionID, err := strconv.ParseInt(sectionIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID секции",
		})
	}

	tables, err := h.tableUC.GetBySection(c.Request().Context(), sectionID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, tables)
}

// Update godoc
// @Summary Обновить данные столика
// @Description Обновляет данные существующего столика
// @Tags tables
// @Accept json
// @Produce json
// @Param id path int true "ID столика"
// @Param table body models.Table true "Обновленные данные столика"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/tables/{id} [put]
func (h *TableHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID столика",
		})
	}

	var table models.Table
	if err := c.Bind(&table); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректные данные столика",
		})
	}

	table.ID = id
	if err := h.tableUC.Update(c.Request().Context(), &table); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "столик успешно обновлен",
	})
}

// Delete godoc
// @Summary Удалить столик
// @Description Удаляет столик по его ID
// @Tags tables
// @Accept json
// @Produce json
// @Param id path int true "ID столика"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/tables/{id} [delete]
func (h *TableHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID столика",
		})
	}

	if err := h.tableUC.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "столик успешно удален",
	})
}

// GenerateQR godoc
// @Summary Сгенерировать QR-код для столика
// @Description Генерирует и возвращает QR-код для указанного столика
// @Tags tables
// @Accept json
// @Produce json
// @Param id path int true "ID столика"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/tables/{id}/qr [post]
func (h *TableHandler) GenerateQR(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID столика",
		})
	}

	qr, err := h.tableUC.GenerateQR(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "QR-код успешно сгенерирован",
		"qr":      qr,
	})
}
