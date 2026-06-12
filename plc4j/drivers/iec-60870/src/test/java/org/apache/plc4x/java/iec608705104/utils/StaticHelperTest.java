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
package org.apache.plc4x.java.iec608705104.utils;

import org.apache.plc4x.java.iec608705104.readwrite.utils.StaticHelper;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class StaticHelperTest {

    @Test
    void finishedWhenLessThanOneByteRemains() {
        // 'finished' is used as the array terminator for the APDUs helper
        // type — it must report 'done' as soon as the buffer can no longer
        // hold another APDU byte (parsing past the end blows up).
        ReadBuffer buf = mock(ReadBuffer.class);
        when(buf.getRemainingBits()).thenReturn(7);
        assertTrue(StaticHelper.finished(buf));

        when(buf.getRemainingBits()).thenReturn(0);
        assertTrue(StaticHelper.finished(buf));
    }

    @Test
    void notFinishedWhenAtLeastOneByteRemains() {
        ReadBuffer buf = mock(ReadBuffer.class);
        when(buf.getRemainingBits()).thenReturn(8);
        assertFalse(StaticHelper.finished(buf));
    }

}
