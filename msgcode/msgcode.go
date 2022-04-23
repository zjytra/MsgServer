/*
创建时间: 2021/6/7 20:25
作者: zjy
功能介绍:

*/

package msgcode



const (
	AccountCode_None                = 0    //没有
	Succeed             			= 1    //成功
	AccountCode_UserNameShort       = 2    // 账号过短 >=4
	AccountCode_UserNameLong        = 3    // 账号超长 <= 11
	AccountCode_PasswordShort       = 4    // 密码过短 >=6
	AccountCode_PasswordLong        = 5    // 密码过短 <=11
	AccountCode_UserNameFormatErro  = 6    // 账号格式错误
	AccountCode_SqlZhuRu              = 7  // sql注入
	AccountCode_IsExsist              = 8  // 账号已经存在
	AccountCode_MACAccountNumIsMore   = 9  // 单台机器超过注册账号数量
	AccountCode_IsRegistering         = 10 // 正在注册中
	AccountCode_IsLogining            = 11 // 正在登录中
	AccountCode_IsLogined             = 12 //已经登录过
	AccountCode_IsLoginPassWordIsErro = 13 //登录密码错误
	AccountCode_NotExsist             = 14 // 账号不存在
	AccountCode_PassWordError         = 15 // 密码错误
	Login_DBERRO                  = 16     // 登录查询数据库错误
	LoginGate_RedisErro           = 17     // redis错误
	LoginGate_CacheNotFind        = 18     // redis错误
	LoginGate_DataCenterClose     = 19     // 数据中心关闭
	SeverCfg_NotFind              = 20     // 服务器配置数据未找到
	LoginQueryRolesDbError        = 21     // 登录数据库错误
	AccountCode_IsLoginDataCenter = 22     // 正在登录中心服
	MsgNotHandler                 = 23     // 无消息处理
	PlayerNotFind                 = 24     // 玩家没有找到
	RoleNameHasSpace              = 25     // 玩家名不能有空格
	RoleNameRepeat                = 26     // 玩家名不能有空格
	RoleIsMax                     = 27     // 玩家角色超过最大数
	ParseReqError                 = 28     // 解析请求出错
	RegisterRedisError            = 29     // 注册账户redis异常
	RoleIsExists                  = 30     // 玩家已经存在
	RoleNameError		 		  = 31    // 角色名错误
	SysError		 		  	  = 32    // 系统错误
	RoleCreateIng		          = 33 //
	CreateRoleRedisError          = 34     // 创建角色redis错误
	CreateRoleDBError             = 34     // 创建角色数据库错误
	RoomNotFind		              = 35     // 房间未找到
	UserNotFind		              = 36     // 玩家未找到
)
