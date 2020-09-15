package org.apache.plc4x.nifi.util;

public enum MODBUS_DATATYPE {   
	//TODO maybe there is a better way to do the mapping within PLC4X libraries
	//got from  https://plc4x.apache.org/users/protocols/modbus.html;
	
	//TODO missing sintax 1x, 3x, 6x, etc. See documentation.

	COIL ("coil", PLC4X_DATA_TYPE.BOOL),
	DISCRETE_INPUT ("discrete-input", PLC4X_DATA_TYPE.BOOL),
	HOLDING_REGISTER ("holding-register", PLC4X_DATA_TYPE.SHORT),
	INPUT_REGISTER ("input-register", PLC4X_DATA_TYPE.SHORT),
	EXTENDED_REGISTER ("extended-register", PLC4X_DATA_TYPE.SHORT);

	private String addressId;
	private PLC4X_DATA_TYPE plc4xDataTypeEnumMap;

	private MODBUS_DATATYPE (String addressId, PLC4X_DATA_TYPE map) {
		this.addressId = addressId;
		this.plc4xDataTypeEnumMap = map;	
	}

	public PLC4X_DATA_TYPE getPlc4xDataTypeEnumMap() {
		return plc4xDataTypeEnumMap;
	}

	public void setPlc4xDataTypeEnumMap(PLC4X_DATA_TYPE plc4xDataTypeEnumMap) {
		this.plc4xDataTypeEnumMap = plc4xDataTypeEnumMap;
	}

	public String getAddressId() {
		return addressId;
	}

	public void setAddressId(String addressId) {
		this.addressId = addressId;
	}
	
	public static PLC4X_DATA_TYPE getPlc4xDatatypeById(String addressId) {
		for (MODBUS_DATATYPE element: MODBUS_DATATYPE.values()) {
			if (addressId.equals(element.getAddressId()))
				return element.getPlc4xDataTypeEnumMap();
		}
		return null;
	}
	
}
