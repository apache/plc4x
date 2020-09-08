package org.apache.plc4x.nifi.record;

import java.io.IOException;

import org.apache.plc4x.java.api.messages.PlcReadResponse;

public class Plc4xFullReadResponseRowCallback implements Plc4xReadResponseRowCallback {

	@Override
	public void processRow(PlcReadResponse result) throws IOException {
		// do nothing
	}

}
