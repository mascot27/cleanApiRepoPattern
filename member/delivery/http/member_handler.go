package http

import (
	"context"
	"github.com/labstack/echo"
	memberUcase "github.com/mascot27/cleanApiRepoPattern/member"
	models "github.com/mascot27/cleanApiRepoPattern/models"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

type ResponseError struct {
	Message string `json:"message"`
}

type HttpMemberHandler struct {
	MUsecase memberUcase.MemberUsecase
}

func isRequestValid(m *models.Member) (bool, error) {

	validate := validator.New()

	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case models.INTERNAL_SERVER_ERROR:

		return http.StatusInternalServerError
	case models.NOT_FOUND_ERROR:
		return http.StatusNotFound
	case models.CONFLICT_ERROR:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func (m *HttpMemberHandler) Fetch(c echo.Context) error {

	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	listMr, nextCursor, err := m.MUsecase.Fetch(ctx, cursor, int64(num))

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listMr)
}

func (m *HttpMemberHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	id := int64(idP)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	memb, err := m.MUsecase.GetByID(ctx, id)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, memb)
}

func (m *HttpMemberHandler) Store(c echo.Context) error {
	var member models.Member
	err := c.Bind(&member)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&member); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	memb, err := m.MUsecase.Store(ctx, &member)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, memb)
}

func (m *HttpMemberHandler) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	id := int64(idP)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	_, err = m.MUsecase.Delete(ctx, id)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func NewMemberHttpHandler(e *echo.Echo, us memberUcase.MemberUsecase) {
	handler := &HttpMemberHandler{
		MUsecase: us,
	}
	e.GET("/member", handler.Fetch)
	e.POST("/member", handler.Store)
	e.GET("/member/:id", handler.GetByID)
	e.DELETE("/member/:id", handler.Delete)
}
