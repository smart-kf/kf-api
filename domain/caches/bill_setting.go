package caches

type BillSettingCache struct {
}

// IsVerifyCodeEnable 是否启用前台登录验证码
func (c *BillSettingCache) IsVerifyCodeEnable() bool {
	return false
}

func (c *BillSettingCache) GetNotice() string {
	return "测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告"
}

func (s *BillSettingCache) OneDayCardPrice() int {
	return 15 * 1000 // 15 U / 天
}

// NextQRCodeDomain 获取下一个二维码的域名
func (s *BillSettingCache) NextQRCodeDomain() string {
	return "https://qr.smart-kf.top"
}
