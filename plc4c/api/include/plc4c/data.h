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
#ifndef PLC4C_API_INCLUDE_PLC4C_DATA_H_
#define PLC4C_API_INCLUDE_PLC4C_DATA_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include <stdint.h>

#include "plc4c/utils/list.h"
#include "plc4c/types.h"

typedef void (*plc4c_data_custom_destroy)(plc4c_data *data);

typedef void (*plc4c_data_custom_printf)(plc4c_data *data);

/*
 * Functions for working with plc4c_data
 */

/**
 * Creates a plc4c_data with bool
 * @param b the bool value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_bool_data(bool b);
plc4c_data *plc4c_data_create_bool_array(bool *b, int nItem);
/**
 * Creates a plc4c_data with unsigned 8 bit integer
 * @param ui the unsigned 8 bit integer value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_usint_data(uint8_t ui);
plc4c_data *plc4c_data_create_usint_array(uint8_t *ui, int nItem);
/**
 * Creates a plc4c_data with a signed 8 bit integer
 * @param i signed 8 bit integer
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_sint_data(int8_t i);
plc4c_data *plc4c_data_create_sint_array(int8_t *i, int nItem);
/**
 * Creates a plc4c_data with unsigned 16 bit integer
 * @param ui the unsigned short value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_uint_data(uint16_t ui);
plc4c_data *plc4c_data_create_uint_array(uint16_t *ui, int nItem);
/**
 * Creates a plc4c_data with signed 16 bit integer
 * @param i the short value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_int_data(int16_t i);
plc4c_data *plc4c_data_create_int_array(int16_t *i, int nItem);
/**
 * Creates a plc4c_data with unsigned 32 bit integer
 * @param ui the unsigned 32 bit integer value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_udint_data(uint32_t ui);
plc4c_data *plc4c_data_create_udint_array(uint32_t *ui, int nItem);
/**
 * Creates a plc4c_data with signed 32 bit integer
 * @param i the signed 32 bit integer value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_dint_data(int32_t i);
plc4c_data *plc4c_data_create_dint_array(int32_t *i, int nItem);
/**
 * Creates a plc4c_data with unsigned 64 bit integer
 * @param ui the 64 bit integer value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_ulint_data(uint64_t ui);
plc4c_data *plc4c_data_create_ulint_array(uint64_t *ui, int nItem);
/**
 * Creates a plc4c_data with signed 64 bit integer
 * @param i the signed 64 bit integer value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_lint_data(int64_t i);
plc4c_data *plc4c_data_create_lint_array(int64_t *i, int nItem);

/**
 * Creates a plc4c_data with float
 * @param f the float value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_real_data(float f);
plc4c_data *plc4c_data_create_real_array(float *f, int nItem);
/**
 * Creates a plc4c_data with float
 * @param f the float value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_lreal_data(double d);
plc4c_data *plc4c_data_create_lreal_array(double *d, int nItems);

plc4c_data *plc4c_data_create_date_data(uint16_t d);
plc4c_data *plc4c_data_create_date_array(uint16_t *d, int nItems);

plc4c_data *plc4c_data_create_ldate_data(uint32_t nanosecondsSinceEpoch);
plc4c_data *plc4c_data_create_ldate_array(uint32_t *nanosecondsSinceEpoch, int nItems);

plc4c_data *plc4c_data_create_time_data(uint32_t t);
plc4c_data *plc4c_data_create_time_array(uint32_t *t, int nItems);

plc4c_data *plc4c_data_create_ltime_data(uint64_t lt);
plc4c_data *plc4c_data_create_ltime_array(uint64_t *lt, int nItems);

plc4c_data *plc4c_data_create_time_of_day_data(uint32_t tod);
plc4c_data *plc4c_data_create_time_of_day_array(uint32_t *tod, int nItems);

plc4c_data *plc4c_data_create_ltime_of_day_data(uint64_t nanosecondsSinceMidnight);
plc4c_data *plc4c_data_create_ltime_of_day_array(uint64_t *nanosecondsSinceMidnight, int nItems);

plc4c_data *plc4c_data_create_date_and_time_data(uint32_t tad);
plc4c_data *plc4c_data_create_date_and_time_array(uint32_t *tad, int nItems);

plc4c_data *plc4c_data_create_ldate_and_time_data(uint64_t nanosecondsSinceEpoch);
plc4c_data *plc4c_data_create_ldate_and_time_array(uint64_t *nanosecondsSinceEpoch, int nItems);

/**
 * Creates a plc4c_data with char*
 * @param size the size of the string
 * @param s the char* value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_string_data(unsigned int size, char *s);

/**
 * Creates a plc4c_data with a constant char
 * @param s the char value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_char_data(char s);

/**
 * Creates a plc4c_data with a constant 16bit wchar
 * @param s the wchar value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_wchar_data(wchar_t s);

/**
 * Creates a plc4c_data which contains a list of values
 * @param list the list value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_list_data(plc4c_list* list);

/**
 * Create a plc4c_data which contains a bit-string of 8 boolean values
 * @param ui unsigned 8 bit integer value.
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_byte_data(uint8_t ui);

/**
 * Create a plc4c_data which contains a bit-string of 16 boolean values
 * @param ui unsigned 16 bit integer value.
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_word_data(uint16_t ui);

/**
 * Create a plc4c_data which contains a bit-string of 32 boolean values
 * @param ui unsigned 32 bit integer value.
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_dword_data(uint32_t ui);

/**
 * Create a plc4c_data which contains a bit-string of 64 boolean values
 * @param ui unsigned 64 bit integer value.
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_lword_data(uint64_t ui);

/**
 * Set a custom function to call when destroying this data.  Typically when the
 * type is a PLC4C_VOID_POINTER
 * @param data pointer to plc4c_data
 * @param data_custom_destroy the function to call
 */
void plc4c_data_set_custom_destroy(
    plc4c_data *data, plc4c_data_custom_destroy data_custom_destroy);

/**
 * Set a custom function to call when plc4c_data_printf is called.
 * @param data pointer to plc4c_data
 * @param data_custom_printf the function to call
 */
void plc4c_data_set_custom_printf(plc4c_data *data,
                                  plc4c_data_custom_printf data_custom_printf);

/**
 * Print the value of the plc4c_data to std out as printf does, formatted by the
 * type. If plc4c_data_set_custom_printf has been called and set, then that
 * function will print the value.
 * @param data
 */
void plc4c_data_printf(plc4c_data *data);

/**
 * Delete the plc4c_data correctly accounting for the type.
 * If plc4c_data_set_custom_destroy has been called and set, then that function
 * will handle the destruction of the data.
 * @param data pointer to plc4c_data
 */
void plc4c_data_destroy(plc4c_data *data);


void* plc4c_data_update_values(plc4c_data *data, void *newData);

#ifdef __cplusplus
}
#endif

#endif  // PLC4C_API_INCLUDE_PLC4C_DATA_H_
