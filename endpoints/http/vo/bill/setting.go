package bill

type SettingRequest struct {
	DailyPackage   Package `json:"dailyPackage" doc:"日卡套餐"`
	WeeklyPackage  Package `json:"weeklyPackage" doc:"周卡套餐"`
	MonthlyPackage Package `json:"monthlyPackage" doc:"月卡套餐"`
	Payment        Payment `json:"payment" doc:"支付配置"`
	Notice         Notice  `json:"notice" doc:"公告"`
	DomainPrice    int     `json:"domainPrice" doc:"域名支付配置"`
}

type Package struct {
	Id    string  `json:"id" doc:"套餐id, daily=日卡, weekly=周卡, monthly=月卡"`
	Days  int     `json:"days" doc:"天数"`
	Price float64 `json:"price" doc:"价格不含小数"`
}

type Payment struct {
	PayUrl string `json:"payUrl" doc:"支付地址域名(带https)"`
	Token  string `json:"token" doc:"token"`
	AppId  string `json:"appId" doc:"appId"`
	Email  string `json:"email" doc:"邮件发送地址"`
}

type Notice struct {
	Content string `json:"content" doc:"公告内容"`
	Enable  bool   `json:"enable" doc:"是否启用,启用之后会在客服后台登录的时候展示"`
}
