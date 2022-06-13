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

#include <ctype.h>
#include <plc4c/driver_plc4x.h>
#include <plc4c/spi/types_private.h>
#include <stdlib.h>
#include <string.h>

#include "plc4c/driver_plc4x_packets.h"
#include "plc4x_message.h"

plc4c_return_code plc4c_driver_plc4x_send_packet(plc4c_connection* connection,
                                                 plc4c_plc4x_read_write_plc4x_message* packet) {
  uint16_t packet_size;
  plc4c_spi_write_buffer* write_buffer;
  plc4c_return_code return_code;

  // Get the size required to contain the serialized form of this packet.
  packet_size = plc4c_plc4x_read_write_plc4x_message_length_in_bytes(packet);

  // Serialize this message to a byte-array.
  return_code = plc4c_spi_write_buffer_create(packet_size, &write_buffer);
  if (return_code != OK) {
    return return_code;
  }
  return_code = plc4c_plc4x_read_write_plc4x_message_serialize(write_buffer, packet);
  if(return_code != OK) {
    return return_code;
  }

  // Now send this to the recipient.
  return_code = connection->transport->send_message(connection->transport_configuration, write_buffer);
  if (return_code != OK) {
    return return_code;
  }

  /* TODO: free when relevant, here works I think*/
  free(write_buffer->data);
  free(write_buffer);
  return OK;
}

int16_t plc4c_driver_plc4x_select_message_function(uint8_t* buffer_data,
                                                   uint16_t buffer_length) {
  if (buffer_length >= 3) {
    uint16_t length = ((uint16_t) *(buffer_data + 1)) << 8 | ((uint16_t) *(buffer_data + 2));
    return length;
  }

  // In all other cases, we'll just have to wait for the next time.
  return 0;
}

plc4c_return_code plc4c_driver_plc4x_receive_packet(plc4c_connection* connection,
                                                    plc4c_plc4x_read_write_plc4x_message** packet) {
  // Check with the transport if there is a packet available.
  // If it is, get a read_buffer for reading it.
  plc4c_spi_read_buffer* read_buffer;
  plc4c_return_code return_code;

  return_code = connection->transport->select_message(
      connection->transport_configuration, 3,
      plc4c_driver_plc4x_select_message_function, &read_buffer);

  // OK is only returned if a packet is available.
  if (return_code != OK) {
    return return_code;
  }

  // Parse the packet by consuming the read_buffer data.
  *packet = NULL;
  return_code = plc4c_plc4x_read_write_plc4x_message_parse(read_buffer, packet);
  if (return_code != OK) {
    return return_code;
  }

  // TODO: verify and maybe move
  free(read_buffer->data);
  free(read_buffer);

  // In this case a packet was available and parsed.
  return OK;
}

void destroy_plc4x_read_request_fields(plc4c_list_element *element) {
  plc4c_plc4x_read_write_plc4x_field_request *item;
  item = element->value;

  free(item->field->name);
  free(item->field->field_query);
  free(item->field);
}

void destroy_plc4x_read_response_fields(plc4c_list_element *element) {
  plc4c_plc4x_read_write_plc4x_field_value_response *item;
  item = element->value;

  free(item->field->name);
  free(item->field->field_query);
  free(item->field);
  // TODO: Free PlcValue
}

void destroy_plc4x_write_request_fields(plc4c_list_element *element) {
  plc4c_plc4x_read_write_plc4x_field_value_request *item;
  item = element->value;

  free(item->field->name);
  free(item->field->field_query);
  free(item->field);
  // TODO: Free PlcValue
}

void destroy_plc4x_write_response_fields(plc4c_list_element *element) {
  plc4c_plc4x_read_write_plc4x_field_response *item;
  item = element->value;
  // Nothing to do here.
}

void plc4c_driver_plc4x_destroy_packet(plc4c_plc4x_read_write_plc4x_message* packet) {
  if (packet) {
    switch (packet->_type) {
      case plc4c_plc4x_read_write_plc4x_message_type_plc4c_plc4x_read_write_plc4x_connect_request:
        free(packet->plc4x_connect_request_connection_string);
        break;

      case plc4c_plc4x_read_write_plc4x_message_type_plc4c_plc4x_read_write_plc4x_connect_response:
        break;

      case plc4c_plc4x_read_write_plc4x_message_type_plc4c_plc4x_read_write_plc4x_read_request:
        plc4c_utils_list_delete_elements(packet->plc4x_read_request_fields,
                                         destroy_plc4x_read_request_fields);
        free(packet->plc4x_read_request_fields);
        break;

      case plc4c_plc4x_read_write_plc4x_message_type_plc4c_plc4x_read_write_plc4x_read_response:
        plc4c_utils_list_delete_elements(packet->plc4x_read_response_fields,
                                         destroy_plc4x_read_response_fields);
        free(packet->plc4x_read_request_fields);
        break;

      case plc4c_plc4x_read_write_plc4x_message_type_plc4c_plc4x_read_write_plc4x_write_request:
        plc4c_utils_list_delete_elements(packet->plc4x_write_request_fields,
                                         destroy_plc4x_write_request_fields);
        free(packet->plc4x_read_request_fields);
        break;

      case plc4c_plc4x_read_write_plc4x_message_type_plc4c_plc4x_read_write_plc4x_write_response:
        plc4c_utils_list_delete_elements(packet->plc4x_write_response_fields,
                                         destroy_plc4x_write_response_fields);
        free(packet->plc4x_read_request_fields);
        break;
    }
    free(packet);
  }
}

plc4c_return_code plc4c_driver_plc4x_create_connection_request(plc4c_driver_plc4x_config* configuration,
                                                               plc4c_plc4x_read_write_plc4x_message** connect_request_packet) {
  *connect_request_packet = malloc(sizeof(plc4c_plc4x_read_write_plc4x_message));
  if (*connect_request_packet == NULL) {
    return NO_MEMORY;
  }

  (*connect_request_packet)->_type = plc4c_plc4x_read_write_plc4x_message_type_plc4c_plc4x_read_write_plc4x_connect_request;
  (*connect_request_packet)->request_id = configuration->request_id++;

  // Copy the connection string.
  size_t size = sizeof(char) * (strlen(configuration->remote_connection_string) + 1);
  (*connect_request_packet)->plc4x_connect_request_connection_string = malloc(size);
  if((*connect_request_packet)->plc4x_connect_request_connection_string == NULL) {
    return NO_MEMORY;
  }
  strncpy((*connect_request_packet)->plc4x_connect_request_connection_string,
          configuration->remote_connection_string, size);

  return OK;
}

void plc4c_driver_plc4x_destroy_connection_request(plc4c_plc4x_read_write_plc4x_message* packet) {
    plc4c_driver_plc4x_destroy_packet(packet);
}

plc4c_return_code plc4c_driver_plc4x_create_plc4x_read_request(plc4c_read_request* read_request,
                                                               plc4c_plc4x_read_write_plc4x_message** plc4x_read_request_packet) {
  plc4c_driver_plc4x_config* configuration = read_request->connection->configuration;

  *plc4x_read_request_packet = malloc(sizeof(plc4c_plc4x_read_write_plc4x_message));
  if (*plc4x_read_request_packet == NULL) {
    return NO_MEMORY;
  }

  (*plc4x_read_request_packet)->_type = plc4c_plc4x_read_write_plc4x_message_type_plc4c_plc4x_read_write_plc4x_connect_request;
  (*plc4x_read_request_packet)->request_id = configuration->request_id++;
  (*plc4x_read_request_packet)->plc4x_read_request_connection_id = configuration->connection_id;
  plc4c_utils_list_create(&((*plc4x_read_request_packet)->plc4x_read_request_fields));

  plc4c_list_element* cur_item = read_request->items->tail;
  while (cur_item != NULL) {
    plc4c_item* item;
    item = cur_item->value;

    // Create a new field item.
    plc4c_plc4x_read_write_plc4x_field* field;
    field = malloc(sizeof(plc4c_plc4x_read_write_plc4x_field));
    if (field == NULL) {
      return NO_MEMORY;
    }

    // Copy the name.
    size_t size = sizeof(char) * (strlen(item->name) + 1);
    field->name = malloc(size);
    if(field->name == NULL) {
      return NO_MEMORY;
    }
    strncpy(field->name,item->name, size);

    // Copy the query.
    size = sizeof(char) * (strlen(item->address) + 1);
    field->field_query = malloc(size);
    if(field->field_query == NULL) {
      return NO_MEMORY;
    }
    strncpy(field->field_query,item->address, size);

    // Add the item to the list.
    plc4c_utils_list_insert_tail_value(
        (*plc4x_read_request_packet)->plc4x_read_request_fields, field);

    cur_item = cur_item->next;
  }

  return OK;
}

void plc4c_driver_plc4x_destroy_plc4x_read_request(plc4c_plc4x_read_write_plc4x_message* packet) {
  plc4c_driver_plc4x_destroy_packet(packet);
}

plc4c_return_code plc4c_driver_plc4x_create_plc4x_write_request(plc4c_write_request* write_request,
                                                                plc4c_plc4x_read_write_plc4x_message** plc4x_write_request_packet) {
  plc4c_driver_plc4x_config* configuration = write_request->connection->configuration;

  *plc4x_write_request_packet = malloc(sizeof(plc4c_plc4x_read_write_plc4x_message));
  if (*plc4x_write_request_packet == NULL) {
    return NO_MEMORY;
  }

  (*plc4x_write_request_packet)->_type = plc4c_plc4x_read_write_plc4x_message_type_plc4c_plc4x_read_write_plc4x_connect_request;
  (*plc4x_write_request_packet)->request_id = configuration->request_id++;
  (*plc4x_write_request_packet)->plc4x_write_request_connection_id = configuration->connection_id;
  plc4c_utils_list_create(&((*plc4x_write_request_packet)->plc4x_write_request_fields));

  plc4c_list_element* cur_item = write_request->items->tail;
  while (cur_item != NULL) {
    plc4c_item* item;
    item = cur_item->value;

    // Create a new field item.
    plc4c_plc4x_read_write_plc4x_field* field;
    field = malloc(sizeof(plc4c_plc4x_read_write_plc4x_field));
    if (field == NULL) {
      return NO_MEMORY;
    }

    // Copy the name.
    size_t size = sizeof(char) * (strlen(item->name) + 1);
    field->name = malloc(size);
    if(field->name == NULL) {
      return NO_MEMORY;
    }
    strncpy(field->name,item->name, size);

    // Copy the query.
    size = sizeof(char) * (strlen(item->address) + 1);
    field->field_query = malloc(size);
    if(field->field_query == NULL) {
      return NO_MEMORY;
    }
    strncpy(field->field_query,item->address, size);

    plc4c_plc4x_read_write_plc4x_field_value_request* value_request;
    value_request = malloc(sizeof(plc4c_plc4x_read_write_plc4x_field_value_request));
    if(value_request == NULL) {
      return NO_MEMORY;
    }
    value_request->field = field;
    // TODO: Set the plc-value
    //value_request->

    // Add the item to the list.
    plc4c_utils_list_insert_tail_value(
        (*plc4x_write_request_packet)->plc4x_write_request_fields, value_request);

    cur_item = cur_item->next;
  }

  return OK;
}

void plc4c_driver_plc4x_destroy_plc4x_write_request(plc4c_plc4x_read_write_plc4x_message* packet) {
  plc4c_driver_plc4x_destroy_packet(packet);
}










