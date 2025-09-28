package enum

type EmailCodeAction string

const (
	ConfirmCode    EmailCodeAction = "CONFIRM_CODE"
	ForgotPassword EmailCodeAction = "FORGOT_PASSWORD"
)

func (v EmailCodeAction) String() string {
	return string(v)
}
