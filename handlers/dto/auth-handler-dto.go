package dto

type AuthLoginReqDTO struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	DeviceId  string `json:"deviceId"`
	UserAgent string `json:"userAgent"`
	IPAddress string `json:"ipAddress"`
}

type AuthLoginResDTO struct {
	Token string `json:"token"`
}

type WhoAmIResDTO struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	DOB       string `json:"dob"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
