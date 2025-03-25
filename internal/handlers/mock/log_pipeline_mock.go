package mock

import (
	"time"

	models "gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type MockLogPipeline struct{}

func NewMockLogPipeline() *MockLogPipeline {
	return &MockLogPipeline{}
}

func (m *MockLogPipeline) LogStatusChange(TS time.Time, ID int64, statusFrom, statusTo models.OrderStatus) {
}
func (m *MockLogPipeline) Shutdown() {}
