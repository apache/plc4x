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
#ifndef PLC4C_DRIVER_S7_PACKETS_H_
#define PLC4C_DRIVER_S7_PACKETS_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <plc4c/types.h>

#include "tpkt_packet.h"
#include "plc4c/driver_s7.h"

plc4c_return_code send_packet(plc4c_connection* connection,
                              plc4c_s7_read_write_tpkt_packet* packet);
plc4c_return_code receive_packet(plc4c_connection* connection,
                                 plc4c_s7_read_write_tpkt_packet** packet);

plc4c_return_code createCOTPConnectionRequest(
    plc4c_driver_s7_config* configuration,
    plc4c_s7_read_write_tpkt_packet** cotp_connect_request_packet);
plc4c_return_code createS7ConnectionRequest(
    plc4c_driver_s7_config* configuration,
    plc4c_s7_read_write_tpkt_packet** s7_connect_request_packet);
plc4c_return_code createS7IdentifyRemoteRequest(
    plc4c_s7_read_write_tpkt_packet** s7_identify_remote_request_packet);
plc4c_return_code createS7ReadRequest(
    plc4c_read_request* read_request,
    plc4c_s7_read_write_tpkt_packet** s7_read_request_packet);
plc4c_return_code createS7WriteRequest(
    plc4c_write_request* write_request,
    plc4c_s7_read_write_tpkt_packet** s7_write_request_packet);
plc4c_return_code parseAddress(
    char* address, plc4c_s7_read_write_s7_var_request_parameter_item** item);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_DRIVER_S7_PACKETS_H_