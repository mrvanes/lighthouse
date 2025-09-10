package lighthouse

import (
	"fmt"
	"time"
	// "encoding/json"

	// "github.com/go-oidfed/lib/oidfedconst"
	"github.com/gofiber/fiber/v2"

	"github.com/go-oidfed/lib"
	// "github.com/go-oidfed/lib/apimodel"
	// "github.com/go-oidfed/lib/unixtime"
	// "github.com/go-oidfed/lib/jwx"

)

type JWT struct {
	Message string `json:"msg"`
}

// AttestationRequest is a request to the attestation endpoint
type JWTRequest struct {
	Nonce       string   `json:"nonce"`
	Subject     string   `json:"sub" form:"sub" query:"sub" url:"sub"`
	Foobar      string   `json:"foobar"`
}

// AddAttestationEndpoint adds an attestation endpoint
func (fed *LightHouse) AddJWTEndpoint(endpoint EndpointConf) {
	// Get endpoint
	fed.server.Get(
		endpoint.Path, func(ctx *fiber.Ctx) error {
			return ctx.Send([]byte("Get"))
		},
	)

	// Get endpoint
	fed.server.Get(
		endpoint.Path + "/:key", func(ctx *fiber.Ctx) error {
			message := []byte(ctx.Params("key"))
			return ctx.Send(message)
		},
	)

	// Post endpoint
	fed.server.Post(
		endpoint.Path, func(ctx *fiber.Ctx) error {
			var req JWTRequest
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
			if req.Subject == "" {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(oidfed.ErrorInvalidRequest("required parameter 'sub' not given"))
			}
			if req.Foobar == "" {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(oidfed.ErrorInvalidRequest("required parameter body not given"))
			}


			var nonce = req.Nonce
			var time time.Time

			time, ok := Nonces[nonce]
			fmt.Println(time)
			if ok {
				Nonces.deleteNonce(nonce)
			} else {
				ctx.Status(403)
				return ctx.JSON(oidfed.ErrorInvalidRequest("invalid nonce"))
			}

			message := &JWT{Message: "Hello " + req.Subject + " " + nonce}

			// jwt, err := fed.GeneralJWTSigner.JWT(message, "oauth-client-attestation+jwt")
			// if err != nil {
			// 	return nil
			// }
			// return ctx.Send(jwt)
			return ctx.JSON(message)
		},
	)
}
