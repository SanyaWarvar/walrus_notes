package smtp

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
	"wn/internal/domain/dto/auth"
	"wn/internal/domain/enum"
	apperrors "wn/internal/errors"
	"wn/pkg/applogger"
	"wn/pkg/util"

	"github.com/resend/resend-go/v3"
)

type Config struct {
	OwnerEmail    string
	OwnerPassword string
	Address       string
	CodeLenght    int
	CodeExp       time.Duration
	MinTTL        time.Duration
	ApiKey        string
}

func NewConfig(ownerEmail, ownerPassword, addres, apikey string, codeLenght int, codeExp, minTTL time.Duration) *Config {
	return &Config{
		OwnerEmail:    ownerEmail,
		OwnerPassword: ownerPassword,
		Address:       addres,
		CodeLenght:    codeLenght,
		CodeExp:       codeExp,
		MinTTL:        minTTL,
		ApiKey:        apikey,
	}
}

type cacheRepo interface {
	GetConfirmCode(ctx context.Context, email string) (*auth.ConfirmationCode, bool, error)
	SaveConfirmCode(ctx context.Context, email string, item auth.ConfirmationCode, ttl *time.Duration) error
}

type Service struct {
	logger applogger.Logger

	cacheRepo cacheRepo
	cfg       *Config
	resend    *resend.Client
}

func NewService(lgr applogger.Logger, cfg *Config, cacheRepo cacheRepo) *Service {
	client := resend.NewClient(cfg.ApiKey)

	return &Service{
		logger:    lgr,
		cfg:       cfg,
		cacheRepo: cacheRepo,
		resend:    client,
	}
}

func (srv *Service) SendConfirmEmailMessage(email, code, action string) error {
	htmlTemplate := `
<!DOCTYPE html>
<html>

<body>
    <div class="email">
        <div class="header">
            <img src="https://walrus-notes-q231.onrender.com/assets/logo.png" 
     alt="Logo" 
     class="logo" 
     style="width: 200px; height:auto; margin-bottom:12px;">
            <h1>Код подтверждения</h1>
        </div>
        
        <div class="content">
            <div class="message">
                <p>Для функции <span class="action">%s</span> требуется подтверждение.</p>
                <p>Используй этот код:<div class="code">%s</div></p>
            </div>  
            <p color: #475569; font-size: 14px; margin: 25px 0;">
                ⏳ Код действителен <strong>5 минут</strong>
            </p>
            <div class="warning">
                ⚠️ Если вы ничего не запрашивали — просто проигнорируйте это письмо.
            </div>
        </div>
        
        <div class="footer">
            <p><strong>Walrus Notes Team</strong></p>
            <p style="margin-top: 8px; font-size: 12px; opacity: 0.8;">
                Автоматическое сообщение • Не отвечать<br>
                support@walrus-notes.ru
            </p>
        </div>
    </div>
</body>
<head>
    <meta charset="UTF-8">
    <style>
        body { margin: 0; padding: 20px; background: #f8fafc; font-family: -apple-system, sans-serif; }
        .email { max-width: 500px; margin: 0 auto; background: white; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 12px rgba(0,0,0,0.08); }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); padding: 25px 20px; text-align: center; }
        .header h1 { color: white; margin: 0; font-size: 20px; font-weight: 600; }
        .content { padding: 30px; }
        .code { font-size: 32px; font-weight: bold; color: #667eea; text-align: center; letter-spacing: 3px; margin: 25px 0; font-family: monospace; }
        .message { background: #f1f5f9; padding: 15px; border-radius: 8px; margin: 20px 0; font-size: 15px; line-height: 1.5; }
        .footer { background: #f8fafc; padding: 20px; text-align: center; color: #64748b; font-size: 13px; border-top: 1px solid #e2e8f0; }
        .warning { background: #fed7d7; color: #742a2a; padding: 12px; border-radius: 6px; margin-top: 20px; font-size: 13px; }
        .action { color: #4f46e5; font-weight: 600; }
    </style>
</head>
</html>
    `

	htmlContent := fmt.Sprintf(htmlTemplate, action, code)
	subject := fmt.Sprintf("Код подтверждения: %s", code)

	return srv.SendMessage(email, htmlContent, subject)
}

func (srv *Service) SendMessage(email, messageText, title string) error {
	toEmail := email
	fromEmail := srv.cfg.OwnerEmail

	params := &resend.SendEmailRequest{
		From:    fromEmail,
		To:      []string{toEmail},
		Subject: title,
		Html:    messageText,
	}

	_, status := srv.resend.Emails.Send(params)
	/*
		subject_body := fmt.Sprintf("Subject:%s\n\n %s", title, messageText)
		status := smtp.SendMail(
			srv.cfg.Address,
			smtp.PlainAuth("", fromEmail, srv.cfg.OwnerPassword, "smtp.gmail.com"),
			fromEmail,
			[]string{toEmail},
			[]byte(subject_body),
		)*/

	return status
}

func (srv *Service) GenerateConfirmCode(action enum.EmailCodeAction) *auth.ConfirmationCode {
	return &auth.ConfirmationCode{
		Code:      fmt.Sprintf("%0*d", srv.cfg.CodeLenght, rand.Intn(int(math.Pow10(srv.cfg.CodeLenght)))),
		Action:    action,
		CreatedAt: util.GetCurrentUTCTime(),
	}
}

func (srv *Service) SendConfirmEmailCode(ctx context.Context, email string, action enum.EmailCodeAction) error {
	cachedCode, ex, err := srv.cacheRepo.GetConfirmCode(ctx, email)
	if ex && (err == nil && srv.cfg.CodeExp-srv.cfg.MinTTL > util.GetCurrentUTCTime().Sub(cachedCode.CreatedAt)) {
		return apperrors.ConfirmCodeAlreadySend
	}

	newCode := srv.GenerateConfirmCode(action)
	err = srv.cacheRepo.SaveConfirmCode(ctx, email, *newCode, &srv.cfg.CodeExp)
	go func() {
		err = srv.SendConfirmEmailMessage(email, newCode.Code, action.String())
		if err != nil {
			srv.logger.Warnf("error while sending confirm email message: %s", err.Error())
		}
	}()

	return err
}

func (srv *Service) ConfirmCode(ctx context.Context, email string, code string) (*auth.ConfirmationCode, error) {
	targetCode, ex, err := srv.cacheRepo.GetConfirmCode(ctx, email)
	if err != nil {
		return nil, err
	}
	if !ex {
		return nil, apperrors.ConfirmCodeNotExist
	}
	if code != targetCode.Code {
		return nil, apperrors.ConfirmCodeIncorrect
	}
	return targetCode, nil
}
