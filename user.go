package main
import "net"
type User struct{
	Name string
	Addr string
	C chan string
	conn net.Conn
	server *Server
}
//创建一个用户的API
func NewUser(conn net.Conn,server *Server) *User{
	userAddr:=conn.RemoteAddr().String()
	user:=&User{
		Name:userAddr,
		Addr:userAddr,
		C:make(chan string),
		conn:conn,
		server:server,
	}
	//启动监听当前的user channel消息的goroutine
	go user.ListenMessage()
	return user
}
//用户的上线业务
func (this *User) Online(){
	//用户上线，将用户加入到onlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name]=this
	this.server.mapLock.Unlock()
	//广播当前用户上线消息
	this.server.BroadCast(this,"已上线")
}
//用户的下线功能
func (this *User) Offline(){
	//用户下线，将用户从OnlineMap中删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap,this.Name)
	this.server.mapLock.Unlock()
	//广播当前用户下线通知
	this.server.BroadCast(this,"下线")
}
//用户处理消息的业务
func (this *User) DoMessage(msg string){
	this.server.BroadCast(this,msg)
}
//监听当前user channel的方法，一旦有消息，就发送给客户端
func (this *User) ListenMessage(){
	for{
		msg:=<-this.C
		this.conn.Write([]byte(msg+"\n"))
	}
}