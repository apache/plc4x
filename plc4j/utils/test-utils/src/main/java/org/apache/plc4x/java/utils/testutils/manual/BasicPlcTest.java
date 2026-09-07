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
package org.apache.plc4x.java.utils.testutils.manual;

import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcRawByteArray;
import org.apache.plc4x.java.spi.values.PlcStruct;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.junit.jupiter.api.Assertions;

import java.util.*;

public class BasicPlcTest {

    private final String connectionString;
    private final PlcAuthentication authentication;
    private final boolean testRead;
    private final boolean testWrite;
    private final boolean enableSingleItemTests;
    private final List<TestCase> testCases;
    private final boolean shuffleMultiItemRequests;
    private final int numRandomMultiItemRequests;

    public BasicPlcTest(String connectionString) {
        this(connectionString, null, true, true, true, true, 100);
    }

    public BasicPlcTest(String connectionString, PlcAuthentication authentication, boolean testRead, boolean testWrite, boolean enableSingleItemTests, boolean shuffleMultiItemRequests, int numRandomMultiItemRequests) {
        this.connectionString = connectionString;
        this.authentication = authentication;
        this.testRead = testRead;
        this.testWrite = testWrite;
        this.enableSingleItemTests = enableSingleItemTests;
        this.shuffleMultiItemRequests = shuffleMultiItemRequests;
        this.numRandomMultiItemRequests = numRandomMultiItemRequests;
        testCases = new ArrayList<>();
    }

    public BasicPlcTest addTestCase(String address, PlcValue expectedReadValue) {
        testCases.add(new TestCase(address, PlcResponseCode.OK, expectedReadValue));
        return this;
    }

    public BasicPlcTest addTestCase(String address, PlcResponseCode expectedResponseCode) {
        testCases.add(new TestCase(address, expectedResponseCode, null));
        return this;
    }

    public BasicPlcTest addTestCase(String address, PlcResponseCode expectedResponseCode, PlcValue expectedReadValue) {
        testCases.add(new TestCase(address, expectedResponseCode, expectedReadValue));
        return this;
    }

    public void run() throws Exception {
        // Timing accumulators for single item tests
        long singleReadTotalMs = 0;
        int singleReadCount = 0;
        long singleWriteTotalMs = 0;
        int singleWriteCount = 0;

        // Timing accumulators for multi item tests
        long multiReadTotalMs = 0;
        int multiReadCount = 0;
        long multiWriteTotalMs = 0;
        int multiWriteCount = 0;

        // Connection metadata for the summary
        PlcConnectionMetadata connectionMetadata = null;
        String protocolCode = "unknown";
        String protocolName = "unknown";
        String transportCode = "unknown";
        String transportName = "unknown";
        try (PlcConnection plcConnection = new DefaultPlcDriverManager().getConnection(connectionString, authentication)) {
            connectionMetadata = plcConnection.getMetadata();
            protocolCode = plcConnection.getProtocolCode();
            protocolName = plcConnection.getProtocolName();
            transportCode = plcConnection.getTransportCode();
            transportName = plcConnection.getTransportName();
            if(enableSingleItemTests) {
                System.out.println("Reading all types in separate requests");

                // Run all entries separately:
                for (TestCase testCase : testCases) {
                    String tagName = testCase.address;
                    System.out.println(" - Single Tag: " + tagName);

                    // Try reading the value from the PLC.
                    if (testRead) {
                        // Prepare the read-request
                        final PlcReadRequest readRequest = plcConnection.readRequestBuilder().addTagAddress(
                            tagName, testCase.address).build();

                        // Execute the read request
                        long startTime = System.currentTimeMillis();
                        final PlcReadResponse readResponse = readRequest.execute().get();
                        long endTime = System.currentTimeMillis();
                        long duration = endTime - startTime;
                        singleReadTotalMs += duration;
                        singleReadCount++;

                        // Check the result
                        Assertions.assertEquals(1, readResponse.getTagNames().size(), tagName);
                        Assertions.assertEquals(tagName, readResponse.getTagNames().iterator().next(), tagName);
                        Assertions.assertEquals(testCase.responseCode, readResponse.getResponseCode(tagName), tagName);
                        Assertions.assertNotNull(readResponse.getPlcValue(tagName), tagName);
                        assertPlcValuesEqual(tagName, testCase.expectedReadValue, readResponse.getPlcValue(tagName));
                        System.out.println("        - Read OK (" + duration + " ms)");
                    }

                    // Try writing the value to the PLC.
                    if (testWrite) {
                        // Prepare the write-request
                        PlcWriteRequest writeRequest = plcConnection.writeRequestBuilder().addTagAddress(
                            tagName, testCase.address, testCase.expectedReadValue).build();

                        // Execute the writ-request
                        long startTime = System.currentTimeMillis();
                        PlcWriteResponse writeResponse = writeRequest.execute().get();
                        long endTime = System.currentTimeMillis();
                        long duration = endTime - startTime;
                        singleWriteTotalMs += duration;
                        singleWriteCount++;

                        // Check the result
                        Assertions.assertEquals(testCase.responseCode, writeResponse.getResponseCode(tagName),
                            String.format("Got status %s for %s",
                                writeResponse.getResponseCode(tagName).name(), testCase.address));

                        System.out.println("        - Write OK (" + duration + " ms)");
                    }
                }
                System.out.println("Success");
            }

            if(numRandomMultiItemRequests > 0) {
                // Read all items in one big request.
                // Shuffle the list of test cases and run the test 10 times.
                System.out.println("Reading all items together in random order");
                for (int i = 0; i < numRandomMultiItemRequests; i++) {
                    System.out.println(" - run number " + i + " of " + numRandomMultiItemRequests + " with " + testCases.size() + " items");
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
                        long duration = endTime - startTime;
                        multiReadTotalMs += duration;
                        multiReadCount++;

                        // Check the result
                        Assertions.assertEquals(shuffledTestcases.size(), readResponse.getTagNames().size());
                        for (TestCase testCase : shuffledTestcases) {
                            String tagName = testCase.address;
                            Assertions.assertEquals(PlcResponseCode.OK, readResponse.getResponseCode(tagName),
                                "Tag: " + tagName);
                            Assertions.assertNotNull(readResponse.getPlcValue(tagName), "Tag: " + tagName);
                            assertPlcValuesEqual(tagName, testCase.expectedReadValue, readResponse.getPlcValue(tagName));
                        }
                        System.out.println("        - Read OK (" + duration + " ms)");
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
                        long duration = endTime - startTime;
                        multiWriteTotalMs += duration;
                        multiWriteCount++;

                        // Check the result
                        Assertions.assertEquals(shuffledTestcases.size(), writeResponse.getTagNames().size());
                        for (TestCase testCase : shuffledTestcases) {
                            String tagName = testCase.address;
                            Assertions.assertEquals(PlcResponseCode.OK, writeResponse.getResponseCode(tagName),
                                "Tag: " + tagName);
                        }
                        System.out.println("        - Write OK (" + duration + " ms)");
                    }
                }
                System.out.println("Success");
            }
        } catch (Exception e) {
            Assertions.fail(e);
        }

        // Print timing summary
        System.out.println();
        System.out.println("=== Test Summary ===");
        System.out.println("Protocol:  " + protocolName + " (" + protocolCode + ")");
        System.out.println("Transport: " + transportName + " (" + transportCode + ")");
        if (singleReadCount > 0 || singleWriteCount > 0) {
            System.out.println("Single Item:");
            if (singleReadCount > 0) {
                double avgSingleRead = (double) singleReadTotalMs / singleReadCount;
                System.out.printf("  Average Read Time:  %.2f ms (%d requests, %d ms total)%n",
                    avgSingleRead, singleReadCount, singleReadTotalMs);
            }
            if (singleWriteCount > 0) {
                double avgSingleWrite = (double) singleWriteTotalMs / singleWriteCount;
                System.out.printf("  Average Write Time: %.2f ms (%d requests, %d ms total)%n",
                    avgSingleWrite, singleWriteCount, singleWriteTotalMs);
            }
            if (singleReadCount > 0 && singleWriteCount > 0) {
                double avgSingle = (double) (singleReadTotalMs + singleWriteTotalMs) / (singleReadCount + singleWriteCount);
                System.out.printf("  Average Time:       %.2f ms%n", avgSingle);
            }
        }
        if (multiReadCount > 0 || multiWriteCount > 0) {
            int multiItems = testCases.size();
            System.out.printf("#Multi Item (each %d items):%n", multiItems);
            if (multiReadCount > 0) {
                double avgMultiRead = (double) multiReadTotalMs / multiReadCount;
                System.out.printf("  Average Read Time:  %.2f ms (%d requests, %d ms total)%n",
                    avgMultiRead, multiReadCount, multiReadTotalMs);
            }
            if (multiWriteCount > 0) {
                double avgMultiWrite = (double) multiWriteTotalMs / multiWriteCount;
                System.out.printf("  Average Write Time: %.2f ms (%d requests, %d ms total)%n",
                    avgMultiWrite, multiWriteCount, multiWriteTotalMs);
            }
            if (multiReadCount > 0 && multiWriteCount > 0) {
                double avgMulti = (double) (multiReadTotalMs + multiWriteTotalMs) / (multiReadCount + multiWriteCount);
                System.out.printf("  Average Time:       %.2f ms%n", avgMulti);
            }
        }
        System.out.println("======================");
    }

    public record TestCase(String address, PlcResponseCode responseCode, Object expectedReadValue) {
    }

    /**
     * Compare expected and actual PlcValues with proper handling for PlcList, PlcStruct, PlcRawByteArray.
     *
     * @param tagName the tag name for error messages
     * @param expected the expected value
     * @param actual the actual value read from the PLC
     */
    private void assertPlcValuesEqual(String tagName, Object expected, PlcValue actual) {
        if (actual instanceof PlcList plcList) {
            List<?> expectedValues;
            if (expected instanceof PlcList expectedPlcList) {
                expectedValues = expectedPlcList.getList();
            } else if (expected instanceof List) {
                expectedValues = (List<?>) expected;
            } else {
                Assertions.fail("Tag " + tagName + ": Got a list of values, but only expected one.");
                return;
            }
            Assertions.assertEquals(expectedValues.size(), plcList.getLength(),
                "Tag " + tagName + ": List size mismatch");
            for (int j = 0; j < expectedValues.size(); j++) {
                Object expectedElement = expectedValues.get(j);
                PlcValue actualElement = plcList.getIndex(j);
                if (expectedElement instanceof PlcValue expectedPlcValue) {
                    // Recursively compare PlcValues (handles nested structs/lists)
                    assertPlcValuesEqual(tagName + "[" + j + "]", expectedPlcValue, actualElement);
                } else {
                    Assertions.assertEquals(expectedElement, actualElement.getObject(),
                        tagName + "[" + j + "]");
                }
            }
        } else if (expected instanceof PlcStruct expectedStruct) {
            Assertions.assertInstanceOf(PlcStruct.class, actual, "Tag " + tagName + ": Expected PlcStruct but got " + actual.getClass().getSimpleName());
            PlcStruct actualStruct = (PlcStruct) actual;
            diffPlcStructs(tagName, expectedStruct, actualStruct);
        } else if (expected instanceof PlcRawByteArray plcRawByteArray) {
            byte[] expectedRawByteArray = plcRawByteArray.getRaw();
            byte[] readRawByteArray = actual.getRaw();
            Assertions.assertArrayEquals(expectedRawByteArray, readRawByteArray, "Tag " + tagName);
        } else if (expected instanceof PlcValue plcValue) {
            Assertions.assertEquals(plcValue.getObject(), actual.getObject(), "Tag " + tagName);
        } else if (expected != null) {
            Assertions.assertEquals(expected.toString(), actual.getObject().toString(), "Tag " + tagName);
        }
    }

    /**
     * Compare two PlcStruct objects with detailed diff output.
     *
     * @param tagName the tag name for error messages
     * @param expected the expected PlcStruct
     * @param actual the actual PlcStruct
     */
    private void diffPlcStructs(String tagName, PlcStruct expected, PlcStruct actual) {
        Map<String, ? extends PlcValue> expectedMap = expected.getStruct();
        Map<String, ? extends PlcValue> actualMap = actual.getStruct();

        StringBuilder sb = new StringBuilder();

        // Missing / extra keys
        var missing = new TreeSet<>(expectedMap.keySet());
        missing.removeAll(actualMap.keySet());
        var extra = new TreeSet<>(actualMap.keySet());
        extra.removeAll(expectedMap.keySet());
        if (!missing.isEmpty()) sb.append("Missing keys in actual: ").append(missing).append('\n');
        if (!extra.isEmpty()) sb.append("Extra keys in actual: ").append(extra).append('\n');

        // Value diffs for common keys
        var common = new TreeSet<>(expectedMap.keySet());
        common.retainAll(actualMap.keySet());
        for (String key : common) {
            PlcValue ev = expectedMap.get(key);
            PlcValue av = actualMap.get(key);

            // Handle nested structs recursively
            if (ev instanceof PlcStruct expectedNested && av instanceof PlcStruct actualNested) {
                try {
                    diffPlcStructs(tagName + "." + key, expectedNested, actualNested);
                } catch (AssertionError e) {
                    sb.append("Nested struct '").append(key).append("' differs:\n")
                        .append("  ").append(e.getMessage()).append('\n');
                }
            } else if (ev instanceof PlcList expectedList && av instanceof PlcList actualList) {
                // Compare lists element by element
                if (expectedList.getLength() != actualList.getLength()) {
                    sb.append("Key '").append(key).append("': list length mismatch\n")
                        .append("  expected: ").append(expectedList.getLength()).append('\n')
                        .append("  actual  : ").append(actualList.getLength()).append('\n');
                } else {
                    for (int i = 0; i < expectedList.getLength(); i++) {
                        Object evObj = expectedList.getIndex(i).getObject();
                        Object avObj = actualList.getIndex(i).getObject();
                        if (!Objects.equals(evObj, avObj)) {
                            sb.append("Key '").append(key).append("[").append(i).append("]':\n")
                                .append("  expected: ").append(pretty(evObj)).append('\n')
                                .append("  actual  : ").append(pretty(avObj)).append('\n');
                        }
                    }
                }
            } else {
                Object evObj = ev != null ? ev.getObject() : null;
                Object avObj = av != null ? av.getObject() : null;
                if (!Objects.equals(evObj, avObj)) {
                    sb.append("Key '").append(key).append("':\n")
                        .append("  expected: ").append(pretty(evObj)).append('\n')
                        .append("  actual  : ").append(pretty(avObj)).append('\n');

                    // If numerics look equal but types differ, call that out explicitly
                    if (evObj instanceof Number en && avObj instanceof Number an &&
                        Double.compare(en.doubleValue(), an.doubleValue()) == 0 &&
                        !evObj.getClass().equals(avObj.getClass())) {
                        sb.append("  note    : numeric values equal, types differ\n");
                    }
                }
            }
        }

        Assertions.assertEquals(0, sb.length(),
            "Tag " + tagName + ": Diff detected in PlcStruct:\n" + sb);
    }

    private static String pretty(Object v) {
        return (v == null) ? "null"
            : v + " (" + v.getClass().getName() + ")";
    }

}
