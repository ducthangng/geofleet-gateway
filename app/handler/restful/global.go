package restful

import (
	"log"
	"net/http"

	"github.com/ducthangng/geofleet/gateway/app/handler/presenter"
	"github.com/ducthangng/geofleet/gateway/app/usecase"
	"github.com/ducthangng/geofleet/gateway/app/usecase/user_service"
	"github.com/ducthangng/geofleet/gateway/service/copier"
	"github.com/gin-gonic/gin"
)

type GlobalHandler struct {
	UserUsecase usecase.UserServiceUsecase
}

func NewGlobalHandler() *GlobalHandler {
	return &GlobalHandler{
		UserUsecase: user_service.NewUserService(),
	}
}

func (gh *GlobalHandler) Register(ctx *gin.Context) {
	var (
		err     error
		userId  int
		req     presenter.UserCreation
		userDTO user_service.UserCreation
	)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		return
	}

	copier.MustCopy(&userDTO, &req)

	userId, err = gh.UserUsecase.CreateUserProfile(ctx, userDTO)
	if err != nil {
		log.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, userId)
}

// Login can mean get online - or log in with account info
func (gh *GlobalHandler) Login(ctx *gin.Context) {

}
