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
package org.apache.plc4x.java.s7.netty.model.types;

public enum ParameterError {
    NO_ERROR((short) 0x0000),
    INVALID_BLOCK_TYPE_NUMBER((short) 0x0110),
    INVALID_PARAMETER((short) 0x0112),
    PG_RESSOURCE_ERROR((short) 0x011A),
    PLC_RESSOURCE_ERROR((short) 0x011B),
    PROTOCOL_ERROR((short) 0x011C),
    USER_BUFFER_TOO_SHORT((short) 0x011F),
    REQUEST_ERROR((short) 0x0141),
    VERSION_MISMATCH((short) 0x01C0),
    NOT_IMPLEMENTED((short) 0x01F0),
    L7_INVALID_CPU_STATE((short) 0x8001),
    L7_PDU_SIZE_ERROR((short) 0x8500),
    L7_INVALID_SZL_ID((short) 0xD401),
    L7_INVALID_INDEX((short) 0xD402),
    L7_DGS_CONNECTION_ALREADY_ANNOUNCED((short) 0xD403),
    L7_MAX_USER_NB((short) 0xD404),
    L7_DGS_FUNCTION_PARAMETER_SYNTAX_ERROR((short) 0xD405),
    L7_NO_INFO((short) 0xD406),
    L7_PRT_FUNCTION_PARAMETER_SYNTAX_ERROR((short) 0xD601),
    L7_INVALID_VARIABLE_ADDRESS((short) 0xD801),
    L7_UNKNOWN_REQUEST((short) 0xD802),
    L7_INVALID_REQUEST_STATUS((short) 0xD803);

    private final short code;

    ParameterError(short code) {
        this.code = code;
    }

    public short getCode() {
        return code;
    }

}
