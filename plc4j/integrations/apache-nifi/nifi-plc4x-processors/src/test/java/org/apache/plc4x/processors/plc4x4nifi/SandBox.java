package org.apache.plc4x.processors.plc4x4nifi;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.utils.connectionpool.PooledPlcDriverManager;

public class SandBox {
	public static void main(String[] args) throws PlcConnectionException {
    	PlcConnection a = new PooledPlcDriverManager().getConnection("s7://10.105.143.1:102?remote-rack=0&remote-slot=0&controller-type=S7_300");
		
	}
}
