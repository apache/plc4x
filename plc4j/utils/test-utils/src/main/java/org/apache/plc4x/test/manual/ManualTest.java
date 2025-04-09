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
package org.apache.plc4x.test.manual;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcRawByteArray;
import org.apache.plc4x.java.spi.values.PlcStruct;
import org.junit.jupiter.api.Assertions;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

public abstract class ManualTest {

    private final String connectionString;
    private final boolean testRead;
    private final boolean testWrite;
    private final boolean enableSingleIteTests;
    private final List<TestCase> testCases;
    private final boolean shuffleMultiItemRequests;
    private final int numRandomMultiItemRequests;

    public ManualTest(String connectionString) {
        this(connectionString, true, true, true, true, 100);
    }

    public ManualTest(String connectionString, boolean testRead, boolean testWrite, boolean enableSingleIteTests, boolean shuffleMultiItemRequests, int numRandomMultiItemRequests) {
        this.connectionString = connectionString;
        this.testRead = testRead;
        this.testWrite = testWrite;
        this.enableSingleIteTests = enableSingleIteTests;
        this.shuffleMultiItemRequests = shuffleMultiItemRequests;
        this.numRandomMultiItemRequests = numRandomMultiItemRequests;
        testCases = new ArrayList<>();
    }

    public ManualTest addTestCase(String address, PlcValue expectedReadValue) {
        testCases.add(new TestCase(address, PlcResponseCode.OK, expectedReadValue));
        return this;
    }

    public ManualTest addTestCase(String address, PlcResponseCode expectedResponseCode) {
        testCases.add(new TestCase(address, expectedResponseCode, null));
        return this;
    }

    public ManualTest addTestCase(String address, PlcResponseCode expectedResponseCode, PlcValue expectedReadValue) {
        testCases.add(new TestCase(address, expectedResponseCode, expectedReadValue));
        return this;
    }

    public void run() throws Exception {
        try (PlcConnection plcConnection = new DefaultPlcDriverManager().getConnection(connectionString)) {
            System.out.println("Reading all types in separate requests");

            if(enableSingleIteTests) {
                // Run all entries separately:
                for (TestCase testCase : testCases) {
                    String tagName = testCase.address;

                    // Try reading the value from the PLC.
                    if (testRead) {
                        System.out.println(" - Reading: " + tagName);
                        // Prepare the read-request
                        final PlcReadRequest readRequest = plcConnection.readRequestBuilder().addTagAddress(
                            tagName, testCase.address).build();

                        // Execute the read request
                        long startTime = System.currentTimeMillis();
                        final PlcReadResponse readResponse = readRequest.execute().get();
                        long endTime = System.currentTimeMillis();

                        // Check the result
                        Assertions.assertEquals(1, readResponse.getTagNames().size(), tagName);
                        Assertions.assertEquals(tagName, readResponse.getTagNames().iterator().next(), tagName);
                        Assertions.assertEquals(testCase.responseCode, readResponse.getResponseCode(tagName), tagName);
                        Assertions.assertNotNull(readResponse.getPlcValue(tagName), tagName);
                        if (readResponse.getPlcValue(tagName) instanceof PlcList) {
                            PlcList plcList = (PlcList) readResponse.getPlcValue(tagName);
                            List<?> expectedValues;
                            if (testCase.expectedReadValue instanceof PlcList) {
                                PlcList expectedPlcList = (PlcList) testCase.expectedReadValue;
                                expectedValues = expectedPlcList.getList();
                            } else if (testCase.expectedReadValue instanceof List) {
                                expectedValues = (List<?>) testCase.expectedReadValue;
                            } else {
                                Assertions.fail("Got a list of values, but only expected one.");
                                return;
                            }
                            for (int j = 0; j < expectedValues.size(); j++) {
                                if (expectedValues.get(j) instanceof PlcValue) {
                                    Assertions.assertEquals(((PlcValue) expectedValues.get(j)).getObject(),
                                        plcList.getIndex(j).getObject(), tagName + "[" + j + "]");
                                } else {
                                    Assertions.assertEquals(expectedValues.get(j), plcList.getIndex(j).getObject(),
                                        tagName + "[" + j + "]");
                                }
                            }
                        } else {
                            if (testCase.expectedReadValue instanceof PlcStruct) {
                                Assertions.assertEquals(testCase.expectedReadValue.toString(),
                                    readResponse.getPlcValue(tagName).toString(), tagName);
                            } else if (testCase.expectedReadValue instanceof PlcRawByteArray) {
                                byte[] expectedRawByteArray = ((PlcRawByteArray) testCase.expectedReadValue).getRaw();
                                byte[] readRawByteArray = ((PlcRawByteArray) readResponse.getPlcValue(tagName)).getRaw();
                                Assertions.assertArrayEquals(expectedRawByteArray, readRawByteArray);
                            } else if (testCase.expectedReadValue instanceof PlcValue) {
                                Assertions.assertEquals(
                                    ((PlcValue) testCase.expectedReadValue).getObject(),
                                    readResponse.getPlcValue(tagName).getObject(), tagName);
                            } else {
                                Assertions.assertEquals(
                                    testCase.expectedReadValue.toString(),
                                    readResponse.getPlcValue(tagName).getObject().toString(), tagName);
                            }
                        }
                    }

                    // Try writing the value to the PLC.
                    if (testWrite) {
                        System.out.println(" - Writing: " + tagName);

                        // Prepare the write request
                        PlcWriteRequest writeRequest = plcConnection.writeRequestBuilder().addTagAddress(
                            tagName, testCase.address, testCase.expectedReadValue).build();

                        // Execute the write request
                        long startTime = System.currentTimeMillis();
                        PlcWriteResponse writeResponse = writeRequest.execute().get();
                        long endTime = System.currentTimeMillis();

                        // Check the result
                        Assertions.assertEquals(testCase.responseCode, writeResponse.getResponseCode(tagName),
                            String.format("Got status %s for %s",
                                writeResponse.getResponseCode(tagName).name(), testCase.address));
                    }
                }
                System.out.println("Success");
            }

            if(numRandomMultiItemRequests > 0) {
                // Read all items in one big request.
                // Shuffle the list of test cases and run the test 10 times.
                System.out.println("Reading all items together in random order");
                for (int i = 0; i < numRandomMultiItemRequests; i++) {
                    System.out.println(" - run number " + i + " of " + numRandomMultiItemRequests);
                    final List<TestCase> shuffledTestcases = new ArrayList<>(testCases);
                    if(shuffleMultiItemRequests) {
                        Collections.shuffle(shuffledTestcases);
                    }

                    StringBuilder sb = new StringBuilder();
                    for (TestCase testCase : shuffledTestcases) {
                        sb.append(testCase.address).append(", ");
                    }
                    System.out.println("       using order: " + sb);

                    if(testRead) {
                        final PlcReadRequest.Builder builder = plcConnection.readRequestBuilder();
                        for (TestCase testCase : shuffledTestcases) {
                            String tagName = testCase.address;
                            builder.addTagAddress(tagName, testCase.address);
                        }
                        final PlcReadRequest readRequest = builder.build();

                        // Execute the read request
                        long startTime = System.currentTimeMillis();
                        final PlcReadResponse readResponse = readRequest.execute().get();
                        long endTime = System.currentTimeMillis();

                        // Check the result
                        Assertions.assertEquals(shuffledTestcases.size(), readResponse.getTagNames().size());
                        for (TestCase testCase : shuffledTestcases) {
                            String tagName = testCase.address;
                            Assertions.assertEquals(PlcResponseCode.OK, readResponse.getResponseCode(tagName),
                                "Tag: " + tagName);
                            Assertions.assertNotNull(readResponse.getPlcValue(tagName), "Tag: " + tagName);
                            if (readResponse.getPlcValue(tagName) instanceof PlcList) {
                                PlcList plcList = (PlcList) readResponse.getPlcValue(tagName);
                                PlcList expectedValues = (PlcList) testCase.expectedReadValue;
                                for (int j = 0; j < expectedValues.getLength(); j++) {
                                    Assertions.assertEquals(expectedValues.getIndex(j).toString(), plcList.getIndex(j).toString(),
                                        "Tag: " + tagName);
                                }
                            } else {
                                Assertions.assertEquals(
                                    testCase.expectedReadValue.toString(), readResponse.getPlcValue(tagName).toString(),
                                    "Tag: " + tagName);
                            }
                        }
                        System.out.println("        - Read OK (" + (endTime - startTime) + " ms)");
                    }

                    if (testWrite) {
                        final PlcWriteRequest.Builder writeBuilder = plcConnection.writeRequestBuilder();
                        for (TestCase testCase : shuffledTestcases) {
                            String tagName = testCase.address;
                            writeBuilder.addTagAddress(tagName, testCase.address, testCase.expectedReadValue);
                        }
                        final PlcWriteRequest writeRequest = writeBuilder.build();

                        // Execute the read request
                        long startTime = System.currentTimeMillis();
                        final PlcWriteResponse writeResponse = writeRequest.execute().get();
                        long endTime = System.currentTimeMillis();

                        // Check the result
                        Assertions.assertEquals(shuffledTestcases.size(), writeResponse.getTagNames().size());
                        for (TestCase testCase : shuffledTestcases) {
                            String tagName = testCase.address;
                            Assertions.assertEquals(PlcResponseCode.OK, writeResponse.getResponseCode(tagName),
                                "Tag: " + tagName);
                        }
                        System.out.println("        - Write OK (" + (endTime - startTime) + " ms)");
                    }
                }
                System.out.println("Success");
            }
        } catch (Exception e) {
            Assertions.fail(e);
        }
    }

    public static class TestCase {
        private final String address;
        private final PlcResponseCode responseCode;
        private final Object expectedReadValue;

        public TestCase(String address, PlcResponseCode responseCode, Object expectedReadValue) {
            this.address = address;
            this.responseCode = responseCode;
            this.expectedReadValue = expectedReadValue;
        }

        public String getAddress() {
            return address;
        }

        public PlcResponseCode getResponseCode() {
            return responseCode;
        }

        public Object getExpectedReadValue() {
            return expectedReadValue;
        }
    }


}
