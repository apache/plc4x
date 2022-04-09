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

#ifndef PLC4C_MQTT_READ_WRITE_MQT_T__CONTROL_PACKET_TYPE_H_
#define PLC4C_MQTT_READ_WRITE_MQT_T__CONTROL_PACKET_TYPE_H_

#include <stdbool.h>
#include <stdint.h>
#include <plc4c/driver_mqtt_static.h>
#include <plc4c/spi/read_buffer.h>
#include <plc4c/spi/write_buffer.h>

// Code generated by code-generation. DO NOT EDIT.


#ifdef __cplusplus
extern "C" {
#endif

enum plc4c_mqtt_read_write_mqt_t__control_packet_type {
  plc4c_mqtt_read_write_mqt_t__control_packet_type_RESERVED = 0x0,
  plc4c_mqtt_read_write_mqt_t__control_packet_type_CONNECT = 0x1,
  plc4c_mqtt_read_write_mqt_t__control_packet_type_CONNACK = 0x2,
  plc4c_mqtt_read_write_mqt_t__control_packet_type_PUBLISH = 0x3,
  plc4c_mqtt_read_write_mqt_t__control_packet_type_PUBACK = 0x4,
  plc4c_mqtt_read_write_mqt_t__control_packet_type_PUBREC = 0x5,
  plc4c_mqtt_read_write_mqt_t__control_packet_type_PUBREL = 0x6,
  plc4c_mqtt_read_write_mqt_t__control_packet_type_PUBCOMP = 0x7,
  plc4c_mqtt_read_write_mqt_t__control_packet_type_SUBSCRIBE = 0x8,
  plc4c_mqtt_read_write_mqt_t__control_packet_type_SUBACK = 0x9,
  plc4c_mqtt_read_write_mqt_t__control_packet_type_UNSUBSCRIBE = 0xA,
  plc4c_mqtt_read_write_mqt_t__control_packet_type_UNSUBACK = 0xB,
  plc4c_mqtt_read_write_mqt_t__control_packet_type_PINGREQ = 0xC,
  plc4c_mqtt_read_write_mqt_t__control_packet_type_PINGRESP = 0xD,
  plc4c_mqtt_read_write_mqt_t__control_packet_type_DISCONNECT = 0xE,
  plc4c_mqtt_read_write_mqt_t__control_packet_type_AUTH = 0xF
};
typedef enum plc4c_mqtt_read_write_mqt_t__control_packet_type plc4c_mqtt_read_write_mqt_t__control_packet_type;

// Get an empty NULL-struct
plc4c_mqtt_read_write_mqt_t__control_packet_type plc4c_mqtt_read_write_mqt_t__control_packet_type_null();

plc4c_return_code plc4c_mqtt_read_write_mqt_t__control_packet_type_parse(plc4c_spi_read_buffer* readBuffer, plc4c_mqtt_read_write_mqt_t__control_packet_type** message);

plc4c_return_code plc4c_mqtt_read_write_mqt_t__control_packet_type_serialize(plc4c_spi_write_buffer* writeBuffer, plc4c_mqtt_read_write_mqt_t__control_packet_type* message);

plc4c_mqtt_read_write_mqt_t__control_packet_type plc4c_mqtt_read_write_mqt_t__control_packet_type_value_of(char* value_string);

int plc4c_mqtt_read_write_mqt_t__control_packet_type_num_values();

plc4c_mqtt_read_write_mqt_t__control_packet_type plc4c_mqtt_read_write_mqt_t__control_packet_type_value_for_index(int index);

uint16_t plc4c_mqtt_read_write_mqt_t__control_packet_type_length_in_bytes(plc4c_mqtt_read_write_mqt_t__control_packet_type* message);

uint16_t plc4c_mqtt_read_write_mqt_t__control_packet_type_length_in_bits(plc4c_mqtt_read_write_mqt_t__control_packet_type* message);

#ifdef __cplusplus
}
#endif

#endif  // PLC4C_MQTT_READ_WRITE_MQT_T__CONTROL_PACKET_TYPE_H_
