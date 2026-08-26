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

#ifndef PLC4C_CONTEXT_H
#define PLC4C_CONTEXT_H

#include <stdint.h>
#include <stdbool.h>

#include <plc4c/types.h>

/**
 * How deeply a message may nest its types. Generated parsers are plain recursive
 * descent, so a type that contains itself turns the sender's byte count into our
 * call depth - and a C stack that runs out takes the process with it rather than
 * the message.
 *
 * The deepest message in the project's own test suites nests 36 levels, and the
 * widest generated parse frame measures 272 bytes, so the whole budget costs
 * under 300 KB of stack against a default thread's 8 MB. It is the same number
 * the other bindings use, so that one setting of PLC4C_MAX_NESTING_DEPTH_ENV
 * means the same thing whichever of them is reading.
 */
#define PLC4C_DEFAULT_MAX_NESTING_DEPTH 1024

/**
 * Names the environment variable that moves the bound for the device whose
 * messages genuinely nest deeper than anything we have seen. It carries the same
 * meaning as in the other bindings: a positive depth, and anything else leaves
 * the default in place.
 */
#define PLC4C_MAX_NESTING_DEPTH_ENV "PLC4X_MAX_NESTING_DEPTH"

struct plc4x_spi_context {
  uint16_t numItems;
  uint16_t curItem;
  uint16_t depth;
};
typedef struct plc4x_spi_context plc4x_spi_context;

plc4x_spi_context plc4x_spi_context_background();

plc4x_spi_context plc4x_spi_context_create_array_context(plc4x_spi_context context, uint16_t numItems, uint16_t curItem);

/**
 * Descends one type deeper, or reports a parse error once the nesting has gone
 * further than the bound allows.
 *
 * Every parse function takes its context by value and hands its own copy to the
 * types it contains, so raising the depth on that copy bounds the whole subtree
 * below it and needs no matching call on the way out - the count unwinds when the
 * function returns, whether it returned a value or an error.
 */
plc4c_return_code plc4x_spi_context_enter_type(plc4x_spi_context* context);

/** How many types deep this context sits. */
uint16_t plc4x_spi_context_get_depth(plc4x_spi_context ctx);

/** The bound in force, which is the default unless the environment says otherwise. */
uint16_t plc4x_spi_context_get_max_depth();

/**
 * Reads a bound out of the value of PLC4C_MAX_NESTING_DEPTH_ENV.
 *
 * @param configured the value of the variable, or NULL when it is unset
 * @return the depth it asks for, or PLC4C_DEFAULT_MAX_NESTING_DEPTH when it asks
 * for nothing usable
 */
uint16_t plc4x_spi_context_resolve_max_depth(const char* configured);

uint16_t plc4x_spi_context_get_num_items_from_context(plc4x_spi_context ctx);

uint16_t plc4x_spi_context_get_cur_item_from_context(plc4x_spi_context ctx);

bool plc4x_spi_context_get_last_item_from_context(plc4x_spi_context ctx);

#endif  // PLC4C_CONTEXT_H
