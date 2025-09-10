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
	"github.com/go-oidfed/lighthouse/storage"
)

// type Attestation struct {
// 	KeyAttestation storage.KeyAttestation
// }

type WalletAttestation struct {
	Format  string `json:"format"`
	// Attestation Attestation `json:"attestation"`
	Attestation storage.KeyAttestation `json:"attestation"`
}

type Attestations struct {
	WalletAttestations[] WalletAttestation `json:"wallet_attestations"`
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
func (fed *LightHouse) AddAttestationEndpoint(
	endpoint EndpointConf,
	store storage.WalletInstanceStorageBackend,
) {
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

			var nonce = req.Nonce
			var hardware_key_tag = req.HardwareKeyTag
			var time time.Time

			time, ok := Nonces[nonce]
			fmt.Println(time)

			if ok {
				Nonces.deleteNonce(nonce)
			} else {
				ctx.Status(403)
				return ctx.JSON(oidfed.ErrorInvalidRequest("invalid nonce"))
			}

			info, err := store.WalletInstance(hardware_key_tag)
			if err != nil {
				return ctx.JSON(oidfed.ErrorInvalidRequest("hardware key tag not found"))
			}

			var key_attestion = info.KeyAttestation

			// Build the Wallet Attestation, containing multiple Wallet Attestations that have an Attestation
			// attestation := &Attestation{KeyAttestation: key_attestion}
			// wallet_attestation := &WalletAttestation{Format: "json", Attestation: *attestation}
			wallet_attestation := &WalletAttestation{Format: "json", Attestation: key_attestion}

			attestations := &Attestations{WalletAttestations: []WalletAttestation{*wallet_attestation}}

			// message := &JWT{Message: "Hello " + nonce + " " + req.Cnf}
			// jwt, err := fed.GeneralJWTSigner.JWT(message, "oauth-client-attestation+jwt")
			// if err != nil {
			// 	return nil
			// }
			// return ctx.Send(jwt)

			return ctx.JSON(attestations)
		},
	)
}
