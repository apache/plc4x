package ads

import (
	"context"
	"testing"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
)

func TestDiscoverer(t *testing.T) {
	discoverer := NewDiscoverer()
	discoverer.Discover(context.Background(), func(event apiModel.PlcDiscoveryItem) {
		print(event)
	})
	time.Sleep(time.Second * 30)
}
