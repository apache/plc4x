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

import org.apache.plc4x.java.s7.netty.model.types.CpuServicesParameterSubFunctionGroup;
import org.apache.plc4x.java.s7.netty.model.types.CpuServicesParameterFunctionGroup;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;

public abstract class CpuServicesParameter implements S7Parameter {

    private CpuServicesParameterFunctionGroup functionGroup;
    private CpuServicesParameterSubFunctionGroup subFunctionGroup;
    private byte sequenceNumber;

    public CpuServicesParameter(CpuServicesParameterFunctionGroup functionGroup, CpuServicesParameterSubFunctionGroup subFunctionGroup, byte sequenceNumber) {
        this.functionGroup = functionGroup;
        this.subFunctionGroup = subFunctionGroup;
        this.sequenceNumber = sequenceNumber;
    }

    @Override
    public ParameterType getType() {
        return ParameterType.CPU_SERVICES;
    }

    public CpuServicesParameterFunctionGroup getFunctionGroup() {
        return functionGroup;
    }

    public CpuServicesParameterSubFunctionGroup getSubFunctionGroup() {
        return subFunctionGroup;
    }

    public byte getSequenceNumber() {
        return sequenceNumber;
    }

}
