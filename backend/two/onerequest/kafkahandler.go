package onerequest

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/kartpop/cruncan/backend/pkg/model"
)

type KafkaHandler struct {
	logger *slog.Logger
}

func NewKafkaHandler() *KafkaHandler {
	return &KafkaHandler{
		logger: slog.Default(),
	}
}

func (h *KafkaHandler) Handle(message []byte, topic string) error {

	var oneRequest model.OneRequest
	err := json.Unmarshal(message, &oneRequest)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to unmarshal message: %v", err))
		return err
	}

	h.logger.Info(fmt.Sprintf("Unmarshaled message: %v", oneRequest))

	return nil
}
