package caches

type BillSettingCache struct{}

// IsVerifyCodeEnable 是否启用前台登录验证码
func (c *BillSettingCache) IsVerifyCodeEnable() bool {
	return false
}

func (c *BillSettingCache) GetNotice() string {
	return "测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告测试公告"
}
