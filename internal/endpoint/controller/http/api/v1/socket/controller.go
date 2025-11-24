package socket

import (
	"net/http"
	"wn/internal/domain/services/socket"
	"wn/pkg/apperror"
	"wn/pkg/applogger"
	"wn/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Controller struct {
	lgr      applogger.Logger
	builder  *response.Builder
	services *socket.Service
	upgrader websocket.Upgrader
}

func NewController(
	logger applogger.Logger,
	builder *response.Builder,
	services *socket.Service,
) *Controller {
	return &Controller{
		lgr:      logger,
		builder:  builder,
		services: services,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (h *Controller) ConnectionController(api *gin.RouterGroup) {
	routerGroup := api.Group("")
	{
		routerGroup.GET("/connection", h.createConnection)
		routerGroup.POST("/secret", h.generateSecret)
	}
}

// createConnection - HTTP хендлер для установки вебсокет соединения
func (h *Controller) createConnection(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		_ = c.Error(apperror.NewBadRequestError("empty user_id", "bind_query"))
		return
	}
	userId, err := uuid.Parse(userID)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), "bind_query"))
		return
	}
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		_ = c.Error(err)
		return
	}
	// Создаем доменное представление соединения
	wsConn := socket.NewWSConnection(conn, userId)
	// Передаем управление доменному сервису
	ctx := c.Request.Context()
	h.services.HandleConnection(ctx, wsConn)

	// Не отправляем ответ через builder, т.к. соединение уже установлено
	// и управление передано доменному сервису
}

// generateSecret - пример другого хендлера
func (h *Controller) generateSecret(c *gin.Context) {
	// Логика генерации секрета для аутентификации соединения
	type Request struct {
		UserID string `json:"user_id"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	// Генерация токена/секрета для вебсокета
	/*secret, err := h.services.GenerateConnectionSecret(c.Request.Context(), req.UserID)
	if err != nil {
		h.builder.BuildError(c, err)
		return
	}*/

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(c.Request.Context(), gin.H{
		"secret":  uuid.New().String(),
		"user_id": req.UserID,
	}))
}
