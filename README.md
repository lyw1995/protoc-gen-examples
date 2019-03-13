### github/track/protoc-gen-examples
### 简单的protobuff pulgin 扩展功能, 用上了Go的Template,
### 有空学学 go/ast

##  执行  MacOs
 -  国内免费代理配置 
 - export GOPROXY=https://athens.azurefd.net 
 -  生成`protoc-gen-tmpl`并拷贝到GOBIN,后生成tmpl代码
 - ./install_gen_plugin.sh 
 - ` go run tmpl/cmd/main.go` 执行
 - ` curl http://localhost:8000/userinfo` 

## 说明: 1. proto为定义文件测试文件
```
package user;
//option go_package="track/example/model";
import "google/protobuf/timestamp.proto";
import "google/protobuf/any.proto";
import "google/protobuf/descriptor.proto";

// 测试生成 user_msg.proto
extend google.protobuf.MethodOptions {
    string method = 24245;
}
// 用户控制器
service UserController {
    // 用户登录 传用户名,密码
    rpc Login (LoginReq) returns (LoginResp) {
        option (method) = "POST";
    }
    // 获取用户信息 传用户名即可
    rpc UserInfo (IdReq) returns (LoginResp) {
        option (method) = "GET";
    }
}

// 用户权限
enum Role {
    // 管理员
    Admin = 0;
    // 普通用户
    User = 1;
}

// 登录请求
message LoginReq {
    // 用户昵称
    string name = 1;
    // 用户密码
    string password = 2;
}
// 用户ID请求
message IdReq {
    // 用户id
    int64 id = 1;
}

// 登录响应
message LoginResp {
    // 用户id
    string user_id = 1;
    // 消息
    string msg = 2;
    // 列表
    repeated google.protobuf.Any list = 3;
    // 创建时间
    google.protobuf.Timestamp created_at = 4;
    // 权限
    Role role=5;
}

```
## 说明: 2. 生成请求,响应结构体
```

// 用户权限
type Role int32

const (
	// 管理员
	Role_Admin Role = 0
	// 普通用户
	Role_User Role = 1
)

// 登录请求
type LoginReq struct {
	// 用户昵称
	Name string `json:"name,omitempty"`
	// 用户密码
	Password string `json:"password,omitempty"`
}

// 用户ID请求
type IdReq struct {
	// 用户id
	Id int64 `json:"id,omitempty"`
}

// 登录响应
type LoginResp struct {
	// 用户id
	UserId string `json:"userId,omitempty"`
	// 消息
	Msg string `json:"msg,omitempty"`
	// 列表
	List []interface{} `json:"list,omitempty"`
	// 创建时间
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	// 权限
	Role Role `json:"role,omitempty"`
}

```
## 说明: 3. 生成http请求 cmd/main.go
```
import (
	"github/track/protoc-gen-examples/tmpl"
	"log"
	"net/http"
)
	

func Login(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	if method != "POST"  {
		writer.WriteHeader(http.StatusInternalServerError)
		_,_ = writer.Write([]byte("request method un support"))
	}
	var _ user.LoginReq
	var _ user.LoginResp

	writer.WriteHeader(http.StatusOK)
	_,_=writer.Write([]byte("did "+method+" Login successful\n"))
}

func UserInfo(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	if method != "GET"  {
		writer.WriteHeader(http.StatusInternalServerError)
		_,_ = writer.Write([]byte("request method un support"))
	}
	var _ user.IdReq
	var _ user.LoginResp

	writer.WriteHeader(http.StatusOK)
	_,_=writer.Write([]byte("did "+method+" UserInfo successful\n"))
}


	
func main() {
	
	http.HandleFunc("/login", Login)
	
	http.HandleFunc("/userinfo", UserInfo)
	
	log.Fatal(http.ListenAndServe(":8000", nil))
}

```