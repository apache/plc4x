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

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;
import org.junit.jupiter.params.provider.ValueSource;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * How GSDML XML maps onto the model. These are the rules the former jackson-xml mapping followed;
 * device descriptions in the field rely on them loading the same way.
 */
class GsdmlXmlTest {

    private static ProfinetISO15745Profile parse(String profileBody) throws IOException {
        String xml = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
            + "<ISO15745Profile xmlns=\"http://www.profibus.com/GSDML/2003/11/DeviceProfile\""
            + " xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\">"
            + "<ProfileHeader><ProfileIdentification>PROFINET Device Profile</ProfileIdentification>"
            + "<ISO15745Reference><ISO15745Part> 4 </ISO15745Part></ISO15745Reference></ProfileHeader>"
            + "<ProfileBody>" + profileBody + "</ProfileBody></ISO15745Profile>";
        return ProfinetISO15745Profile.parse(new ByteArrayInputStream(xml.getBytes(StandardCharsets.UTF_8)));
    }

    private static ProfinetDeviceAccessPointItem accessPoint(String attributes, String content) throws IOException {
        return parse("<ApplicationProcess><DeviceAccessPointList><DeviceAccessPointItem " + attributes + ">"
            + content + "</DeviceAccessPointItem></DeviceAccessPointList></ApplicationProcess>")
            .getProfileBody().getApplicationProcess().getDeviceAccessPointList().get(0);
    }

    private static List<ProfinetDataItem> inputDataItems(String inputContent) throws IOException {
        return parse("<ApplicationProcess><ModuleList><ModuleItem><VirtualSubmoduleList><VirtualSubmoduleItem>"
            + "<IOData><Input>" + inputContent + "</Input></IOData>"
            + "</VirtualSubmoduleItem></VirtualSubmoduleList></ModuleItem></ModuleList></ApplicationProcess>")
            .getProfileBody().getApplicationProcess().getModuleList().get(0)
            .getVirtualSubmoduleList().get(0).getIoData().getInput().get(0).getDataItemList();
    }

    @Test
    void elementTextAndAttributesFeedTheSameProperties() throws IOException {
        ProfinetISO15745Profile profile = parse("<DeviceIdentity VendorID=\"0x1\"><VendorName><Value>v</Value></VendorName></DeviceIdentity>");

        assertEquals("PROFINET Device Profile", profile.getProfileHeader().getProfileIdentification());
        assertEquals(4, profile.getProfileHeader().getIso15745Reference().getIso15745Part());
        assertEquals("0x1", profile.getProfileBody().getDeviceIdentity().getVendorId());
        assertEquals("v", profile.getProfileBody().getDeviceIdentity().getVendorName().getValue());
    }

    @Test
    void laterOccurrencesReplaceEarlierOnes() throws IOException {
        ProfinetDeviceAccessPointItem accessPoint = accessPoint("ID=\"attribute\"",
            "<ID>element</ID><ModuleInfo><Name TextId=\"first\"/></ModuleInfo><ModuleInfo><InfoText TextId=\"second\"/></ModuleInfo>");

        assertEquals("element", accessPoint.getId());
        assertNull(accessPoint.getModuleInfo().getName(), "a repeated element replaces, it does not merge");
        assertEquals("second", accessPoint.getModuleInfo().getInfoText().getTextId());
    }

    @Test
    void wrappedListsTakeEveryChildElement() throws IOException {
        ProfinetApplicationProcess process = parse("<ApplicationProcess><ModuleList><ModuleItem ID=\"a\"/><!-- c -->"
            + "<OtherItem ID=\"b\"/></ModuleList><GraphicsList/></ApplicationProcess>").getProfileBody().getApplicationProcess();

        assertEquals(2, process.getModuleList().size());
        assertEquals("b", process.getModuleList().get(1).getId());
        assertTrue(process.getGraphicsList().isEmpty());
        assertNull(process.getDeviceAccessPointList());
    }

    @Test
    void unwrappedListsAreOneRunOfSiblings() throws IOException {
        ProfinetPrimaryLanguage language = parse("<ApplicationProcess><ExternalTextList><PrimaryLanguage>"
            + "<Text TextId=\"a\"/><Text TextId=\"b\"/><Language/><Text TextId=\"c\"/>"
            + "</PrimaryLanguage></ExternalTextList></ApplicationProcess>")
            .getProfileBody().getApplicationProcess().getExternalTextList().getPrimaryLanguage();

        assertEquals(1, language.getText().size(), "a second run replaces the first");
        assertEquals("c", language.getText().get(0).getTextId());
    }

    @Test
    void scalarsReadLikeBefore() throws IOException {
        ProfinetDeviceAccessPointItem accessPoint = accessPoint(
            "MinDeviceInterval=\" +32 \" FixedInSlots=\"\" DeviceAccessSupported=\"TRUE\" MultipleWriteSupported=\"\"", "");

        assertEquals(32, accessPoint.getMinDeviceInterval());
        assertEquals(0, accessPoint.getFixedInSlots());
        assertEquals(0, accessPoint.getObjectUUIDLocalIndex());
        assertTrue(accessPoint.isDeviceAccessSupported());
        assertFalse(accessPoint.isMultipleWriteSupported());
        assertNull(accessPoint.getPhysicalSlots());
    }

    @ParameterizedTest
    @CsvSource({"MinDeviceInterval, 0x10", "MinDeviceInterval, 2147483648", "MinDeviceInterval, 1..3",
        "DeviceAccessSupported, 1", "DeviceAccessSupported, yes"})
    void malformedScalarsFailTheFile(String attribute, String value) {
        assertThrows(IOException.class, () -> accessPoint(attribute + "=\"" + value + "\"", ""));
    }

    @Test
    void nilMakesAPropertyNull() throws IOException {
        assertNull(parse("<DeviceIdentity><VendorName xsi:nil=\"true\" Value=\"v\"/></DeviceIdentity>")
            .getProfileBody().getDeviceIdentity().getVendorName());
    }

    @Test
    void textWhereAnElementIsExpectedFailsTheFile() {
        assertThrows(IOException.class, () -> parse("<DeviceIdentity><VendorName>v</VendorName></DeviceIdentity>"));
        assertThrows(IOException.class, () -> parse("<ApplicationProcess><LogBookEntryList><LogBookEntryItem>"
            + "<ErrorCode2Value Name=\"n\"/></LogBookEntryItem></LogBookEntryList></ApplicationProcess>"));
    }

    @Test
    void strictClassesStillReadWhatTheyKnow() throws IOException {
        List<ProfinetDataItem> items = inputDataItems("\n  <DataItem DataType=\"Unsigned8\" TextId=\"t\" UseAsBits=\"true\">"
            + "<BitDataItem BitOffset=\"0\" TextId=\"b0\"/><BitDataItem BitOffset=\"1\" TextId=\"b1\"/></DataItem>\n");

        assertEquals("Unsigned8", items.get(0).getDataType());
        assertEquals(1, items.get(0).getBitDataItem().get(1).getBitOffset());
    }

    @ParameterizedTest
    @ValueSource(strings = {
        "<DataItem DataType=\"Unsigned8\" Subordinate=\"true\"/>",
        "<DataItem DataType=\"Unsigned8\"><Unknown/></DataItem>",
        "<DataItem DataType=\"Unsigned8\">text</DataItem>",
        "<DataItem DataType=\"Unsigned8\"/><Channel/>"})
    void strictClassesRejectWhatTheyDoNotKnow(String inputContent) {
        assertThrows(IOException.class, () -> inputDataItems(inputContent));
    }

    @ParameterizedTest
    @ValueSource(strings = {"", "not xml", "<ISO15745Profile>", "<?xml version=\"1.0\"?><!DOCTYPE ISO15745Profile><ISO15745Profile/>"})
    void unreadableDocumentsAreIoErrors(String xml) {
        assertThrows(IOException.class,
            () -> ProfinetISO15745Profile.parse(new ByteArrayInputStream(xml.getBytes(StandardCharsets.UTF_8))));
    }

    @Test
    void externalEntitiesAreNeverResolved(@TempDir Path dir) throws IOException {
        Path secret = Files.writeString(dir.resolve("secret.txt"), "secret");
        Path gsdml = Files.writeString(dir.resolve("gsdml.xml"), "<?xml version=\"1.0\"?>"
            + "<!DOCTYPE ISO15745Profile [<!ENTITY xxe SYSTEM \"" + secret.toUri() + "\">]>"
            + "<ISO15745Profile><ProfileHeader><ProfileName>&xxe;</ProfileName></ProfileHeader></ISO15745Profile>");

        IOException thrown = assertThrows(IOException.class, () -> ProfinetISO15745Profile.parse(gsdml));
        assertFalse(String.valueOf(thrown.getMessage()).contains("secret"));
    }

    @Test
    void anyRootElementIsRead(@TempDir Path dir) throws IOException {
        Path file = Files.writeString(dir.resolve("other.xml"), "<Other><ProfileBody/></Other>");

        ProfinetISO15745Profile profile = ProfinetISO15745Profile.parse(file);

        assertNull(profile.getProfileHeader());
        assertNull(profile.getProfileBody().getDeviceIdentity());
    }
}
