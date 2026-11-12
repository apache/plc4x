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

import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

import java.io.IOException;
import java.io.InputStream;
import java.util.List;
import java.util.Objects;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
public class ProfinetGSDMLParseTest {

    private ProfinetISO15745Profile gsdml = null;

    @BeforeAll
    public void setUp() throws IOException {
        try (InputStream in = Objects.requireNonNull(getClass().getClassLoader().getResourceAsStream("gsdml.xml"))) {
            this.gsdml = ProfinetISO15745Profile.parse(in);
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

    @Test
    public void readGsdmlFileDeepValues() {
        ProfinetApplicationProcess process = this.gsdml.getProfileBody().getApplicationProcess();
        ProfinetDeviceAccessPointItem accessPoint = process.getDeviceAccessPointList().get(0);
        assertEquals(2, process.getDeviceAccessPointList().size());
        assertEquals("PLC4X_1", accessPoint.getId());
        assertEquals("0..32", accessPoint.getPhysicalSlots());
        assertEquals(32, accessPoint.getMinDeviceInterval());
        assertEquals(1, accessPoint.getObjectUUIDLocalIndex());
        assertEquals("", accessPoint.getCertificationInfo().getApplicationClass());
        assertEquals("PLC4X-Icon", accessPoint.getGraphics().getGraphicItemRef().getGraphicItemTarget());
        ProfinetInterfaceSubmoduleItem interfaceModule = accessPoint.getSystemDefinedSubmoduleList().getInterfaceSubmodules().get(0);
        assertEquals(32768, interfaceModule.getSubslotNumber());
        assertTrue(interfaceModule.isDcpBoundarySupported());
        assertFalse(interfaceModule.isNetworkComponentDiagnosisSupported());
        assertEquals("1 2 4 8 16 32 64 128 256 512", interfaceModule.getApplicationRelations().getTimingProperties().getReductionRatio());
        ProfinetPortSubmoduleItem port = accessPoint.getSystemDefinedSubmoduleList().getPortSubmodules().get(0);
        assertEquals(3, port.getMauTypeItem().size());
        assertEquals("2", port.getMauTypeItem().get(1).getValue());

        assertEquals(3, process.getModuleList().size());
        List<ProfinetDataItem> dataItems = process.getModuleList().get(0).getVirtualSubmoduleList().get(0)
            .getIoData().getInput().get(0).getDataItemList();
        assertEquals(2, dataItems.size());
        assertTrue(dataItems.get(0).isUseAsBits());
        assertEquals(8, dataItems.get(0).getBitDataItem().size());
        assertEquals(7, dataItems.get(0).getBitDataItem().get(7).getBitOffset());
        assertEquals("Float32", dataItems.get(1).getDataType());
        assertNull(dataItems.get(1).getBitDataItem());
        assertEquals(10, process.getModuleList().get(2).getInputDataLength());

        assertEquals("1", process.getLogBookEntryList().get(0).getStatus());
        assertEquals("PLC4X-Icon", process.getGraphicsList().get(0).getId());
        assertEquals(16, process.getExternalTextList().getPrimaryLanguage().getText().size());
    }

    /**
     * Real vendor files: a few values from deep inside each, as the former jackson-based parser read them.
     */
    @ParameterizedTest
    @CsvSource({
        "GSDML-V2.2-Advantech-ADAM6100-20111216.xml, 0x01C6, 0x6100, 1, 7, 57, 0..1",
        "GSDML-V2.3-SIEMENS-SIMOCODEproVPN-20201104.xml, 0x002A, 0x0904, 3, 4, 101, 0..2",
        "GSDML-V2.33-wago-series750_753-20230115.xml, 0x011D, 0x02EE, 48, 448, 3663, 0..250",
        "GSDML-V2.43-SIEMENS-SIMOCODEproVPN-20220817.xml, 0x002A, 0x0904, 6, 7, 109, 0..2"})
    public void readVendorGsdmlFiles(String file, String vendorId, String deviceId, int accessPoints, int modules,
                                     int texts, String physicalSlots) throws IOException {
        ProfinetISO15745Profile profile;
        try (InputStream in = Objects.requireNonNull(getClass().getClassLoader().getResourceAsStream(file))) {
            profile = ProfinetISO15745Profile.parse(in);
        }
        assertEquals("Device", profile.getProfileHeader().getProfileClassID());
        assertEquals(vendorId, profile.getProfileBody().getDeviceIdentity().getVendorId());
        assertEquals(deviceId, profile.getProfileBody().getDeviceIdentity().getDeviceID());
        ProfinetApplicationProcess process = profile.getProfileBody().getApplicationProcess();
        assertEquals(accessPoints, process.getDeviceAccessPointList().size());
        assertEquals(modules, process.getModuleList().size());
        assertEquals(texts, process.getExternalTextList().getPrimaryLanguage().getText().size());
        assertEquals(physicalSlots, process.getDeviceAccessPointList().get(0).getPhysicalSlots());
    }

}
