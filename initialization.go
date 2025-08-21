package lighthouse

import (
	// "time"
	// "encoding/json"

	// "github.com/go-oidfed/lib/oidfedconst"
	"github.com/gofiber/fiber/v2"

	"github.com/go-oidfed/lib"
	// "github.com/go-oidfed/lib/apimodel"
	// "github.com/go-oidfed/lib/unixtime"
	// "github.com/go-oidfed/lib/jwx"

)

type Initialization struct {
	Message string `json:"msg"`
}

// InitializationRequest is a request to the initialization endpoint
type InitializationRequest struct {
	Nonce          string   `json:"nonce"`
	HardwareKeyTag string   `json:"hardware_key_tag"`
	KeyAttestation string   `json:"key_attestation"`
}

// AddInitializationEndpoint adds an initialization endpoint
func (fed *LightHouse) AddInitializationEndpoint(endpoint EndpointConf) {
	fed.server.Post(
		endpoint.Path, func(ctx *fiber.Ctx) error {
			var req InitializationRequest
			if err := ctx.QueryParser(&req); err != nil {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(oidfed.ErrorInvalidRequest("could not parse request parameters: " + err.Error()))
			}
			if err := ctx.BodyParser(&req); err != nil {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(oidfed.ErrorInvalidRequest("could not parse body parameters: " + err.Error()))
			}
			if req.Nonce == "" {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(oidfed.ErrorInvalidRequest("required parameter 'nonce' not given"))
			}
			if req.HardwareKeyTag == "" {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(oidfed.ErrorInvalidRequest("required parameter 'hardware_key_tag' not given"))
			}
			if req.KeyAttestation == "" {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(oidfed.ErrorInvalidRequest("required parameter 'key_attestation' not given"))
			}

			return ctx.SendStatus(204)
		},
	)
}
