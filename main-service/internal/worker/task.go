package worker

const TaskSendMail = "task:send_mail"

type PayloadSendMail struct {
	DestAddr string `json:"destAddr"`
	Content  string `json:"content"`
}
