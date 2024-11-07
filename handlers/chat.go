package handlers

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Conversation struct {
	user1 *User
	user2 *User
}

type User struct {
	conn  *websocket.Conn
	mutex sync.Mutex
}

var (
	upgrader           = websocket.Upgrader{}
	users              = map[string]*User{}
	usersMutex         sync.Mutex
	conversations      = map[string]*Conversation{}
	conversationsMutex sync.Mutex
)

func WsHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	conversationId := c.QueryParam("conversationId")
	userId := c.QueryParam("userId")

	// Save user
	usersMutex.Lock()
	user, exists := users[userId]
	if !exists {
		user = &User{
			conn: ws,
		}
		users[userId] = user
	}
	usersMutex.Unlock()

	// Save conversation
	conversationsMutex.Lock()
	conversation, exists := conversations[conversationId]
	if !exists {
		conversation = &Conversation{
			user1: user,
		}
		conversations[conversationId] = conversation
	} else {
		conversation.user2 = user
	}
	conversationsMutex.Unlock()

	defer func() {
		ws.Close()
		usersMutex.Lock()
		delete(users, userId)
		usersMutex.Unlock()

		conversationsMutex.Lock()
		if conversation.user1 == user {
			conversation.user1 = nil
		} else if conversation.user2 == user {
			conversation.user2 = nil
		}
		// Delete conversation if both users have disconnected
		if conversation.user1 == nil && conversation.user2 == nil {
			delete(conversations, conversationId)
		}
		conversationsMutex.Unlock()
	}()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			// Handle different types of errors
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Logger().Error("Unexpected closure:", err)
			} else {
				c.Logger().Info("Connection closed:", err) // Normal closure
			}
			break // Exit the loop on error
		}

		// Determine who should receive the message
		if user == conversation.user1 && conversation.user2 != nil {
			conversation.user2.mutex.Lock()
			err = conversation.user2.conn.WriteMessage(websocket.TextMessage, msg)
			conversation.user2.mutex.Unlock()
		} else if user == conversation.user2 && conversation.user1 != nil {
			conversation.user1.mutex.Lock()
			err = conversation.user1.conn.WriteMessage(websocket.TextMessage, msg)
			conversation.user1.mutex.Unlock()
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}

	return err
}
