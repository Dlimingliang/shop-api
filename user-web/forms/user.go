package forms

type PasswordLoginForm struct {
	Mobile   string `json:"mobile" binding:"required"`
	Password string `json:"password" binding:"required,min=3,max=16"`
}
