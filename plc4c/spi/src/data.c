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
  data->data.bool_value = b;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_bool_array(bool *b, int nItems) {
  plc4c_list *list = malloc(sizeof(plc4c_list));
  plc4c_data *elem;

  plc4c_utils_list_init(list);
  for (size_t idx = 0; idx < nItems; idx++) {
    elem = plc4c_data_create_bool_data(*(b + idx));
    plc4c_utils_list_insert_head_value(list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_usint_data(uint8_t uc) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_USINT;
  data->size = sizeof(uint8_t);
  data->data.usint_value = uc;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_usint_array(uint8_t *uc, int nItems) {
  plc4c_list *list = malloc(sizeof(plc4c_list));
  plc4c_data *elem;
  plc4c_utils_list_init(list);
  for (size_t idx = 0; idx < nItems; idx++) {
    elem = plc4c_data_create_usint_data(*(uc + idx));
    plc4c_utils_list_insert_head_value(list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_sint_data(int8_t uc) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_SINT;
  data->size = sizeof(int8_t);
  data->data.usint_value = uc;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_sint_array(int8_t *uc, int nItems) {
  plc4c_list *list = malloc(sizeof(plc4c_list));
  plc4c_data *elem;
  plc4c_utils_list_init(list);
  for (size_t idx = 0; idx < nItems; idx++) {
    elem = plc4c_data_create_sint_data(*(uc + idx));
    plc4c_utils_list_insert_head_value(list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_uint_data(uint16_t us) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_UINT;
  data->size = sizeof(uint16_t);
  data->data.uint_value = us;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_uint_array(uint16_t *us, int nItems) {
  plc4c_list *list = malloc(sizeof(plc4c_list));
  plc4c_data *elem;
  plc4c_utils_list_init(list);
  for (size_t idx = 0; idx < nItems; idx++) {
    elem = plc4c_data_create_uint_data(*(us + idx));
    plc4c_utils_list_insert_head_value(list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_int_data(int16_t s) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_INT;
  data->size = sizeof(int16_t);
  data->data.int_value = s;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_int_array(int16_t *s, int nItems) {
  plc4c_list *list = malloc(sizeof(plc4c_list));
  plc4c_data *elem;
  plc4c_utils_list_init(list);
  for (size_t idx = 0; idx < nItems; idx++) {
    elem = plc4c_data_create_int_data(*(s + idx));
    plc4c_utils_list_insert_head_value(list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_udint_data(uint32_t ui) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_UDINT;
  data->size = sizeof(uint32_t);
  data->data.uint_value = ui;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_udint_array(uint32_t *ui, int nItems) {
  plc4c_list *list = malloc(sizeof(plc4c_list));
  plc4c_data *elem;
  plc4c_utils_list_init(list);
  for (size_t idx = 0; idx < nItems; idx++) {
    elem = plc4c_data_create_udint_data(*(ui + idx));
    plc4c_utils_list_insert_head_value(list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_dint_data(int32_t i) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_DINT;
  data->size = sizeof(int32_t);
  data->data.dint_value = i;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_dint_array(int32_t *i, int nItems) {
  plc4c_list *list = malloc(sizeof(plc4c_list));
  plc4c_data *elem;
  plc4c_utils_list_init(list);
  for (size_t idx = 0; idx < nItems; idx++) {
    elem = plc4c_data_create_dint_data(*(i + idx));
    plc4c_utils_list_insert_head_value(list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_ulint_data(uint64_t ui) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_ULINT;
  data->size = sizeof(uint64_t);
  data->data.ulint_value = ui;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_ulint_array(uint64_t *ui, int nItems) {
  plc4c_list *list = malloc(sizeof(plc4c_list));
  plc4c_data *elem;
  plc4c_utils_list_init(list);
  for (size_t idx = 0; idx < nItems; idx++) {
    elem = plc4c_data_create_ulint_data(*(ui + idx));
    plc4c_utils_list_insert_head_value(list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_lint_data(int64_t i) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_LINT;
  data->size = sizeof(int64_t);
  data->data.lint_value = i;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_lint_array(int64_t *i, int nItems) {
  plc4c_list *list = malloc(sizeof(plc4c_list));
  plc4c_data *elem;
  plc4c_utils_list_init(list);
  for (size_t idx = 0; idx < nItems; idx++) {
    elem = plc4c_data_create_lint_data(*(i + idx));
    plc4c_utils_list_insert_head_value(list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_real_data(float f) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_REAL;
  data->size = sizeof(float);
  data->data.real_value = f;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_real_array(float *f, int nItems) {
  plc4c_list *list = malloc(sizeof(plc4c_list));
  plc4c_data *elem;
  plc4c_utils_list_init(list);
  for (size_t idx = 0; idx < nItems; idx++) {
    elem = plc4c_data_create_real_data(*(f + idx));
    plc4c_utils_list_insert_head_value(list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_lreal_data(double d) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_LREAL;
  data->size = sizeof(double);
  data->data.lreal_value = d;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_lreal_array(double *d, int nItem) {
  plc4c_list *list = malloc(sizeof(plc4c_list));
  plc4c_data *elem;
  plc4c_utils_list_init(list);
  for (size_t idx = 0; idx < nItem; idx++) {
    elem = plc4c_data_create_lreal_data(*(d + idx));
    plc4c_utils_list_insert_head_value(list, elem);
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_date_data(uint16_t d) {
  // TODO: Implement
  return NULL;
}

plc4c_data *plc4c_data_create_date_array(uint16_t *d, int nItems) {
  // TODO: Implement
  return NULL;
}

plc4c_data *plc4c_data_create_time_data(uint32_t t) {
  // TODO: Implement
  return NULL;
}

plc4c_data *plc4c_data_create_time_array(uint32_t *t, int nItems) {
  // TODO: Implement
  return NULL;
}

plc4c_data *plc4c_data_create_ltime_data(uint64_t lt) {
  // TODO: Implement
  return NULL;
}

plc4c_data *plc4c_data_create_ltime_array(uint64_t *lt, int nItems) {
  // TODO: Implement
  return NULL;
}

plc4c_data *plc4c_data_create_time_of_day_data(uint32_t tod) {
  // TODO: Implement
  return NULL;
}

plc4c_data *plc4c_data_create_time_of_day_array(uint32_t *tod, int nItems) {
  // TODO: Implement
  return NULL;
}

plc4c_data *plc4c_data_create_date_and_time_data(uint32_t tad) {
  // TODO: Implement
  return NULL;
}

plc4c_data *plc4c_data_create_date_and_time_array(uint32_t *tad, int nItems) {
  // TODO: Implement
  return NULL;
}

plc4c_data *plc4c_data_create_string_data(unsigned int size, char *s) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_STRING;
  data->size = size;
  data->data.string_value = s;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_char_data(char s) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_CHAR;
  data->size = 1;
  data->data.char_value = s;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_wchar_data(wchar_t s) {
  plc4c_data *data = malloc(sizeof(plc4c_data));
  data->data_type = PLC4C_CHAR;
  data->size = 2;
  data->data.wchar_value = s;
  data->custom_destroy = NULL;
  data->custom_printf = NULL;
  return data;
}

plc4c_data *plc4c_data_create_list_data(plc4c_list *list) {
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

plc4c_data *plc4c_data_create_byte_data(uint8_t uc) {
  uint8_t cur_bit = ((uint8_t)1) << 7;
  plc4c_list *list = malloc(sizeof(plc4c_list));
  plc4c_utils_list_init(list);
  for (int i = 0; i < 8; i++) {
    plc4c_data *item = plc4c_data_create_bool_data((uc & cur_bit) != 0);
    plc4c_utils_list_insert_head_value(list, item);
    cur_bit = cur_bit >> 1;
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_word_data(uint16_t us) {
  uint16_t cur_bit = ((uint16_t)1) << 15;
  plc4c_list *list = malloc(sizeof(plc4c_list));
  plc4c_utils_list_init(list);
  for (int i = 0; i < 16; i++) {
    plc4c_data *item = plc4c_data_create_bool_data((us & cur_bit) != 0);
    plc4c_utils_list_insert_head_value(list, item);
    cur_bit = cur_bit >> 1;
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_dword_data(uint32_t ui) {
  uint32_t cur_bit = ((uint32_t)1) << 31;
  plc4c_list *list = malloc(sizeof(plc4c_list));
  plc4c_utils_list_init(list);
  for (int i = 0; i < 32; i++) {
    plc4c_data *item = plc4c_data_create_bool_data((ui & cur_bit) != 0);
    plc4c_utils_list_insert_head_value(list, item);
    cur_bit = cur_bit >> 1;
  }
  return plc4c_data_create_list_data(list);
}

plc4c_data *plc4c_data_create_lword_data(uint64_t ui) {
  uint64_t cur_bit = ((uint64_t)1) << 63;
  plc4c_list *list = malloc(sizeof(plc4c_list));
  plc4c_utils_list_init(list);
  for (int i = 0; i < 64; i++) {
    plc4c_data *item = plc4c_data_create_bool_data((ui & cur_bit) != 0);
    plc4c_utils_list_insert_head_value(list, item);
    cur_bit = cur_bit >> 1;
  }
  return plc4c_data_create_list_data(list);
}

void plc4c_data_printf(plc4c_data *data) {
  if (data == NULL) {
    printf("NULL");
    return;
  }
  switch (data->data_type) {
    case PLC4C_BOOL:
      printf("%s", data->data.bool_value ? "true" : "false");
      break;
    case PLC4C_BYTE:
      printf("%hhu", data->data.byte_value);
      break;
    case PLC4C_WORD:
      printf("%d", data->data.word_value);
      break;
    case PLC4C_DWORD:
      printf("%d", data->data.dword_value);
      break;
    case PLC4C_LWORD:
      printf("%llu", data->data.lword_value);
      break;

    case PLC4C_USINT:
      printf("%u", data->data.usint_value);
      break;
    case PLC4C_SINT:
      printf("%d", data->data.sint_value);
      break;
    case PLC4C_UINT:
      printf("%u", data->data.uint_value);
      break;
    case PLC4C_INT:
      printf("%d", data->data.int_value);
      break;
    case PLC4C_UDINT:
      printf("%u", data->data.udint_value);
      break;
    case PLC4C_DINT:
      printf("%d", data->data.dint_value);
      break;
    case PLC4C_ULINT:
      printf("%llu", data->data.ulint_value);
      break;
    case PLC4C_LINT:
      printf("%lld", data->data.lint_value);
      break;

    case PLC4C_REAL:
      printf("%f", data->data.real_value);
      break;
    case PLC4C_LREAL:
      printf("%.20f", data->data.lreal_value);
      break;

      /*    case PLC4C_TIME:
            printf("%d", data->data.time_value);
            break;
            case PLC4C_LTIME:
              printf("%d", data->data.ltime_value);
              break;
          case PLC4C_DATE:
            printf("%d", data->data._value);
            break;
          case PLC4C_LDATE:
            printf("%d", data->data._value);
            break;
          case PLC4C_TIME_OF_DAY:
            printf("%d", data->data._value);
            break;
          case PLC4C_LTIME_OF_DAY:
            printf("%d", data->data._value);
            break;
          case PLC4C_DATE_AND_TIME:
            printf("%d", data->data._value);
            break;
          case PLC4C_LDATE_AND_TIME:
            printf("%d", data->data._value);
            break;*/

    case PLC4C_CHAR:
      printf("%d", data->data.char_value);
      break;
    case PLC4C_WCHAR:
      printf("%d", data->data.wchar_value);
      break;
    case PLC4C_STRING:
      printf("%s", data->data.string_value);
      break;
    case PLC4C_WSTRING:
      printf("%ws", data->data.wstring_value);
      break;

    case PLC4C_LIST:
      printf("[");
      plc4c_list_element *cur_element =
          plc4c_utils_list_tail(data->data.list_value);
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
  plc4c_data *value;
  value = elem->value;
  plc4c_data_destroy(value);
}

void plc4c_data_destroy(plc4c_data *data) {
  assert(data != NULL);
  if (data->custom_destroy != NULL) {
    data->custom_destroy(data);
  } else {
    switch (data->data_type) {
      case PLC4C_STRING:
        free(data->data.string_value);
        break;
      case PLC4C_LIST:
        plc4c_utils_list_delete_elements(data->data.list_value, list_data_rm);
        // free(&data->data.list_value);
      default:
        break;
    }
  }
  free(data);
}

void *plc4c_data_update_values(plc4c_data *data, void *newData) {
  switch (data->data_type) {
      // Boolean/Bit-String Types
    case PLC4C_BOOL:
      data->data.bool_value = *(bool *)newData;
      newData += sizeof(bool);
      break;
    case PLC4C_BYTE:
      data->data.byte_value = *(uint8_t *)newData;
      newData += sizeof(uint8_t);
      break;
    case PLC4C_WORD:
      data->data.word_value = *(uint16_t *)newData;
      newData += sizeof(uint16_t);
      break;
    case PLC4C_DWORD:
      data->data.dword_value = *(uint32_t *)newData;
      newData += sizeof(uint32_t);
      break;
    case PLC4C_LWORD:
      data->data.lword_value = *(uint64_t *)newData;
      newData += sizeof(uint64_t);
      break;
    // Integer Types (Unsigned/Signed)
    // 1-byte
    case PLC4C_USINT:
      data->data.usint_value = *(uint8_t *)newData;
      newData += sizeof(uint8_t);
      break;
    case PLC4C_SINT:
      data->data.sint_value = *(int8_t *)newData;
      newData += sizeof(int8_t);
      break;
    // 2-byte
    case PLC4C_UINT:
      data->data.uint_value = *(uint16_t *)newData;
      newData += sizeof(uint16_t);
      break;
    case PLC4C_INT:
      data->data.int_value = *(int16_t *)newData;
      newData += sizeof(int16_t);
      break;
    // 4-byte
    case PLC4C_UDINT:
      data->data.udint_value = *(uint32_t *)newData;
      newData += sizeof(uint32_t);
      break;
    case PLC4C_DINT:
      data->data.dint_value = *(int32_t *)newData;
      newData += sizeof(int32_t);
      break;
      // 8-byte
    case PLC4C_ULINT:
      data->data.ulint_value = *(uint64_t *)newData;
      newData += sizeof(uint64_t);
      break;
    case PLC4C_LINT:
      data->data.lint_value = *(int64_t *)newData;
      newData += sizeof(int64_t);
      break;
    // Floating-Point Types
    case PLC4C_REAL:
      data->data.real_value = *(float *)newData;
      newData += sizeof(float);
      break;
    case PLC4C_LREAL:
      data->data.lreal_value = *(double *)newData;
      newData += sizeof(double);
      break;
      // Time Types
      /*    case PLC4C_TIME:
            data->data.time_value = *(*)newData;
            newData += sizeof();
            break;
          case PLC4C_LTIME:
            data->data.ltime_value = *(*)newData;
            newData += sizeof();
              break;
          // Date Types
          case PLC4C_DATE:
            data->data. = *(*)newData;
            newData += sizeof();
            break;
          case PLC4C_LDATE:
            data->data. = *(*)newData;
            newData += sizeof();
            break;
          // Time of day Types
          case PLC4C_TIME_OF_DAY:
            data->data. = *(*)newData;
            newData += sizeof();
            break;
          case PLC4C_LTIME_OF_DAY:
            data->data. = *(*)newData;
            newData += sizeof();
            break;
          // Date and Time Types
          case PLC4C_DATE_AND_TIME:
            data->data. = *(*)newData;
            newData += sizeof();
            break;
          case PLC4C_LDATE_AND_TIME:
            data->data. = *(*)newData;
            newData += sizeof();
            break;*/
    // Char and String Types
    case PLC4C_CHAR:
      data->data.char_value = *(char *)newData;
      newData += sizeof(char);
      break;
    case PLC4C_WCHAR:
      data->data.wchar_value = *(wchar_t *)newData;
      newData += sizeof(char *);
      break;
    case PLC4C_STRING:
      data->data.string_value = *(char **)newData;
      newData += sizeof(char *);
      break;
    case PLC4C_WSTRING:
      data->data.wstring_value = *(wchar_t **)newData;
      newData += sizeof(wchar_t *);
      break;

    case PLC4C_LIST: {
      plc4c_list_element *list_element;
      list_element = plc4c_utils_list_tail(data->data.list_value);
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

    default:
      printf("unknown");
      break;
  }
  return newData;
}
