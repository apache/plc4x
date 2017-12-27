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
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;
import org.apache.plc4x.java.s7.netty.model.types.SpecificationType;
import org.apache.plc4x.java.s7.netty.model.types.TransportSize;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;

import static org.junit.jupiter.api.Assertions.*;

class S7ParameterTests {

    @Test
    @Tag("fast")
    void varParameter() {
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
        assertTrue(varParameter.getType() == ParameterType.READ_VAR, "Unexpected parameter type");
        assertTrue(varParameter.getItems().containsAll(parameterItems), "Unexpected paramater items");
    }

    @Test
    @Tag("fast")
    void cpuServicesParameter() {
        CpuServicesParameter cpuParameter = new CpuServicesParameter();
        assertTrue(cpuParameter.getType() == ParameterType.CPU_SERVICES, "Unexpected parameter type");
    }
    
    @Test
    @Tag("fast")
    void setupCommunicationsParameter() {
        short maxAmqCaller = 4;
        short maxAmqCallee = 8;
        short pduLength = 512;

        SetupCommunicationParameter setupParameter = new SetupCommunicationParameter(maxAmqCaller, maxAmqCallee, pduLength);
        assertTrue(setupParameter.getType() == ParameterType.SETUP_COMMUNICATION, "Unexpected parameter type");
        assertTrue(setupParameter.getMaxAmqCallee() == maxAmqCallee, "Unexpected value for maxAmqCallee");
        assertTrue(setupParameter.getMaxAmqCaller() == maxAmqCaller, "Unexpected value for maxAmqCaller");
        assertTrue(setupParameter.getPduLength() == pduLength, "Unexpected value for pduLength");
    }



}