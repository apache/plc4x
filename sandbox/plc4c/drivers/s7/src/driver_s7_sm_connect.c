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

#include <ctype.h>
#include <plc4c/plc4c.h>
#include <plc4c/spi/types_private.h>
#include <stdlib.h>
#include <string.h>
#include "plc4c/driver_s7.h"
#include "plc4c/driver_s7_encode_decode.h"
#include "plc4c/driver_s7_packets.h"
#include "cotp_protocol_class.h"
#include "tpkt_packet.h"

enum plc4c_driver_s7_connect_states {
  PLC4C_DRIVER_S7_CONNECT_INIT,
  PLC4C_DRIVER_S7_CONNECT_TRANSPORT_CONNECT,
  PLC4C_DRIVER_S7_CONNECT_SEND_COTP_CONNECT_REQUEST,
  PLC4C_DRIVER_S7_CONNECT_RECEIVE_COTP_CONNECT_RESPONSE,
  PLC4C_DRIVER_S7_CONNECT_SEND_S7_CONNECT_REQUEST,
  PLC4C_DRIVER_S7_CONNECT_RECEIVE_S7_CONNECT_RESPONSE,
  PLC4C_DRIVER_S7_CONNECT_SEND_S7_IDENTIFICATION_REQUEST,
  PLC4C_DRIVER_S7_CONNECT_RECEIVE_S7_IDENTIFICATION_RESPONSE,
  PLC4C_DRIVER_S7_CONNECT_FINISHED
};

/**
 * State machine function for establishing a connection to a remote S7 device.
 * @param task the current system task
 * @return return code of the current state machine step execution
 */
plc4c_return_code plc4c_driver_s7_connect_machine_function(plc4c_system_task* task) {

  plc4c_connection* connection;
  plc4c_driver_s7_config* configuration;
  plc4c_return_code return_code;
  plc4c_s7_read_write_tpkt_packet* packet;

  connection = task->connection;
  if (connection == NULL) {
    return INTERNAL_ERROR;
  }
  // If we were already connected, return an error
  if (plc4c_connection_get_connected(connection)) {
    return ALREADY_CONNECTED;
  }
  configuration = connection->configuration;

  // Initialize the pdu id (The first messages are hard-coded)
  configuration->pdu_id = 4;

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

      task->state_id = PLC4C_DRIVER_S7_CONNECT_TRANSPORT_CONNECT;
      break;
    }
    case PLC4C_DRIVER_S7_CONNECT_TRANSPORT_CONNECT: {
      return_code = connection->transport->open(connection->transport_configuration);
      if(return_code != OK) {
        return return_code;
      }

      task->state_id = PLC4C_DRIVER_S7_CONNECT_SEND_COTP_CONNECT_REQUEST;
      break;
    }
    // Send a COTP connection request.
    case PLC4C_DRIVER_S7_CONNECT_SEND_COTP_CONNECT_REQUEST: {
      // Get a COTP connection response for the settings in the config.
      return_code = plc4c_driver_s7_create_cotp_connection_request(configuration, &packet);
      if (return_code != OK) {
        return return_code;
      }

      // Send the packet to the remote.
      return_code = plc4c_driver_s7_send_packet(connection, packet);
      if (return_code != OK) {
        return return_code;
      }

      task->state_id = PLC4C_DRIVER_S7_CONNECT_RECEIVE_COTP_CONNECT_RESPONSE;
      break;
    }
    // Receive a COTP connection response.
    case PLC4C_DRIVER_S7_CONNECT_RECEIVE_COTP_CONNECT_RESPONSE: {
      // Read a response packet.
      return_code = plc4c_driver_s7_receive_packet(connection, &packet);
      // If we haven't read enough to process a full message, just try again
      // next time.
      if (return_code == UNFINISHED) {
        return OK;
      } else if (return_code != OK) {
        return return_code;
      }

      // Check if the packet has the right type
      if (packet->payload->_type !=
          plc4c_s7_read_write_cotp_packet_type_plc4c_s7_read_write_cotp_packet_connection_response) {
        return INTERNAL_ERROR;
      }

      // Extract the information for: called-tsap-id, calling-tsap-id and
      // tpdu-size.
      plc4c_list_element* parameter_element;
      plc4c_s7_read_write_cotp_parameter* parameter;

      parameter_element = packet->payload->parameters->tail;
      while (parameter_element != NULL) {
        parameter = parameter_element->value;
        switch (parameter->_type) {
          case plc4c_s7_read_write_cotp_parameter_type_plc4c_s7_read_write_cotp_parameter_tpdu_size: {
            configuration->cotp_tpdu_size = parameter->cotp_parameter_tpdu_size_tpdu_size;
            break;
          }
          case plc4c_s7_read_write_cotp_parameter_type_plc4c_s7_read_write_cotp_parameter_calling_tsap: {
            configuration->calling_tsap_id = parameter->cotp_parameter_calling_tsap_tsap_id;
            break;
          }
          case plc4c_s7_read_write_cotp_parameter_type_plc4c_s7_read_write_cotp_parameter_called_tsap: {
            configuration->called_tsap_id = parameter->cotp_parameter_called_tsap_tsap_id;
            break;
          }
          default: {
            break;
          }
        }
        parameter_element = parameter_element->next;
      }

      // If we got the expected response, continue with the next higher level
      // of connection.
      task->state_id = PLC4C_DRIVER_S7_CONNECT_SEND_S7_CONNECT_REQUEST;
      break;
    }
    // Send a S7 connection request.
    case PLC4C_DRIVER_S7_CONNECT_SEND_S7_CONNECT_REQUEST: {
      return_code = plc4c_driver_s7_create_s7_connection_request(configuration, &packet);
      if (return_code != OK) {
        return return_code;
      }

      // Send the packet to the remote.
      return_code = plc4c_driver_s7_send_packet(connection, packet);
      if (return_code != OK) {
        return return_code;
      }

      task->state_id = PLC4C_DRIVER_S7_CONNECT_RECEIVE_S7_CONNECT_RESPONSE;
      break;
    }
    // Receive a S7 connection response.
    case PLC4C_DRIVER_S7_CONNECT_RECEIVE_S7_CONNECT_RESPONSE: {
      // Read a response packet.
      
      plc4c_s7_read_write_s7_parameter* s7_parameter;
      
      return_code = plc4c_driver_s7_receive_packet(connection, &packet);
      // If we haven't read enough to process a full message, just try again
      // next time.
      if (return_code == UNFINISHED) {
        return OK;
      } else if (return_code != OK) {
        return return_code;
      }

      // Check if the packet has the right type
      if (packet->payload->_type !=
          plc4c_s7_read_write_cotp_packet_type_plc4c_s7_read_write_cotp_packet_data) {
        return INTERNAL_ERROR;
      }
      if (packet->payload->payload->_type !=
          plc4c_s7_read_write_s7_message_type_plc4c_s7_read_write_s7_message_response_data) {
        return INTERNAL_ERROR;
      }
      if (packet->payload->payload->parameter->_type !=
          plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_setup_communication) {
        return INTERNAL_ERROR;
      }

      // Extract and save the information for:
      // max-amq-caller, max-amq-callee, pdu-size.
      s7_parameter  = packet->payload->payload->parameter;
      
      configuration->pdu_size = s7_parameter->s7_parameter_setup_communication_pdu_length;
      configuration->max_amq_caller = s7_parameter->s7_parameter_setup_communication_max_amq_caller;
      configuration->max_amq_callee = s7_parameter->s7_parameter_setup_communication_max_amq_callee;

      // If no controller is explicitly selected, detect it.
      if (configuration->controller_type == PLC4C_DRIVER_S7_CONTROLLER_TYPE_ANY) {
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
      return_code = plc4c_driver_s7_create_s7_identify_remote_request(
              &packet);
      if (return_code != OK) {
        return return_code;
      }

      // Send the packet to the remote.
      return_code = plc4c_driver_s7_send_packet(
          connection, packet);
      if (return_code != OK) {
        return return_code;
      }

      task->state_id = PLC4C_DRIVER_S7_CONNECT_RECEIVE_S7_IDENTIFICATION_RESPONSE;
      break;
    }
    // Receive a S7 identification response.
    case PLC4C_DRIVER_S7_CONNECT_RECEIVE_S7_IDENTIFICATION_RESPONSE: {
      // Read a response packet.
      return_code = plc4c_driver_s7_receive_packet(connection, &packet);
      // If we haven't read enough to process a full message, just try again
      // next time.
      if (return_code == UNFINISHED) {
        return OK;
      } else if (return_code != OK) {
        return return_code;
      }

      // Check if the packet has the right type
      if (packet->payload->_type !=
          plc4c_s7_read_write_cotp_packet_type_plc4c_s7_read_write_cotp_packet_data) {
        return INTERNAL_ERROR;
      }
      if (packet->payload->payload->_type !=
          plc4c_s7_read_write_s7_message_type_plc4c_s7_read_write_s7_message_user_data) {
        return INTERNAL_ERROR;
      }
      if (packet->payload->payload->payload
              ->_type !=
          plc4c_s7_read_write_s7_payload_type_plc4c_s7_read_write_s7_payload_user_data) {
        return INTERNAL_ERROR;
      }

      plc4c_list_element* cur_item =
          packet->payload->payload->payload
              ->s7_payload_user_data_items->tail;
      while (cur_item != NULL) {
        plc4c_s7_read_write_s7_payload_user_data_item* item = cur_item->value;
        if (item->_type ==
            plc4c_s7_read_write_s7_payload_user_data_item_type_plc4c_s7_read_write_s7_payload_user_data_item_cpu_function_read_szl_response) {
          plc4c_list_element* szl_item =
              item->s7_payload_user_data_item_cpu_function_read_szl_response_items
                  ->tail;
          while (szl_item != NULL) {
            plc4c_s7_read_write_szl_data_tree_item* data_tree_item =
                szl_item->value;
            if (data_tree_item->item_index == 0x0001) {
              char* article_number = list_to_string(data_tree_item->mlfb);
              if (article_number != NULL) {
                configuration->controller_type =
                    decode_controller_type(article_number);
                free(article_number);
                break;
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
