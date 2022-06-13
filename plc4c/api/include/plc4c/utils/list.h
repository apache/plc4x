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
#ifndef PLC4C_UTILS_LIST_H_
#define PLC4C_UTILS_LIST_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

typedef struct plc4c_list plc4c_list;
typedef struct plc4c_list_element plc4c_list_element;

typedef void (*plc4c_list_delete_element_callback)(plc4c_list_element *element);

struct plc4c_list {
  plc4c_list_element *head;
  plc4c_list_element *tail;
};

struct plc4c_list_element {
  plc4c_list_element *next;
  plc4c_list_element *previous;
  void *value;
};

void plc4c_utils_list_init(plc4c_list *list);

void plc4c_utils_list_create(plc4c_list **list);

size_t plc4c_utils_list_size(plc4c_list *list);

bool plc4c_utils_list_empty(plc4c_list *list);

bool plc4c_utils_list_contains(plc4c_list *list, plc4c_list_element *element);

void *plc4c_utils_list_get_value(plc4c_list *list, size_t element_index);

void plc4c_utils_list_insert_head(plc4c_list *list,
                                  plc4c_list_element *element);

void plc4c_utils_list_insert_head_value(plc4c_list *list, void *value);

void plc4c_utils_list_insert_tail(plc4c_list *list,
                                  plc4c_list_element *element);

void plc4c_utils_list_insert_tail_value(plc4c_list *list, void *value);

void plc4c_utils_list_remove(plc4c_list *list, plc4c_list_element *element);

plc4c_list_element *plc4c_utils_list_remove_head(plc4c_list *list);

plc4c_list_element *plc4c_utils_list_remove_tail(plc4c_list *list);

plc4c_list_element *plc4c_utils_list_head(plc4c_list *list);

plc4c_list_element *plc4c_utils_list_tail(plc4c_list *list);

void plc4c_utils_list_delete_elements(plc4c_list *list,
                                      plc4c_list_delete_element_callback);

plc4c_list_element *plc4c_utils_list_find_element_by_item(plc4c_list *list,
                                                          void *item);

uint8_t *plc4c_list_to_byte_array(plc4c_list *list);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_UTILS_LIST_H_
