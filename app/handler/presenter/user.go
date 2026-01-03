package presenter

type UserCreation struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

type UserCreationResponse struct {
	UserId string `json:"userId"`
}
