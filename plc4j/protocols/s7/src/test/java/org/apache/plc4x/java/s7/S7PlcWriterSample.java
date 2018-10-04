/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.s7;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class S7PlcWriterSample {

    private static final Logger logger = LoggerFactory.getLogger(S7PlcWriterSample.class);

    /**
     * Example code do demonstrate using the S7 Plc Driver.
     *
     * @param args ignored.
     */
/*    public static void main(String[] args) {
        // Create a connection to the S7 PLC (s7://{hostname/ip}/{racknumber}/{slotnumber})
        logger.info("Connecting");
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection("s7://192.168.0.1/0/0")) {
            logger.info("Connected");

            Optional<PlcWriter> writer = plcConnection.getWriter();
            // Check if this connection support reading of data.
            if (writer.isPresent()) {
                PlcWriter plcWriter = writer.get();
                PlcField inputs = plcConnection.prepareField("DATA_BLOCKS/1/2");
                //////////////////////////////////////////////////////////
                // Write synchronously ...
                // NOTICE: the ".get()" immediately lets this thread pause till
                // the response is processed and available.
                TypeSafePlcWriteResponse<Float> plcWriteResponse = plcWriter.write(
                    new TypeSafePlcWriteRequest<>(Float.class, inputs, 2.0f)).get();
                System.out.println("Written: " + plcWriteResponse.getResponseItem()
                    .orElseThrow(() -> new IllegalStateException("No response available"))
                    .getResponseCode().name());
            }
        }
        // Catch any exception or the application won't be able to finish if something goes wrong.
        catch (Exception e) {
            e.printStackTrace();
        }
        // The application would cleanly terminate after several seconds ... this just speeds things up.
        System.exit(0);
    }*/

}
