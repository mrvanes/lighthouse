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

type Attestation struct {
	Message string `json:"msg"`
}

// AttestationRequest is a request to the attestation endpoint
type AttestationRequest struct {
	Nonce               string   `json:"nonce"`
	IntegrityAssertion  string   `json:"integrity_assertion"`
	HardwareSignature   string   `json:"hardware_signature"`
	HardwareKeyTag      string   `json:"hardware_key_tag"`
	Cnf                 string   `json:"cnf"`
}

// AddAttestationEndpoint adds an attestation endpoint
func (fed *LightHouse) AddAttestationEndpoint(endpoint EndpointConf) {
	fed.server.Post(
		endpoint.Path, func(ctx *fiber.Ctx) error {
			var req AttestationRequest
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
			if req.IntegrityAssertion == "" {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(oidfed.ErrorInvalidRequest("required parameter 'integrity_assertion' not given"))
			}
			if req.HardwareSignature == "" {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(oidfed.ErrorInvalidRequest("required parameter 'hardware_signature' not given"))
			}
			if req.HardwareKeyTag == "" {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(oidfed.ErrorInvalidRequest("required parameter 'hardware_key_tag' not given"))
			}
			if req.Cnf == "" {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(oidfed.ErrorInvalidRequest("required parameter 'cnf' not given"))
			}

			message := &JWT{Message: "Hello " + req.Nonce + " " + req.Cnf}

			jwt, err := fed.GeneralJWTSigner.JWT(message, "oauth-client-attestation+jwt")
			if err != nil {
				return nil
			}
			return ctx.Send(jwt)
		},
	)
}
