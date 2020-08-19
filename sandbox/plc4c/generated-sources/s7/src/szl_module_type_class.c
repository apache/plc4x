/*
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License"); you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing,
  software distributed under the License is distributed on an
  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  KIND, either express or implied.  See the License for the
  specific language governing permissions and limitations
  under the License.
*/

#include "szl_module_type_class.h"


// Create an empty NULL-struct
static const plc4c_s7_read_write_szl_module_type_class plc4c_s7_read_write_szl_module_type_class_null_const;

plc4c_s7_read_write_szl_module_type_class plc4c_s7_read_write_szl_module_type_class_null() {
  return plc4c_s7_read_write_szl_module_type_class_null_const;
}

int plc4c_s7_read_write_szl_module_type_class_num_values() {
  return 4;
}

