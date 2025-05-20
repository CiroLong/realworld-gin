package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github/CiroLong/realworld-gin/internal/common"
	"github/CiroLong/realworld-gin/internal/models"
	"github/CiroLong/realworld-gin/internal/service"
	"net/http"
)

type UserController interface {
	RegisterUsers(context *gin.Context)

	Login(context *gin.Context)
}
type userController struct {
	userService service.UserService
}

// 使用依赖注入
func NewUserController() UserController {
	return &userController{userService: service.NewUserService()}
}

type UserModelValidator struct {
	User struct {
		Username string `form:"username" json:"username" binding:"exists,alphanum,min=4,max=255"`
		Email    string `form:"email" json:"email" binding:"exists,email"`
		Password string `form:"password" json:"password" binding:"exists,min=8,max=255"`
		Bio      string `form:"bio" json:"bio" binding:"max=1024"`
		Image    string `form:"image" json:"image" binding:"omitempty,url"`
	} `json:"user"`
	userModel models.UserModel `json:"-"`
}

func NewUserModelValidator() UserModelValidator {
	return UserModelValidator{}
}

func (self *UserModelValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	self.userModel.Username = self.User.Username
	self.userModel.Email = self.User.Email
	self.userModel.Bio = self.User.Bio

	if self.User.Password != common.NBRandomPassword {
		// TODO: SetPassword
	}
	if self.User.Image != "" {
		self.userModel.Image = &self.User.Image
	}
	return nil
}

type UserSerializer struct {
	c *gin.Context
}

type UserResponse struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Bio      string  `json:"bio"`
	Image    *string `json:"image"`
	Token    string  `json:"token"`
}

func (self *UserSerializer) Response() UserResponse {
	myUserModel := self.c.MustGet("my_user_model").(models.UserModel)
	user := UserResponse{
		Username: myUserModel.Username,
		Email:    myUserModel.Email,
		Bio:      myUserModel.Bio,
		Image:    myUserModel.Image,
		Token:    common.GenToken(myUserModel.ID),
	}
	return user
}

func (uc userController) RegisterUsers(c *gin.Context) {
	userModelValidator := NewUserModelValidator()
	if err := userModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	if err := uc.userService.SaveOneUser(&userModelValidator.userModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	c.Set("my_user_model", userModelValidator.userModel)

	serializer := UserSerializer{c}
	c.JSON(http.StatusCreated, gin.H{"user": serializer.Response()})
}

type LoginValidator struct {
	User struct {
		Email    string `form:"email" json:"email" binding:"exists,email"`
		Password string `form:"password"json:"password" binding:"exists,min=8,max=255"`
	} `json:"user"`
	userModel models.UserModel `json:"-"`
}

func (lv *LoginValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, lv)
	if err != nil {
		return err
	}

	lv.userModel.Email = lv.User.Email
	return nil
}

func NewLoginValidator() LoginValidator {
	loginValidator := LoginValidator{}
	return loginValidator
}

func (uc userController) Login(c *gin.Context) {
	loginValidator := NewLoginValidator()
	if err := loginValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	// 使用email查询
	userModel, err := uc.userService.FindOneUser(&models.UserModel{Email: loginValidator.userModel.Email})
	if err != nil {
		c.JSON(http.StatusForbidden, common.NewError("login", errors.New("Not Registered email or invalid password")))
		return
	}

	if uc.userService.CheckPassword(userModel, loginValidator.User.Password) != nil {
		c.JSON(http.StatusForbidden, common.NewError("login", errors.New("Not Registered email or invalid password")))
		return
	}

	// 把UserId和 model 存储在上下文中

	c.Set("my_user_id", userModel.ID)
	c.Set("my_user_model", userModel)

	serializer := UserSerializer{c}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}
