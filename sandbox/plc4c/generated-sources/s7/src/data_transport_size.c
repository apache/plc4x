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

#include "data_transport_size.h"


// Create an empty NULL-struct
static const plc4c_s7_read_write_data_transport_size plc4c_s7_read_write_data_transport_size_null_const;

plc4c_s7_read_write_data_transport_size plc4c_s7_read_write_data_transport_size_null() {
  return plc4c_s7_read_write_data_transport_size_null_const;
}


bool plc4c_s7_read_write_data_transport_size_get_size_in_bits(plc4c_s7_read_write_data_transport_size value) {
  switch(value) {
    case 0: { /* '0x00' */
      return false;
    }
    case 3: { /* '0x03' */
      return true;
    }
    case 4: { /* '0x04' */
      return true;
    }
    case 5: { /* '0x05' */
      return true;
    }
    case 6: { /* '0x06' */
      return false;
    }
    case 7: { /* '0x07' */
      return false;
    }
    case 9: { /* '0x09' */
      return false;
    }
    default: {
      return 0;
    }
  }
}
