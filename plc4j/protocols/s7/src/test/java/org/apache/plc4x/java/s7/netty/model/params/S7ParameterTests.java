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

package org.apache.plc4x.java.s7.netty.model.params;

import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem;
import org.apache.plc4x.java.s7.netty.model.params.items.VarParameterItem;
import org.apache.plc4x.java.s7.netty.model.types.*;
import org.apache.plc4x.test.FastTests;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import java.util.ArrayList;

import static org.hamcrest.collection.IsIterableContainingInOrder.contains;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.junit.Assert.assertThat;

public class S7ParameterTests {

    @Test
    @Category(FastTests.class)
    public void varParameter() {
        ParameterType parameterType = ParameterType.READ_VAR;
        ArrayList<VarParameterItem> parameterItems = new ArrayList<>();
        SpecificationType specificationType = SpecificationType.VARIABLE_SPECIFICATION;
        MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
        TransportSize transportSize = TransportSize.INT;
        short numElements = 1;
        byte dataBlock = (byte) 0x1;
        byte byteOffset = (byte) 0x10;
        byte bitOffset = (byte) 0x0;

        parameterItems.add(new S7AnyVarParameterItem(specificationType, memoryArea, transportSize, numElements, dataBlock, byteOffset, bitOffset));

        VarParameter varParameter = new VarParameter(parameterType, parameterItems);
        assertThat("Unexpected parameter type", varParameter.getType(), equalTo(ParameterType.READ_VAR));
        assertThat("Unexpected parameter items", varParameter.getItems(), contains(parameterItems.toArray()));
    }

    @Test
    @Category(FastTests.class)
    public void cpuServicesParameter() {
        CpuServicesParameter cpuParameter = new CpuServicesRequestParameter(
            CpuServicesParameterFunctionGroup.CPU_FUNCTIONS, CpuServicesParameterSubFunctionGroup.READ_SSL, (byte) 0);
        assertThat("Unexpected parameter type", cpuParameter.getType(), equalTo(ParameterType.CPU_SERVICES));
    }
    
    @Test
    @Category(FastTests.class)
    public void setupCommunicationsParameter() {
        short maxAmqCaller = 4;
        short maxAmqCallee = 8;
        short pduLength = 512;

        SetupCommunicationParameter setupParameter = new SetupCommunicationParameter(maxAmqCaller, maxAmqCallee, pduLength);
        assertThat("Unexpected parameter type", setupParameter.getType(), equalTo(ParameterType.SETUP_COMMUNICATION));
        assertThat("Unexpected value for maxAmqCallee", setupParameter.getMaxAmqCallee(), equalTo(maxAmqCallee));
        assertThat("Unexpected value for maxAmqCaller", setupParameter.getMaxAmqCaller(), equalTo(maxAmqCaller));
        assertThat("Unexpected value for pduLength", setupParameter.getPduLength(), equalTo(pduLength));
    }

}