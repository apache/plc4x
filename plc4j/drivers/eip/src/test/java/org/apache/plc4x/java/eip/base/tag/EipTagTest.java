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

package org.apache.plc4x.java.eip.base.tag;

import org.apache.plc4x.java.eip.readwrite.CIPDataTypeCode;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

public class EipTagTest {

    @Test
    public void testTagParse() {
        EipTag eipTag = EipTag.of("%A0[0..1]");

        Assertions.assertNotNull(eipTag);
        Assertions.assertEquals("%A0[0..1]", eipTag.getTag());
        Assertions.assertEquals(eipTag.getType(), CIPDataTypeCode.DINT);
        Assertions.assertEquals(eipTag.getElementNb(), 2);
    }

    /**
     * Every address form documented in the EtherNet/IP page (website/asciidoc, "Address Format")
     * has to parse the way it is documented - see GH-1481, where the docs had the data type and
     * the element count the wrong way round. If a form changes here, update the page too.
     */
    @Test
    public void testDocumentedAddressForms() {
        assertTag("myTag", "myTag", CIPDataTypeCode.DINT, 1);
        assertTag("myTag:REAL", "myTag", CIPDataTypeCode.REAL, 1);
        assertTag("myTag[0..3]", "myTag[0..3]", CIPDataTypeCode.DINT, 4);
        assertTag("myArray[3]:DINT", "myArray[3]", CIPDataTypeCode.DINT, 1);
        assertTag("myArray[0..3]:DINT", "myArray[0..3]", CIPDataTypeCode.DINT, 4);
        // The '%' prefix is optional.
        assertTag("%myTag:REAL", "%myTag", CIPDataTypeCode.REAL, 1);
    }

    private void assertTag(String address, String tag, CIPDataTypeCode type, int elementNb) {
        EipTag eipTag = EipTag.of(address);
        Assertions.assertNotNull(eipTag, address);
        Assertions.assertEquals(tag, eipTag.getTag(), address);
        Assertions.assertEquals(type, eipTag.getType(), address);
        Assertions.assertEquals(elementNb, eipTag.getElementNb(), address);
    }

}
