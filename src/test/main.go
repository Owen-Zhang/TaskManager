package main

import 
(
	"app/redis"
	"time"
	"fmt"
	//import "github.com/garyburd/redigo/redis"
 	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/process"
	"net/smtp"
	"strings"
	//"os"
	"runtime"
)

//测试 
func main() {
	//testComputerMem()
	//testComputerCpu()
	//testComputerDisk()
	//testComputerProcess()
	
	//sendMail()
	
	fmt.Println(runtime.GOOS)
	
	for {
		
	}
}

func sendMail() {
	user := "zhangxiaohui.cn@163.com"
	password := "thankmygod"
	host := "smtp.163.com:25"
	to := "1034992601@qq.com"

	subject := "使用Golang发送邮件"

	body := `
		<html>
		<body>
		<h3>
		"Test send to email"
		</h3>
		</body>
		</html>
		`
	fmt.Println("send email")
	err := SendToMail(user, password, host, to, subject, body, "html")
	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("Send mail success!")
	}
}

func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func testComputerMem() {
	v, _ := mem.VirtualMemory()
	fmt.Printf("总内存: %v, 可供程序分配内存: %v, 已用内存: %v,  Free:%v, UsedPercent:%f%%\n, Active: %v, Inactive: %v, Wired: %v\n", 
				v.Total, v.Available, v.Used, v.Free, v.UsedPercent, v.Active, v.Inactive, v.Wired)
}

func testComputerCpu() {
	c, _ := cpu.Info()
	cc,_ := cpu.Percent(time.Second,false)
	i := 0
	for _, sub_cup := range c {
		modelName := sub_cup.ModelName
		cores := sub_cup.Cores
		fmt.Printf("cpu: %s, %d, 使用率%f%%\n", modelName, cores, cc[i])
		i++
	}
}

func testComputerDisk() {
	d, _ := disk.Usage("/")
	
	fmt.Printf("硬盘总共: %v,  用了: %v, 剩下：%v\n", d.Total/1024/1024/1024, d.Used/1024/1024/1024, d.Free/1024/1024/1024)
}

func testComputerProcess() {
	p, _ := process.Processes()
	for _, val := range p {
		fmt.Println(val.Pid)
		//fmt.Println(val.Name(), val.Pid, val.MemoryInfo().Data, val.CPUPercent())
	}
}

//读写redis
func testRedis() {
	
	/*
	c, err := redis.Dial("tcp", "10.10.6.8:8501")
	if err != nil {
		fmt.Println(err)
		return
	}
	//密码授权
	c.Do("AUTH", "kjt@123")
	c.Do("SELECT", 15)
	c.Do("SET", "a", "1223456789")
	a, err := redis.String(c.Do("GET", "a"))

	fmt.Println(a)

	defer c.Close()
	*/
	
	redis.Set("123456789", "sdfasfdasdfasdfadsfasdf")
}