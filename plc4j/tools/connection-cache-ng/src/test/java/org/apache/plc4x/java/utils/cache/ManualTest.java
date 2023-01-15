package org.apache.plc4x.java.utils.cache;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadResponse;

public class ManualTest {

    public static void main(String[] args) {
        CachedPlcConnectionManager cachedPlcConnectionManager = CachedPlcConnectionManager.getBuilder(new PlcDriverManager()).build();
        for (int i = 0; i < 100; i++){
            try (PlcConnection connection = cachedPlcConnectionManager.getConnection("ads:tcp://192.168.23.20?sourceAmsNetId=192.168.23.200.1.1&sourceAmsPort=65534&targetAmsNetId=192.168.23.20.1.1&targetAmsPort=851")) {
                PlcReadResponse plcReadResponse = connection.readRequestBuilder().addTagAddress("var", "MAIN.hurz_REAL").build().execute().get();
                System.out.printf("Run %d: Value: %f%n", i, plcReadResponse.getFloat("var"));
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        }
    }

}
