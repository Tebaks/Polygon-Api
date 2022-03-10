package handler

import (
	"app/log"
	"app/polygon"
	"app/repository"
	"app/service"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PolygonHandler interface {
	CreateNewPolygonRequest(c echo.Context) error
	GetPolygonByName(c echo.Context) error
}

type polygonHandler struct {
	polygonService    service.PolygonService
	polygonRepository repository.PolygonRepository
	logger            log.Logger
}

func NewPolygonHandler(polygonService service.PolygonService, polygonRepository repository.PolygonRepository, logger log.Logger) PolygonHandler {
	return &polygonHandler{
		polygonService:    polygonService,
		polygonRepository: polygonRepository,
		logger:            logger,
	}
}

// CreateNewPolygonRequest godoc
// @Description Create new polygon
// @Tags polygon
// @Accept  json
// @Produce  json
// @Param CreateNewPolygonRequest body CreateNewPolygonRequest true "Create new polygon info"
// @Success 200 {object} Response{error=string,success=bool,data=polygon.Polygon}
// @Failure 500 {object} Response
// @Router /polygon/ [post]
func (h *polygonHandler) CreateNewPolygonRequest(c echo.Context) error {
	var (
		req     CreateNewPolygonRequest
		polygon polygon.Polygon
		err     error
	)

	if err = c.Bind(&req); err != nil {
		h.logger.Error("invalid request body x", err)
		return ErrorWithCodeResponse(c, http.StatusBadRequest, errors.New("invalid request body"))
	}

	if err = c.Validate(req); err != nil {
		h.logger.Error("invalid request body", err)
		return ErrorWithCodeResponse(c, http.StatusBadRequest, errors.New("invalid request body"))
	}

	if len(req.Vertices) < 3 {
		return ErrorWithCodeResponse(c, http.StatusBadRequest, errors.New("polygon must have at least 3 vertices"))
	}

	if polygon, err = h.polygonService.CreateNewPolygon(c.Request().Context(), req.Vertices); err != nil {
		h.logger.ErrorWithFields("failed to create new polygon", err, log.Fields{"vertices": req.Vertices})
		return ErrorWithCodeResponse(c, http.StatusInternalServerError, errors.New("failed to create new polygon"))
	}

	return SuccessWithCodeResponse(c, http.StatusCreated, polygon)
}

// GetPolygonByName
// @Description Get polygon by name
// @Tags polygon
// @Accept  json
// @Produce  json
// @Param name path string true "name"
// @Success 200 {object} Response{error=string,success=bool,data=polygon.Polygon}
// @Failure 500 {object} Response
// @Router /polygon/{name} [get]
func (h *polygonHandler) GetPolygonByName(c echo.Context) error {
	var (
		req     GetPolygonByNameRequest
		polygon polygon.Polygon
		err     error
	)

	if err = c.Bind(&req); err != nil {
		h.logger.Error("invalid request body", err)
		return ErrorWithCodeResponse(c, http.StatusBadRequest, err)
	}

	if polygon, err = h.polygonRepository.FindByName(c.Request().Context(), req.Name); err != nil {
		if err == repository.ErrPolygonNotFound {
			return ErrorWithCodeResponse(c, http.StatusNotFound, err)
		}
		h.logger.ErrorWithFields("failed to get polygon by name", err, log.Fields{"name": req.Name})
		return ErrorWithCodeResponse(c, http.StatusInternalServerError, err)
	}

	return SuccessWithCodeResponse(c, http.StatusOK, polygon)
}
