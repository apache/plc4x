/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.profinet.utils;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.profinet.gsdml.ProfinetDataItem;

public class ProfinetDataTypeMapper {

    public static class DataTypeInformation {
        private final PlcValueType plcValueType;
        private final int numElements;

        public DataTypeInformation(PlcValueType plcValueType, int numElements) {
            this.plcValueType = plcValueType;
            this.numElements = numElements;
        }

        public PlcValueType getPlcValueType() {
            return plcValueType;
        }

        public int getNumElements() {
            return numElements;
        }
    }

    public static DataTypeInformation getPlcValueType(ProfinetDataItem dataItem) {
        if(dataItem.isUseAsBits()) {
            switch (dataItem.getDataType()) {
                case "Integer8":
                case "Unsigned8":
                    return new DataTypeInformation(PlcValueType.BOOL, 8);
                case "Integer16":
                case "Unsigned16":
                    return new DataTypeInformation(PlcValueType.BOOL, 16);
                case "Integer32":
                case "Unsigned32":
                    return new DataTypeInformation(PlcValueType.BOOL, 32);
                case "Integer64":
                case "Unsigned64":
                    return new DataTypeInformation(PlcValueType.BOOL, 64);
                case "Float32":
                case "Float64":
                case "Date":
                case "TimeOfDay with date indication":
                case "TimeOfDay without date indication":
                case "TimeDifference with date indication":
                case "TimeDifference without date indication":
                case "NetworkTime":
                case "NetworkTimeDifference":
                case "VisibleString":
                    break;
                case "OctetString":
                    return new DataTypeInformation(PlcValueType.BOOL, dataItem.getLength());
                case "Unsigned8+Unsigned8":
                case "Float32+Unsigned8":
                case "Float32+Status8":
                case "F_MessageTrailer4Byte":
                case "F_MessageTrailer5Byte":
            }
        }
        switch (dataItem.getDataType()) {
            case "Integer8":
                return new DataTypeInformation(PlcValueType.SINT, 1);
            case "Integer16":
                return new DataTypeInformation(PlcValueType.INT, 1);
            case "Integer32":
                return new DataTypeInformation(PlcValueType.DINT, 1);
            case "Integer64":
                return new DataTypeInformation(PlcValueType.LINT, 1);
            case "Unsigned8":
                return new DataTypeInformation(PlcValueType.USINT, 1);
            case "Unsigned16":
                return new DataTypeInformation(PlcValueType.UINT, 1);
            case "Unsigned32":
                return new DataTypeInformation(PlcValueType.UDINT, 1);
            case "Unsigned64":
                return new DataTypeInformation(PlcValueType.ULINT, 1);
            case "Float32":
                return new DataTypeInformation(PlcValueType.REAL, 1);
            case "Float64":
                return new DataTypeInformation(PlcValueType.LREAL, 1);
            case "Date":
            case "TimeOfDay with date indication":
            case "TimeOfDay without date indication":
            case "TimeDifference with date indication":
            case "TimeDifference without date indication":
            case "NetworkTime":
            case "NetworkTimeDifference":
                break;
            case "VisibleString":
                return new DataTypeInformation(PlcValueType.STRING, 1);
            case "OctetString":
                // This should actually never happen, as the OctetString should always be used with "useAsBits".
                return new DataTypeInformation(PlcValueType.BOOL, dataItem.getLength());
            case "Unsigned8+Unsigned8":
            case "Float32+Unsigned8":
            case "Float32+Status8":
            case "F_MessageTrailer4Byte":
            case "F_MessageTrailer5Byte":
        }

        throw new PlcRuntimeException("Unable to find PlcValueType for dataType " + dataItem.getDataType());
    }

}
