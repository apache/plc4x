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
package org.apache.plc4x.test.manual;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.values.PlcList;
import org.junit.jupiter.api.Assertions;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

public abstract class ManualTest {

    private final String connectionString;
    private final List<TestCase> testCases;

    public ManualTest(String connectionString) {
        this.connectionString = connectionString;
        testCases = new ArrayList<>();
    }

    public void addTestCase(String address, Object expectedReadValue) {
        testCases.add(new TestCase(address, expectedReadValue, null));
    }

    public void run() throws Exception {
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection(connectionString)) {
            System.out.println("Reading all types in separate requests");
            // Run all entries separately:
            for (TestCase testCase : testCases) {
                String fieldName = testCase.address;
                // Prepare the read-request
                final PlcReadRequest readRequest = plcConnection.readRequestBuilder().addItem(fieldName, testCase.address).build();

                // Execute the read request
                final PlcReadResponse readResponse = readRequest.execute().get();

                // Check the result
                Assertions.assertEquals(1, readResponse.getFieldNames().size(), fieldName);
                Assertions.assertEquals(fieldName, readResponse.getFieldNames().iterator().next(), fieldName);
                Assertions.assertEquals(PlcResponseCode.OK, readResponse.getResponseCode(fieldName), fieldName);
                Assertions.assertNotNull(readResponse.getPlcValue(fieldName), fieldName);
                if(readResponse.getPlcValue(fieldName) instanceof PlcList) {
                    PlcList plcList = (PlcList) readResponse.getPlcValue(fieldName);
                    List<Object> expectedValues = (List<Object>) testCase.expectedReadValue;
                    for (int j = 0; j < expectedValues.size(); j++) {
                        Assertions.assertEquals(expectedValues.get(j), plcList.getIndex(j).getObject(), fieldName + "[" + j + "]");
                    }
                } else {
                    Assertions.assertEquals(
                        testCase.expectedReadValue.toString(), readResponse.getPlcValue(fieldName).getObject().toString(), fieldName);
                }
            }
            System.out.println("Success");


            // Read all items in one big request.
            // Shuffle the list of test cases and run the test 10 times.
            System.out.println("Reading all items together in random order");
            for (int i = 0; i < 100; i++) {
                System.out.println(" - run number " + i + " of " + 100);
                final List<TestCase> shuffledTestcases = new ArrayList<>(testCases);
                Collections.shuffle(shuffledTestcases);

                StringBuilder sb = new StringBuilder();
                for (TestCase testCase : shuffledTestcases) {
                    sb.append(testCase.address).append(", ");
                }
                System.out.println("       using order: " + sb.toString());

                final PlcReadRequest.Builder builder = plcConnection.readRequestBuilder();
                for (TestCase testCase : shuffledTestcases) {
                    String fieldName = testCase.address;
                    builder.addItem(fieldName, testCase.address);
                }
                final PlcReadRequest readRequest = builder.build();

                // Execute the read request
                final PlcReadResponse readResponse = readRequest.execute().get();

                // Check the result
                Assertions.assertEquals(shuffledTestcases.size(), readResponse.getFieldNames().size());
                for (TestCase testCase : shuffledTestcases) {
                    String fieldName = testCase.address;
                    Assertions.assertEquals(PlcResponseCode.OK, readResponse.getResponseCode(fieldName));
                    Assertions.assertNotNull(readResponse.getPlcValue(fieldName));
                    if (readResponse.getPlcValue(fieldName) instanceof PlcList) {
                        PlcList plcList = (PlcList) readResponse.getPlcValue(fieldName);
                        List<Object> expectedValues = (List<Object>) testCase.expectedReadValue;
                        for (int j = 0; j < expectedValues.size(); j++) {
                            Assertions.assertEquals(expectedValues.get(j), plcList.getIndex(j).getObject());
                        }
                    } else {
                        Assertions.assertEquals(
                            testCase.expectedReadValue.toString(), readResponse.getPlcValue(fieldName).getObject().toString());
                    }
                }
            }
            System.out.println("Success");
        } catch (Exception e) {
            Assertions.fail(e);
        }
    }

    public static class TestCase {
        private final String address;
        private final Object expectedReadValue;
        private final Object writeValue;

        public TestCase(String address, Object expectedReadValue, Object writeValue) {
            this.address = address;
            this.expectedReadValue = expectedReadValue;
            this.writeValue = writeValue;
        }

        public String getAddress() {
            return address;
        }

        public Object getExpectedReadValue() {
            return expectedReadValue;
        }

        public Object getWriteValue() {
            return writeValue;
        }
    }


}
