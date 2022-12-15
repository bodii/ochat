package bootstrap

import (
	"ochat/comm"
)

var (
	PROJECT_DIR = comm.GetProjectDIR()
)

func Init() {
	// 初始化配置项
	InitConfig()

	// 创建数据库链接
	DBOnceInit()
}
