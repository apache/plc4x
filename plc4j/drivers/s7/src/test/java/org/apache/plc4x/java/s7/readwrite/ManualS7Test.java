package org.apache.plc4x.java.s7.readwrite;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.value.PlcValue;

import java.util.List;

public class ManualS7Test {

    public static void main(String[] args) throws Exception {
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection("s7-light://192.168.23.30")) {
            final PlcReadRequest.Builder readBuilder = connection.readRequestBuilder();
            readBuilder.addTagAddress("test", "%DB4:1:BYTE[100]");
            PlcReadResponse plcReadResponse = readBuilder.build().execute().get();
            PlcValue test = plcReadResponse.getPlcValue("test");
            List<? extends PlcValue> list = test.getList();
        }
    }

}
