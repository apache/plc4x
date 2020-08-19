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

#include <cotp_protocol_class.h>
#include <plc4c/driver_s7.h>
#include <plc4c/plc4c.h>
#include <plc4c/spi/types_private.h>
#include <stdlib.h>
#include <string.h>
#include <tpkt_packet.h>

// State definitions
enum plc4c_driver_s7_connect_states {
  PLC4C_DRIVER_S7_CONNECT_INIT,
  PLC4C_DRIVER_S7_CONNECT_SEND_COTP_CONNECT_REQUEST,
  PLC4C_DRIVER_S7_CONNECT_RECEIVE_COTP_CONNECT_RESPONSE,
  PLC4C_DRIVER_S7_CONNECT_SEND_S7_CONNECT_REQUEST,
  PLC4C_DRIVER_S7_CONNECT_RECEIVE_S7_CONNECT_RESPONSE,
  PLC4C_DRIVER_S7_CONNECT_SEND_S7_IDENTIFICATION_REQUEST,
  PLC4C_DRIVER_S7_CONNECT_RECEIVE_S7_IDENTIFICATION_RESPONSE,
  PLC4C_DRIVER_S7_CONNECT_FINISHED
};

enum plc4c_driver_s7_disconnect_states {
  PLC4C_DRIVER_S7_DISCONNECT_INIT,
  PLC4C_DRIVER_S7_DISCONNECT_WAIT_TASKS_FINISHED,
  PLC4C_DRIVER_S7_DISCONNECT_FINISHED
};

enum plc4c_driver_s7_read_states {
  PLC4C_DRIVER_S7_READ_INIT,
  PLC4C_DRIVER_S7_READ_FINISHED
};

enum plc4c_driver_s7_write_states {
  PLC4C_DRIVER_S7_WRITE_INIT,
  PLC4C_DRIVER_S7_WRITE_FINISHED
};

int16_t plc4c_driver_s7_select_message_function(uint8_t* buffer_data,
                                                uint16_t buffer_length) {
  // If the packet doesn't start with 0x03, it's a corrupt package.
  if (buffer_length >= 1) {
    // The buffer seems to be corrupt, try to find a sequence of 0x03 0x00
    // and return the negative number of bytes needed to find that or the
    // number of bytes in the buffer so it will simply clean the buffer
    // completely.
    if (*buffer_data != 0x03) {
      for (int i = 1; i < (buffer_length - 1); i++) {
        buffer_data++;
        if ((*buffer_data == 0x03) && (*(buffer_data + 1) == 0x00)) {
          // We've found a potential new packet start.
          return -(i - 1);
        }
      }
      // We didn't find a new start, delete the entire content except the last
      // byte (as this could be the start of the next frame and we couldn't
      // confirm this.
      return -(buffer_length - 1);
    }
  }
  // The length information is located in bytes 3 and 4
  if (buffer_length >= 4) {
    uint16_t packet_length =
        (((uint16_t) * (buffer_data + (uint8_t)2)) << (uint16_t)8) |
        ((uint16_t) * (buffer_data + 3));
    if (buffer_length >= packet_length) {
      return packet_length;
    }
    // 8192 is the maximum pdu size, so if the value is larger, the packet is
    // probably corrupt.
    if (packet_length > 8192) {
      for (int i = 1; i < (buffer_length - 1); i++) {
        buffer_data++;
        if ((*buffer_data == 0x03) && (*(buffer_data + 1) == 0x00)) {
          // We've found a potential new packet start.
          return -(i - 1);
        }
      }
      return -(buffer_length - 1);
    }
  }
  // In all other cases, we'll just have to wait for the next time.
  return 0;
}

uint16_t plc4c_driver_s7_encode_tsap_id(
    plc4c_driver_s7_device_group device_group, uint8_t rack, uint8_t slot) {
  return (((uint16_t)device_group) << 8) | (((uint16_t)rack & 0x000F) << 4) |
         ((uint16_t)slot & 0x000F);
}

uint16_t plc4c_driver_s7_get_nearest_matching_tpdu_size(uint16_t pdu_size) {
  for(int i = 0; i < plc4c_s7_read_write_cotp_tpdu_size_num_values(); i++) {
    uint16_t cur_value = plc4c_s7_read_write_cotp_tpdu_size_get_size_in_bytes(i);
    if(cur_value >= pdu_size) {
      return cur_value;
    }
  }
  return 0;
}

plc4c_driver_s7_controller_type decode_controller_type(char* article_number) {
  char* prefix = "6ES7 ";
  // If this article-number doesn't start with this prefix, we can't decode it.
  if (strncmp(prefix, article_number, strlen(prefix)) != 0) {
    return PLC4C_DRIVER_S7_CONTROLLER_TYPE_ANY;
  }
  char model = *(article_number + 5);
  switch (model) {
    case 2:
      return PLC4C_DRIVER_S7_CONTROLLER_TYPE_S7_1200;
    case 5:
      return PLC4C_DRIVER_S7_CONTROLLER_TYPE_S7_1500;
    case 3:
      return PLC4C_DRIVER_S7_CONTROLLER_TYPE_S7_300;
    case 4:
      return PLC4C_DRIVER_S7_CONTROLLER_TYPE_S7_400;
    default:
      return PLC4C_DRIVER_S7_CONTROLLER_TYPE_ANY;
  }
}

plc4c_return_code send_packet(plc4c_connection* connection,
                              plc4c_s7_read_write_tpkt_packet* packet) {
  // Get the size required to contain the serialized form of this packet.
  uint16_t packet_size =
      plc4c_s7_read_write_tpkt_packet_length_in_bytes(packet);

  // Serialize this message to a byte-array.
  plc4c_spi_write_buffer* write_buffer;
  plc4c_return_code return_code =
      plc4c_spi_write_buffer_create(packet_size, &write_buffer);
  if (return_code != OK) {
    return return_code;
  }
  plc4c_s7_read_write_tpkt_packet_serialize(write_buffer, packet);

  // Now send this to the recipient.
  return_code = connection->transport->send_message(write_buffer);
  if (return_code != OK) {
    return return_code;
  }

  return OK;
}

plc4c_return_code receive_packet(plc4c_connection* connection,
                                 plc4c_s7_read_write_tpkt_packet** packet) {
  // Get a response from the transport.
  plc4c_spi_read_buffer* read_buffer;
  plc4c_return_code return_code = connection->transport->select_message(
      plc4c_driver_s7_select_message_function, &read_buffer);
  if (return_code != OK) {
    return return_code;
  }

  // Parse the given data.
  *packet = NULL;
  return_code = plc4c_s7_read_write_tpkt_packet_parse(
      read_buffer, packet);
  if (return_code != OK) {
    return return_code;
  }

  return OK;
}

// Declare the functions for which the definition will come later ...
plc4c_return_code createCOTPConnectionRequest(
    plc4c_driver_s7_config* configuration,
    plc4c_s7_read_write_tpkt_packet** cotp_connect_request_packet);
plc4c_return_code createS7ConnectionRequest(
    plc4c_driver_s7_config* configuration,
    plc4c_s7_read_write_tpkt_packet** s7_connect_request_packet);
plc4c_return_code createS7IdentifyRemoteRequest(
    plc4c_s7_read_write_tpkt_packet** s7_identify_remote_request_packet);

/**
 * State machine function for establishing a connection to a remote S7 device.
 * @param task the current system task
 * @return return code of the current state machine step execution
 */
plc4c_return_code plc4c_driver_s7_connect_machine_function(
    plc4c_system_task* task) {
  plc4c_connection* connection = task->context;
  if (connection == NULL) {
    return INTERNAL_ERROR;
  }
  // If we were already connected, return an error
  if (plc4c_connection_get_connected(connection)) {
    return ALREADY_CONNECTED;
  }
  plc4c_driver_s7_config* configuration = connection->configuration;

  switch (task->state_id) {
    // Initialize some internal data-structures.
    case PLC4C_DRIVER_S7_CONNECT_INIT: {
      // Calculate some internal settings from the values provided in the
      // configuration.
      configuration->calling_tsap_id = plc4c_driver_s7_encode_tsap_id(
          PLC4C_DRIVER_S7_DEVICE_GROUP_OTHERS, configuration->local_rack,
          configuration->local_slot);
      configuration->called_tsap_id = plc4c_driver_s7_encode_tsap_id(
          PLC4C_DRIVER_S7_DEVICE_GROUP_PG_OR_PC, configuration->remote_rack,
          configuration->remote_slot);
      configuration->cotp_tpdu_size =
          plc4c_driver_s7_get_nearest_matching_tpdu_size(
              configuration->pdu_size);
      configuration->pdu_size = configuration->cotp_tpdu_size - 16;

      task->state_id = PLC4C_DRIVER_S7_CONNECT_SEND_COTP_CONNECT_REQUEST;
      break;
    }
    // Send a COTP connection request.
    case PLC4C_DRIVER_S7_CONNECT_SEND_COTP_CONNECT_REQUEST: {
      // Get a COTP connection response for the settings in the config.
      plc4c_s7_read_write_tpkt_packet* cotp_connect_request_packet;
      plc4c_return_code return_code = createCOTPConnectionRequest(
          configuration, &cotp_connect_request_packet);
      if(return_code != OK) {
        return return_code;
      }

      // Send the packet to the remote.
      return_code = send_packet(
          connection, cotp_connect_request_packet);
      if(return_code != OK) {
        return return_code;
      }

      task->state_id = PLC4C_DRIVER_S7_CONNECT_RECEIVE_COTP_CONNECT_RESPONSE;
      break;
    }
    // Receive a COTP connection response.
    case PLC4C_DRIVER_S7_CONNECT_RECEIVE_COTP_CONNECT_RESPONSE: {
      // Read a response packet.
      plc4c_s7_read_write_tpkt_packet* cotp_connect_response_packet;
      plc4c_return_code return_code = receive_packet(connection, &cotp_connect_response_packet);
      // If we haven't read enough to process a full message, just try again
      // next time.
      if(return_code == UNFINISHED) {
        return OK;
      } else if(return_code != OK) {
        return return_code;
      }

      // Check if the packet has the right type
      if (cotp_connect_response_packet->payload->_type !=
          plc4c_s7_read_write_cotp_packet_type_plc4c_s7_read_write_cotp_packet_connection_response) {
        return INTERNAL_ERROR;
      }

      // Extract the information for: called-tsap-id, calling-tsap-id and
      // tpdu-size.
      plc4c_list_element* parameter_element =
          cotp_connect_response_packet->payload->parameters->tail;
      do {
        plc4c_s7_read_write_cotp_parameter* parameter =
            parameter_element->value;
        switch (parameter->_type) {
          case plc4c_s7_read_write_cotp_parameter_type_plc4c_s7_read_write_cotp_parameter_tpdu_size: {
            configuration->cotp_tpdu_size =
                parameter->cotp_parameter_tpdu_size_tpdu_size;
            break;
          }
          case plc4c_s7_read_write_cotp_parameter_type_plc4c_s7_read_write_cotp_parameter_calling_tsap: {
            configuration->calling_tsap_id =
                parameter->cotp_parameter_calling_tsap_tsap_id;
            break;
          }
          case plc4c_s7_read_write_cotp_parameter_type_plc4c_s7_read_write_cotp_parameter_called_tsap: {
            configuration->called_tsap_id =
                parameter->cotp_parameter_called_tsap_tsap_id;
            break;
          }
          default: {
            break;
          }
        }
      } while (parameter_element != NULL);

      // If we got the expected response, continue with the next higher level
      // of connection.
      task->state_id = PLC4C_DRIVER_S7_CONNECT_SEND_S7_CONNECT_REQUEST;
      break;
    }
    // Send a S7 connection request.
    case PLC4C_DRIVER_S7_CONNECT_SEND_S7_CONNECT_REQUEST: {
      plc4c_s7_read_write_tpkt_packet* s7_connect_request_packet;
      plc4c_return_code return_code = createS7ConnectionRequest(
          configuration, &s7_connect_request_packet);
      if(return_code != OK) {
        return return_code;
      }

      // Send the packet to the remote.
      return_code = send_packet(
          connection, s7_connect_request_packet);
      if(return_code != OK) {
        return return_code;
      }

      task->state_id = PLC4C_DRIVER_S7_CONNECT_RECEIVE_S7_CONNECT_RESPONSE;
      break;
    }
    // Receive a S7 connection response.
    case PLC4C_DRIVER_S7_CONNECT_RECEIVE_S7_CONNECT_RESPONSE: {
      // Read a response packet.
      plc4c_s7_read_write_tpkt_packet* s7_connect_response_packet;
      plc4c_return_code return_code = receive_packet(
          connection, &s7_connect_response_packet);
      // If we haven't read enough to process a full message, just try again
      // next time.
      if(return_code == UNFINISHED) {
        return OK;
      } else if(return_code != OK) {
        return return_code;
      }

      // Check if the packet has the right type
      if (s7_connect_response_packet->payload->_type !=
          plc4c_s7_read_write_cotp_packet_type_plc4c_s7_read_write_cotp_packet_data) {
        return INTERNAL_ERROR;
      }
      if (s7_connect_response_packet->payload->payload->_type !=
          plc4c_s7_read_write_s7_message_type_plc4c_s7_read_write_s7_message_response) {
        return INTERNAL_ERROR;
      }
      if (s7_connect_response_packet->payload->payload->parameter->_type !=
          plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_setup_communication) {
        return INTERNAL_ERROR;
      }

      // Extract and save the information for:
      // max-amq-caller, max-amq-callee, pdu-size.
      configuration->pdu_size =
          s7_connect_response_packet->payload->payload->parameter
              ->s7_parameter_setup_communication_pdu_length;
      configuration->max_amq_caller =
          s7_connect_response_packet->payload->payload->parameter
              ->s7_parameter_setup_communication_max_amq_caller;
      configuration->max_amq_callee =
          s7_connect_response_packet->payload->payload->parameter
              ->s7_parameter_setup_communication_max_amq_callee;

      // If no controller is explicitly selected, detect it.
      if (configuration->controller_type !=
          PLC4C_DRIVER_S7_CONTROLLER_TYPE_ANY) {
        task->state_id = PLC4C_DRIVER_S7_CONNECT_SEND_S7_IDENTIFICATION_REQUEST;
      }
      // If a controller is explicitly selected, we're done connecting.
      else {
        task->state_id = PLC4C_DRIVER_S7_CONNECT_FINISHED;
      }
      break;
    }
    // Send a S7 identification request.
    case PLC4C_DRIVER_S7_CONNECT_SEND_S7_IDENTIFICATION_REQUEST: {
      plc4c_s7_read_write_tpkt_packet* s7_identify_remote_request_packet;
      plc4c_return_code return_code = createS7IdentifyRemoteRequest(
          &s7_identify_remote_request_packet);
      if(return_code != OK) {
        return return_code;
      }

      // Send the packet to the remote.
      return_code = send_packet(
          connection, s7_identify_remote_request_packet);
      if(return_code != OK) {
        return return_code;
      }

      task->state_id =
          PLC4C_DRIVER_S7_CONNECT_RECEIVE_S7_IDENTIFICATION_RESPONSE;
      break;
    }
    // Receive a S7 identification response.
    case PLC4C_DRIVER_S7_CONNECT_RECEIVE_S7_IDENTIFICATION_RESPONSE: {
      // Read a response packet.
      plc4c_s7_read_write_tpkt_packet* s7_identify_remote_response_packet;
      plc4c_return_code return_code = receive_packet(
          connection, &s7_identify_remote_response_packet);
      // If we haven't read enough to process a full message, just try again
      // next time.
      if(return_code == UNFINISHED) {
        return OK;
      } else if(return_code != OK) {
        return return_code;
      }

      // Check if the packet has the right type
      if (s7_identify_remote_response_packet->payload->_type !=
          plc4c_s7_read_write_cotp_packet_type_plc4c_s7_read_write_cotp_packet_data) {
        return INTERNAL_ERROR;
      }
      if (s7_identify_remote_response_packet->payload->payload->_type !=
          plc4c_s7_read_write_s7_message_type_plc4c_s7_read_write_s7_message_user_data) {
        return INTERNAL_ERROR;
      }
      if (s7_identify_remote_response_packet->payload->payload->payload->_type !=
          plc4c_s7_read_write_s7_payload_type_plc4c_s7_read_write_s7_payload_user_data) {
        return INTERNAL_ERROR;
      }

      plc4c_list_element* cur_item =
          s7_identify_remote_response_packet->payload->payload->payload
              ->s7_payload_user_data_items->tail;
      while (cur_item != NULL) {
        plc4c_s7_read_write_s7_payload_user_data_item* item = cur_item->value;
        if (item->_type ==
            plc4c_s7_read_write_s7_payload_user_data_item_type_plc4c_s7_read_write_s7_payload_user_data_item_cpu_function_read_szl_response) {
          plc4c_list_element* szl_item =
              item->s7_payload_user_data_item_cpu_function_read_szl_response_items
                  ->head;
          while (szl_item != NULL) {
            plc4c_s7_read_write_szl_data_tree_item* data_tree_item =
                szl_item->value;
            if (data_tree_item->item_index == 0x0001) {
              char* article_number = list_to_string(data_tree_item->mlfb);
              if (article_number != NULL) {
                configuration->controller_type =
                    decode_controller_type(article_number);
                free(article_number);
              } else {
                configuration->controller_type =
                    PLC4C_DRIVER_S7_CONTROLLER_TYPE_ANY;
              }
            }
            szl_item = szl_item->next;
          }
        }

        cur_item = cur_item->next;
      }

      task->state_id = PLC4C_DRIVER_S7_CONNECT_FINISHED;
      break;
    }
    // Clean up some internal data-structures.
    case PLC4C_DRIVER_S7_CONNECT_FINISHED: {
      plc4c_connection_set_connected(connection, true);
      task->completed = true;
      break;
    }
    // If an unexpected state id was received, this is not really something
    // we can recover from.
    default: {
      task->completed = true;
      return INTERNAL_ERROR;
    }
  }
  return OK;
}

plc4c_return_code plc4c_driver_s7_disconnect_machine_function(
    plc4c_system_task* task) {
  plc4c_connection* connection = task->context;
  if (connection == NULL) {
    return INTERNAL_ERROR;
  }

  switch (task->state_id) {
    case PLC4C_DRIVER_S7_DISCONNECT_INIT: {
      plc4c_connection_set_disconnect(connection, true);
      task->state_id = PLC4C_DRIVER_S7_DISCONNECT_WAIT_TASKS_FINISHED;
      break;
    }
    case PLC4C_DRIVER_S7_DISCONNECT_WAIT_TASKS_FINISHED: {
      // The disconnect system-task also counts.
      if (plc4c_connection_get_running_tasks_count(connection) == 1) {
        plc4c_connection_set_connected(connection, false);
        task->completed = true;
        task->state_id = PLC4C_DRIVER_S7_DISCONNECT_FINISHED;
      }
      break;
    }
    case PLC4C_DRIVER_S7_DISCONNECT_FINISHED: {
      // Do nothing
      break;
    }
    default: {
      return INTERNAL_ERROR;
    }
  }
  return OK;
}

plc4c_return_code plc4c_driver_s7_read_machine_function(
    plc4c_system_task* task) {
  plc4c_read_request_execution* read_request_execution = task->context;
  if (read_request_execution == NULL) {
    return INTERNAL_ERROR;
  }
  if (task->connection == NULL) {
    return INTERNAL_ERROR;
  }

  switch (task->state_id) {
    case PLC4C_DRIVER_S7_READ_INIT: {
      task->completed = true;
      // TODO: Implement this ...
      task->state_id = PLC4C_DRIVER_S7_READ_FINISHED;
      break;
    }
    case PLC4C_DRIVER_S7_READ_FINISHED: {
      break;
    }
  }
  return OK;
}

plc4c_return_code plc4c_driver_s7_write_machine_function(
    plc4c_system_task* task) {
  plc4c_write_request_execution* write_request_execution = task->context;
  if (write_request_execution == NULL) {
    return INTERNAL_ERROR;
  }
  if (task->connection == NULL) {
    return INTERNAL_ERROR;
  }

  switch (task->state_id) {
    case PLC4C_DRIVER_S7_WRITE_INIT: {
      task->completed = true;
      // TODO: Implement this ...
      task->state_id = PLC4C_DRIVER_S7_WRITE_FINISHED;
      break;
    }
    case PLC4C_DRIVER_S7_WRITE_FINISHED: {
      break;
    }
  }
  return OK;
}

plc4c_return_code plc4c_driver_s7_connect_function(plc4c_connection* connection,
                                                   plc4c_system_task** task) {
  plc4c_system_task* new_task = malloc(sizeof(plc4c_system_task));
  // There's nothing to do here, so no need for a state-machine.
  new_task->state_id = PLC4C_DRIVER_S7_CONNECT_INIT;
  new_task->state_machine_function = &plc4c_driver_s7_connect_machine_function;
  new_task->completed = false;
  new_task->context = connection;
  new_task->connection = connection;
  *task = new_task;
  return OK;
}

plc4c_return_code plc4c_driver_s7_disconnect_function(
    plc4c_connection* connection, plc4c_system_task** task) {
  plc4c_system_task* new_task = malloc(sizeof(plc4c_system_task));
  new_task->state_id = PLC4C_DRIVER_S7_DISCONNECT_INIT;
  new_task->state_machine_function =
      &plc4c_driver_s7_disconnect_machine_function;
  new_task->completed = false;
  new_task->context = connection;
  new_task->connection = connection;
  *task = new_task;
  return OK;
}

plc4c_return_code plc4c_driver_s7_read_function(
    plc4c_read_request_execution* read_request_execution,
    plc4c_system_task** task) {
  plc4c_system_task* new_task = malloc(sizeof(plc4c_system_task));
  new_task->state_id = PLC4C_DRIVER_S7_READ_INIT;
  new_task->state_machine_function = &plc4c_driver_s7_read_machine_function;
  new_task->completed = false;
  new_task->context = read_request_execution;
  new_task->connection = read_request_execution->system_task->connection;
  *task = new_task;
  return OK;
}

void plc4c_driver_s7_free_read_response_item(
    plc4c_list_element* read_item_element) {
  plc4c_response_value_item* value_item =
      (plc4c_response_value_item*)read_item_element->value;
  plc4c_data_destroy(value_item->value);
  value_item->value = NULL;
}

void plc4c_driver_s7_free_read_response(plc4c_read_response* response) {
  // the request will be cleaned up elsewhere
  plc4c_utils_list_delete_elements(response->items,
                                   &plc4c_driver_s7_free_read_response_item);
}

plc4c_return_code plc4c_driver_s7_write_function(
    plc4c_write_request_execution* write_request_execution,
    plc4c_system_task** task) {
  plc4c_system_task* new_task = malloc(sizeof(plc4c_system_task));
  new_task->state_id = PLC4C_DRIVER_S7_WRITE_INIT;
  new_task->state_machine_function = &plc4c_driver_s7_write_machine_function;
  new_task->completed = false;
  new_task->context = write_request_execution;
  new_task->connection = write_request_execution->system_task->connection;
  *task = new_task;
  return OK;
}

void plc4c_driver_s7_free_write_response_item(
    plc4c_list_element* write_item_element) {
  plc4c_response_value_item* value_item =
      (plc4c_response_value_item*)write_item_element->value;
  // do not delete the plc4c_item
  // we also, in THIS case don't delete the random value which isn't really
  // a pointer
  // free(value_item->value);
  value_item->value = NULL;
}

void plc4c_driver_s7_free_write_response(plc4c_write_response* response) {
  // the request will be cleaned up elsewhere
  plc4c_utils_list_delete_elements(response->response_items,
                                   &plc4c_driver_s7_free_write_response_item);
}

plc4c_driver* plc4c_driver_s7_create() {
  plc4c_driver* driver = (plc4c_driver*)malloc(sizeof(plc4c_driver));
  driver->protocol_code = "s7";
  driver->protocol_name = "Siemens S7 (Basic)";
  driver->default_transport_code = "tcp";
  driver->parse_address_function = NULL;
  driver->connect_function = &plc4c_driver_s7_connect_function;
  driver->disconnect_function = &plc4c_driver_s7_disconnect_function;
  driver->read_function = &plc4c_driver_s7_read_function;
  driver->write_function = &plc4c_driver_s7_write_function;
  driver->subscribe_function = NULL;
  driver->unsubscribe_function = NULL;
  driver->free_read_response_function = &plc4c_driver_s7_free_read_response;
  driver->free_write_response_function = &plc4c_driver_s7_free_write_response;
  driver->free_subscription_response_function = NULL;
  driver->free_unsubscription_response_function = NULL;
  return driver;
}

////////////////////////////////////////////////////////////////////////////////
// Helpers for preparing the different packets
////////////////////////////////////////////////////////////////////////////////

plc4c_return_code createCOTPConnectionRequest(
    plc4c_driver_s7_config* configuration,
    plc4c_s7_read_write_tpkt_packet** cotp_connect_request_packet) {
  *cotp_connect_request_packet =
      malloc(sizeof(plc4c_s7_read_write_tpkt_packet));
  if (*cotp_connect_request_packet == NULL) {
    return NO_MEMORY;
  }
  (*cotp_connect_request_packet)->payload =
      malloc(sizeof(plc4c_s7_read_write_cotp_packet));
  if ((*cotp_connect_request_packet)->payload == NULL) {
    return NO_MEMORY;
  }
  (*cotp_connect_request_packet)->payload->_type =
      plc4c_s7_read_write_cotp_packet_type_plc4c_s7_read_write_cotp_packet_connection_request;
  (*cotp_connect_request_packet)->payload
      ->cotp_packet_connection_request_destination_reference = 0x0000;
  (*cotp_connect_request_packet)->payload
      ->cotp_packet_connection_request_source_reference = 0x000F;
  (*cotp_connect_request_packet)->payload
      ->cotp_packet_connection_request_protocol_class =
      plc4c_s7_read_write_cotp_protocol_class_CLASS_0;

  // Add the COTP parameters: Called TSAP, Calling TSAP and TPDU Size.
  plc4c_utils_list_create(&((*cotp_connect_request_packet)->payload->parameters));
  plc4c_s7_read_write_cotp_parameter* called_tsap_parameter =
      malloc(sizeof(plc4c_s7_read_write_cotp_parameter));
  if (called_tsap_parameter == NULL) {
    return NO_MEMORY;
  }
  called_tsap_parameter->_type =
      plc4c_s7_read_write_cotp_parameter_type_plc4c_s7_read_write_cotp_parameter_called_tsap;
  called_tsap_parameter->cotp_parameter_called_tsap_tsap_id =
      configuration->called_tsap_id;

  plc4c_utils_list_insert_head_value(
      (*cotp_connect_request_packet)->payload->parameters, called_tsap_parameter);
  plc4c_s7_read_write_cotp_parameter* calling_tsap_parameter =
      malloc(sizeof(plc4c_s7_read_write_cotp_parameter));
  if (calling_tsap_parameter == NULL) {
    return NO_MEMORY;
  }
  calling_tsap_parameter->_type =
      plc4c_s7_read_write_cotp_parameter_type_plc4c_s7_read_write_cotp_parameter_calling_tsap;
  calling_tsap_parameter->cotp_parameter_calling_tsap_tsap_id =
      configuration->calling_tsap_id;

  plc4c_utils_list_insert_head_value(
      (*cotp_connect_request_packet)->payload->parameters, calling_tsap_parameter);
  plc4c_s7_read_write_cotp_parameter* tpdu_size_parameter =
      malloc(sizeof(plc4c_s7_read_write_cotp_parameter));
  if (tpdu_size_parameter == NULL) {
    return NO_MEMORY;
  }
  tpdu_size_parameter->_type =
      plc4c_s7_read_write_cotp_parameter_type_plc4c_s7_read_write_cotp_parameter_tpdu_size;
  tpdu_size_parameter->cotp_parameter_tpdu_size_tpdu_size =
      configuration->cotp_tpdu_size;

  plc4c_utils_list_insert_head_value(
      (*cotp_connect_request_packet)->payload->parameters, tpdu_size_parameter);

  // For a COTP connection request, there is no payload.
  (*cotp_connect_request_packet)->payload->payload = NULL;

  return OK;
}

plc4c_return_code createS7ConnectionRequest(
    plc4c_driver_s7_config* configuration,
    plc4c_s7_read_write_tpkt_packet** s7_connect_request_packet) {
  *s7_connect_request_packet = malloc(sizeof(plc4c_s7_read_write_tpkt_packet));
  if (*s7_connect_request_packet == NULL) {
    return NO_MEMORY;
  }
  (*s7_connect_request_packet)->payload =
      malloc(sizeof(plc4c_s7_read_write_cotp_packet));
  if ((*s7_connect_request_packet)->payload == NULL) {
    return NO_MEMORY;
  }
  (*s7_connect_request_packet)->payload->_type =
      plc4c_s7_read_write_cotp_packet_type_plc4c_s7_read_write_cotp_packet_data;
  (*s7_connect_request_packet)->payload->parameters = NULL;
  (*s7_connect_request_packet)->payload->cotp_packet_data_eot = true;
  (*s7_connect_request_packet)->payload->cotp_packet_data_tpdu_ref = 1;

  (*s7_connect_request_packet)->payload->payload =
      malloc(sizeof(plc4c_s7_read_write_s7_message));
  if ((*s7_connect_request_packet)->payload->payload == NULL) {
    return NO_MEMORY;
  }
  (*s7_connect_request_packet)->payload->payload->_type =
      plc4c_s7_read_write_s7_message_type_plc4c_s7_read_write_s7_message_request;
  (*s7_connect_request_packet)->payload->payload->parameter =
      malloc(sizeof(plc4c_s7_read_write_s7_parameter));
  if ((*s7_connect_request_packet)->payload->payload->parameter == NULL) {
    return NO_MEMORY;
  }
  (*s7_connect_request_packet)->payload->payload->parameter->_type =
      plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_setup_communication;
  (*s7_connect_request_packet)->payload->payload->parameter
      ->s7_parameter_setup_communication_max_amq_callee =
      configuration->max_amq_callee;
  (*s7_connect_request_packet)->payload->payload->parameter
      ->s7_parameter_setup_communication_max_amq_caller =
      configuration->max_amq_caller;
  (*s7_connect_request_packet)->payload->payload->parameter
      ->s7_parameter_setup_communication_pdu_length =
      configuration->pdu_size;

  return OK;
}

plc4c_return_code createS7IdentifyRemoteRequest(
    plc4c_s7_read_write_tpkt_packet** s7_identify_remote_request_packet) {

  *s7_identify_remote_request_packet =
      malloc(sizeof(plc4c_s7_read_write_tpkt_packet));
  if (*s7_identify_remote_request_packet == NULL) {
    return NO_MEMORY;
  }
  (*s7_identify_remote_request_packet)->payload =
      malloc(sizeof(plc4c_s7_read_write_cotp_packet));
  if ((*s7_identify_remote_request_packet)->payload == NULL) {
    return NO_MEMORY;
  }
  (*s7_identify_remote_request_packet)->payload->_type =
      plc4c_s7_read_write_cotp_packet_type_plc4c_s7_read_write_cotp_packet_data;
  (*s7_identify_remote_request_packet)->payload->parameters = NULL;
  (*s7_identify_remote_request_packet)->payload->cotp_packet_data_eot = true;
  (*s7_identify_remote_request_packet)->payload->cotp_packet_data_tpdu_ref = 2;

  (*s7_identify_remote_request_packet)->payload->payload =
      malloc(sizeof(plc4c_s7_read_write_s7_message));
  if ((*s7_identify_remote_request_packet)->payload->payload == NULL) {
    return NO_MEMORY;
  }
  (*s7_identify_remote_request_packet)->payload->payload->_type =
      plc4c_s7_read_write_s7_message_type_plc4c_s7_read_write_s7_message_user_data;
  (*s7_identify_remote_request_packet)->payload->payload->parameter =
      malloc(sizeof(plc4c_s7_read_write_s7_parameter));
  if ((*s7_identify_remote_request_packet)->payload->payload->parameter ==
      NULL) {
    return NO_MEMORY;
  }
  (*s7_identify_remote_request_packet)->payload->payload->parameter->_type =
      plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_user_data;

  plc4c_utils_list_create(
      &((*s7_identify_remote_request_packet)->payload->payload->parameter
          ->s7_parameter_user_data_items));
  plc4c_s7_read_write_s7_parameter_user_data_item* parameter_item =
      malloc(sizeof(plc4c_s7_read_write_s7_parameter_user_data_item));
  parameter_item->_type =
      plc4c_s7_read_write_s7_parameter_user_data_item_type_plc4c_s7_read_write_s7_parameter_user_data_item_cpu_functions;
  parameter_item->s7_parameter_user_data_item_cpu_functions_method = 0x11;
  parameter_item
      ->s7_parameter_user_data_item_cpu_functions_cpu_function_type = 0x4;
  parameter_item
      ->s7_parameter_user_data_item_cpu_functions_cpu_function_group = 0x4;
  parameter_item
      ->s7_parameter_user_data_item_cpu_functions_cpu_subfunction = 0x01;
  parameter_item
      ->s7_parameter_user_data_item_cpu_functions_sequence_number = 0x00;
  parameter_item
      ->s7_parameter_user_data_item_cpu_functions_data_unit_reference_number =
      NULL;
  parameter_item->s7_parameter_user_data_item_cpu_functions_last_data_unit =
      NULL;
  parameter_item->s7_parameter_user_data_item_cpu_functions_error_code =
      NULL;
  plc4c_utils_list_insert_head_value(
      (*s7_identify_remote_request_packet)->payload->payload->parameter
          ->s7_parameter_user_data_items,
      parameter_item);

  plc4c_utils_list_create(
      &((*s7_identify_remote_request_packet)->payload->payload->payload
          ->s7_payload_user_data_items));
  plc4c_s7_read_write_s7_payload_user_data_item* payload_item =
      malloc(sizeof(plc4c_s7_read_write_s7_payload_user_data_item));
  payload_item->_type =
      plc4c_s7_read_write_s7_payload_user_data_item_type_plc4c_s7_read_write_s7_payload_user_data_item_cpu_function_read_szl_request;
  payload_item->return_code =
      plc4c_s7_read_write_data_transport_error_code_OK;
  payload_item->transport_size =
      plc4c_s7_read_write_data_transport_size_OCTET_STRING;
  payload_item->szl_index = 0x0000;
  payload_item->szl_id = malloc(sizeof(plc4c_s7_read_write_szl_id));
  if (payload_item->szl_id == NULL) {
    return NO_MEMORY;
  }
  payload_item->szl_id->type_class =
      plc4c_s7_read_write_szl_module_type_class_CPU;
  payload_item->szl_id->sublist_extract = 0x00;
  payload_item->szl_id->sublist_list =
      plc4c_s7_read_write_szl_sublist_MODULE_IDENTIFICATION;

  return OK;
}
