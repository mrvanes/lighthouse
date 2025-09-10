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
	// "github.com/pkg/errors"

	"github.com/go-oidfed/lighthouse/storage"
)

type Msg struct {
	Message string `json:"msg"`
}

// type KeyAttestation struct {
// 	Key string `json:"key"`
// }

// RevocationRequest is a request to the initialization endpoint
type RevocationRequest struct {
	Nonce          string         `json:"nonce"`
	Status         string         `json:"status"`
	HardwareKeyTag string         `json:"hardware_key_tag"`
}

// InitializationRequest is a request to the initialization endpoint
type InitializationRequest struct {
	Nonce          string         `json:"nonce"`
	HardwareKeyTag string         `json:"hardware_key_tag"`
	KeyAttestation storage.KeyAttestation `json:"key_attestation"`
}

// AddInitializationEndpoint adds an initialization endpoint
func (fed *LightHouse) AddInitializationEndpoint(
	endpoint EndpointConf,
	store storage.WalletInstanceStorageBackend,
) {
	// Get endpoint
	fed.server.Patch(
		endpoint.Path, func(ctx *fiber.Ctx) error {
			var req RevocationRequest
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
			if req.Status == "" {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(oidfed.ErrorInvalidRequest("required parameter 'hardware_key_tag' not given"))
			}

			var hardware_key_tag = req.HardwareKeyTag
			var status = req.Status

			if status == "REVOKED" {
				err := store.Delete(hardware_key_tag)
				if err != nil {
					return ctx.JSON(oidfed.ErrorInvalidRequest("hardware key tag not found"))
				}
			}

			return ctx.SendStatus(204) // No Content
		},
	)

	// Get endpoint
	fed.server.Get(
		endpoint.Path, func(ctx *fiber.Ctx) error {
			return ctx.JSON(oidfed.ErrorInvalidRequest("missing path parameter"))
		},
	)

	// Get endpoint
	fed.server.Get(
		endpoint.Path + "/:key", func(ctx *fiber.Ctx) error {
			var hardware_key_tag = ctx.Params("key")

			info, err := store.WalletInstance(hardware_key_tag)
			if err != nil {
				return ctx.JSON(oidfed.ErrorInvalidRequest("hardware key tag not found"))
			}

			// return ctx.Send([]byte(info.KeyAttestation))
			return ctx.JSON(info.KeyAttestation)
		},
	)

	// Post endpoint
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
			if req.KeyAttestation == (storage.KeyAttestation{}) {
				ctx.Status(fiber.StatusBadRequest)
				return ctx.JSON(oidfed.ErrorInvalidRequest("required parameter 'key_attestation' not given"))
			}

			var nonce = req.Nonce
			var hardware_key_tag = req.HardwareKeyTag
			var key_attestation = storage.KeyAttestation{Key: req.KeyAttestation.Key}
			var time time.Time

			time, ok := Nonces[nonce]
			fmt.Println(time)
			if ok {
				Nonces.deleteNonce(nonce)
			} else {
				ctx.Status(403)
				return ctx.JSON(oidfed.ErrorInvalidRequest("invalid nonce"))
			}

			info := storage.WalletInstanceInfo{
				HardwareKeyTag: hardware_key_tag,
				KeyAttestation: key_attestation,
			}

			if err := store.Write(
				hardware_key_tag, info,
			); err != nil {
				ctx.Status(fiber.StatusInternalServerError)
				return ctx.JSON(oidfed.ErrorServerError(err.Error()))
			}

			return ctx.SendStatus(204) // No Content
		},
	)
}
