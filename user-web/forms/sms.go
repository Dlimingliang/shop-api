package forms

type SendSmsForm struct {
	Mobile string `json:"mobile" binding:"required,mobile"`
}
