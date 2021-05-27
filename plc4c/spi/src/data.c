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

#include <assert.h>
#include <plc4c/data.h>
#include <plc4c/spi/types_private.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>

plc4c_data *plc4c_data_create_bool_data(bool b) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_BOOL;
  data->size = sizeof(bool);
  data->data.boolean_value = b;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data* plc4c_data_create_bool_array(bool *b, int nItems) {
  plc4c_list list;
  plc4c_data* elem;
  plc4c_utils_list_init(&list);
  for (size_t idx = 0; idx < nItems ; idx++) {
    elem = plc4c_data_create_bool_data(*(b + idx));
    plc4c_utils_list_insert_head_value(&list, elem);
  }
  return plc4c_data_create_list_data(list);
}


plc4c_data *plc4c_data_create_int8_t_data(int8_t c) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_CHAR;
  data->size = sizeof(char);
  data->data.char_value = c;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data* plc4c_data_create_int8_t_array(int8_t *c, int nItems) {
  plc4c_list list;
  plc4c_data* elem;
  plc4c_utils_list_init(&list);
  for (size_t idx = 0; idx < nItems ; idx++) {
    elem = plc4c_data_create_int8_t_data(*(c + idx));
    plc4c_utils_list_insert_head_value(&list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_uint8_t_data(uint8_t uc) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_UCHAR;
  data->size = sizeof(unsigned char);
  data->data.uchar_value = uc;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data* plc4c_data_create_uint8_t_array(uint8_t *uc, int nItems) {
  plc4c_list list;
  plc4c_data* elem;
  plc4c_utils_list_init(&list);
  for (size_t idx = 0; idx < nItems ; idx++) {
    elem = plc4c_data_create_uint8_t_data(*(uc + idx));
    plc4c_utils_list_insert_head_value(&list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_int16_t_data(int16_t s) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_SHORT;
  data->size = sizeof(short);
  data->data.short_value = s;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data* plc4c_data_create_int16_t_array(int16_t *s, int nItems) {
  plc4c_list list;
  plc4c_data* elem;
  plc4c_utils_list_init(&list);
  for (size_t idx = 0; idx < nItems ; idx++) {
    elem = plc4c_data_create_int16_t_data(*(s + idx));
    plc4c_utils_list_insert_head_value(&list, elem);
  }
  return plc4c_data_create_list_data(list);
}


plc4c_data *plc4c_data_create_uint16_t_data(uint16_t us) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_USHORT;
  data->size = sizeof(unsigned short);
  data->data.ushort_value = us;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data* plc4c_data_create_uint16_t_array(uint16_t *us, int nItems) {
  plc4c_list list;
  plc4c_data* elem;
  plc4c_utils_list_init(&list);
  for (size_t idx = 0; idx < nItems ; idx++) {
    elem = plc4c_data_create_uint16_t_data(*(us + idx));
    plc4c_utils_list_insert_head_value(&list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_int32_t_data(int32_t i) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_INT;
  data->size = sizeof(int);
  data->data.int_value = i;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data* plc4c_data_create_int32_t_array(int32_t *i, int nItems) {
  plc4c_list list;
  plc4c_data* elem;
  plc4c_utils_list_init(&list);
  for (size_t idx = 0; idx < nItems ; idx++) {
    elem = plc4c_data_create_int32_t_data(*(i + idx));
    plc4c_utils_list_insert_head_value(&list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_uint32_t_data(uint32_t ui) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_UINT;
  data->size = sizeof(unsigned int);
  data->data.uint_value = ui;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data* plc4c_data_create_uint32_t_array(uint32_t *ui, int nItems) {
  plc4c_list list;
  plc4c_data* elem;
  plc4c_utils_list_init(&list);
  for (size_t idx = 0; idx < nItems ; idx++) {
    elem = plc4c_data_create_uint32_t_data(*(ui + idx));
    plc4c_utils_list_insert_head_value(&list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_int64_t_data(int64_t i) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_LINT;
  data->size = sizeof(int);
  data->data.lint_value = i;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data* plc4c_data_create_int64_t_array(int64_t *i, int nItems) {
  plc4c_list list;
  plc4c_data* elem;
  plc4c_utils_list_init(&list);
  for (size_t idx = 0; idx < nItems ; idx++) {
    elem = plc4c_data_create_int64_t_data(*(i + idx));
    plc4c_utils_list_insert_head_value(&list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_uint64_t_data(uint64_t ui) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_ULINT;
  data->size = sizeof(unsigned int);
  data->data.ulint_value = ui;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data* plc4c_data_create_uint64_t_array(uint64_t *ui, int nItems) {
  plc4c_list list;
  plc4c_data* elem;
  plc4c_utils_list_init(&list);
  for (size_t idx = 0; idx < nItems ; idx++) {
    elem = plc4c_data_create_uint64_t_data(*(ui + idx));
    plc4c_utils_list_insert_head_value(&list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_float_data(float f) {
  plc4c_data *data = calloc(1,sizeof(plc4c_data));
  data->data_type = PLC4C_FLOAT;
  data->size = sizeof(float);
  data->data.float_value = f;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data* plc4c_data_create_float_array(float *f, int nItems) {
  plc4c_list list;
  plc4c_data* elem;
  plc4c_utils_list_init(&list);
  for (size_t idx = 0; idx < nItems ; idx++) {
    elem = plc4c_data_create_float_data(*(f + idx));
    plc4c_utils_list_insert_head_value(&list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_double_data(double d) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_DOUBLE;
  data->size = sizeof(float);
  data->data.double_value = d;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data* plc4c_data_create_double_array(double *d, int nItem) {
  plc4c_list list;
  plc4c_data* elem;
  plc4c_utils_list_init(&list);
  for (size_t idx = 0; idx < nItem; idx++) {
    elem = plc4c_data_create_double_data(*(d + idx));
    plc4c_utils_list_insert_head_value(&list, elem);
  }
  return plc4c_data_create_list_data(list);
}


plc4c_data *plc4c_data_create_string_data(unsigned int size, char *s) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_STRING_POINTER;
  data->size = size;
  data->data.pstring_value = s;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_constant_string_data(unsigned int size, char *s) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_CONSTANT_STRING;
  data->size = size;
  data->data.const_string_value = s;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_char_data(char* s) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_CONSTANT_STRING;
  data->size = 1;
  data->data.const_string_value = s;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_list_data(plc4c_list list) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_LIST;
  // TODO: Perhaps the list size makes more sense here
  data->size = 1;
  data->data.list_value = list;
  // TODO: Add a destroy function for lists
  data->custom_destroy = NULL;
  // TODO: Add a print function for lists
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_uint8_t_bit_string_data(uint8_t uc) {
  uint8_t cur_bit = ((uint8_t) 1) << 7;
  plc4c_list list ;
  plc4c_utils_list_init(&list);
  for(int i = 0; i < 8; i++) {
    plc4c_data *item = plc4c_data_create_bool_data((uc & cur_bit) != 0);
    plc4c_utils_list_insert_head_value(&list, item);
    cur_bit = cur_bit >> 1;
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_uint16_t_bit_string_data(uint16_t us) {
  uint16_t cur_bit = ((uint16_t) 1) << 15;
  plc4c_list list ;
  plc4c_utils_list_init(&list);
  for(int i = 0; i < 16; i++) {
    plc4c_data *item = plc4c_data_create_bool_data((us & cur_bit) != 0);
    plc4c_utils_list_insert_head_value(&list, item);
    cur_bit = cur_bit >> 1;
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_uint32_t_bit_string_data(uint32_t ui) {
  uint32_t cur_bit = ((uint32_t) 1) << 31;
  plc4c_list list;
  plc4c_utils_list_init(&list);
  for(int i = 0; i < 32; i++) {
    plc4c_data *item = plc4c_data_create_bool_data((ui & cur_bit) != 0);
    plc4c_utils_list_insert_head_value(&list, item);
    cur_bit = cur_bit >> 1;
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_uint64_t_bit_string_data(uint64_t ui) {
  uint64_t cur_bit = ((uint64_t) 1) << 63;
  plc4c_list list;
  plc4c_utils_list_init(&list);
  for(int i = 0; i < 64; i++) {
    plc4c_data *item = plc4c_data_create_bool_data((ui & cur_bit) != 0);
    plc4c_utils_list_insert_head_value(&list, item);
    cur_bit = cur_bit >> 1;
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_void_pointer_data(void *v) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_VOID_POINTER;
  data->size = 0;
  data->data.pvoid_value = v;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

void plc4c_data_printf(plc4c_data *data) {
  if (data == NULL) {
    printf("NULL");
    return;
  }
  switch (data->data_type) {
    case PLC4C_BOOL:
      printf("%s", data->data.boolean_value ? "true" : "false");
      break;
    case PLC4C_CHAR:
      printf("%d", data->data.char_value);
      break;
    case PLC4C_UCHAR:
      printf("%u", data->data.uchar_value);
      break;
    case PLC4C_SHORT:
      printf("%d", data->data.short_value);
      break;
    case PLC4C_USHORT:
      printf("%u", data->data.ushort_value);
      break;
    case PLC4C_INT:
      printf("%d", data->data.int_value);
      break;
    case PLC4C_UINT:
      printf("%u", data->data.uint_value);
      break;
    case PLC4C_LINT:
      printf("%ll", data->data.lint_value);
      break;
    case PLC4C_ULINT:
      printf("%llu", data->data.ulint_value);
      break;
    case PLC4C_FLOAT:
      printf("%f", data->data.float_value);
      break;
    case PLC4C_DOUBLE:
      printf("%.20f", data->data.double_value);
      break;
    case PLC4C_STRING_POINTER:
      printf("%s", data->data.pstring_value);
      break;
    case PLC4C_CONSTANT_STRING:
      printf("%s", data->data.const_string_value);
      break;
    case PLC4C_LIST:
      printf("[");
      plc4c_list_element *cur_element =
          plc4c_utils_list_tail(&data->data.list_value);
      while (cur_element != NULL) {
        plc4c_data *data_item = cur_element->value;
        plc4c_data_printf(data_item);
        cur_element = cur_element->next;
        if (cur_element != NULL) {
          printf(", ");
        }
      }
      printf("]");
      break;
    case PLC4C_VOID_POINTER:
      if (data->custom_printf != NULL) {
        data->custom_printf(data);
      } else {
        printf("%p", data->data.pvoid_value);
      }
      break;
    default:
      printf("unknown");

      break;
  }
}

void plc4c_data_set_custom_destroy(
    plc4c_data *data, plc4c_data_custom_destroy data_custom_destroy) {
  data->custom_destroy = data_custom_destroy;
}

void plc4c_data_set_custom_printf(plc4c_data *data,
                                  plc4c_data_custom_printf data_custom_printf) {
  data->custom_printf = data_custom_printf;
}

void list_data_rm(plc4c_list_element *elem) {
  plc4c_data* value;
  value = elem->value;
  plc4c_data_destroy(value);
}

void plc4c_data_destroy(plc4c_data *data) {
  assert(data != NULL);
  if (data->custom_destroy != NULL) {
    data->custom_destroy(data);
  } else {
    switch (data->data_type) {
      case PLC4C_VOID_POINTER:
        free(data->data.pvoid_value);
        break;
      case PLC4C_STRING_POINTER:
        free(data->data.pstring_value);
        break;
      case PLC4C_LIST:
        plc4c_utils_list_delete_elements(&data->data.list_value,
            list_data_rm);
        //free(&data->data.list_value);
      default:
        break;
    }
  }
  free(data);
}


#include <string.h>
void* plc4c_data_update_values(plc4c_data *data, void *newData) {

  switch (data->data_type) {
      case PLC4C_BOOL:
        data->data.boolean_value = *(bool*)newData;
        newData += sizeof(bool);
        break;
      case PLC4C_CHAR:
        data->data.char_value = *(int8_t*)newData;
        newData += sizeof(int8_t);
        break;
      case PLC4C_UCHAR:
        data->data.char_value = *(int8_t*)newData;
        newData += sizeof(int8_t);
        break;
      case PLC4C_SHORT:
        data->data.short_value = *(int16_t*)newData;
        newData += sizeof(int16_t);
        break;
      case PLC4C_USHORT:
        data->data.ushort_value = *(uint16_t*)newData;
        newData += sizeof(uint16_t);
        break;
      case PLC4C_INT:
        data->data.int_value = *(int32_t*)newData;
        newData += sizeof(int32_t);
        break;
      case PLC4C_UINT:
        data->data.uint_value = *(uint32_t*)newData;
        newData += sizeof(uint32_t);
        break;
      case PLC4C_LINT:
        data->data.lint_value = *(int64_t*)newData;
        newData += sizeof(int64_t);
        break;
      case PLC4C_ULINT:
        data->data.ulint_value = *(uint64_t*)newData;
        newData += sizeof(uint64_t);
        break;
      case PLC4C_FLOAT:
        data->data.float_value = *(float*)newData;
        newData += sizeof(float);
        break;
      case PLC4C_DOUBLE:
        data->data.double_value = *(double*)newData;
        newData += sizeof(double);
        break;
      case PLC4C_STRING_POINTER:
        data->data.pstring_value = (char*) newData;
        newData += strlen(newData) + 1;
        break;
      case PLC4C_CONSTANT_STRING:
        data->data.const_string_value = (char *)newData;
        newData += strlen(newData) + 1;
        break;
      case PLC4C_VOID_POINTER:
        data->data.pvoid_value = newData;
        newData += sizeof(void*);
        break;
      case PLC4C_LIST: {
        plc4c_list_element *list_element;
        plc4c_data *list_data;
        list_element = plc4c_utils_list_tail(&data->data.list_value);
        do {
          plc4c_data *list_data = list_element->value;
          newData = plc4c_data_update_values(list_data, newData);
          list_element = list_element->next;
        } while (list_element != NULL);
        break;
      }
      case PLC4C_STRUCT:
        // TODO
        break;
  }
  return newData;
}

