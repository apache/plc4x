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

#ifndef PLC4C_DRIVER_MODBUS_STATIC_H_
#define PLC4C_DRIVER_MODBUS_STATIC_H_

#include "modbus_pdu.h"

uint16_t plc4c_modbus_read_write_rtu_crc_check(uint8_t address, plc4c_modbus_read_write_modbus_pdu* pdu);

uint8_t plc4c_modbus_read_write_ascii_lrc_check(uint8_t address, plc4c_modbus_read_write_modbus_pdu* pdu);

#endif  // PLC4C_DRIVER_MODBUS_STATIC_H_
