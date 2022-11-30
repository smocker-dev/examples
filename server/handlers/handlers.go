package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"example/server/services"
	"example/types"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	Reservations *reservations
}

func Setup(services *services.Services) *Handlers {
	return &Handlers{
		Reservations: &reservations{
			service: services.Reservations,
		},
	}
}

type reservations struct {
	service services.Reservations
}

func (h *reservations) GetReservations(c echo.Context) error {
	reservations, err := h.service.GetReservations(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, reservations)
}

func (h *reservations) GetReservationByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	reservation, err := h.service.GetReservationByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, reservation)
}

func (h *reservations) CreateReservation(c echo.Context) error {
	var body types.CreateReservationBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	reservation, err := h.service.CreateReservation(c.Request().Context(), body)
	if err != nil {
		if errors.Is(err, services.ErrInsufficientCapacity) {
			return c.JSON(http.StatusConflict, echo.Map{
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, reservation)
}
