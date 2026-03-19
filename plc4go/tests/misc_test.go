package tests

import (
	"encoding/binary"
	"testing"

	"github.com/apache/plc4x/plc4go/spi/codegen"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestByteOrderUpcast(t *testing.T) {
	t.Run("little endian", func(t *testing.T) {
		rw := utils.UpcastReaderArgs(codegen.WithByteOrder(binary.LittleEndian))
		fc := codegen.FieldCommons[any]{}
		byteOrder := fc.ExtractByteOrder(rw...)
		require.NotNil(t, byteOrder)
		assert.Equal(t, binary.LittleEndian, *byteOrder)
	})
	t.Run("big endian", func(t *testing.T) {
		rw := utils.UpcastReaderArgs(codegen.WithByteOrder(binary.BigEndian))
		fc := codegen.FieldCommons[any]{}
		byteOrder := fc.ExtractByteOrder(rw...)
		require.NotNil(t, byteOrder)
		assert.Equal(t, binary.BigEndian, *byteOrder)
	})
}
