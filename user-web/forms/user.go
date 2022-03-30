package forms

type PasswordLoginForm struct {
	Mobile    string `json:"mobile" binding:"required,mobile"`
	Password  string `json:"password" binding:"required,min=3,max=16"`
	Captcha   string `json:"captcha" binding:"required"`
	CaptchaID string `json:"captchaID" binding:"required"`
}
