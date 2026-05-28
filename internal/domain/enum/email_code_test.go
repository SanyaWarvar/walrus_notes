package enum

import "testing"

func TestEmailCodeAction_String(t *testing.T) {
	tests := []struct {
		action EmailCodeAction
		want   string
	}{
		{ConfirmCode, "CONFIRM_CODE"},
		{ForgotPassword, "FORGOT_PASSWORD"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.action.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
