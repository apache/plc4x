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

#include <plc4c/types.h>
#include "plc4c/spi/types_private.h"

char *plc4c_return_code_to_message(return_code err) {
    return "hurz";
}

char *plc4c_response_code_to_message(plc4c_response_code response_code) {
    switch (response_code) {
        case PLC4C_RESPONSE_CODE_OK: {
            return "OK";
        }
        case PLC4C_RESPONSE_CODE_NOT_FOUND: {
            return "NOT FOUND";
        }
        case PLC4C_RESPONSE_CODE_ACCESS_DENIED: {
            return "ACCESS DENIED";
        }
        case PLC4C_RESPONSE_CODE_INVALID_ADDRESS: {
            return "INVALID ADDRESS";
        }
        case PLC4C_RESPONSE_CODE_INVALID_DATATYPE: {
            return "INVALID_DATATYPE";
        }
        case PLC4C_RESPONSE_CODE_INTERNAL_ERROR: {
            return "INTERNAL ERROR";
        }
        case PLC4C_RESPONSE_CODE_RESPONSE_PENDING: {
            return "RESPONSE PENDING";
        }
        default: {
            return "UNKNOWN RESPONSE CODE";
        }
    }
}
