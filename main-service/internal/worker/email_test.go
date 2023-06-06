package worker

import (
	"fmt"
	"net/smtp"
	"testing"

	"github.com/Banana-Boat/terryminal/main-service/internal/util"
	"github.com/jordan-wright/email"
)

func TestCreateBash(t *testing.T) {
	// github action 不执行该测试
	if testing.Short() {
		t.Skip("Skipping test in github action.")
	}

	/* 加载配置 */
	config, err := util.LoadConfig("../..")
	if err != nil {
		t.Error(err)
	}

	e := &email.Email{
		To:      []string{"937499953@qq.com"},
		From:    config.EmailFromAddr,
		Subject: "test",
		HTML:    []byte("<h1>test</h1>"),
	}

	if err = e.Send(
		fmt.Sprintf("%s:%s", config.SmtpHost, config.SmtpPort),
		smtp.PlainAuth("", config.EmailFromAddr, config.EmailAuthCode, config.SmtpHost),
	); err != nil {
		t.Error(err)
	}

}
