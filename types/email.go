package types

type EmailPayload struct {
	MailTo  string      `json:"mailTo"`
	Subject string      `json:"subject"`
	Body    interface{} `json:"body"`
}
