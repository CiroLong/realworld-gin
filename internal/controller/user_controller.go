package controller

import "github/CiroLong/realworld-gin/internal/service"

type UserController struct {
	userService service.UserService
}

// 使用依赖注入
func NewUserController(s service.UserService) *UserController {
	return &UserController{userService: s}
}
