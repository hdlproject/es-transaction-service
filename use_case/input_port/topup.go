package input_port

import "context"

type (
	TopUpRequest struct {
		UserID uint
		Amount uint64
	}

	TopUpResponse struct {
		Ok      bool
		Message string
	}

	TopUpUseCase interface {
		TopUp(ctx context.Context, request TopUpRequest) (response TopUpResponse, err error)
	}
)
