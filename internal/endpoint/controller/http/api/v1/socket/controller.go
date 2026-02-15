package socket

import (
	"net/http"
	"wn/internal/domain/services/multyplayer"
	"wn/internal/domain/services/socket"
	"wn/pkg/apperror"
	"wn/pkg/applogger"
	"wn/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Controller struct {
	lgr         applogger.Logger
	builder     *response.Builder
	services    *socket.Service
	multyplayer *multyplayer.Service
	upgrader    websocket.Upgrader
}

func NewController(
	logger applogger.Logger,
	builder *response.Builder,
	services *socket.Service,
	multyplayer *multyplayer.Service,
) *Controller {
	return &Controller{
		lgr:         logger,
		builder:     builder,
		services:    services,
		multyplayer: multyplayer,
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
		routerGroup.GET("/room/:id", h.createRoomConnection)
		routerGroup.POST("/secret", h.generateSecret)
	}
}

func (h *Controller) createRoomConnection(c *gin.Context) {
	// Валидация user_id
	userID := c.Query("user_id")
	if userID == "" {
		_ = c.Error(apperror.NewBadRequestError("empty user_id", "bind_query"))
		return
	}

	userId, err := uuid.Parse(userID)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError("invalid user_id: "+err.Error(), "bind_query"))
		return
	}

	// Валидация note_id (ID комнаты)
	noteID := c.Param("id")
	if noteID == "" {
		_ = c.Error(apperror.NewBadRequestError("empty id", "bind_query"))
		return
	}

	noteId, err := uuid.Parse(noteID)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError("invalid note_id: "+err.Error(), "bind_query"))
		return
	}

	// Upgrade HTTP to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.lgr.Errorf("failed to upgrade websocket connection: %s", err.Error())
		_ = c.Error(err)
		return
	}

	// Создаем connection object для мультиплеер сервиса
	wsConn := &multyplayer.Connection{
		Send: make(chan []byte, 512),
		Close: func() error {
			// При закрытии соединения удаляем пользователя из комнаты
			h.multyplayer.Disconnect(userId.String(), noteId.String())
			return conn.Close()
		},
	}

	// Подключаем пользователя к комнате
	h.multyplayer.Connect(userId.String(), noteId.String(), wsConn)

	// Запускаем горутину для отправки сообщений клиенту
	go func() {
		defer func() {
			h.multyplayer.Disconnect(userId.String(), noteId.String())
			conn.Close()
		}()

		for {
			select {
			case message, ok := <-wsConn.Send:
				if !ok {
					// Канал закрыт, закрываем соединение
					conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				// Отправляем сообщение клиенту
				if err := conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
					h.lgr.Errorf("failed to write message to client: %s",
						err.Error())
					return
				}
			}
		}
	}()

	// Основной цикл чтения сообщений от клиента
	go func() {
		defer func() {
			h.multyplayer.Disconnect(userId.String(), noteId.String())
			conn.Close()
		}()

		for {
			// Читаем сообщение от клиента
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				// Клиент отключился или ошибка
				h.lgr.Errorf("Клиент отключился или ошибка: %s",
					err.Error())
				break
			}

			// Поддерживаем только бинарные сообщения для проксирования
			if messageType != websocket.BinaryMessage {
				h.lgr.Warnf("received non-binary message: %v", message)
				continue
			}

			// Проксируем сообщение всем участникам комнаты
			h.multyplayer.HandleMessage(userId.String(), noteId.String(), message)
		}
	}()

	// Не отправляем ответ через builder, т.к. соединение уже установлено
	// и управление передано доменному сервису
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
	go h.services.HandleConnection(ctx, wsConn)

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
