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
#ifndef PLC4C_S7_READ_WRITE_S7_MESSAGE_USER_DATA_H_
#define PLC4C_S7_READ_WRITE_S7_MESSAGE_USER_DATA_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include <stdint.h>
#include <plc4c/utils/list.h>
#include "s7_parameter.h"
#include "s7_payload.h"

struct plc4c_s7_read_write_s7_message_user_data {
  uint16_t tpdu_reference;
  plc4c_s7_read_write_s7_parameter parameter;
  plc4c_s7_read_write_s7_payload payload;
};
typedef struct plc4c_s7_read_write_s7_message_user_data plc4c_s7_read_write_s7_message_user_data;

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_S7_READ_WRITE_S7_MESSAGE_USER_DATA_H_
