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
package org.apache.plc4x.java.transport.test;

import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.test.config.TestTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;

/**
 * Comprehensive tests for TestTransportInstance.
 */
class TestTransportInstanceTest {

    private TestTransportInstance instance;
    private TestTransportConfiguration config;
    private AuditLog mockAuditLog;

    @BeforeEach
    void setUp() throws TransportException {
        config = new TestTransportConfiguration();
        config.receiveBufferSize = 1024;
        mockAuditLog = mock(AuditLog.class);
        instance = new TestTransportInstance(config, mockAuditLog);
    }

    // === Happy Path Tests ===

    @Test
    void testInstanceCreation() {
        assertNotNull(instance);
        assertTrue(instance.isOpen());
        assertEquals(config, instance.getConfiguration());
    }

    @Test
    void testInitialState() throws TransportException {
        assertEquals(0, instance.getNumBytesAvailable());
        assertEquals(0, instance.getNumBytesWritten());
        assertTrue(instance.isOpen());
    }

    @Test
    void testInjectAndReadData() throws TransportException {
        byte[] testData = new byte[]{0x01, 0x02, 0x03, 0x04};
        int written = instance.injectTestData(testData);

        assertEquals(4, written);
        assertEquals(4, instance.getNumBytesAvailable());

        byte[] readData = instance.read(4);
        assertArrayEquals(testData, readData);
        assertEquals(0, instance.getNumBytesAvailable());
    }

    @Test
    void testPeekData() throws TransportException {
        byte[] testData = new byte[]{0x10, 0x20, 0x30, 0x40};
        instance.injectTestData(testData);

        byte[] peekedData = instance.peekReadableBytes(4);
        assertArrayEquals(testData, peekedData);

        // Data should still be available after peek
        assertEquals(4, instance.getNumBytesAvailable());

        // Reading should return the same data
        byte[] readData = instance.read(4);
        assertArrayEquals(testData, readData);
        assertEquals(0, instance.getNumBytesAvailable());
    }

    @Test
    void testPartialPeek() throws TransportException {
        byte[] testData = new byte[]{0x01, 0x02, 0x03, 0x04};
        instance.injectTestData(testData);

        byte[] peekedData = instance.peekReadableBytes(2);
        assertEquals(2, peekedData.length);
        assertEquals(0x01, peekedData[0]);
        assertEquals(0x02, peekedData[1]);

        // All data should still be available
        assertEquals(4, instance.getNumBytesAvailable());
    }

    @Test
    void testWriteAndRetrieve() throws TransportException {
        byte[] testData = new byte[]{0x05, 0x06, 0x07, 0x08};
        instance.write(testData);

        assertEquals(4, instance.getNumBytesWritten());

        byte[] writtenData = instance.getAllWrittenData();
        assertArrayEquals(testData, writtenData);
        assertEquals(0, instance.getNumBytesWritten());
    }

    @Test
    void testPartialRetrieveWrittenData() throws TransportException {
        byte[] testData = new byte[]{0x11, 0x22, 0x33, 0x44, 0x55};
        instance.write(testData);

        byte[] partial = instance.getWrittenData(3);
        assertEquals(3, partial.length);
        assertEquals(0x11, partial[0]);
        assertEquals(0x22, partial[1]);
        assertEquals(0x33, partial[2]);

        // Remaining data should still be available
        assertEquals(2, instance.getNumBytesWritten());
    }

    @Test
    void testMultipleWritesAndReads() throws TransportException {
        // Write first chunk
        instance.write(new byte[]{0x01, 0x02});
        assertEquals(2, instance.getNumBytesWritten());

        // Write second chunk
        instance.write(new byte[]{0x03, 0x04});
        assertEquals(4, instance.getNumBytesWritten());

        // Retrieve all
        byte[] allData = instance.getAllWrittenData();
        assertArrayEquals(new byte[]{0x01, 0x02, 0x03, 0x04}, allData);
    }

    @Test
    void testMultipleInjectsAndReads() throws TransportException {
        instance.injectTestData(new byte[]{0x10, 0x20});
        instance.injectTestData(new byte[]{0x30, 0x40});

        assertEquals(4, instance.getNumBytesAvailable());

        byte[] data = instance.read(4);
        assertArrayEquals(new byte[]{0x10, 0x20, 0x30, 0x40}, data);
    }

    @Test
    void testCloseTransport() throws TransportException {
        assertTrue(instance.isOpen());

        instance.close();
        assertFalse(instance.isOpen());

        // Closing again should not throw
        instance.close();
        assertFalse(instance.isOpen());
    }

    @Test
    void testGetNumBytesAvailableAfterClose() throws TransportException {
        instance.close();
        assertEquals(0, instance.getNumBytesAvailable());
    }

    // === Unhappy Path Tests ===

    @Test
    void testReadMoreBytesThanAvailable() throws TransportException {
        instance.injectTestData(new byte[]{0x01, 0x02});

        TransportException exception = assertThrows(TransportException.class, () ->
            instance.read(5)
        );

        assertTrue(exception.getMessage().contains("Requested 5 bytes but only 2 available"));
    }

    @Test
    void testPeekMoreBytesThanAvailable() throws TransportException {
        instance.injectTestData(new byte[]{0x01, 0x02, 0x03});

        TransportException exception = assertThrows(TransportException.class, () ->
            instance.peekReadableBytes(10)
        );

        assertTrue(exception.getMessage().contains("Requested 10 bytes but only 3 available"));
    }

    @Test
    void testReadAfterClose() throws TransportException {
        instance.injectTestData(new byte[]{0x01, 0x02});
        instance.close();

        TransportException exception = assertThrows(TransportException.class, () ->
            instance.read(2)
        );

        assertTrue(exception.getMessage().contains("Transport is closed"));
    }

    @Test
    void testPeekAfterClose() throws TransportException {
        instance.injectTestData(new byte[]{0x01, 0x02});
        instance.close();

        TransportException exception = assertThrows(TransportException.class, () ->
            instance.peekReadableBytes(2)
        );

        assertTrue(exception.getMessage().contains("Transport is closed"));
    }

    @Test
    void testWriteAfterClose() throws TransportException {
        instance.close();

        TransportException exception = assertThrows(TransportException.class, () ->
            instance.write(new byte[]{0x01, 0x02})
        );

        assertTrue(exception.getMessage().contains("Transport is closed"));
    }

    @Test
    void testReadZeroBytes() throws TransportException {
        byte[] result = instance.read(0);
        assertNotNull(result);
        assertEquals(0, result.length);
    }

    @Test
    void testPeekZeroBytes() throws TransportException {
        byte[] result = instance.peekReadableBytes(0);
        assertNotNull(result);
        assertEquals(0, result.length);
    }

    @Test
    void testWriteNullBytes() throws TransportException {
        // Should not throw
        instance.write(null);
        assertEquals(0, instance.getNumBytesWritten());
    }

    @Test
    void testWriteEmptyArray() throws TransportException {
        instance.write(new byte[]{});
        assertEquals(0, instance.getNumBytesWritten());
    }

    @Test
    void testInjectNullData() {
        int written = instance.injectTestData(null);
        assertEquals(0, written);
    }

    @Test
    void testInjectEmptyData() {
        int written = instance.injectTestData(new byte[]{});
        assertEquals(0, written);
    }

    @Test
    void testBufferOverflow() throws TransportException {
        // Fill buffer to capacity
        byte[] largeData = new byte[config.receiveBufferSize];
        for (int i = 0; i < largeData.length; i++) {
            largeData[i] = (byte) (i % 256);
        }

        int written = instance.injectTestData(largeData);
        assertEquals(config.receiveBufferSize, written);

        // Try to inject more data - should only write what fits
        byte[] extraData = new byte[]{0x01, 0x02};
        int extraWritten = instance.injectTestData(extraData);
        assertEquals(0, extraWritten); // Buffer is full
    }

    @Test
    void testWriteBufferOverflow() throws TransportException {
        // Fill write buffer to capacity
        byte[] largeData = new byte[config.receiveBufferSize];
        for (int i = 0; i < largeData.length; i++) {
            largeData[i] = (byte) (i % 256);
        }

        instance.write(largeData);

        // Try to write more - should throw exception
        TransportException exception = assertThrows(TransportException.class, () ->
            instance.write(new byte[]{0x01, 0x02})
        );

        assertTrue(exception.getMessage().contains("Could only write"));
    }

    @Test
    void testReadNegativeBytes() throws TransportException {
        byte[] result = instance.read(-1);
        assertNotNull(result);
        assertEquals(0, result.length);
    }

    @Test
    void testPeekNegativeBytes() throws TransportException {
        byte[] result = instance.peekReadableBytes(-1);
        assertNotNull(result);
        assertEquals(0, result.length);
    }

    @Test
    void testGetAllWrittenDataWhenEmpty() {
        byte[] data = instance.getAllWrittenData();
        assertNotNull(data);
        assertEquals(0, data.length);
    }

    @Test
    void testGetWrittenDataWhenEmpty() {
        byte[] data = instance.getWrittenData(10);
        assertNotNull(data);
        assertEquals(0, data.length);
    }

    @Test
    void testGetWrittenDataMoreThanAvailable() throws TransportException {
        instance.write(new byte[]{0x01, 0x02, 0x03});

        // Request more than available
        byte[] data = instance.getWrittenData(10);
        assertEquals(3, data.length);
        assertArrayEquals(new byte[]{0x01, 0x02, 0x03}, data);
    }

    @Test
    void testConcurrentReadAndWrite() throws Exception {
        // Test thread safety with concurrent operations
        Thread writeThread = new Thread(() -> {
            try {
                for (int i = 0; i < 10; i++) {
                    instance.write(new byte[]{(byte) i});
                    Thread.sleep(1);
                }
            } catch (Exception e) {
                fail("Write thread failed: " + e.getMessage());
            }
        });

        Thread injectThread = new Thread(() -> {
            try {
                for (int i = 0; i < 10; i++) {
                    instance.injectTestData(new byte[]{(byte) i});
                    Thread.sleep(1);
                }
            } catch (Exception e) {
                fail("Inject thread failed: " + e.getMessage());
            }
        });

        writeThread.start();
        injectThread.start();

        writeThread.join();
        injectThread.join();

        // Verify both operations completed
        assertEquals(10, instance.getNumBytesWritten());
        assertEquals(10, instance.getNumBytesAvailable());
    }

}
