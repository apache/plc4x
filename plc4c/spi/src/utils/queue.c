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

#include <plc4c/utils/queue.h>

void plc4c_utils_queue_create(plc4c_queue **queue) {
  plc4c_queue *new_queue = malloc(sizeof(plc4c_queue));
  new_queue->head = NULL;
  new_queue->tail = NULL;
  *queue = new_queue;
}

size_t plc4c_utils_queue_size(plc4c_queue *queue) {
  if (queue->head == NULL) {
    return 0;
  }
  plc4c_queue_element *cur_element = queue->head;
  int count = 1;
  while (cur_element->next != NULL) {
    count++;
    cur_element = cur_element->next;
  }
  return count;
}

bool plc4c_utils_queue_empty(plc4c_queue *queue) { return queue->head == NULL; }

void plc4c_utils_queue_push(plc4c_queue *queue, plc4c_queue_element *element) {
  if (queue->tail != NULL) {
    queue->tail->next = element;
  } else {
    queue->head = element;
  }
  queue->tail = element;
}

void plc4c_utils_queue_push_value(plc4c_queue *queue, void *value) {
  plc4c_queue_element *new_element = malloc(sizeof(plc4c_queue_element));
  new_element->value = value;
  new_element->next = NULL;
  plc4c_utils_queue_push(queue, new_element);
}

plc4c_queue_element *plc4c_utils_queue_pop(plc4c_queue *queue) {
  plc4c_queue_element *head_element = queue->head;
  queue->head = queue->head->next;
  head_element->next = NULL;
  return head_element;
}

plc4c_queue_element *plc4c_utils_queue_head(plc4c_queue *queue) {
  return queue->head;
}

plc4c_queue_element *plc4c_utils_queue_tail(plc4c_queue *queue) {
  return queue->tail;
}
