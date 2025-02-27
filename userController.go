package controller

import (
    "log"
    "context"
    "net/http"
    "skyfox/bookings/model"
    ae "skyfox/error"

    "github.com/gin-gonic/gin"
)

type UserService interface {
    UserDetails(context.Context, string) (model.User, error)
    ChangePassword(ctx context.Context, username, newPassword string) error
}

type UserController struct {
    userService UserService
}

func NewUserController(userService UserService) *UserController {
    return &UserController{
        userService: userService,
    }
}

// login godoc
//
//		@Summary		Login
//		@Description	login
//		@Tags			login
//		@Accept			json
//		@Produce		json
//	 @param Authorization header string true "Enter basic auth"
//		@Success		200	{string}	string
//		@Failure		401	{object}	ae.AppError
//		@Failure		404	{object}	ae.AppError
//		@Failure		500	{object}	ae.AppError
//		@Router			/login [get]
func (uh *UserController) Login(c *gin.Context) {

    username, _, _ := c.Request.BasicAuth()

    user, err := uh.userService.UserDetails(c.Request.Context(), username)
    if err != nil {
        appError := err.(*ae.AppError)
        c.AbortWithStatusJSON(appError.HTTPCode(), appError)
        return
    }

    c.JSON(http.StatusOK, user.Username)
}

// changePassword godoc
//
//		@Summary		Change Password
//		@Description	change password
//		@Tags			user
//		@Accept			json
//		@Produce		json
//		@Param			body	body		map[string]string	true	"Change Password Request"
//		@Success		200	{string}	string
//		@Failure		400	{object}	ae.AppError
//		@Failure		500	{object}	ae.AppError
//		@Router			/change-password [post]
func (uh *UserController) ChangePassword(c *gin.Context) {
    var req struct {
        Username    string `json:"username"`
        NewPassword string `json:"newpassword"`
        OldPassword string `json:"oldpassword"`
    }
    log.Printf("Recieved rqst %+v", req)
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := uh.userService.ChangePassword(c.Request.Context(), req.Username, req.NewPassword)
    if err != nil {
        appError := err.(*ae.AppError)
        c.JSON(appError.HTTPCode(), appError)
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
