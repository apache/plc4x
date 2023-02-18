package model

import (
	"context"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

func NoMorePathSegments(readBuffer utils.ReadBuffer, oder IntegerEncoding) bool {
	initialPos := readBuffer.GetPos()
	defer readBuffer.Reset(initialPos)

	_, err := PathSegmentParseWithBuffer(context.Background(), readBuffer, oder)
	return err != nil
}
