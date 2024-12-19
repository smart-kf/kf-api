package bill

import "github.com/clearcodecn/swaggos"

func SwaggerDoc(group *swaggos.Group) {
	bg := group.Group("/bill").Tag("计费后台")
	bg.Post("/login").Body(LoginRequest{}).JSON(LoginResponse{})
}
