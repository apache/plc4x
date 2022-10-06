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
package org.apache.plc4x.java.utils.connectionpool2;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;

class PooledDriverManagerTest implements WithAssertions {

    @Test
    void getCachedDriverManager() throws PlcConnectionException {
        CachedDriverManager mock = Mockito.mock(CachedDriverManager.class, Mockito.RETURNS_DEEP_STUBS);
        PooledDriverManager driverManager = new PooledDriverManager(key -> mock);

        assertThat(driverManager.getCachedManagers().size()).isEqualTo(0);
        PlcConnection connection = driverManager.getConnection("abc");

        assertThat(driverManager.getCachedManagers())
            .containsValue(mock)
            .containsKey("abc")
            .hasSize(1);

        verify(mock, times(1)).getConnection("abc");
    }
}