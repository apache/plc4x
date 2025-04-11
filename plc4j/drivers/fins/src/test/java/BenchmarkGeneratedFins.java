/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.fins.readwrite.FinsTcpDriver;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.TimeUnit;

public class BenchmarkGeneratedFins {

    private static final Logger log = LoggerFactory.getLogger(BenchmarkGeneratedFins.class);

    public static void main(String[] args) throws Exception {

        Map<String, String> resultMap = new HashMap<>();
        PlcConnection mockConnection =  new DefaultPlcDriverManager().getConnection("finstcp:tcp://127.0.0.1:9600");

        PlcReadRequest.Builder requestBuilder = mockConnection.readRequestBuilder();
        String s1 = "D100";
        requestBuilder.addTagAddress( "test" , s1);

        PlcReadRequest readRequest = requestBuilder.build();
        PlcReadResponse response = readRequest.execute().get(Long.parseLong("2000"), TimeUnit.MILLISECONDS);
        for (String tagName : response.getTagNames()) {
            if (response.getResponseCode(tagName) == PlcResponseCode.OK) {
                int numValues = response.getNumberOfValues(tagName);
                // If it's just one element, output just one single line.
                log.info("{}: {}", tagName, response.getPlcValue(tagName));
                if (numValues == 1) {
                    resultMap.put(tagName, response.getPlcValue(tagName).toString());
                }
                // If it's more than one element, output each in a single row.
                else {
                    for (int i = 0; i < numValues; i++) {
                        resultMap.put(tagName + "-" + i, response.getObject(tagName, i).toString());
                    }
                }
            } else {
                log.error("Error[{}]: {}", tagName, response.getResponseCode(tagName).name());
            }
        }

    }

}
