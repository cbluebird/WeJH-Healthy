package emailService

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"healthy/config/email"
	"log"
)

func SendEmail(service string) {
	email.MailConf.Title = "微精弘警报"
	email.MailConf.RecipientList = []string{email.MailConf.Sender}
	html := fmt.Sprintf(`<div>
        <div>
            经检测微精弘的%s服务瘫痪，请及时修复。
        </div> 
    </div>`, service)
	m := gomail.NewMessage()
	// 第三个参数是我们发送者的名称，但是如果对方有发送者的好友，优先显示对方好友备注名
	m.SetHeader(`From`, email.MailConf.Sender)
	m.SetHeader(`To`, email.MailConf.RecipientList...)
	m.SetHeader(`Subject`, email.MailConf.Title)
	m.SetBody(`text/html`, html)
	// m.Attach("./Dockerfile") //添加附件
	d := gomail.NewDialer(email.MailConf.SMTPAddr, email.MailConf.SMTPPort, email.MailConf.Sender, email.MailConf.SPassword)
	err := d.DialAndSend(m)
	if err != nil {
		log.Println(err)
	}
}
