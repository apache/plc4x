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

import org.apache.plc4x.java.s7.netty.model.messages.S7PushMessage;
import org.apache.plc4x.java.s7.netty.model.types.CpuCurrentModeType;
import org.apache.plc4x.java.s7.netty.model.types.CpuUserDataMethodType;
import org.apache.plc4x.java.s7.netty.model.types.CpuUserDataParameterFunctionGroupType;
import org.apache.plc4x.java.s7.netty.model.types.CpuUserDataParameterType;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;

/**
 *
 * @author cgarcia
 */
public class CpuDiagnosticPushParameter implements S7Parameter, S7PushMessage {
    
    private CpuUserDataMethodType Method;
    private CpuUserDataParameterType parameterType;
    private CpuUserDataParameterFunctionGroupType ParameterFunction;
    private CpuCurrentModeType CurrentMode;
    private byte sequenceNumber;

    public CpuDiagnosticPushParameter(
            CpuUserDataMethodType Method, 
            CpuUserDataParameterType ParameterType, 
            CpuUserDataParameterFunctionGroupType ParameterFunction, 
            CpuCurrentModeType CurrentMode, 
            byte sequenceNumber) {
        this.Method = Method;
        this.parameterType = ParameterType;
        this.ParameterFunction = ParameterFunction;
        this.CurrentMode = CurrentMode;
        this.sequenceNumber = sequenceNumber;
    }

    public CpuUserDataMethodType getMethod() {
        return Method;
    }

    public CpuUserDataParameterType getParameterType() {
        return parameterType;
    }

    public CpuUserDataParameterFunctionGroupType getParameterFunction() {
        return ParameterFunction;
    }

    public CpuCurrentModeType getCurrentMode() {
        return CurrentMode;
    }

    @Override
    public ParameterType getType() {
        return ParameterType.MODE_TRANSITION;
    }

    @Override
    public String toString() {
        return "CpuDiagnosticPushParameter{" + "Method=" + Method 
                + ", parameterType=" + parameterType 
                + ", ParameterFunction=" + ParameterFunction 
                + ", CurrentMode=" + CurrentMode 
                + ", sequenceNumber=" + sequenceNumber 
                + '}';
    }



    
}
