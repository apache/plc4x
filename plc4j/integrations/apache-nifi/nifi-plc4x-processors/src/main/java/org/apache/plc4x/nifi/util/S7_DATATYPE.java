package org.apache.plc4x.nifi.util;

public enum S7_DATATYPE {   
	//TODO maybe there is a better way to do the mapping within PLC4X libraries
	//got from  org.apache.plc4x.java.s7.readwrite.types.TransportSize;

	BOOL ("BOOL", PLC4X_DATA_TYPE.BOOL),
	BYTE ("BYTE", PLC4X_DATA_TYPE.ARRAY),
	WORD ("WORD", PLC4X_DATA_TYPE.ARRAY),
	DWORD ("DWORD", PLC4X_DATA_TYPE.ARRAY),
	INT ("INT", PLC4X_DATA_TYPE.SHORT),
	LWORD ("LWORD", PLC4X_DATA_TYPE.ARRAY),
	UINT ("UINT", PLC4X_DATA_TYPE.SHORT),
	SINT ("SINT", PLC4X_DATA_TYPE.BYTE),
	USINT ("USINT", PLC4X_DATA_TYPE.BYTE),
	DINT ("DINT", PLC4X_DATA_TYPE.INT),
	UDINT ("UDINT", PLC4X_DATA_TYPE.INT),
	LINT ("LINT", PLC4X_DATA_TYPE.LONG),
	ULINT ("ULINT", PLC4X_DATA_TYPE.LONG),
	REAL ("REAL", PLC4X_DATA_TYPE.DOUBLE),
	LREAL ("LREAL", PLC4X_DATA_TYPE.DOUBLE),
	CHAR ("CHAR", PLC4X_DATA_TYPE.STRING),
	WCHAR ("WCHAR", PLC4X_DATA_TYPE.STRING),
	STRING ("STRING", PLC4X_DATA_TYPE.STRING),
	WSTRING ("WSTRING", PLC4X_DATA_TYPE.STRING),
	TIME ("TIME", PLC4X_DATA_TYPE.STRING),
	LTIME ("LTIME", PLC4X_DATA_TYPE.STRING),
	DATE ("DATE", PLC4X_DATA_TYPE.STRING),
	TIME_OF_DAY ("TIME_OF_DAY", PLC4X_DATA_TYPE.STRING),
	DATE_AND_TIME ("DATE_AND_TIME", PLC4X_DATA_TYPE.STRING);

	private String addressId;
	private PLC4X_DATA_TYPE plc4xDataTypeEnumMap;

	private S7_DATATYPE (String addressId, PLC4X_DATA_TYPE map) {
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
		for (S7_DATATYPE element: S7_DATATYPE.values()) {
			if (addressId.equals(element.getAddressId()))
				return element.getPlc4xDataTypeEnumMap();
		}
		return null;
	}
	
}
