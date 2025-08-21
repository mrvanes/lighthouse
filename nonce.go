package lighthouse

import (
	"github.com/gofiber/fiber/v2"
)

type Nonce struct {
	Message string `json:"nonce"`
}

// AddNonceEndpoint adds a nonce endpoint
func (fed *LightHouse) AddNonceEndpoint(endpoint EndpointConf) {
	fed.server.Get(
		endpoint.Path, func(ctx *fiber.Ctx) error {
			nonce := Nonces.addNonce()
			message := &Nonce{Message: nonce}

			return ctx.JSON(message)
		},
	)
}
