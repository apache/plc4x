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

package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.io.IOException;
import java.io.InputStreamReader;
import java.util.Objects;

import static org.junit.jupiter.api.Assertions.assertEquals;

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
public class ProfinetGSDMLParseTest {

    private ProfinetISO15745Profile gsdml = null;

    @BeforeAll
    public void setUp() {
        try {
            XmlMapper xmlMapper = new XmlMapper();
            this.gsdml = xmlMapper.readValue(new InputStreamReader(Objects.requireNonNull(getClass().getClassLoader().getResourceAsStream("gsdml.xml"))), ProfinetISO15745Profile.class);
        } catch(IOException e) {
            assert false;
        }
    }

    @Test
    public void readGsdmlFile()  {
        assertEquals(this.gsdml.getProfileBody().getDeviceIdentity().getVendorName().getValue(), "Apache PLC4X");
    }

    @Test
    public void readGsdmlFileStartupMode()  {
        ProfinetInterfaceSubmoduleItem interfaceModule = (ProfinetInterfaceSubmoduleItem) this.gsdml.getProfileBody().getApplicationProcess().getDeviceAccessPointList().get(0).getSystemDefinedSubmoduleList().getInterfaceSubmodules().get(0);
        assertEquals(interfaceModule.getApplicationRelations().getStartupMode(), "Advanced");
    }

}
