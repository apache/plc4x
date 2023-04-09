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
package org.apache.plc4x.examples.plc4j.s7event;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import java.util.Collection;
import java.util.concurrent.TimeUnit;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.simple.SimpleLogger;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.s7.readwrite.utils.StaticHelper.SZL;
import org.apache.plc4x.java.s7.utils.S7ParamErrorCode;
import org.apache.commons.codec.binary.Hex;
import org.apache.commons.lang3.ArrayUtils;
import static io.netty.buffer.Unpooled.wrappedBuffer;

/**
 * Example for capturing events generated from a Siemens S7-300, S7-400 or VIPA PLC.
 * Support for mode events ("MODE"), system events ("SYS"), user events ("USR")
 * and alarms ("ALM").
 * Each consumer shows the tags and associated values of the "map" containing
 * the event parameters.
 */
public class SzlRequest {

    private static final Logger logger = LoggerFactory.getLogger(SzlRequest.class);    
    
    /**
     * @param args the command line arguments
     */
    public static void main(String[] args) throws Exception {
        System.setProperty(SimpleLogger.DEFAULT_LOG_LEVEL_KEY, "trace");
        
        System.out.println("******************************************************************************************");
        System.out.println("Before using, take a look at:");
        System.out.println("System Software for S7-300/400.\r\nSystem and Standard Functions - Volume 1/2");
        System.out.println("Document: A5E02789976-01");
        System.out.println("Chapter 34 System Status Lists (SSL).");
        System.out.println("URL: https://cache.industry.siemens.com/dl/files/604/44240604/att_67003/v1/s7sfc_en-EN.pdf");
        System.out.println("******************************************************************************************");        
        System.out.println("* +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+\n" +
                           "* |15|14|13|12|11|10| 9| 8| 7| 6| 5| 4| 3| 2| 1|\n" +
                           "* +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+\n" +
                           "* \\__________/\\__________/\\_________________/\n" +
                           "*    Module      Number of    Number of the\n" +
                           "*    class       the partial  partial list\n" +
                           "*    list\n" +
                           "*    extract\n" +
                           "*\n" +
                           "* <b>Module Class:</b>\n" +
                           "* +--------------+-----------------+\n" +
                           "* | Module class | Coding (Binary) |\n" +
                           "* +--------------|-----------------+\n" +
                           "* |     CPU      |      0000       |\n" +
                           "* +--------------|-----------------+\n" +
                           "* |     IM       |      0100       |\n" +
                           "* +--------------|-----------------+\n" +
                           "* |     FM       |      1000       |\n" +
                           "* +--------------|-----------------+\n" +
                           "* |     CP       |      1100       |\n" +
                           "* +--------------|-----------------+");
        System.out.println("******************************************************************************************\r\n");    
        
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection("s7://10.10.1.33?remote-rack=0&remote-slot=3&controller-type=S7_400")) {
            final PlcReadRequest.Builder subscription = connection.readRequestBuilder();
            
            System.out.println("Request: SZL_ID=16#0000;INDEX=16#0000");
            
            subscription.addTagAddress("MySZL", "SZL_ID=16#0000;INDEX=16#0000");
            
            final PlcReadRequest sub = subscription.build();
            final PlcReadResponse szlresponse = sub.execute().get();

            if (szlresponse.getResponseCode("MySZL") == PlcResponseCode.OK){
                Collection<Byte>  data = szlresponse.getAllBytes("MySZL");
                byte[] dbytes = ArrayUtils.toPrimitive(data.toArray(new Byte[data.size()]));
                
                //System.out.println("DATA: \r\n" + Hex.encodeHexString(dbytes));
                //System.out.println("");
                
                SZL szl = SZL.valueOf(0x0000);
                ByteBuf wb = wrappedBuffer(dbytes);
                StringBuilder sb =  szl.execute(wb);
                System.out.println(sb.toString());                
            } else if (szlresponse.getResponseCode("MySZL") == PlcResponseCode.NOT_FOUND){
//                System.out.println("Service not found."); 
//                ByteBuf data = ((S7ByteReadResponse) szlresponse).getByteBufValues("MySSL");
//                System.out.println(S7ParamErrorCode.valueOf(data.getShort(0)) + 
//                        " : " + 
//                        S7ParamErrorCode.valueOf(data.getShort(0)).getEvent());
            }

            Thread.sleep(2000);
            //connection.close();


            System.out.println("Bye...");

        }
    }

}
