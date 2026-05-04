package layout

import (
	"context"
	"wn/internal/domain/dto"
	"wn/internal/domain/events"
	"wn/pkg/applogger"

	"github.com/google/uuid"
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

func (p *Producer) SendToAssociatedUsers(ctx context.Context, targetId uuid.UUID, recipients []uuid.UUID, event events.Event) error {

	for _, recipientId := range recipients {
		err := p.socketService.SendTo(dto.ConnectionID(recipientId.String()), event.ToSocketEvent())
		if err != nil {
			p.lgr.Errorf("SendToAssociatedUsers for id %s: SendTo: %s", targetId.String(), err.Error())
		}
	}
	return nil
}
