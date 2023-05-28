package forms

type RegistrationForm struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	//some info
}

type AuthForm struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
