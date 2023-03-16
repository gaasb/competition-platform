package auth

type User struct {
}
type Session struct {
}

func GetSession(token string) Session {
	return Session{}
}
func SetSession() {}
