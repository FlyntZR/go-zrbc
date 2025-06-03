package main

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"net/smtp"
	"testing"
	"time"
)

// VerificationCode 验证码结构体
type VerificationCode struct {
	Code string
}

// 生成随机验证码
func generateCode(length int) string {
	var letters = []rune("1234567890")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// 发送邮件
func sendEmail(to, subject, body string) error {
	// 配置邮件服务器信息
	from := "flyntlee2024@gmail.com"
	password := "qqwjtqkaacdqpjxx"
	host := "smtp.gmail.com"
	port := "587"

	// 创建邮件正文
	t := template.Must(template.New("mailBody").Parse(body))
	buffer := bytes.NewBufferString("")
	code := generateCode(5) // 生成5位长度的验证码
	vc := VerificationCode{Code: code}
	t.Execute(buffer, vc)

	// 创建邮件消息
	msg := []byte("To: " + to + "\r\n" +
		"From: " + from + "<" + from + ">\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		buffer.String())

	// 连接邮件服务器
	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(host+":"+port, auth, from, []string{to}, msg)
	if err != nil {
		fmt.Printf("SendMail failed: %+v\n", err)
		return err
	}
	fmt.Printf("Email sent successfully to %s with verification code: %s\n", to, code)
	return nil
}

func TestSend_email(t *testing.T) {
	to := "matthew828486@gmail.com"
	subject := "Email Verification"
	body := `【德学院】您的验证码是 {{.Code}}，30分钟内有效。`
	err := sendEmail(to, subject, body)
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}
}
