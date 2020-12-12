package knxnetip

import (
	driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	"strconv"
)

func GroupAddressToString(groupAddress *driverModel.KnxGroupAddress) string {
	if groupAddress != nil {
		switch groupAddress.Child.(type) {
		case *driverModel.KnxGroupAddress3Level:
			level3 := driverModel.CastKnxGroupAddress3Level(groupAddress)
			return strconv.Itoa(int(level3.MainGroup)) + "/" + strconv.Itoa(int(level3.MiddleGroup)) + "/" + strconv.Itoa(int(level3.SubGroup))
		case *driverModel.KnxGroupAddress2Level:
			level2 := driverModel.CastKnxGroupAddress2Level(groupAddress)
			return strconv.Itoa(int(level2.MainGroup)) + "/" + strconv.Itoa(int(level2.SubGroup))
		case *driverModel.KnxGroupAddressFreeLevel:
			level1 := driverModel.CastKnxGroupAddressFreeLevel(groupAddress)
			return strconv.Itoa(int(level1.SubGroup))
		}
	}
	return ""
}
