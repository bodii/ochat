package service

import (
	"ochat/bootstrap"

	"xorm.io/xorm"
)

type GroupService struct {
	DB *xorm.Engine
}

func NewGroupServ() *GroupService {
	return &GroupService{
		DB: bootstrap.DB_Engine,
	}
}

func (g *GroupService) Add() {

}
