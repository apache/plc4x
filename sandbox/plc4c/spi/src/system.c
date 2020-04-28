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

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <plc4c/system.h>
#include "plc4c/spi/types_private.h"
#include "plc4c/spi/system_private.h"

return_code plc4c_system_create(plc4c_system **system) {
    plc4c_system *new_system = malloc(sizeof(plc4c_system));
    new_system->driver_list_head = NULL;
    new_system->connection_list_head = NULL;
    *system = new_system;
    return OK;
}

void plc4c_system_destroy(plc4c_system *system) {
    // TODO: So some more cleaning up ...
    free(system);
}

void plc4c_system_set_on_driver_load_success_callback(plc4c_system *system,
                                                      plc4c_system_on_driver_load_success_callback callback) {
    system->on_driver_load_success_callback = callback;
}

void plc4c_system_set_on_driver_load_failure_callback(plc4c_system *system,
                                                      plc4c_system_on_driver_load_failure_callback callback) {
    system->on_driver_load_failure_callback = callback;
}

void plc4c_system_set_on_connect_success_callback(plc4c_system *system,
                                                  plc4c_system_on_connect_success_callback callback) {
    system->on_connect_success_callback = callback;
}

void plc4c_system_set_on_connect_failure_callback(plc4c_system *system,
                                                  plc4c_system_on_connect_failure_callback callback) {
    system->on_connect_failure_callback = callback;
}

void plc4c_system_set_on_disconnect_success_callback(plc4c_system *system,
                                                     plc4c_system_on_disconnect_success_callback callback) {
    system->on_disconnect_success_callback = callback;
}

void plc4c_system_set_on_disconnection_failure_callback(plc4c_system *system,
                                                        plc4c_system_on_disconnect_failure_callback callback) {
    system->on_disconnect_failure_callback = callback;
}

void plc4c_system_set_on_loop_failure_callback(plc4c_system *system,
                                               plc4c_system_on_loop_failure_callback callback) {
    system->on_loop_failure_callback = callback;
}

return_code plc4c_system_add_driver(plc4c_system *system,
                                    plc4c_driver *driver) {
    // If the system is not initialized, return an error.
    // There is nothing we can do here.
    if (system == NULL) {
        return INTERNAL_ERROR;
    }

    // Get the first element of the driver list.
    plc4c_driver_list_item *cur_driver = system->driver_list_head;

    // If the driver list is empty. Start a new list of drivers.
    if (cur_driver == NULL) {
        system->driver_list_head = malloc(sizeof(plc4c_driver_list_item));
        system->driver_list_head->driver = driver;
        system->driver_list_head->next = NULL;
    }
        // Drivers are already listed, add the currentdriver to the end of the list.
    else {
        // Go to the last driver in the list.
        while (cur_driver->next != NULL) {
            cur_driver = cur_driver->next;
        }

        // Add a new driver item to the end of the list.
        cur_driver->next = malloc(sizeof(plc4c_driver_list_item));
        cur_driver->next->driver = driver;
        cur_driver->next->next = NULL;
    }
    return OK;
}

return_code plc4c_system_init(plc4c_system *system) {
    // Nothing to really do at the moment.

    return OK;
}

void plc4c_system_shutdown(plc4c_system *system) {

}

enum connection_string_parser_state {
    PROTOCOL_CODE,
    TRANSPORT_CODE,
    TRANSPORT_CONNECTION_INFORMATION,
    PARAMETERS,
    FINISHED
};

return_code plc4c_system_create_connection(const char *connection_string, plc4c_connection** connection) {
    // Count the number of colons and question-marks so we know which pattern to use for
    // matching and how large the arrays for containing the different segments should be.
    int num_colons = 0;
    int num_question_marks = 0;
    char* protocol_code = NULL;
    char* transport_code = NULL;
    char* transport_connect_information = NULL;
    char* parameters = NULL;
    int start_segment_index = 0;
    const char* start_segment = connection_string;
    // The connection string has two parts ... the first, where colons are the delimiters
    // and the second where a question mark is the delimiter.
    enum mode {
        SEARCHING_FOR_COLONS,
        SEARCHING_FOR_QUESTION_MARKS
    } mode = SEARCHING_FOR_COLONS;
    for(int i = 0; i <= strlen(connection_string); i++) {
        // If we're in the first part of the connection string ... watch out for colons.
        if(mode == SEARCHING_FOR_COLONS) {
            // If we encounter a colon, depending on the number of colons already found, save the information in
            // either the protocol code or transport code variable.
            if (*(connection_string + i) == ':') {
                num_colons++;
                // The first colon delimits the protocol-code.
                if (num_colons == 1) {
                    // Allocate enough memory to hold the sub-string.
                    protocol_code = malloc(sizeof(char) * ((i - start_segment_index) + 1));
                    // Copy the sub-string to the freshly allocated memory area.
                    strlcpy(protocol_code, start_segment, (i - start_segment_index) + 1);

                    // Set the start of the next segment to directly after the colon.
                    start_segment_index = i + 1;
                    start_segment = connection_string + start_segment_index;

                    // If the following character would be a slash, we're probably finished and no transport code is
                    // provided. If this is the case, ensure it's actually a double-slash and if this is the case
                    // switch to the question-mark searching mode.
                    if(*start_segment == '/') {
                        if(*(start_segment + 1) == '/') {
                            mode = SEARCHING_FOR_QUESTION_MARKS;
                            start_segment_index += 2;
                            start_segment += 2;
                        } else {
                            return INVALID_CONNECTION_STRING;
                        }
                    }
                }
                // If we encountered a second colon, this is the transport code.
                else if (num_colons == 2) {
                    // Allocate enough memory to hold the sub-string.
                    transport_code = malloc(sizeof(char) * ((i - start_segment_index) + 1));
                    // Copy the sub-string to the freshly allocated memory area.
                    strlcpy(transport_code, start_segment, (i - start_segment_index) + 1);

                    // Set the start of the next segment to directly after the colon.
                    start_segment_index = i + 1;
                    start_segment = connection_string + start_segment_index;

                    // The transport code is allways followed by "://". So check if this is the case.
                    // If it is, switch to question-mark searching mode.
                    if((*start_segment != '/') || (*(start_segment + 1) != '/')) {
                        return INVALID_CONNECTION_STRING;
                    }
                    mode = SEARCHING_FOR_QUESTION_MARKS;

                    // Bump the start of the segment to after the double slashes.
                    start_segment_index += 2;
                    start_segment += 2;
                } else {
                    return INVALID_CONNECTION_STRING;
                }
            }
        }
        // If we're in the second part, look for question marks.
        else {
            // The question-mark separates the transport connect information from the parameters.
            if (*(connection_string + i) == '?') {
                num_question_marks++;
                // Only one question-mark is allowed in connection strings.
                if (num_question_marks > 1) {
                    return INVALID_CONNECTION_STRING;
                }

                // Allocate enough memory to hold the sub-string.
                transport_connect_information = malloc(sizeof(char) * ((i - start_segment_index) + 1));
                // Copy the sub-string to the freshly allocated memory area.
                strlcpy(transport_connect_information, start_segment, (i - start_segment_index) + 1);

                // Set the start of the next segment to directly after the question-mark.
                start_segment_index = i + 1;
                start_segment = connection_string + start_segment_index;
            }
            // This is the last character ... finish up the last loose end.
            if (i == strlen(connection_string)) {
                // If no question-mark has been encountered, this connection string doesn't have one and the
                // remaining part is simply the transport connect information.
                if (num_question_marks == 0) {
                    transport_connect_information = malloc(sizeof(char) * ((i - start_segment_index) + 1));
                    strlcpy(transport_connect_information, start_segment, (i - start_segment_index) + 1);
                }
                // I a question-mark was found, this is the paramters section.
                else {
                    parameters = malloc(sizeof(char) * (i - start_segment_index));
                    strlcpy(parameters, start_segment, (i - start_segment_index) + 1);
                }
            }
        }
    }
    if(num_colons == 0) {
        return INVALID_CONNECTION_STRING;
    }

    // Initialize a new connection data-structure with the parsed information.
    plc4c_connection* new_connection = malloc(sizeof(plc4c_connection));
    new_connection->connection_string = connection_string;
    new_connection->protocol_code = protocol_code;
    new_connection->transport_code = transport_code;
    new_connection->transport_connect_information = transport_connect_information;
    new_connection->parameters = parameters;
    *connection = new_connection;

    return OK;
}

return_code plc4c_system_connect(plc4c_system *system,
                                 const char *connection_string,
                                 plc4c_connection **connection) {

    // Parse the connection string and initialize some of the connection field variables from this.
    plc4c_connection* new_connection = NULL;
    return_code result = plc4c_system_create_connection(connection_string, &new_connection);
    if(result != OK) {
        return result;
    }

    // Find a matching driver from the driver-list
    plc4c_driver_list_item* cur_driver_list_item = system->driver_list_head;
    // If no driver is available at all this is devinitely a developer error,
    // so we output a special error code for this case
    if(cur_driver_list_item == NULL) {
        return NO_DRIVER_AVAILABLE;
    }
    do {
        if (strcmp(cur_driver_list_item->driver->protocol_code, new_connection->protocol_code) == 0) {
            // Set the driver reference in the new connection.
            new_connection->driver = cur_driver_list_item->driver;
            // If no transport was selected, use the drivers default transport (if it exists).
            if(new_connection->transport_code == NULL) {
                if(cur_driver_list_item->driver->default_transport_code != NULL) {
                    new_connection->transport_code = cur_driver_list_item->driver->default_transport_code;
                }
            }
            break;
        }
        cur_driver_list_item = cur_driver_list_item->next;
    } while (cur_driver_list_item != NULL);

    // If the driver property is still NULL, the desired driver was not found.
    if(new_connection->driver == NULL) {
        return UNKNOWN_DRIVER;
    }

    // Return an error if the user didn't specify a transport and the driver doesn't have a default one.
    if(new_connection->transport_code == NULL) {
        return UNSPECIFIED_TRANSPORT;
    }

    // Pass the new connection back.
    *connection = new_connection;

    return OK;
}

return_code plc4c_system_loop() {
    return OK;
}


