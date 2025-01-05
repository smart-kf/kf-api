package bill

import "github.com/smart-fm/kf-api/pkg/xerrors"

type LoginRequest struct {
	Username string `json:"username" binding:"required" doc:"用户名"`
	Password string `json:"password" binding:"required" doc:"密码"`
}

func (r *LoginRequest) Validate() error {
	if len(r.Username) == 0 || len(r.Password) == 0 {
		return xerrors.NewParamsErrors("账号或密码不能为空")
	}
	return nil
}

type LoginResponse struct {
	Token string `json:"token"`
}
