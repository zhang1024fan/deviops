package result

// Codes 定义状态码和消息
type Codes struct {
	SUCCESS                                 uint
	FAILED                                  uint
	Message                                 map[uint]string
	NOAUTH                                  uint
	AUTHFORMATERROR                         uint
	INVALIDTOKEN                            uint
	MissingLoginParameter                   uint
	VerificationCodeHasExpired              uint
	CAPTCHANOTTRUE                          uint
	PASSWORDNOTTRUE                         uint
	STATUSISENABLE                          uint
	POSTALREADYEXISTS                       uint
	DEPTISEXIST                             uint
	DEPTISDISTRIBUTE                        uint
	MENUISEXIST                             uint
	DELSYSMENUFAILED                        uint
	ROLENAMEALREADYEXISTS                   uint
	MissingNewAdminParameter                uint
	USERNAMEALREADYEXISTS                   uint
	FILEUPLOADERROR                         uint
	MissingModificationOfPersonalParameters uint
	MissingChangePasswordParameter          uint
	RESETPASSWORD                           uint
	KUBEClUSTEREXIST                        uint
	KUBEClUSTERNOTEXIST                     uint
	ValidationParameterError                uint
	KUBEClUSTERNAMENOTEXIST                 uint
	WEBSOCKETERROR                          uint
}

// ApiCode 状态码
var ApiCode = &Codes{
	SUCCESS:                                 200,
	FAILED:                                  501,
	NOAUTH:                                  403,
	AUTHFORMATERROR:                         405,
	INVALIDTOKEN:                            406,
	MissingLoginParameter:                   407,
	VerificationCodeHasExpired:              408,
	CAPTCHANOTTRUE:                          409,
	PASSWORDNOTTRUE:                         410,
	STATUSISENABLE:                          411,
	POSTALREADYEXISTS:                       412,
	DEPTISEXIST:                             413,
	DEPTISDISTRIBUTE:                        414,
	MENUISEXIST:                             415,
	DELSYSMENUFAILED:                        416,
	ROLENAMEALREADYEXISTS:                   417,
	MissingNewAdminParameter:                418,
	USERNAMEALREADYEXISTS:                   419,
	FILEUPLOADERROR:                         420,
	MissingModificationOfPersonalParameters: 421,
	MissingChangePasswordParameter:          422,
	RESETPASSWORD:                           423,
	KUBEClUSTEREXIST:                        424,
	ValidationParameterError:                425,
	KUBEClUSTERNOTEXIST:                     426,
	KUBEClUSTERNAMENOTEXIST:                 427,
	WEBSOCKETERROR:                          428,
}

// 状态信息
func init() {
	ApiCode.Message = map[uint]string{
		ApiCode.SUCCESS:                                 "成功",
		ApiCode.FAILED:                                  "失败",
		ApiCode.NOAUTH:                                  "请求头中的auth为空",
		ApiCode.AUTHFORMATERROR:                         "请求头中的auth格式有错误",
		ApiCode.INVALIDTOKEN:                            "无效的token，请重新登录",
		ApiCode.MissingLoginParameter:                   "缺少登录参数",
		ApiCode.VerificationCodeHasExpired:              "验证码已失效",
		ApiCode.CAPTCHANOTTRUE:                          "验证码不正确，请重新输入",
		ApiCode.PASSWORDNOTTRUE:                         "密码不正确",
		ApiCode.STATUSISENABLE:                          "您的账号已被停用,请联系管理员",
		ApiCode.POSTALREADYEXISTS:                       "岗位名称或岗位编码已存在，请重新输入",
		ApiCode.DEPTISEXIST:                             "部门名称已存在，请重新输入",
		ApiCode.DEPTISDISTRIBUTE:                        "部门已分配，不能删除",
		ApiCode.MENUISEXIST:                             "菜单名称已存在，请重新输入",
		ApiCode.DELSYSMENUFAILED:                        "菜单已分配，不能删除",
		ApiCode.ROLENAMEALREADYEXISTS:                   "角色名称或角色KEY已存在，请重新输入",
		ApiCode.MissingNewAdminParameter:                "缺少新增参数",
		ApiCode.USERNAMEALREADYEXISTS:                   "用户名已存在，请重新输入",
		ApiCode.FILEUPLOADERROR:                         "文件上传错误",
		ApiCode.MissingModificationOfPersonalParameters: "缺少修改个人信息参数",
		ApiCode.MissingChangePasswordParameter:          "修改密码参数不能为空",
		ApiCode.RESETPASSWORD:                           "两次输入的密码不一致，请重新输入",
		ApiCode.KUBEClUSTEREXIST:                        "集群名称已存在，请重新输入",
		ApiCode.ValidationParameterError:                "参数验证错误",
		ApiCode.KUBEClUSTERNOTEXIST:                     "集群ID不存在",
		ApiCode.KUBEClUSTERNAMENOTEXIST:                 "集群名称不存在",
		ApiCode.WEBSOCKETERROR:                          "WebSocket连接错误",
	}
}

// 供外部调用
func (c *Codes) GetMessage(code uint) string {
	message, ok := c.Message[code]
	if !ok {
		return ""
	}
	return message
}
