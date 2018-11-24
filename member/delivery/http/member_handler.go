package http

import (
	"context"
	memberUcase "github.com/mascot27/cleanApiRepoPattern/member"
	"github.com/labstack/echo"
	models "github.com/mascot27/cleanApiRepoPattern/models"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	validator "gopkg.in/go-playground/validator.v9"
)

type ResponseError struct {
	Message string `json:"message"`
}

type HttpMemberHandler struct {
	MUsecase memberUcase.MemberUsecase
}

func(m *HttpMemberHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	id := int64(idP)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	art, err := m.MUsecase.GetById(ctx, id)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, art)
}
func isRequestValid(m *models.Member) (bool, error) {

	validate := validator.New()

	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
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

	ar, err := m.MUsecase.Store(ctx, &member)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, ar)
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


func NewMemberHttpHandler(e *echo.Echo, us memberUcase.MemberUsecase) {
	handler := &HttpMemberHandler{
		MUsecase: us,
	}
	e.POST("/member", handler.Store)
	e.GET("/member/:id", handler.GetByID)

}