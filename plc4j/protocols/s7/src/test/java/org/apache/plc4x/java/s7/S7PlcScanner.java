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

import org.junit.jupiter.api.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class S7PlcScanner {

    private static final Logger logger = LoggerFactory.getLogger(S7PlcScanner.class);

    @Test
    public void hurz() {

    }


    /**
     * Example code do demonstrate using the S7 Plc Driver.
     *
     * @param args ignored.
     */
/*    public static void main(String[] args) {
        // Create a connection to the S7 PLC (s7://{hostname/ip}/{racknumber}/{slotnumber})
        logger.info("Connecting");
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection("s7://10.10.64.20/1/1")) {
            logger.info("Connected");

            Optional<PlcReader> reader = plcConnection.getReader();
            // Check if this connection support reading of data.
            if (reader.isPresent()) {
                PlcReader plcReader = reader.get();

                for (MemoryArea memoryArea : MemoryArea.values()) {
                    System.out.println(memoryArea);
                    System.out.println("------------------------------------------");
                    for (int i = 0; i < 8959; i++) {
                        try {
                            PlcField field;
                            if (memoryArea == MemoryArea.DATA_BLOCKS) {
                                field = plcConnection.prepareField("DATA_BLOCKS/1/" + i);
                            } else {
                                field = plcConnection.prepareField(memoryArea.name() + "/" + i);
                            }
                            TypeSafePlcReadResponse<Byte> plcReadResponse = plcReader.read(
                                new TypeSafePlcReadRequest<>(Byte.class, field)).get();
                            Byte data = plcReadResponse.getResponseItem()
                                .orElseThrow(() -> new IllegalStateException("No response available"))
                                .getValues().get(0);
                            if (data != null && data != 0) {
                                System.out.println(String.format(
                                    "Response: Memory Area: %s Index: %d Value: %02X", memoryArea.name(), i, data));
                            }
                        } catch (Exception e) {
                            e.printStackTrace();
                        }
                    }
                }
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
