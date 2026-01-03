package restful

import (
	"log"
	"strconv"

	"github.com/ducthangng/geofleet/gateway/app/usecase"
	"github.com/ducthangng/geofleet/gateway/app/usecase/user_service"
	"github.com/ducthangng/geofleet/gateway/service/ehandler"
	"github.com/gin-gonic/gin"
)

type UserHanlder struct {
	UserUsecase usecase.UserServiceUsecase
}

func NewUserHandler() *UserHanlder {
	return &UserHanlder{}
}

func (uh *UserHanlder) GetUser(ctx *gin.Context) {
	var (
		err         error
		userIdParam string
		userId      int
		res         user_service.UserDTO
	)

	defer func() {
		// do somethingg
		if err != nil {
			return
		}
	}()

	userIdParam = ctx.Query("userId")
	if len(userIdParam) == 0 {
		err = ehandler.ERROR_MISSING_USER
		return
	}

	userId, err = strconv.Atoi(userIdParam)
	if err != nil {
		return
	}

	// call to the user service
	res, err = uh.UserUsecase.GetUserProfile(ctx, userId)
	if err != nil {
		err = ehandler.ERROR_USER_SERVICE
		return
	}

	// call to userService
	log.Println(res)
}

// user ping for online time
func (uh *UserHanlder) Ping(ctx *gin.Context) {

}
