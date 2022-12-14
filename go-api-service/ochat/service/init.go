package service

import (
	"ochat/bootstrap"

	"xorm.io/xorm"
)

var DB *xorm.Engine = bootstrap.DBOnceInit()
