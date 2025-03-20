package converter

import (
	"encoding/json"
	"fmt"

	"github.com/tbe-team/raybot/internal/controller/http/oas/gen"
	"github.com/tbe-team/raybot/internal/model"
)

func ConvertCommandToResponse(command model.Command) (gen.CommandResponse, error) {
	inputsJSON, err := json.Marshal(command.Inputs)
	if err != nil {
		return gen.CommandResponse{}, fmt.Errorf("marshal command inputs: %w", err)
	}

	var inputs gen.CommandResponse_Inputs
	if err := inputs.UnmarshalJSON(inputsJSON); err != nil {
		return gen.CommandResponse{}, fmt.Errorf("unmarshal command inputs: %w", err)
	}

	return gen.CommandResponse{
		Id:          command.ID,
		Type:        command.Type.String(),
		Status:      command.Status.String(),
		Source:      command.Source.String(),
		Inputs:      inputs,
		Error:       command.Error,
		CreatedAt:   command.CreatedAt,
		CompletedAt: command.CompletedAt,
	}, nil
}
