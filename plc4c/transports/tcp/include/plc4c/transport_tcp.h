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
#ifndef PLC4C_TRANSPORT_TCP_H_
#define PLC4C_TRANSPORT_TCP_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <plc4c/types.h>
#include <string.h>

struct plc4c_transport_tcp_config {
  char* address;
  uint16_t port;

  int sockfd;
};
typedef struct plc4c_transport_tcp_config plc4c_transport_tcp_config;

#if defined(_WIN32) || defined(_WIN64)
/* We are on Windows */
# define strtok_r strtok_s
#endif

plc4c_transport *plc4c_transport_tcp_create();

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_TRANSPORT_TCP_H_