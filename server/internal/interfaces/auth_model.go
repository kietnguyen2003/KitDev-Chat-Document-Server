package interfaces

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
}

type signinRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
