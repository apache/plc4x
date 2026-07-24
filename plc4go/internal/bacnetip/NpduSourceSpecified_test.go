package bacnetip

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

// Reproduces the routed-reply rejection: a source-specified NPDU carrying a
// ComplexACK must round-trip with its APDU intact.
func TestNPDU_SourceSpecified_APDUSurvives(t *testing.T) {
	ack := readWriteModel.NewAPDUSimpleAck(0, readWriteModel.BACnetConfirmedServiceChoice_WRITE_PROPERTY)
	control := readWriteModel.NewNPDUControl(false, false, true, false, readWriteModel.NPDUNetworkPriority_NORMAL_MESSAGE)
	snet := uint16(3001)
	slen := uint8(6)
	sadr := []uint8{192, 168, 102, 2, 0xBA, 0xC0}
	npdu := readWriteModel.NewNPDU(1, control, nil, nil, nil, &snet, &slen, sadr, nil, nil, ack)
	bvlc := readWriteModel.NewBVLCOriginalUnicastNPDU(npdu)
	data, err := bvlc.Serialize()
	require.NoError(t, err)
	t.Logf("frame: % x", data)

	reparsed, err := readWriteModel.BVLCParse[readWriteModel.BVLC](context.Background(), data)
	require.NoError(t, err)
	renpdu := reparsed.(interface{ GetNpdu() readWriteModel.NPDU }).GetNpdu()
	require.NotNil(t, renpdu.GetSourceNetworkAddress())
	assert.Equal(t, uint16(3001), *renpdu.GetSourceNetworkAddress())
	require.NotNil(t, renpdu.GetApdu(), "APDU must survive a source-specified NPDU parse")
	require.NoError(t, err)
	_, err2 := getInvokeIdFromApdu(renpdu.GetApdu())
	assert.NoError(t, err2)
}
