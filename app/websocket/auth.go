package websocket

func CheckAccount(user string, pwd string) bool {
	if user != "13657134800" || pwd != "17671422338" {
		return false
	}
	return true
}
