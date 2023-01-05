package bootstrap

import "ochat/comm/funcs"

var (
	PROJECT_DIR = funcs.GetProjectDIR()
)

func Init() {
	// 初始化系统配置项
	InitSysConfig()

	// 创建数据库链接
	RedisOnceInit()

	// 创建数据库链接
	DBOnceInit()
}
