package Cmd

const (
	 C2L_LoginConnectAck = 500 //连接登录服成功后发送确认
	 L2C_LoginConnectAck = 501 //连接登录服成功后登录服回复
	 C2L_LoginMsg = 502 //登陆
	 L2C_LoginMsg = 503 //登录返回角色
	 C2L_CreateRoleReq = 504 //创建角色
	 L2C_CreateRoleRes = 505 //返回创建的角色id与房间号

	 C2L_SendMsgReq = 600 //发送消息
	 L2C_SendMsgRes = 601 //发送消息
	 C2L_GetRoomMsgListReq = 602 //获取消息列表
	 L2C_GetRoomMsgListRes = 603 //获取消息列表回复
	 C2L_GetRoleInfoReq = 604 //获取玩家信息
	 L2C_GetUserInfoRes = 605 //获取玩家信息回复
	 C2L_PopularMsgReq = 606 //获取最近10分钟发送最频繁的消息
	 L2C_PopularMsgRes = 607 //
	 C2L_SwitchRoomReq = 608 //切换房间
	 L2C_SwitchRoomRes = 609 //切换房间回复

)