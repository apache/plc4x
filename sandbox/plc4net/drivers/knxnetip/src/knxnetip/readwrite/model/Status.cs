//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

namespace org.apache.plc4net.drivers.knxnetip.readwrite.model
{

    public enum Status : byte
    {

        NO_ERROR = 0x00,
        PROTOCOL_TYPE_NOT_SUPPORTED = 0x01,
        UNSUPPORTED_PROTOCOL_VERSION = 0x02,
        OUT_OF_ORDER_SEQUENCE_NUMBER = 0x04,
        INVALID_CONNECTION_ID = 0x21,
        CONNECTION_TYPE_NOT_SUPPORTED = 0x22,
        CONNECTION_OPTION_NOT_SUPPORTED = 0x23,
        NO_MORE_CONNECTIONS = 0x24,
        NO_MORE_UNIQUE_CONNECTIONS = 0x25,
        DATA_CONNECTION = 0x26,
        KNX_CONNECTION = 0x27,
        TUNNELLING_LAYER_NOT_SUPPORTED = 0x29,

    }

}

