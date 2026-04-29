package layout

import (
	"context"
	"wn/internal/domain/dto"
	"wn/internal/domain/events"
	"wn/pkg/applogger"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type permissionsService interface {
	GetAssociatedUsersByLayout(ctx context.Context, layoutId uuid.UUID) ([]uuid.UUID, error)
}

type socketService interface {
	SendTo(connID dto.ConnectionID, msg *dto.SocketMessage) error
}

type Producer struct {
	lgr                applogger.Logger
	permissionsService permissionsService
	socketService      socketService
}

func NewProducer(
	lgr applogger.Logger,
	permissionsService permissionsService,
	socketService socketService,
) *Producer {
	return &Producer{
		lgr:                lgr,
		permissionsService: permissionsService,
		socketService:      socketService,
	}
}

func (p *Producer) SendToAssociatedUsers(ctx context.Context, layoutId uuid.UUID, event events.Event) error {
	recipients, err := p.permissionsService.GetAssociatedUsersByLayout(ctx, layoutId)
	if err != nil {
		p.lgr.Errorf("SendToAssociatedUsers for layout %s: p.permissionsService.GetAssociatedUsersByLayout: %s", layoutId.String(), err.Error())
		return errors.Wrap(err, "p.permissionsService.GetAssociatedUsersByLayout")
	}
	for _, recipientId := range recipients {
		err = p.socketService.SendTo(dto.ConnectionID(recipientId.String()), event.ToSocketEvent())
		if err != nil {
			p.lgr.Errorf("SendToAssociatedUsers for layout %s: SendTo: %s", layoutId.String(), err.Error())
		}
	}
	return nil
}
