package kffrontend

import "github.com/clearcodecn/swaggos"

func SwaggerDoc(group *swaggos.Group) {
	bg := group.Group("/kf-fe").Tag("客服前台")
	_ = bg
}
