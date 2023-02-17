package service

import (
	"fmt"
	"ochat/bootstrap"
	"ochat/models"

	"xorm.io/xorm"
)

type MessageService struct {
	DB *xorm.Engine
}

func NewMessageServ() *MessageService {
	return &MessageService{
		DB: bootstrap.DB_Engine,
	}
}

func (m *MessageService) Save(data *models.Message) (
	message *models.Message, err error) {

	_, err = m.DB.Insert(data)
	if err != nil {
		return message, err
	}

	fmt.Println("\nsave data: ", data)

	return data, err
}
