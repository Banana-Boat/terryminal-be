package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/smtp"

	"github.com/hibiken/asynq"
	"github.com/jordan-wright/email"
	"github.com/rs/zerolog/log"
)

/* SendMail任务 */
func (processor *TaskProcessor) processTaskSendMail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendMail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal payload")
		return err
	}

	e := &email.Email{
		To:      []string{payload.To},
		From:    processor.config.EmailFromAddr,
		Subject: payload.Subject,
		Text:    []byte(payload.Text),
		HTML:    []byte(payload.Html),
	}

	e.Send(
		fmt.Sprintf("%s:%s", processor.config.SmtpHost, processor.config.SmtpPort),
		smtp.PlainAuth("", processor.config.EmailFromAddr, processor.config.EmailAuthCode, processor.config.SmtpHost),
	)

	log.Info().Msgf("task processed: type=%s, payload=%v", task.Type(), payload)
	return nil
}
