/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.s7.readwrite;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;

public class ManualS7CounterTest {

    public static void main(String[] args) throws Exception {

        //Maybe a S7-1200
        String URL1 = "s7://192.168.23.30";
        
        //S7-400
        String URL2 = "s7://192.168.0.47?remote-rack=0&"
                + "remote-slot=3&"
                + "controller-type=S7_400&read-timeout=8&"
                + "ping=false&ping-time=2&retry-time=3";
        
        
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(URL2)) {
            final PlcReadRequest.Builder readBuilder = connection.readRequestBuilder();
            
            final PlcWriteRequest.Builder writeBuilder = connection.writeRequestBuilder();

            writeBuilder.addTagAddress("counter-0", "%C0:COUNTER", Integer.decode("0x0123"));            
            writeBuilder.addTagAddress("counter-1", "%C9:COUNTER", Integer.decode("0x0456"));            
            writeBuilder.addTagAddress("counter-2", "%C12:COUNTER", Integer.decode("0x0789"));
            writeBuilder.addTagAddress("counter-3", "%C18:COUNTER", Integer.decode("0x0012"));            
            
            final PlcWriteRequest writeRequest = writeBuilder.build();
            final PlcWriteResponse writeResposne = writeRequest.execute().get();
            
            if ( writeResposne.getResponseCode("counter-3") == PlcResponseCode.OK ){
                System.out.println("Write the counter");
            } else {
                System.out.println("Problems....");
            }
                                    
            readBuilder.addTagAddress("counter-0", "%C0:COUNTER"); // Set this counter to 123
            readBuilder.addTagAddress("counter-1", "%C9:COUNTER"); // Set this counter to 456
            readBuilder.addTagAddress("counter-2", "%C12:COUNTER"); // Set this counter to 789
            readBuilder.addTagAddress("counter-3", "%C18:COUNTER"); // Set this counter to 012
            readBuilder.addTagAddress("counters",  "%C0:COUNTER[20]");  // Set this counter to 000
            
            final PlcReadRequest readRequest = readBuilder.build();

            final PlcReadResponse readResponse = readRequest.execute().get();

            System.out.println(readResponse);
            
            byte[] responseBytes = readResponse.getPlcValue("counter-0").getRaw();
            
            short bcd_0 =  readResponse.getShort("counter-0");
            short bcd_1 =  readResponse.getShort("counter-1");
            short bcd_2 =  readResponse.getShort("counter-2");
            short bcd_3 =  readResponse.getShort("counter-3");
            System.out.println("counter-0 = " + convertShortToBcd(bcd_0));
            System.out.println("counter-1 = " + convertShortToBcd(bcd_1));
            System.out.println("counter-2 = " + convertShortToBcd(bcd_2));
            System.out.println("counter-3 = " + convertShortToBcd(bcd_3));            
            
            PlcValue plcValues = readResponse.getPlcValue("counters");

            System.out.println(plcValues.toString());
            
            if (plcValues.isList()) {            
                System.out.println(plcValues.getIndex(0).getShort() == bcd_0);
                System.out.println(plcValues.getIndex(9).getShort() == bcd_1);
                System.out.println(plcValues.getIndex(12).getShort() == bcd_2);                
                System.out.println(plcValues.getIndex(18).getShort() == bcd_3);                                                            
            }
          
            
        }
    }
    
    
    private static short convertShortToBcd(short incomingShort) {
        return (short) (((incomingShort >> 8) & 0x0F) * 100 +
            ((incomingShort >> 4) & 0x0F) * 10 +
            (incomingShort & 0x000f));
    }    

}
