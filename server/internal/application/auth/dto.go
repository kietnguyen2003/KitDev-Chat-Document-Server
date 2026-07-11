package auth

type RegisterRequest struct {
	Username string
	Password string
	Fullname string
}

type Token struct {
	AccessToken  string `json:"access_token"`
	ReFreshToken string `json:"refresh_token"`
	ExpireIn     int    `json:"expire_in"`
}

type User struct {
	FullName string `json:"fullname"`
	Role     string `json:"role"`
}

type Storage struct {
	CurrenSize uint64 `json:"current_size"`
	LimitSize  uint64 `json:"limit_size"`
}

type RegisterResponse struct {
	User    User    `json:"user"`
	Token   Token   `json:"token"`
	Storage Storage `json:"storage"`
}

type SignInRequest struct {
	Username string
	Password string
}

func DomainToDTo(accessToken string, refresheToken string, ttl int, fullName string, role string, currentStorage uint64, limitStorage uint64) *RegisterResponse {
	return &RegisterResponse{
		User: User{
			FullName: fullName,
			Role:     role,
		},
		Token: Token{
			AccessToken:  accessToken,
			ReFreshToken: refresheToken,
			ExpireIn:     ttl,
		},
		Storage: Storage{
			CurrenSize: currentStorage,
			LimitSize:  limitStorage,
		},
	}
}
