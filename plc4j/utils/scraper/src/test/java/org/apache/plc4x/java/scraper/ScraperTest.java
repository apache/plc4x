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

package org.apache.plc4x.java.scraper;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.items.DefaultStringFieldItem;
import org.apache.plc4x.java.mock.MockDevice;
import org.apache.plc4x.java.mock.PlcMockConnection;
import org.junit.Test;
import org.mockito.Mockito;

import java.util.Map;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.when;

public class ScraperTest {

    @Test
    public void scrape() throws PlcConnectionException {
        PlcDriverManager driverManager = new PlcDriverManager();
        PlcMockConnection connection = (PlcMockConnection) driverManager.getConnection("mock:scraper");
        MockDevice mockDevice = Mockito.mock(MockDevice.class);
        connection.setDevice(mockDevice);
        when(mockDevice.read(any())).thenReturn(Pair.of(PlcResponseCode.OK, new DefaultStringFieldItem("hallo")));

        Scraper scraper = new Scraper(driverManager, "mock:scraper", 1_000, new Scraper.ResultHandler() {
            @Override
            public void handle(Map<String, Object> result) {
                System.out.println(result);
            }

            @Override
            public void handleException(Exception e) {
                System.err.println(e);
            }
        });

        scraper.run();
    }

    @Test
    public void scrape_badResponseCode_shouldHandleException() throws PlcConnectionException {
        PlcDriverManager driverManager = new PlcDriverManager();
        PlcMockConnection connection = (PlcMockConnection) driverManager.getConnection("mock:scraper");
        MockDevice mockDevice = Mockito.mock(MockDevice.class);
        connection.setDevice(mockDevice);
        when(mockDevice.read(any())).thenReturn(Pair.of(PlcResponseCode.NOT_FOUND, new DefaultStringFieldItem("hallo")));

        Scraper.ResultHandler handler = Mockito.mock(Scraper.ResultHandler.class);

        Scraper scraper = new Scraper(driverManager, "mock:scraper", 1_000, null);

        scraper.run();
    }
}