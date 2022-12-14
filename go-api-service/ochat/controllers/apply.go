package controllers

import "net/http"

func ApplyFriendInitiate(w http.ResponseWriter, r *http.Request) {
	// TODO：发起申请添加好友
}

func ApplyFriendList(w http.ResponseWriter, r *http.Request) {
	// TODO：查看向我申请好友的列表信息
}

func ApplyCommunityList(w http.ResponseWriter, r *http.Request) {
	// TODO：查看向我申请群成员的列表信息
}

func ApplyOperation(w http.ResponseWriter, r *http.Request) {
	// TODO: 处理别人的申请操作
}
