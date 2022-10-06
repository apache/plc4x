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
#ifndef PLC4C_UTILS_QUEUE_H_
#define PLC4C_UTILS_QUEUE_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include <stdlib.h>

typedef struct plc4c_queue plc4c_queue;
typedef struct plc4c_queue_element plc4c_queue_element;

struct plc4c_queue {
  plc4c_queue_element *head;
  plc4c_queue_element *tail;
};

struct plc4c_queue_element {
  plc4c_queue_element *next;
  void *value;
};

void plc4c_utils_queue_create(plc4c_queue **queue);

size_t plc4c_utils_queue_size(plc4c_queue *queue);

bool plc4c_utils_queue_empty(plc4c_queue *queue);

void plc4c_utils_queue_push(plc4c_queue *queue, plc4c_queue_element *element);

void plc4c_utils_queue_push_value(plc4c_queue *queue, void *value);

plc4c_queue_element *plc4c_utils_queue_pop(plc4c_queue *queue);

plc4c_queue_element *plc4c_utils_queue_head(plc4c_queue *queue);

plc4c_queue_element *plc4c_utils_queue_tail(plc4c_queue *queue);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_UTILS_QUEUE_H_
