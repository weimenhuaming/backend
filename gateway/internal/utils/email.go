package utils

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"math/big"
	"net/smtp"
	"strings"

	"github.com/jordan-wright/email"
)

func GenerateCode() (string, error) {
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n), nil
}

// Email 发送电子邮件
// 可以给多个人发送，逗号隔开就行
func Email(To, subject string, body string) error {
	to := strings.Split(To, ",") // 将收件人邮箱地址按逗号拆分成多个地址
	return send(to, subject, body)
}

// send 执行邮件发送操作
func send(to []string, subject string, body string) error {
	from := "492730753@qq.com"
	nickname := "chena7"
	secret := "ympuazjpvcxccagh"
	host := "smtp.qq.com"
	port := 465
	isSSL := true

	// 使用 PlainAuth 创建认证信息
	auth := smtp.PlainAuth("", from, secret, host)

	// 创建新的电子邮件对象
	e := email.NewEmail()
	if nickname != "" {
		// 如果设置了昵称，则格式化发件人地址为 "昵称 <邮箱>"
		e.From = fmt.Sprintf("%s <%s>", nickname, from)
	} else {
		// 否则直接使用发件人邮箱
		e.From = from
	}

	// 设置收件人、主题和邮件内容
	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)

	// 定义错误变量
	var err error
	// 构建邮件服务器的地址，格式为 host:port
	hostAddr := fmt.Sprintf("%s:%d", host, port)

	// 根据配置的是否使用 SSL 来选择邮件发送方法
	if isSSL {
		// 使用带 TLS 的邮件发送
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: host})
	} else {
		// 使用普通的邮件发送
		err = e.Send(hostAddr, auth)
	}

	return err
}

func SendEmailVerificationCode(to, captcha string) error {
	// 生成一个6位数的指定的验证吗，返回一个界面
	subject := "您的邮箱验证码"
	body := `亲爱的用户[` + to + `]，<br/>
<br/>
感谢您注册` + "chena7" + `的个人官网！为了确保您的邮箱安全，请使用以下验证码进行验证：<br/>
<br/>
验证码：[<font color="blue"><u>` + captcha + `</u></font>]<br/>
该验证码在 1 分钟内有效，请尽快使用。<br/>
<br/>
如果您没有请求此验证码，请忽略此邮件。
<br/>
如有任何疑问，请联系我们的支持团队：<br/>
邮箱：` + "492730753@qq.com" + `<br/>
<br/>
祝好，<br/>` + "陈阿七" + `<br/>
<br/>`

	_ = Email(to, subject, body)

	return nil
}
