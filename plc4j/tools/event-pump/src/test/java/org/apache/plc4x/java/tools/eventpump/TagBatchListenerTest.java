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

package org.apache.plc4x.java.tools.eventpump;

import org.apache.plc4x.java.api.PlcConnectionFactory;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.tools.eventpump.triggers.TimerTrigger;
import org.junit.jupiter.api.Test;

import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;

class TagBatchListenerTest {

    @Test
    void testDefaultOnErrorMethod() {
        PlcConnectionFactory mockManager = mock(PlcConnectionFactory.class);

        // Create a real batch for testing
        TagBatch batch = TagBatch.builder()
                .withBatchId("test-batch")
                .withConnectionFactory(mockManager)
                .withConnectionString("test://localhost")
                .addTagAddress("test", "MAIN.test")
                .withTrigger(new TimerTrigger(1, TimeUnit.SECONDS))
                .build();

        Exception testException = new RuntimeException("Test error");

        // Create a listener that uses the default onError method
        TagBatch.TagBatchListener listener = new TagBatch.TagBatchListener() {
            @Override
            public void onTagsFetched(TagBatch b, PlcReadResponse response) {
                // Not used in this test
            }
        };

        // The default onError method should just log and not throw an exception
        assertDoesNotThrow(() -> {
            listener.onError(batch, testException);
        });

        batch.close();
    }

    @Test
    void testListenerCanOverrideOnError() {
        PlcConnectionFactory mockManager = mock(PlcConnectionFactory.class);

        // Create a real batch for testing
        TagBatch batch = TagBatch.builder()
                .withBatchId("test-batch")
                .withConnectionFactory(mockManager)
                .withConnectionString("test://localhost")
                .addTagAddress("test", "MAIN.test")
                .withTrigger(new TimerTrigger(1, TimeUnit.SECONDS))
                .build();

        Exception testException = new RuntimeException("Test error");
        final boolean[] onErrorCalled = {false};

        TagBatch.TagBatchListener listener = new TagBatch.TagBatchListener() {
            @Override
            public void onTagsFetched(TagBatch b, PlcReadResponse response) {
                // Not used in this test
            }

            @Override
            public void onError(TagBatch b, Throwable error) {
                onErrorCalled[0] = true;
                assertEquals("test-batch", b.getBatchId());
                assertEquals(testException, error);
            }
        };

        listener.onError(batch, testException);

        assertTrue(onErrorCalled[0], "Custom onError should have been called");

        batch.close();
    }
}
