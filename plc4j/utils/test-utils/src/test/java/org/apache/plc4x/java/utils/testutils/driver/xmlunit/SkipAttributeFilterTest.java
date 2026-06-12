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
package org.apache.plc4x.java.utils.testutils.driver.xmlunit;

import org.junit.jupiter.api.Test;
import org.w3c.dom.Attr;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class SkipAttributeFilterTest {

    @Test
    void testFilterIgnoresSkipAttribute() {
        SkipAttributeFilter filter = new SkipAttributeFilter();
        Attr attr = mock(Attr.class);
        when(attr.getName()).thenReturn(SkipAttributeFilter.IGNORE_ATTRIBUTE_NAME);

        assertFalse(filter.test(attr));
    }

    @Test
    void testFilterAllowsOtherAttributes() {
        SkipAttributeFilter filter = new SkipAttributeFilter();
        Attr attr = mock(Attr.class);
        when(attr.getName()).thenReturn("other-attribute");

        assertTrue(filter.test(attr));
    }
}
