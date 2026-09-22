package cache_schema



type EmailVerificationCache struct{
	UserId int64 `json:"userid"`
	VerificationCode string `json:"verificationcode"`
}
