package tools

import (
	"clothesStore_go_Learn/DB"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"net/url"
	"time"
)

var (
	appID = "2021000119670683"
	//aliPublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu6wLDfKLr7pdPNW+E4t+ELtSfvk7GCXOpJIr+kq62E9OPIj2KXnBZgDewNNiXymDvyy6elTk4+wkhHGcj6GJjSLJpcYgfrJwrYvONja//mqrXRHG04FmGC2LsA5zfP/lLOUUyYGADNKwM2l2XeYYaMq8M+R8Jf9Q5KKnbGG/fN4h/YXmzDRAxZB679Y9ki5nWX+K6Yu2so3ho2V/5mtbkBnq5b4lD1RXrk8DtfFSB/QFhjt3cFarga1UE25hHBXL46WgpMPhO0cIezT6Ed4qjsPqoJg13hKi0EolOenZblcrcz5bug7O2sOPMmYaP8nsf/RdLKOhN3GQoyY48sfvBQIDAQAB"
	privateKey = "MIIEowIBAAKCAQEAkdOyuEZxG2MBHsyVxdaOJDwa4Rr74rkzMk2p8S1mUB9ZK1xoFy2kRCEM2tHL4viHHBj8LLPmDHmTbPlbOXAZ7L1Q2bWxgvzS4yCqpauMfT7eP1lcbvJ+RHptfwHZBYhx+0FLfar3cAZy3ZJBxQwKQjWw5TlpXFYJErxC7j5CxmCdNGFyz7BrGZyN2WDZPpAd6hitKJLNhRiU1xqUWoMk/EgxJ7c5gSchIZcqB8X4FrAmdbxE6ClE5HxCDXzdtBWuxitTy+S5h399Jk0vGzBUpepDc9ioXTPB7/cdUUE3QFUE9YBEIC1D2NsFN7YoTeLH/btdJ26+LJxLJnF0MJflsQIDAQABAoIBACvmwb08Z7zI94NgMA7ZYv2Bos32I7LD8qfIPcs/0bd5WIz3Stb/hJ6GHKqb0nfIPlS1KOYEWtOSnlGGWHJYT1W4QOjqDEDVAGAka3tow+jIznvf2TYFhwHyoZhE5CMISthLdgClQczWBCq0Z1x9HXGFXHYF7LRBqoWba8Lxt4SlBr2oNL3uUXiI+u19qZzbe0ogHJfVHAbQRASNe9YBoCGn1ALgzCSqTi6LakoxKKTu5+zB21i5AJXk0cc3+jxhw+oJQfl9GIpLxEkBuPRoF2RuNTugcZOYZbYLJ8PXq6bo1kzAaDMm98Lie22P+PENMkm5FMIER3N955p4pvnmbykCgYEA7DPOmmOWlWhCDT9dlbiYyqhmAYKOIPLKpeLdsDscEzAjmqJrLUdqD1S+dSODtprPRBMNnNacCVX/OQRUyE9qh6aCHM3fCEVAr7c1FrFiGaALFok87A0HAy/BPyDOYbf82Pfuz3EzUiaoD7WfKPaCn/qq21UUXOQsscVQAY4RTSsCgYEAngy0wd6hxTQIlYmUGqu1q0JhYn7JBCEAYg8g/U0VCer48Z+FuNqmleb71SCFyyzW7yhL/fTl41BYWIVo8MXM9lbyKTqaNK3be6eJNp4o51OL2nLMxB+siJY2J+a0frihswqWCWJ1qddvQ3ZeOXuTFFRmp+2w9QwO0O8nSzWKwpMCgYBnf8f/FLZOH5IZ1fM/ANVKsAGKldeLjnfHuqIjb7M8oTJottS50Xoi36JZF8fGQw2hKawkVlGnMZyVMlWoNExcxlRrJLafHCFdHa1QlUeELQHOzTH5yTeSaOGHtOtaHFHaDMIC+fpf+/pWb+IfA+13BlLJqv0yOvVurCQDmmnwYwKBgF1c8CJeG33c4P1FCkI/ENAcJF8EukZAIHPMsBYxxK3ZKjnBnEK4lxOSIU2jKqX81PLuAQYB9xMy0R1pobYpgow6jE6imZlo4nDHZRzojQ0po0Hl8uQgOdFtuowTkqgQ9SRIqpzcltk/tDBL6hlW0Gl/+ixVEuWOu+ncfH/HHzMVAoGBAKpcEj6ocBxz5eFCJRQrk3VwevVzf/p/ygsDi4onHt0IdgZpeE3qAjsK0lawuSE+7Vi17TFvLhKPr+auoy9Pf0unMFE1PcdlXkVEZPKzSa3fb6hBk/ZXmqT10qxAfG1lOoWMf0t2B1Qsw/HzIu75jUakKRduczHC9jIRFqObfNvH"
	client, _  = alipay.New(appID, privateKey, false)
)

func init() {
	//err := client.LoadAppPublicCertFromFile("conf/appCertPublicKey_2021000119670683.crt")
	//if err != nil {
	//	return
	//} // 加载应用公钥证书
	//err = client.LoadAliPayRootCertFromFile("conf/alipayRootCert.crt")
	//if err != nil {
	//	return
	//} // 加载支付宝根证书
	//err = client.LoadAliPayPublicCertFromFile("conf/alipayCertPublicKey_RSA2.crt")
	//if err != nil {
	//	return
	//} // 加载支付宝公钥证书
	client.LoadAliPayPublicKey("MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAgdmLEdTNS6zyj+X9mtZmwo6bN7Fu5G497PySW+y/PAYwCzA5OTHJeMymp8oqqKUAH4GjIAzW9oZOyrY1w8UvfK9DoRayu5wx+WJSu09xjQQaz180FZImPUS02TPtDT0Xz4qDisq3YOKpSKeDOsOU824mJvhID0qFpYfEqhFII2ohI46DaiBEL8yWyP5OWX+/ij8K9qMdKYMF8ilJy0MB2dYRd1RklY+UAKWMoAlumol+jURyS33aETnZuqiCLVI1a5f1yxHq7ymPE8fYsntppgHvBYQuUQvOwyvN6gTdOEqkzvkhymyHy11T5nYdq618xiGb2CZanNH0s+cUUn30awIDAQAB")
}

// WebPageAlipay 网页扫码支付
func WebPageAlipay(c *gin.Context) {
	//获取订单号
	orderSn := c.Query("orderSn")
	//查询订单金额
	var price string
	sqlStr := fmt.Sprintf("select pay_money from order_master where order_sn='%s'", orderSn)
	res, ok := DB.Query(sqlStr)
	if ok {
		price = res[0]["pay_money"]
	} else {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "生成付款订单界面时失败",
			"data": "",
		})
		return
	}
	var p = alipay.TradePagePay{}
	//p.NotifyURL = "http://127.0.0.1:8100/notify"
	p.ReturnURL = "http://127.0.0.1:8080/return"
	p.Subject = fmt.Sprintf("订单号：%s", orderSn) //付款标题
	p.OutTradeNo = orderSn                     //商家订单号
	p.TotalAmount = price
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	var url, err = client.TradePagePay(p)
	if err != nil {
		fmt.Println(err)
	}

	var payURL = url.String()
	//这个 payURL 即是用于支付的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	//fmt.Println(payURL)
	//打开默认浏览器
	//payURL = strings.Replace(payURL,"&","^&",-1)
	//exec.Command("cmd", "/c", "start",payURL).Start()
	//返回支付url用于打开支付
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "订单已生成，请尽快付款！",
		"data": payURL,
	})
	return
}

// AliPayNotify 接受通知接口
func AliPayNotify(c *gin.Context) {
	//获取订单号
	orderSn := c.Query("out_trade_no")
	//获取url并转成*URL
	x, _ := url.Parse(c.Request.URL.String())
	//验证是否成功支付
	ok, err := client.VerifySign(x.Query())
	if err != nil {
		fmt.Println("交易状态为:", err)
	}
	if ok {
		//更新订单状态
		nowTime := time.Now().Format("2006-01-02 15:04:05")
		sqlStr := fmt.Sprintf("update order_master set order_status = 3,update_time='%s' where order_sn='%s';", nowTime, orderSn)
		err := DB.Exec(sqlStr)
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "支付成功！",
			"data": "",
		})
	} else {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "支付失败！",
			"data": "",
		})
	}
}

func AliPayTest(c *gin.Context) {
	var noti, _ = client.GetTradeNotification(c.Request)
	if noti != nil {
		fmt.Println("交易状态为:", noti.TradeStatus)
	}
	alipay.AckNotification(c.Writer) // 确认收到通知消息
}
