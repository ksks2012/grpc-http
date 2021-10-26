package errcode

var (
	Success          = NewError(0, "成功")
	Fail             = NewError(10000000, "內部錯誤")
	InvalidParams    = NewError(10000001, "無效參數")
	Unauthorized     = NewError(10000002, "認證錯誤")
	NotFound         = NewError(10000003, "沒有找到")
	Unknown          = NewError(10000004, "未知")
	DeadlineExceeded = NewError(10000005, "超出最後截止期限")
	AccessDenied     = NewError(10000006, "訪問被拒絕")
	LimitExceed      = NewError(10000007, "訪問限制")
	MethodNotAllowed = NewError(10000008, "不支持該方法")
)
