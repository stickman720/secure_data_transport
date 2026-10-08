package cache_schema



type EmailVerificationCache struct{
	UserId int32 `json:"userid"`
	VerificationCode string `json:"verificationcode"`
}


type PassChangeCache struct {
	UserID int32 `json:"userid"`
}