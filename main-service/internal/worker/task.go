package worker

/* 邮件发送 */
const TaskSendMail = "task:send_mail"

type PayloadSendMail struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Html    string `json:"html"`
	Text    string `json:"text"`
}
