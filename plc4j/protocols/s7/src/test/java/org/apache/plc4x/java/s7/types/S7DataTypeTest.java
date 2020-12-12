/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.s7.types;

import org.apache.plc4x.test.FastTests;
import org.junit.Test;
import org.junit.experimental.categories.Category;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.junit.Assert.assertThat;

public class S7DataTypeTest {

/*    @Test
    @Category(FastTests.class)
    public void findMatchingTypeRedundantInformation() {
        S7DataType result = S7DataType.findMatchingType(Integer.class, S7DataType.UINT, "W");
        assertThat(result, equalTo(S7DataType.UINT));
    }

    @Test
    @Category(FastTests.class)
    public void findMatchingTypeMissingSizeCode() {
        S7DataType result = S7DataType.findMatchingType(Integer.class, S7DataType.UINT, null);
        assertThat(result, equalTo(S7DataType.UINT));
    }

    @Test
    @Category(FastTests.class)
    public void findMatchingTypeBaseTypeAndSizeCode() {
        S7DataType result = S7DataType.findMatchingType(Integer.class, S7DataType.INT, "W");
        assertThat(result, equalTo(S7DataType.UINT));
    }

    @Test(expected = IllegalArgumentException.class)
    @Category(FastTests.class)
    public void findMatchingTypeBaseTypeAndWrongSizeCode() {
        S7DataType result = S7DataType.findMatchingType(Integer.class, S7DataType.INT, "X");
        assertThat(result, equalTo(S7DataType.UINT));
    }*/

}