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
 * Example of reading partial status list (SZL).
 * SZL_ID = 0x0011, allows to identify the device (PLC).
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
        Thread.sleep(10000);
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection("s7://10.10.1.33?remote-rack=0&remote-slot=3&controller-type=S7_400")) { //(01)
            
            final PlcReadRequest.Builder readrequest = connection.readRequestBuilder(); //(02)
                       
            readrequest.addTagAddress("MySZL", "SZL_ID=16#0022;INDEX=16#0000"); //(03)
            
            final PlcReadRequest rr = readrequest.build(); //(04)
            final PlcReadResponse szlresponse = rr.execute().get(); //(05)

            if (szlresponse.getResponseCode("MySZL") == PlcResponseCode.OK){ //(06)
                
                Collection<Byte>  data = szlresponse.getAllBytes("MySZL"); //(07)
                byte[] dbytes = ArrayUtils.toPrimitive(data.toArray(new Byte[data.size()])); //(08)
                
                SZL szl = SZL.valueOf(0x0022); //(09)
                ByteBuf wb = wrappedBuffer(dbytes); //(10)
                StringBuilder sb =  szl.execute(wb); //(11)
                System.out.println(sb.toString());  //(12)
                
            } else if (szlresponse.getResponseCode("MySZL") == PlcResponseCode.NOT_FOUND){ //(13)
                System.out.println("SZL is not supported.");
            }

            Thread.sleep(2000);
            System.out.println("Bye...");

        }
    }

}
