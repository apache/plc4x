package org.apache.plc4x.nifi.record;

import java.io.IOException;

import org.apache.plc4x.java.api.messages.PlcReadResponse;

public interface Plc4xReadResponseRowCallback {
	public void processRow(PlcReadResponse result) throws IOException;
}
