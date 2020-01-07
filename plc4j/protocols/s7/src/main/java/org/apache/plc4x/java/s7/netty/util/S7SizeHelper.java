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
package org.apache.plc4x.java.s7.netty.util;

import java.util.List;
import org.apache.plc4x.java.s7.netty.model.params.CpuServicesRequestParameter;
import org.apache.plc4x.java.s7.netty.model.params.S7Parameter;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem;
import org.apache.plc4x.java.s7.netty.model.params.items.VarParameterItem;
import org.apache.plc4x.java.s7.netty.model.payloads.CpuMessageSubscriptionServicePayload;
import org.apache.plc4x.java.s7.netty.model.payloads.CpuServicesPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.S7Payload;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.payloads.ssls.SslDataRecord;
import org.apache.plc4x.java.s7.netty.model.types.VariableAddressingMode;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class S7SizeHelper {

    private static final Logger logger = LoggerFactory.getLogger(S7SizeHelper.class);

    private S7SizeHelper() {
        // We don't want this to be instantiated.
    }

    public static short getParametersLength(List<S7Parameter> parameters) {
        short l = 0;
        if (parameters != null) {
            for (S7Parameter parameter : parameters) {
                l += getParameterLength(parameter);
            }
        }
        return l;
    }

    public static short getPayloadsLength(List<S7Payload> payloads) {
        short l = 0;
        if (payloads == null) {
            return 0;
        }

        for (S7Payload payload : payloads) {
            if(payload instanceof VarPayload) {
                VarPayload varPayload = (VarPayload) payload;
                for (VarPayloadItem payloadItem : varPayload.getItems()) {
                    l += getPayloadLength(payloadItem);
                }
            } else if(payload instanceof CpuServicesPayload) {
                CpuServicesPayload cpuServicesPayload = (CpuServicesPayload) payload;
                if(cpuServicesPayload.getSslDataRecords().isEmpty()) {
                    return 8;
                } else {
                    short length = 0;
                    for (SslDataRecord sslDataRecord : cpuServicesPayload.getSslDataRecords()) {
                        length += sslDataRecord.getLengthInWords() * 2;
                    }
                    return length;
                }
            } else if(payload instanceof CpuMessageSubscriptionServicePayload) {
                CpuMessageSubscriptionServicePayload submsg = (CpuMessageSubscriptionServicePayload) payload;
                if ((submsg.getSubscribedEvents() & 0x80) == 0x00){
                    return 14;
                } else {
                    return 16;                    
                }
  
            }
        }
        return l;
    }

    public static short getParameterLength(S7Parameter parameter) {
        if (parameter == null) {
            return 0;
        }

        switch (parameter.getType()) {
            case MODE_TRANSITION:
                return 8;
            case READ_VAR:
            case WRITE_VAR:
                return getReadWriteVarParameterLength((VarParameter) parameter);
            case SETUP_COMMUNICATION:
                return 8;
            case CPU_SERVICES:
                if(parameter instanceof CpuServicesRequestParameter) {
                    return 8;
                } else {
                    return 12;
                }
            default:
                logger.error("Not implemented");
                return 0;
        }
    }

    public static short getReadWriteVarParameterLength(VarParameter parameter) {
        short length = 2;
        for (VarParameterItem varParameterItem : parameter.getItems()) {
            VariableAddressingMode addressMode = varParameterItem.getAddressingMode();

            if (addressMode == VariableAddressingMode.S7ANY) {
                length += 12;
            } else {
                logger.error("Not implemented");
            }
        }
        return length;
    }

    public static short getPayloadLength(VarParameterItem parameterItem) {
        if (parameterItem == null) {
            return 0;
        }

        if(parameterItem instanceof S7AnyVarParameterItem) {
            S7AnyVarParameterItem anyVarParameterItem = (S7AnyVarParameterItem) parameterItem;
            return (short) (4 + (
                (anyVarParameterItem.getNumElements()) * anyVarParameterItem.getDataType().getSizeInBytes()));
        } else {
            logger.error("Not implemented");
            return 0;
        }
    }

    public static short getPayloadLength(VarPayloadItem payloadItem) {
        if (payloadItem == null) {
            return 0;
        }

        return (short) (4 + payloadItem.getData().length);
    }

}
