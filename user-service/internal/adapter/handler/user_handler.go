package handler

import (
	"net/http"
	"user-service/internal/adapter/handler/request"
	"user-service/internal/adapter/handler/response"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type UserHandlerInterface interface {
	SignIn(ctx echo.Context) error
}

type userHandler struct {
	userService service.UserServiceInterface
}

// SignIn implements [UserHandlerInterface].
func (u *userHandler) SignIn(c echo.Context) error {
	var (
		req      = request.SignInRequest{}
		resp     = response.DefaultResponse{}
		respSign = response.SignInResponse{}
		ctx      = c.Request().Context()
	)

	if err := c.Bind(&req); err != nil {
		log.Errorf("[UserHandler-1] SignIn: %v", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	if err := c.Validate(req); err != nil {
		log.Errorf("{UserHandler-2]: SignIn %v,", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusUnprocessableEntity, resp)
	}

	reqEntity := entity.UserEntity{
		Email:    req.Email,
		Password: req.Password,
	}
	user, token, err := u.userService.SignIn(ctx, reqEntity)
	if err != nil {
		if err.Error() == "404" {
			log.Errorf("{UserHandler-3]: SignIn %v,", err)
			resp.Message = "User not found"
			resp.Data = nil
			return c.JSON(http.StatusNotFound, resp)
		}
		log.Errorf("{UserHandler-4]: SignIn %v,", err)
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(http.StatusInternalServerError, resp)
	}

	respSign.ID = user.ID
	respSign.Name = user.Name
	respSign.Email = user.Email
	respSign.Role = user.RoleName
	respSign.Lat = user.Lat
	respSign.Lng = user.Lng
	respSign.Phone = user.Phone
	respSign.AccessToken = token

	resp.Message = "Success"
	resp.Data = respSign

	return c.JSON(http.StatusOK, resp)
}

func NewUserHandler(e *echo.Echo, userService service.UserServiceInterface) UserHandlerInterface {
	userHandler := &userHandler{userService: userService}

	e.Use(middleware.Recover())
	e.POST("/signin", userHandler.SignIn)

	return userHandler
}
