package ws

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	"github.com/sudoblockio/icon-go-api/config"
	"github.com/sudoblockio/icon-go-api/redis"
)

// BlocksAddHandlers - add fiber endpoint handlers for websocket connections
func BlocksAddHandlers(app *fiber.App) {

	prefix := config.Config.WebsocketPrefix + "/blocks"

	app.Use(prefix, func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get(prefix+"/", websocket.New(handlerGetBlocks))
}

func handlerGetBlocks(c *websocket.Conn) {

	// Add broadcaster
	msgChan := make(chan []byte)
	broadcasterID := redis.GetBroadcaster(config.Config.RedisBlocksChannel).AddBroadcastChannel(msgChan)
	defer func() {
		// Remove broadcaster
		redis.GetBroadcaster(config.Config.RedisBlocksChannel).RemoveBroadcastChannel(broadcasterID)
	}()

	// Read for close
	clientCloseSig := make(chan bool)
	go func() {
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				clientCloseSig <- true
				break
			}
		}
	}()

	for {
		// Read
		msg := <-msgChan

		// Broadcast
		err := c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}

		// check for client close
		select {
		case _ = <-clientCloseSig:
			break
		default:
			continue
		}
	}
}
