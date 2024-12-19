package kfbackend

import "github.com/clearcodecn/swaggos"

func SwaggerDoc(group *swaggos.Group) {
	bg := group.Group("/kf-be").Tag("客服后台")
	_ = bg
}
