package main

//验证码相关文件，包括邮箱验证码和图片验证码

import (
	"fmt"
	"github.com/robfig/cron"
	"math/rand"
	"time"
)

// Code 定义一个map存放邮箱和验证码，key为邮箱，value为验证码
var Code = make(map[string]CodeTime)

type CodeTime struct {
	CodeNum  int
	CodeTime int64
}

func init() {
	go DestroyCode()
	//var c = make(chan bool)
	//<- c
	//time.Sleep(5 * time.Second)
}

func SendEmail(email string) int {
	//smtp.PlainAuth()
	// 参数1：Usually identity should be the empty string, to act as username
	// 参数2：username
	//参数3：password
	//参数4：host
	//auth := smtp.PlainAuth("", "3216300435@qq.com", "sosdeenqiaqfdedg", "smtp.qq.com")
	//to := []string{email}
	//发送随机数为验证码
	// Seed uses the provided seed value to initialize the default Source to a
	// deterministic state. If Seed is not called, the generator behaves as
	// if seeded by Seed(1). Seed values that have the same remainder when
	// divided by 2^31-1 generate the same pseudo-random sequence.
	// Seed, unlike the Rand.Seed method, is safe for concurrent use.
	rand.Seed(time.Now().Unix())
	// Intn returns, as an int, a non-negative pseudo-random number in [0,n)
	num := 9999
	for {
		num = rand.Intn(10000)
		if num > 1000 {
			break
		}
	}
	Code[email] = CodeTime{
		CodeNum:  num,
		CodeTime: time.Now().Unix(),
	}
	//发送内容使用base64 编码，单行不超过80字节，需要插入\r\n进行换行
	//The msg headers should usually include
	// fields such as "From", "To", "Subject", and "Cc".  Sending "Bcc"
	// messages is accomplished by including an email address in the to
	// parameter but not including it in the msg headers.
	//str := fmt.Sprintf("From:3216300435@qq.com\r\nTo:%s\r\nSubject:商城注册\r\n\r\n您的即将注册xx商城，您的验证码是：%d，如若不是您本人操作请忽略。", email,num) //邮件格式
	//msg := []byte(str)
	//err := smtp.SendMail("smtp.qq.com:587", auth, "3216300435@qq.com", to, msg)
	//if err != nil {
	//	log.Fatal(err)
	//}
	return num
}

// DestroyCode 每隔5分钟检查是否有超时的验证码进行删除
func DestroyCode() {
	c := cron.New()
	err := c.AddFunc("@every 5m", func() {
		for k, v := range Code {
			if (time.Now().Unix() - v.CodeTime) > 600 {
				//大于10分钟则销毁验证码
				delete(Code, k)
			}
		}
	})
	if err != nil {
		fmt.Errorf("AddFunc error : %v", err)
		return
	}
	c.Start()
	defer c.Stop()
	select {}
	//time.Sleep(5 * time.Second)
}
