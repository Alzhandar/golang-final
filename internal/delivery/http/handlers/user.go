package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"restaurant-management/internal/models"
	"restaurant-management/internal/usecase"
)

type UserHandler struct {
	userUC usecase.UserUseCase
}

func NewUserHandler(userUC usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUC: userUC,
	}
}

func (h *UserHandler) Register(e *echo.Group) {
	users := e.Group("/users")
	users.POST("", h.Create)
	users.GET("/:id", h.GetByID)
	users.GET("/phone/:phone", h.GetByPhone)
	users.PUT("/:id", h.Update)
	users.DELETE("/:id", h.Delete)
	users.GET("", h.List)
}

// Create godoc
// @Summary Создать нового пользователя
// @Description Создает нового пользователя в системе
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "Данные пользователя"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users [post]
func (h *UserHandler) Create(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректные данные пользователя",
		})
	}

	id, err := h.userUC.Create(c.Request().Context(), &user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "пользователь успешно создан",
	})
}

// GetByID godoc
// @Summary Получить пользователя по ID
// @Description Возвращает пользователя по его ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID пользователя",
		})
	}

	user, err := h.userUC.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}

// GetByPhone godoc
// @Summary Получить пользователя по номеру телефона
// @Description Возвращает пользователя по его номеру телефона
// @Tags users
// @Accept json
// @Produce json
// @Param phone path string true "Номер телефона пользователя"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/phone/{phone} [get]
func (h *UserHandler) GetByPhone(c echo.Context) error {
	phone := c.Param("phone")
	user, err := h.userUC.GetByPhone(c.Request().Context(), phone)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}

// Update godoc
// @Summary Обновить данные пользователя
// @Description Обновляет данные существующего пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param user body models.User true "Обновленные данные пользователя"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [put]
func (h *UserHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID пользователя",
		})
	}

	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректные данные пользователя",
		})
	}

	user.ID = id
	if err := h.userUC.Update(c.Request().Context(), &user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "пользователь успешно обновлен",
	})
}

// Delete godoc
// @Summary Удалить пользователя
// @Description Удаляет пользователя по его ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "некорректный ID пользователя",
		})
	}

	if err := h.userUC.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "пользователь успешно удален",
	})
}

// List godoc
// @Summary Получить список пользователей
// @Description Возвращает список пользователей с пагинацией
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Количество записей на странице" default(10)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]interface{}
// @Router /users [get]
func (h *UserHandler) List(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := 10
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err == nil && o >= 0 {
			offset = o
		}
	}

	users, err := h.userUC.List(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, users)
}
