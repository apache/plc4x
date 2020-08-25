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
  specific language governing permiSchreib' missions and limitations
  under the License.
*/

#include <cotp_protocol_class.h>
#include <ctype.h>
#include <plc4c/driver_s7.h>
#include <plc4c/spi/types_private.h>
#include <stdlib.h>
#include <string.h>
#include <tpkt_packet.h>

#include "plc4c/driver_s7_encode_decode.h"

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
        ((uint16_t) (buffer_data + 2) << 8) |
        ((uint16_t) (buffer_data + 3));
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
  return_code = plc4c_s7_read_write_tpkt_packet_parse(read_buffer, packet);
  if (return_code != OK) {
    return return_code;
  }

  return OK;
}

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
  (*cotp_connect_request_packet)
      ->payload->cotp_packet_connection_request_destination_reference = 0x0000;
  (*cotp_connect_request_packet)
      ->payload->cotp_packet_connection_request_source_reference = 0x000F;
  (*cotp_connect_request_packet)
      ->payload->cotp_packet_connection_request_protocol_class =
      plc4c_s7_read_write_cotp_protocol_class_CLASS_0;

  // Add the COTP parameters: Called TSAP, Calling TSAP and TPDU Size.
  plc4c_utils_list_create(
      &((*cotp_connect_request_packet)->payload->parameters));
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
      (*cotp_connect_request_packet)->payload->parameters,
      called_tsap_parameter);
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
      (*cotp_connect_request_packet)->payload->parameters,
      calling_tsap_parameter);
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
  (*s7_connect_request_packet)
      ->payload->payload->parameter
      ->s7_parameter_setup_communication_max_amq_callee =
      configuration->max_amq_callee;
  (*s7_connect_request_packet)
      ->payload->payload->parameter
      ->s7_parameter_setup_communication_max_amq_caller =
      configuration->max_amq_caller;
  (*s7_connect_request_packet)
      ->payload->payload->parameter
      ->s7_parameter_setup_communication_pdu_length = configuration->pdu_size;

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
      &((*s7_identify_remote_request_packet)
            ->payload->payload->parameter->s7_parameter_user_data_items));
  plc4c_s7_read_write_s7_parameter_user_data_item* parameter_item =
      malloc(sizeof(plc4c_s7_read_write_s7_parameter_user_data_item));
  parameter_item->_type =
      plc4c_s7_read_write_s7_parameter_user_data_item_type_plc4c_s7_read_write_s7_parameter_user_data_item_cpu_functions;
  parameter_item->s7_parameter_user_data_item_cpu_functions_method = 0x11;
  parameter_item->s7_parameter_user_data_item_cpu_functions_cpu_function_type =
      0x4;
  parameter_item->s7_parameter_user_data_item_cpu_functions_cpu_function_group =
      0x4;
  parameter_item->s7_parameter_user_data_item_cpu_functions_cpu_subfunction =
      0x01;
  parameter_item->s7_parameter_user_data_item_cpu_functions_sequence_number =
      0x00;
  parameter_item
      ->s7_parameter_user_data_item_cpu_functions_data_unit_reference_number =
      NULL;
  parameter_item->s7_parameter_user_data_item_cpu_functions_last_data_unit =
      NULL;
  parameter_item->s7_parameter_user_data_item_cpu_functions_error_code = NULL;
  plc4c_utils_list_insert_head_value(
      (*s7_identify_remote_request_packet)
          ->payload->payload->parameter->s7_parameter_user_data_items,
      parameter_item);

  plc4c_utils_list_create(
      &((*s7_identify_remote_request_packet)
            ->payload->payload->payload->s7_payload_user_data_items));
  plc4c_s7_read_write_s7_payload_user_data_item* payload_item =
      malloc(sizeof(plc4c_s7_read_write_s7_payload_user_data_item));
  payload_item->_type =
      plc4c_s7_read_write_s7_payload_user_data_item_type_plc4c_s7_read_write_s7_payload_user_data_item_cpu_function_read_szl_request;
  payload_item->return_code = plc4c_s7_read_write_data_transport_error_code_OK;
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

plc4c_return_code createS7ReadRequest(
    plc4c_read_request* read_request,
    plc4c_s7_read_write_tpkt_packet** s7_read_request_packet) {
  plc4c_driver_s7_config* configuration =
      read_request->connection->configuration;

  *s7_read_request_packet = malloc(sizeof(s7_read_request_packet));
  if (*s7_read_request_packet == NULL) {
    return NO_MEMORY;
  }

  (*s7_read_request_packet)->payload =
      malloc(sizeof(plc4c_s7_read_write_cotp_packet));
  (*s7_read_request_packet)->payload->_type =
      plc4c_s7_read_write_cotp_packet_type_plc4c_s7_read_write_cotp_packet_data;
  (*s7_read_request_packet)->payload->cotp_packet_data_tpdu_ref =
      configuration->pdu_id++;
  (*s7_read_request_packet)->payload->cotp_packet_data_eot = true;
  (*s7_read_request_packet)->payload->payload =
      malloc(sizeof(plc4c_s7_read_write_s7_message));
  (*s7_read_request_packet)->payload->payload->_type =
      plc4c_s7_read_write_s7_message_type_plc4c_s7_read_write_s7_message_request;
  (*s7_read_request_packet)->payload->payload->parameter =
      malloc(sizeof(plc4c_s7_read_write_s7_parameter));
  (*s7_read_request_packet)->payload->payload->parameter->_type =
      plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_read_var_request;
  plc4c_utils_list_create(
      &(*s7_read_request_packet)
           ->payload->payload->parameter->s7_parameter_read_var_request_items);

  plc4c_list_element* item = read_request->items->tail;
  while (item != NULL) {
    // Get the item address from the API request.
    char* itemAddress = item->value;

    // Create the item ...
    plc4c_s7_read_write_s7_var_request_parameter_item* request_item;
    plc4c_return_code return_code = parseAddress(itemAddress, &request_item);
    if (return_code != OK) {
      return return_code;
    }

    // Add the new item to the request.
    plc4c_utils_list_insert_head_value(
        (*s7_read_request_packet)
            ->payload->payload->parameter->s7_parameter_read_var_request_items,
        request_item);

    item = item->next;
  }

  return OK;
}

plc4c_return_code createS7WriteRequest(
    plc4c_write_request* write_request,
    plc4c_s7_read_write_tpkt_packet** s7_write_request_packet) {
  plc4c_driver_s7_config* configuration =
      write_request->connection->configuration;

  *s7_write_request_packet = malloc(sizeof(s7_write_request_packet));
  if (*s7_write_request_packet == NULL) {
    return NO_MEMORY;
  }

  (*s7_write_request_packet)->payload =
      malloc(sizeof(plc4c_s7_read_write_cotp_packet));
  (*s7_write_request_packet)->payload->_type =
      plc4c_s7_read_write_cotp_packet_type_plc4c_s7_read_write_cotp_packet_data;
  (*s7_write_request_packet)->payload->cotp_packet_data_tpdu_ref =
      configuration->pdu_id++;
  (*s7_write_request_packet)->payload->cotp_packet_data_eot = true;
  (*s7_write_request_packet)->payload->payload =
      malloc(sizeof(plc4c_s7_read_write_s7_message));
  (*s7_write_request_packet)->payload->payload->_type =
      plc4c_s7_read_write_s7_message_type_plc4c_s7_read_write_s7_message_request;
  (*s7_write_request_packet)->payload->payload->parameter =
      malloc(sizeof(plc4c_s7_read_write_s7_parameter));
  (*s7_write_request_packet)->payload->payload->parameter->_type =
      plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_write_var_request;
  plc4c_utils_list_create(
      &(*s7_write_request_packet)
          ->payload->payload->parameter->s7_parameter_read_var_request_items);

  plc4c_list_element* item = write_request->items->tail;
  while (item != NULL) {
    // Get the item address from the API request.
    char* itemAddress = item->value;

    // Create the item ...
    plc4c_s7_read_write_s7_var_request_parameter_item* request_item;
    plc4c_return_code return_code = parseAddress(itemAddress, &request_item);
    if (return_code != OK) {
      return return_code;
    }

    // Add the new item to the request.
    plc4c_utils_list_insert_head_value(
        (*s7_write_request_packet)
            ->payload->payload->parameter->s7_parameter_read_var_request_items,
        request_item);

    item = item->next;
  }

  (*s7_write_request_packet)->payload->payload->payload =
      malloc(sizeof(plc4c_s7_read_write_s7_payload));
  if((*s7_write_request_packet)->payload->payload->payload == NULL) {
    return NO_MEMORY;
  }
  (*s7_write_request_packet)->payload->payload->payload->_type =
      plc4c_s7_read_write_s7_payload_type_plc4c_s7_read_write_s7_payload_write_var_request;
  plc4c_utils_list_create(
      &(*s7_write_request_packet)->payload->payload->payload->s7_payload_write_var_request_items);

  // TODO: Implement the value encoding ...
  // TODO: Add all the encoded item values ...

  return OK;
}
