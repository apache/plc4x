/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
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

// fcn added by thomas ~~~~~~~~~
void plc4c_driver_s7_destroy_s7_write_request(
    plc4c_s7_read_write_tpkt_packet* packet);
void plc4c_driver_s7_destroy_s7_read_request(
    plc4c_s7_read_write_tpkt_packet* packet);
void plc4c_driver_s7_time_transport_size(
    plc4c_s7_read_write_transport_size *transport_size);
void plc4c_add_data_to_request(plc4c_data* parsed_value, 
    plc4c_s7_read_write_s7_var_payload_data_item* request_value);
void delete_copt_parameter_list_element(plc4c_list_element *element);
void delete_s7_parameter_list_element(plc4c_list_element *element);
void plc4c_driver_s7_destroy_receive_packet(
    plc4c_s7_read_write_tpkt_packet* packet);
void plc4c_driver_s7_destroy_cotp_connection_request(
    plc4c_s7_read_write_tpkt_packet *packet);
void plc4c_driver_s7_destroy_s7_connection_request(
    plc4c_s7_read_write_tpkt_packet* packet);
void delete_payload_user_data_item_list_element(
    plc4c_list_element *element);
void delete_parameter_user_data_item_list_element(
  plc4c_list_element *element);
void plc4c_driver_s7_destroy_s7_identify_remote_request(
    plc4c_s7_read_write_tpkt_packet* packet);
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

plc4c_return_code plc4c_driver_s7_send_packet(
    plc4c_connection* connection, plc4c_s7_read_write_tpkt_packet* packet);
plc4c_return_code plc4c_driver_s7_receive_packet(
    plc4c_connection* connection, plc4c_s7_read_write_tpkt_packet** packet);

plc4c_return_code plc4c_driver_s7_create_cotp_connection_request(
    plc4c_driver_s7_config* configuration,
    plc4c_s7_read_write_tpkt_packet** cotp_connect_request_packet);
plc4c_return_code plc4c_driver_s7_create_s7_connection_request(
    plc4c_driver_s7_config* configuration,
    plc4c_s7_read_write_tpkt_packet** s7_connect_request_packet);
plc4c_return_code plc4c_driver_s7_create_s7_identify_remote_request(
    plc4c_s7_read_write_tpkt_packet** s7_identify_remote_request_packet);
plc4c_return_code plc4c_driver_s7_create_s7_read_request(
    plc4c_read_request* read_request,
    plc4c_s7_read_write_tpkt_packet** s7_read_request_packet);
plc4c_return_code plc4c_driver_s7_create_s7_write_request(
    plc4c_write_request* write_request,
    plc4c_s7_read_write_tpkt_packet** s7_write_request_packet);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_DRIVER_S7_PACKETS_H_