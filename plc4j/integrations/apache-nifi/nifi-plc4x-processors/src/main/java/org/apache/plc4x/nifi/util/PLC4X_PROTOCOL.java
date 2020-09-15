package org.apache.plc4x.nifi.util;

public enum PLC4X_PROTOCOL {
		S7("s7"),
		MODBUS("modbus");
	
		String protocolId;
		
		private PLC4X_PROTOCOL(String protocolId) {
			this.protocolId = protocolId;
		}
		
		public String getProtocolId() {
			return this.protocolId;
		}
		
		public static PLC4X_PROTOCOL qualifierValueOf(String qualifierValue) {
			for (PLC4X_PROTOCOL v : values()) {
				if (v.getProtocolId().equals(qualifierValue)) {
					return v;
				}
			}
			return null;
		}
}
