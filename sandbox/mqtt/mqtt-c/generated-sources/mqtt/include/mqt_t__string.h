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

#ifndef PLC4C_MQTT_READ_WRITE_MQT_T__STRING_H_
#define PLC4C_MQTT_READ_WRITE_MQT_T__STRING_H_

#include <stdbool.h>
#include <stdint.h>
#include <plc4c/driver_mqtt_static.h>
#include <plc4c/spi/read_buffer.h>
#include <plc4c/spi/write_buffer.h>
#include <plc4c/utils/list.h>

// Code generated by code-generation. DO NOT EDIT.

#ifdef __cplusplus
extern "C" {
#endif


struct plc4c_mqtt_read_write_mqt_t__string {
  /* Properties */
  char* value;
};
typedef struct plc4c_mqtt_read_write_mqt_t__string plc4c_mqtt_read_write_mqt_t__string;

// Create an empty NULL-struct
plc4c_mqtt_read_write_mqt_t__string plc4c_mqtt_read_write_mqt_t__string_null();

plc4c_return_code plc4c_mqtt_read_write_mqt_t__string_parse(plc4c_spi_read_buffer* readBuffer, plc4c_mqtt_read_write_mqt_t__string** message);

plc4c_return_code plc4c_mqtt_read_write_mqt_t__string_serialize(plc4c_spi_write_buffer* writeBuffer, plc4c_mqtt_read_write_mqt_t__string* message);

uint16_t plc4c_mqtt_read_write_mqt_t__string_length_in_bytes(plc4c_mqtt_read_write_mqt_t__string* message);

uint16_t plc4c_mqtt_read_write_mqt_t__string_length_in_bits(plc4c_mqtt_read_write_mqt_t__string* message);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_MQTT_READ_WRITE_MQT_T__STRING_H_
