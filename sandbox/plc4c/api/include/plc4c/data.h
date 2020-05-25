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
#ifndef PLC4C_API_INCLUDE_PLC4C_DATA_H_
#define PLC4C_API_INCLUDE_PLC4C_DATA_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include <stdint.h>

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
plc4c_data *plc4c_data_create_boolean_data(bool b);

/**
 * Creates a plc4c_data with char
 * @param c the char value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_char_data(char c);

/**
 * Creates a plc4c_data with unsigned char
 * @param uc the unsigned char value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_uchar_data(unsigned char uc);

/**
 * Creates a plc4c_data with short
 * @param s the short value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_short_data(short s);

/**
 * Creates a plc4c_data with unsigned short
 * @param us the unsigned short value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_ushort_data(unsigned short us);

/**
 * Creates a plc4c_data with int
 * @param i the int value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_int_data(int i);

/**
 * Creates a plc4c_data with unsigned int
 * @param ui the unsigned int value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_uint_data(unsigned int ui);

/**
 * Creates a plc4c_data with void*
 * @param v the void* value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_void_pointer_data(void *v);

/**
 * Creates a plc4c_data with float
 * @param f the float value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_float_data(float f);

/**
 * Creates a plc4c_data with char*
 * @param size the size of the string
 * @param s the char* value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_string_data(unsigned int size, char *s);

/**
 * Creates a plc4c_data with a constant char*
 * @param size the size of the string
 * @param s the char *value
 * @return pointer to plc4c_data
 */
plc4c_data *plc4c_data_create_constant_string_data(unsigned int size, char *s);

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

#ifdef __cplusplus
}
#endif

#endif  // PLC4C_API_INCLUDE_PLC4C_DATA_H_
