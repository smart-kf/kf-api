package bill

import (
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/pkg/common"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type BatchAddCardRequest struct {
	CardType constant.CardType `json:"cardType" binding:"required" validate:"required,oneof=1 2" doc:"卡密类型: 1=正式卡，2=测试卡"`
	Days     int               `json:"days" doc:"天数,正式卡必填,测试卡忽略"`
	Num      int               `json:"num" doc:"数量,1-100之间的整数" binding:"required" validate:"required,gte=1,lte=100"`
}

type BatchAddResponse struct {
	Num int `json:"num"`
}

func (req *BatchAddCardRequest) Validate() error {
	if req.CardType == constant.CardTypeNormal {
		if req.Days <= 0 {
			return xerrors.NewParamsErrors("请填写天数")
		}
	}
	return nil
}

type ListCardRequest struct {
	common.PageRequest
	SaleStatus         constant.SaleStatus  `json:"saleStatus" doc:"卡片状态,1=销售中，2=下架，3=已出售"`
	LoginStatus        constant.LoginStatus `json:"loginStatus" doc:"登录状态: 1=未登录过，2=登录过"`
	CardType           constant.CardType    `json:"cardType" doc:"卡片类型: 1正式卡, 2测试卡"`
	ExpireStartTime    int64                `json:"expireStartTime" doc:"过期时间-开始时间，秒"`
	ExpireEndTime      int64                `json:"expireEndTime" doc:"过期时间-结束时间，秒"`
	CardID             string               `json:"cardID" doc:"卡密id"`
	LastLoginTimeStart int64                `json:"lastLoginTimeStart" doc:"上次登录时间-开始，秒"`
	LastLoginTimeEnd   int64                `json:"lastLoginTimeEnd" doc:"上次登录时间-开始，秒"`
}

type ListCardResponse struct {
	List  []*KFCardResponse `json:"list" doc:"列表数据"`
	Total int64             `json:"total" doc:"统计"`
}

type KFCardResponse struct {
	ID            uint                 `json:"id" doc:"主键id"`
	CardID        string               `json:"cardId" doc:"卡密id"`
	Password      string               `json:"password" doc:"密码"`
	SaleStatus    constant.SaleStatus  `json:"saleStatus" doc:"销售状态"`
	LoginStatus   constant.LoginStatus `json:"loginStatus" doc:"登录状态"`
	CardType      constant.CardType    `json:"cardType" doc:"卡片类型"`
	Day           int                  `json:"day"  doc:"卡密的天数"`
	ExpireTime    int64                `json:"expireTime" doc:"过期时间"`
	LastLoginTime int64                `json:"lastLoginTime" doc:"上次登录时间"`
}

type UpdateStatusRequest struct {
	ID     uint                `json:"id" binding:"required" doc:"主键id" validate:"required,gt=0"`
	Status constant.SaleStatus `json:"status" binding:"required" doc:"卡片状态,1=销售中，2=下架，3=已出售" validate:"required,oneof=1 2 3"`
}
