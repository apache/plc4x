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
#include <plc4c/driver_s7.h>
#include <plc4c/plc4c.h>
#include <plc4c/transport_tcp.h>
#include <stdio.h>
#include <string.h>

#include "../../../spi/include/plc4c/spi/types_private.h"

#define DEFAULT_CONNECTION_TEST_STRING "s7:tcp://0.0.0.0:102"
#define TEST_DOUBLE_INT
int numOpenConnections = 0;

#ifndef _WIN32
  #include <unistd.h>
  #include <time.h>
  #define S7_LOOPBACK_TIME_IO
#endif

/**
 * Here we could implement something that keeps track of all open connections.
 * For example on embedded devices using the W5100 SPI Network device, this can
 * only handle 4 simultaneous connections.
 *
 * @param connection the connection that was just established
 */
void onGlobalConnect(plc4c_connection *cur_connection) {
  printf("Connected to %s",
         plc4c_connection_get_connection_string(cur_connection));
  numOpenConnections++;
}

void onGlobalDisconnect(plc4c_connection *cur_connection) {
  printf("Disconnected from %s",
         plc4c_connection_get_connection_string(cur_connection));
  numOpenConnections--;
}



enum plc4c_connection_state_t {
  CONNECTING,
  WRITE_REQUEST_CREATE,
  WRITE_REQUEST_SENT,
  WRITE_RESPONSE_RECEIVED,
  READ_REQUEST_CREATE,
  READ_REQUEST_SENT,
  READ_RESPONSE_RECEIVED,
  DISCONNECTING,
  DISCONNECTED
};
typedef enum plc4c_connection_state_t plc4c_connection_state;

//#pragma clang diagnostic push
//#pragma ide diagnostic ignored "hicpp-multiway-paths-covered"


#define CHECK_RESULT(chk, ret, fs) do {if (chk) {printf(fs); return(ret);}} while(0)

void get_user_loopback_values(int argc, char **argv, long *value) {
  int count = 0;
  while (count < argc)
    value[count++] = atol(argv[count]);
}

int main(int argc, char** argv) {
  
#ifdef S7_LOOPBACK_TIME_IO
  struct timespec start, finish;
  long delta_us, diff_s, diff_ns;
#endif

  char* connection_test_string;
  plc4c_return_code result;
  plc4c_connection_state state;
  bool loop = true;
  int idx = 0;
  plc4c_system *system = NULL;
  plc4c_connection *connection = NULL;
  
  plc4c_read_request *read_request = NULL;
  plc4c_read_request_execution *read_request_execution = NULL;
  plc4c_read_response *read_response;

  plc4c_write_request *write_request = NULL;
  plc4c_write_request_execution *write_request_execution = NULL;
  plc4c_write_response *write_response;

  plc4c_list_element *cur_element;
  plc4c_response_value_item *value_item;

  long loopback_value[7] = {0,0,0,0,0};
  plc4c_data *loopback_data;
  
  // Connection string argument and do_write test arg (defaults to off as
  // not currently working)
  if (argc < 2)
    connection_test_string = DEFAULT_CONNECTION_TEST_STRING;
  else
    connection_test_string = argv[1];

  if (argc < 3) {
    loopback_value[0] = 0; // bool
    loopback_value[1] = 44; // uint8
    loopback_value[2] = -55; // int8
    loopback_value[3] = 666; // uint16
    loopback_value[4] = -777; //int16
    loopback_value[5] = 88888; // uint32
    loopback_value[6] = -99999; //int32
  } else {
    get_user_loopback_values(argc - 2, &argv[2], loopback_value);
  }

  // Initialisation and startup sequence
  result = plc4c_system_create(&system);
  CHECK_RESULT(result != OK, result, "plc4c_system_create failed\n");
  plc4c_driver *s7_driver = plc4c_driver_s7_create();
  result = plc4c_system_add_driver(system, s7_driver);
  CHECK_RESULT(result != OK, result, "plc4c_system_add_driver failed\n");
  plc4c_transport *tcp_transport = plc4c_transport_tcp_create();
  result = plc4c_system_add_transport(system, tcp_transport);
  CHECK_RESULT(result != OK, result, "plc4c_system_add_transport failed\n");
  result = plc4c_system_init(system);
  CHECK_RESULT(result != OK, result, "plc4c_system_init failed\n");

  // Register the global callbacks
  plc4c_system_set_on_connect_success_callback(system, &onGlobalConnect);
  plc4c_system_set_on_disconnect_success_callback(system, &onGlobalDisconnect);

  // Establish connections to remote devices
  result = plc4c_system_connect(system, connection_test_string, &connection);
  CHECK_RESULT(result != OK, result, "plc4c_system_connect failed\n");


  // Central program loop ...
  state = CONNECTING;

  while (loop) {
    
    // Give plc4c a chance to do something. This is where all I/O is done.
    if (plc4c_system_loop(system) != OK) {
      printf("ERROR in the system loop\n");
      break;
    }

    // Depending on the current state, implement some logic
    switch (state) {
      case CONNECTING: {
        // Check if the connection is established:
        if (plc4c_connection_get_connected(connection)) {
          state = WRITE_REQUEST_CREATE;
        } else if (plc4c_connection_has_error(connection)) {
          printf("FAILED\n");
          return -1;
        }
        break;
      }

      
      case WRITE_REQUEST_CREATE: {

        result = plc4c_connection_create_write_request(connection, &write_request);
        CHECK_RESULT(result != OK, result,"plc4c_connection_create_write_request failed\n");

        float valuestowrite[] = {1.1,2.2};
        printf("Writing %f %f to %%DB2:4.0:REAL[2] ...\n", valuestowrite[0], valuestowrite[1]);
        loopback_data = plc4c_data_create_float_array(valuestowrite, 2);
        result = plc4c_write_request_add_item(write_request, "%DB2:0.0:REAL[2]", loopback_data);
        
        float valuetowrite = {4.4};
        loopback_data = plc4c_data_create_float_data(valuetowrite);
        result = plc4c_write_request_add_item(write_request, "%DB2:4.0:REAL", loopback_data);
        
        /*
        printf("Writing %d to %%DB2:0.0:BOOL ...\n", (bool) loopback_value[0]);
        loopback_data = plc4c_data_create_bool_data(true);
        result = plc4c_write_request_add_item(write_request, "%DB2:100.0:BOOL", loopback_data);
        
        printf("Writing %d to %%DB2:4.0:USINT ...\n", (uint8_t) loopback_value[1]);
        loopback_data = plc4c_data_create_uint8_t_data((uint8_t) loopback_value[1]);
        result = plc4c_write_request_add_item(write_request, "%DB2:4.0:USINT", loopback_data);
        
        printf("Writing %d to %%DB2:8.0:SINT ...\n", (int8_t) loopback_value[2]);
        loopback_data = plc4c_data_create_int8_t_data((int8_t) loopback_value[2]);
        result = plc4c_write_request_add_item(write_request, "%DB2:8.0:SINT", loopback_data);
                
        printf("Writing %d to %%DB2:12.0:UINT ...\n", (uint16_t) loopback_value[3]);
        loopback_data = plc4c_data_create_uint16_t_data((uint16_t) loopback_value[3]);
        result = plc4c_write_request_add_item(write_request, "%DB2:12.0:UINT", loopback_data);
        
        printf("Writing %d to %%DB2:16.0:INT ...\n", (int16_t) loopback_value[4]);
        loopback_data = plc4c_data_create_int16_t_data((int16_t) loopback_value[4]);
        result = plc4c_write_request_add_item(write_request, "%DB2:16.0:INT", loopback_data);
               
        printf("Writing %d to %%DB2:20.0:UDINT ...\n", (uint32_t) loopback_value[5]);
        loopback_data = plc4c_data_create_uint32_t_data((uint32_t) loopback_value[5]);
        result = plc4c_write_request_add_item(write_request, "%DB2:20.0:UDINT", loopback_data);
        
        printf("Writing %d to %%DB2:24.0:DINT ...\n", (int32_t) loopback_value[6]);
        loopback_data = plc4c_data_create_int32_t_data((int32_t) loopback_value[6]);
        result = plc4c_write_request_add_item(write_request, "%DB2:24.0:DINT", loopback_data);

        printf("Writing %d to %%DB2:28.0:BYTE ...\n", (uint8_t) loopback_value[1]);
        loopback_data = plc4c_data_create_uint8_t_data((uint8_t) loopback_value[1]);
        result = plc4c_write_request_add_item(write_request, "%DB2:28.0:BYTE", loopback_data);
        
        printf("Writing %d to %%DB2:32.0:WORD ...\n", (uint16_t) loopback_value[3]);
        loopback_data = plc4c_data_create_uint16_t_data((uint16_t) loopback_value[3]);
        result = plc4c_write_request_add_item(write_request, "%DB2:32.0:WORD", loopback_data);

        printf("Writing %d to %%DB2:36.0:DWORD ...\n", (uint32_t) loopback_value[5]);
        loopback_data = plc4c_data_create_uint32_t_data((uint32_t) loopback_value[5]);
        result = plc4c_write_request_add_item(write_request, "%DB2:36.0:DWORD", loopback_data);

        printf("Writing %f to %%DB2:40.0:REAL ...\n", (float) loopback_value[5]/3.14);
        loopback_data = plc4c_data_create_float_data((float) loopback_value[5]/3.14);
        result = plc4c_write_request_add_item(write_request, "%DB2:40.0:REAL", loopback_data);
        */
#ifdef S7_LOOPBACK_TIME_IO
        clock_gettime(CLOCK_MONOTONIC,&start);
#endif
        result = plc4c_write_request_execute(write_request, &write_request_execution);
        CHECK_RESULT(result != OK, result,"plc4c_write_request_execute failed\n");
        
        state = WRITE_REQUEST_SENT;
        break;
      }

      case WRITE_REQUEST_SENT: {
        if (plc4c_write_request_check_finished_successfully(write_request_execution)) {
          printf("Success\n");
          state = WRITE_RESPONSE_RECEIVED;
        } else if (plc4c_write_request_execution_check_completed_with_error(write_request_execution)) {
          printf("FAILED\n");
          return -1;
        }
        break;
      }

      case WRITE_RESPONSE_RECEIVED: {
        write_response = plc4c_write_request_execution_get_response(write_request_execution);
        CHECK_RESULT(write_response == NULL, -1,"plc4c_write_request_execution_get_response failed (no responce)\n");
        cur_element = plc4c_utils_list_tail(write_response->response_items);
        idx = 0;
        while (cur_element != NULL) {
          plc4c_response_item *checker = (plc4c_response_item*) cur_element->value;
          printf("Write item %d status: '%s'\n", idx++,
                 plc4c_response_code_to_message(checker->response_code));
          cur_element = cur_element->next;
        }

#ifdef S7_LOOPBACK_TIME_IO
        clock_gettime(CLOCK_MONOTONIC,&finish);
        diff_s = finish.tv_sec - start.tv_sec;
        diff_ns = finish.tv_nsec - start.tv_nsec;
        delta_us = (diff_s * 1000000L) + (diff_ns / 1000L);
        printf("Write took %ld us\n", delta_us);
#endif
        state = READ_REQUEST_CREATE;
        plc4c_write_destroy_write_response(write_response);
        plc4c_write_request_execution_destroy(write_request_execution);
        plc4c_write_request_destroy(write_request);
        break;
      }

      case READ_REQUEST_CREATE: {
        result = plc4c_connection_create_read_request(connection, &read_request);
        CHECK_RESULT(result != OK, result, "plc4c_connection_create_read_request failed\n");
        
        result = plc4c_read_request_add_item(read_request, "WORD", "%DB2:0.0:WORD[163]");
        CHECK_RESULT(result != OK, result, "plc4c_read_request_add_item failed\n");
        
        result = plc4c_read_request_add_item(read_request, "REAL","%DB2:8.0:REAL[63]");
        CHECK_RESULT(result != OK, result, "plc4c_read_request_add_item failed\n");
        
        /*
        result = plc4c_read_request_add_item(read_request, "BOOL", "%DB2:100.0:BOOL");
        CHECK_RESULT(result != OK, result, "plc4c_read_request_add_item failed\n");
        
        result = plc4c_read_request_add_item(read_request, "USINT", "%DB2:4.0:USINT");
        CHECK_RESULT(result != OK, result, "plc4c_read_request_add_item failed\n");
        
        result = plc4c_read_request_add_item(read_request, "SINT", "%DB2:8.0:SINT");
        CHECK_RESULT(result != OK, result, "plc4c_read_request_add_item failed\n");
        
        result = plc4c_read_request_add_item(read_request, "UINT", "%DB2:12.0:UINT");
        CHECK_RESULT(result != OK, result, "plc4c_read_request_add_item failed\n");
        
        result = plc4c_read_request_add_item(read_request, "INT", "%DB2:16.0:INT");
        CHECK_RESULT(result != OK, result, "plc4c_read_request_add_item failed\n");
        
        result = plc4c_read_request_add_item(read_request, "UDINT", "%DB2:20.0:UDINT");
        CHECK_RESULT(result != OK, result, "plc4c_read_request_add_item failed\n");
        
        result = plc4c_read_request_add_item(read_request, "DINT", "%DB2:24.0:DINT");
        CHECK_RESULT(result != OK, result, "plc4c_read_request_add_item failed\n");

        result = plc4c_read_request_add_item(read_request,"BYTE", "%DB2:28.0:BYTE");
        CHECK_RESULT(result != OK, result, "plc4c_read_request_add_item failed\n");

        result = plc4c_read_request_add_item(read_request, "WORD", "%DB2:32.0:WORD");
        CHECK_RESULT(result != OK, result, "plc4c_read_request_add_item failed\n");

        result = plc4c_read_request_add_item(read_request, "DWORD", "%DB2:36.0:DWORD");
        CHECK_RESULT(result != OK, result, "plc4c_read_request_add_item failed\n");
 
        result = plc4c_read_request_add_item(read_request, "REAL", "%DB2:40.0:REAL");
        CHECK_RESULT(result != OK, result, "plc4c_read_request_add_item failed\n");
        */
        result = plc4c_read_request_execute(read_request, &read_request_execution);
        CHECK_RESULT(result != OK, result, "plc4c_read_request_execute failed\n");
        
        state = READ_REQUEST_SENT;

#ifdef S7_LOOPBACK_TIME_IO
        clock_gettime(CLOCK_MONOTONIC,&start);
#endif
        break;
      }

      // Wait until the read-request execution is finished
      case READ_REQUEST_SENT: {
        if (plc4c_read_request_execution_check_finished_successfully(read_request_execution)) {
          state = READ_RESPONSE_RECEIVED;
#ifdef S7_LOOPBACK_TIME_IO
          clock_gettime(CLOCK_MONOTONIC,&finish);
          diff_s = finish.tv_sec - start.tv_sec;
          diff_ns = finish.tv_nsec - start.tv_nsec;
          delta_us = (diff_s * 1000000L) + (diff_ns / 1000L);
          printf("Read took %ld us\n", delta_us);
#endif
        } else if (plc4c_read_request_execution_check_finished_with_error(read_request_execution)) {
          printf("plc4c_read_request_execution_check_finished_with_error FAILED\n");
          return -1;
        }
        break;
      }

      case READ_RESPONSE_RECEIVED: {

        read_response = plc4c_read_request_execution_get_response(read_request_execution);
        CHECK_RESULT(read_response == NULL, -1, "plc4c_read_request_execution_get_response failed (No Response)\n");
        
        // Iterate over all returned items.
        cur_element = plc4c_utils_list_tail(read_response->items);
        while (cur_element != NULL) {
          value_item = cur_element->value;
          printf("Value %s (%s): ", value_item->item->name,
                 plc4c_response_code_to_message(value_item->response_code));
          plc4c_data_printf(value_item->value);
          printf("\n");
          cur_element = cur_element->next;
        }

        plc4c_read_destroy_read_response(read_response);
        plc4c_read_request_execution_destroy(read_request_execution);
        plc4c_read_request_destroy(read_request);

        result = plc4c_connection_disconnect(connection);
        CHECK_RESULT(result != OK, -1,"plc4c_connection_disconnect failed\n");
        state = DISCONNECTING;
        break;
      }
      
      // Wait until the connection is disconnected
      case DISCONNECTING: {
        if (!plc4c_connection_get_connected(connection)) {
          printf("SUCCESS\n");
          plc4c_system_remove_connection(system, connection);
          plc4c_connection_destroy(connection);
          state = DISCONNECTED;
          loop = false;
        }
        break;
      }
      
      // End the loop
      case DISCONNECTED: {
        loop = false;
        break;
      }
      default: {
      }
    }

  }

  // Make sure everything is cleaned up correctly then destroy the 
  // plc4c_system, freeing up all memory allocated by plc4
  plc4c_system_shutdown(system);
  plc4c_system_destroy(system);
  return 0;
}

//#pragma clang diagnostic pop
