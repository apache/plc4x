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

#define _GNU_SOURCE
/*#ifndef CLOCK_MONOTONIC
#define CLOCK_MONOTONIC 1
#endif*/


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


bool syncBoolLoop(plc4c_connection *conn, 
  bool (passCallback)(plc4c_connection* conn),
  bool (failCallback)(plc4c_connection *conn)  ) {

  plc4c_system *system;
  system = plc4c_connection_get_system(conn);
  while (true) {
    if (passCallback(conn))
      return EXIT_SUCCESS;
    else if (failCallback(conn))
      return EXIT_FAILURE;
    if (plc4c_system_loop(system) != OK)
      return EXIT_FAILURE;
  }
}

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
  bool errorFlag = false;
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

  // ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
  // ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
  bool doRead = 1;
  bool doWrite = 1;
  // ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
  // ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

  // Connection string argument and do_write test arg (defaults to off as
  // not currently working)
  connection_test_string = DEFAULT_CONNECTION_TEST_STRING;

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

  unsigned int loopTimes;
  if (argc >= 2)
    loopTimes = (unsigned int) atoi(argv[1]);
  else
    loopTimes = 1;

  unsigned int loopCount = 0;
  #define NREAD 3
  float valuetowrite[NREAD] = {1.1 ,2.2, 3.3};//, 5.5};

  #define ITEM_STR "%DB2:4.0:REAL[3]"

  // Central program loop ...
  syncBoolLoop(connection, plc4c_connection_get_connected, plc4c_connection_has_error);
  bool rtMalloc = false;
  
  while ( (!errorFlag) && (loopCount++ < loopTimes)) {
    if (doWrite) {
      { // Write request create scope
        if (loopCount == 1 || rtMalloc) {
          result = plc4c_connection_create_write_request(connection, &write_request);
          CHECK_RESULT(result != OK, result,"plc4c_connection_create_write_request failed\n");
          loopback_data = plc4c_data_create_real_array(valuetowrite,NREAD);
          valuetowrite[(loopCount-1)%NREAD]++;
          result = plc4c_write_request_add_item(write_request, ITEM_STR, loopback_data);
        } else {
          plc4c_data_update_values(loopback_data, valuetowrite);
          valuetowrite[(loopCount-1)%NREAD]++;
          //valuetowrite++;
        }
        

        #ifdef S7_LOOPBACK_TIME_IO
          clock_gettime(CLOCK_MONOTONIC,&start);
        #endif

        result = plc4c_write_request_execute(write_request, &write_request_execution);
        
        CHECK_RESULT(result != OK, result,"plc4c_write_request_execute failed\n");
        
        if (plc4c_system_loop(system) != OK) {
          printf("ERROR in the system loop\n");
          break;
        }
      }

      { // Write request sent scope
        while(1) {
          if (plc4c_write_request_check_finished_successfully(write_request_execution)) {
            break;
          } else if (plc4c_write_request_execution_check_completed_with_error(write_request_execution)) {
            printf("FAILED\n");
            errorFlag = true;
            break;
          }
          if (plc4c_system_loop(system) != OK) {
            printf("ERROR in the system loop\n");
            errorFlag = true;
            break;
          }
        }
      }

      { // WRITE_RESPONSE_RECEIVED scope
        write_response = plc4c_write_request_execution_get_response(write_request_execution);
        CHECK_RESULT(write_response == NULL, -1,"plc4c_write_request_execution_get_response failed (no response)\n");
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
          printf("Write %ldus%c", delta_us, doRead ? '\t' : '\n');
        #endif
        
        if (rtMalloc)
          plc4c_write_request_destroy(write_request);
  
        plc4c_write_response_destroy(write_response);
        plc4c_write_request_execution_destroy(write_request_execution);

        if (plc4c_system_loop(system) != OK) {
          printf("ERROR in the system loop\n");
          break;
        }
      }
    } // end of doWrite


    if (doRead) {
      { // READ_REQUEST_CREATE scope
        
        result = plc4c_connection_create_read_request(connection, &read_request);
        CHECK_RESULT(result != OK, result, "plc4c_connection_create_read_request failed\n");
        
        result = plc4c_read_request_add_item(read_request, "A_REQUEST", ITEM_STR);
        CHECK_RESULT(result != OK, result, "plc4c_read_request_add_item failed\n");

        result = plc4c_read_request_execute(read_request, &read_request_execution);
        CHECK_RESULT(result != OK, result, "plc4c_read_request_execute failed\n");
        
        #ifdef S7_LOOPBACK_TIME_IO
          clock_gettime(CLOCK_MONOTONIC,&start);
        #endif

        if (plc4c_system_loop(system) != OK) {
          printf("ERROR in the system loop\n");
          break;
        }
      }


      { // READ_REQUEST_SENT scope
        while(true) {
          if (plc4c_read_request_execution_check_finished_successfully(read_request_execution)) {
            #ifdef S7_LOOPBACK_TIME_IO
              clock_gettime(CLOCK_MONOTONIC,&finish);
              diff_s = finish.tv_sec - start.tv_sec;
              diff_ns = finish.tv_nsec - start.tv_nsec;
              delta_us = (diff_s * 1000000L) + (diff_ns / 1000L);
              printf("Read %ldus\n", delta_us);
            #endif
            break;
          } else if (plc4c_read_request_execution_check_finished_with_error(read_request_execution)) {
            printf("plc4c_read_request_execution_check_finished_with_error FAILED\n");
            errorFlag= true;
            break;
          }
          if (plc4c_system_loop(system) != OK) {
            printf("ERROR in the system loop\n");
            errorFlag = true;
            break;
          }
        }
      }


      { // READ_RESPONSE_RECEIVED scope

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

        plc4c_read_response_destroy(read_response);
        plc4c_read_request_execution_destroy(read_request_execution);
        plc4c_read_request_destroy(read_request);

        if (plc4c_system_loop(system) != OK) {
          printf("ERROR in the system loop\n");
          errorFlag = true;
          break;
        }
      }
    } // end of do read
  } 

  if (!rtMalloc && doWrite) 
    plc4c_write_request_destroy(write_request);
  
  // Start disconnecting, break on error or dis-connection
  result = plc4c_connection_disconnect(connection);
  CHECK_RESULT(result != OK, -1,"plc4c_connection_disconnect failed\n");

  while (true) {
    if (plc4c_system_loop(system) != OK) {
      printf("ERROR in the system loop\n");
      return EXIT_FAILURE;
    } else if (!plc4c_connection_get_connected(connection)) {
      plc4c_system_remove_connection(system, connection);
      plc4c_connection_destroy(connection);
      break;
    }
  }

  // Make sure everything is cleaned up correctly then destroy the 
  // plc4c_system, freeing up all memory allocated by plc4
  plc4c_system_shutdown(system);
  plc4c_system_destroy(system);
  return 0;
}

//#pragma clang diagnostic pop
