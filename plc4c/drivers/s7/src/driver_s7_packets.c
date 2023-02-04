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

#include <cotp_protocol_class.h>
#include <plc4c/driver_s7.h>
#include <plc4c/spi/types_private.h>
#include <stdlib.h>
#include <string.h>
#include <tpkt_packet.h>

#include "plc4c/driver_s7_packets.h"
#include "plc4c/driver_s7_encode_decode.h"

// undef to use pointer to plc4c_data item on writes
// probably safe to comment out but on for now.
//#define DEEP_WRITE_DATA_COPY 

// forward declaration for helper function todo: move to header or relocate
    //plc4c_utils_list_insert_head_value(request_value->data, &parsed_value->data);
void plc4c_add_data_to_request(plc4c_data* parsed_value, 
    plc4c_s7_read_write_s7_var_payload_data_item* request_value) {

  uint8_t* data_array;
#ifdef DEEP_WRITE_DATA_COPY
  uint8_t* deep_data_array;
#endif
  plc4c_list_element* list_array;
  plc4c_data* the_data;
  int i, j;
  size_t items;
  size_t item_size;

  // If its a list size doesn't really mean anything we care about list len
  if (parsed_value->data_type == PLC4C_LIST)
    items = plc4c_utils_list_size(parsed_value->data.list_value);
  else
    items = 1; 

  for (i = 0 ; i < items ; i++) {
    if (parsed_value->data_type == PLC4C_LIST) {
      // will not work with nested lists if that's even a thing
      list_array = (i == 0) ? parsed_value->data.list_value->head : list_array->previous;
      the_data = list_array->value;
      item_size = the_data->size;
      data_array = (uint8_t*) &the_data->data;
    } else {
      item_size = parsed_value->size;
      data_array = (uint8_t*) &parsed_value->data;
    }
    // Now add the bytes to a list
    for (j = 0 ; j < item_size ; j++) {
#ifdef DEEP_WRITE_DATA_COPY
        deep_data_array = calloc(1,sizeof(uint8_t));
        memcpy(deep_data_array, data_array + j, sizeof(uint8_t));
        plc4c_utils_list_insert_tail_value(request_value->data, deep_data_array);
#else
        plc4c_utils_list_insert_tail_value(request_value->data, data_array + j);
#endif
    }
  }
}


/**
 * Function used by the driver to tell the transport if there is a full
 * packet the driver can operate on available.
 *
 * @param buffer_data pointer to the buffer
 * @param buffer_length length of the buffer
 * @return positive integer = length of the packet, 0 = not enough,
 * try again later, negative integer = remove the given amount of bytes
 * from the buffer.
 */
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
          return (int16_t) -(i - 1);
        }
      }
      // We didn't find a new start, delete the entire content except the last
      // byte (as this could be the start of the next frame and we couldn't
      // confirm this.
      return (int16_t) -(buffer_length - 1);
    }
  }
  // The length information is located in bytes 3 and 4
  if (buffer_length >= 4) {
    uint16_t packet_length =
        ((uint16_t) *(buffer_data + 2) << 8) |
        ((uint16_t) *(buffer_data + 3));
    if (buffer_length >= packet_length) {
      return (int16_t) packet_length;
    }
    // 8192 is the maximum pdu size, so if the value is larger, the packet is
    // probably corrupt.
    if (packet_length > 8192) {
      for (int i = 1; i < (buffer_length - 1); i++) {
        buffer_data++;
        if ((*buffer_data == 0x03) && (*(buffer_data + 1) == 0x00)) {
          // We've found a potential new packet start.
          return (int16_t) -(i - 1);
        }
      }
      return (int16_t) -(buffer_length - 1);
    }
    return (int16_t) packet_length;
  }
  // In all other cases, we'll just have to wait for the next time.
  return 0;
}

plc4c_return_code plc4c_driver_s7_send_packet(plc4c_connection* connection,
                              plc4c_s7_read_write_tpkt_packet* packet) {
  
  uint16_t packet_size;
  plc4c_spi_write_buffer* write_buffer;
  plc4c_return_code return_code;
  
  // Get the size required to contain the serialized form of this packet.
  packet_size = plc4c_s7_read_write_tpkt_packet_length_in_bytes(packet);

  // Serialize this message to a byte-array.
  return_code = plc4c_spi_write_buffer_create(packet_size, &write_buffer);
  if (return_code != OK) {
    return return_code;
  }
  return_code = plc4c_s7_read_write_tpkt_packet_serialize(write_buffer, packet);
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

plc4c_return_code plc4c_driver_s7_receive_packet(plc4c_connection* connection,
                                 plc4c_s7_read_write_tpkt_packet** packet) {
  // Check with the transport if there is a packet available.
  // If it is, get a read_buffer for reading it.
  plc4c_spi_read_buffer* read_buffer;
  plc4c_return_code return_code;
  
  return_code = connection->transport->select_message(
      connection->transport_configuration, 4, 
      plc4c_driver_s7_select_message_function, &read_buffer);

  // OK is only returned if a packet is available.
  if (return_code != OK) {
    return return_code;
  }

  // Parse the packet by consuming the read_buffer data.
  *packet = NULL;
  return_code = plc4c_s7_read_write_tpkt_packet_parse(read_buffer, packet);
  if (return_code != OK) {
    return return_code;
  }

  // TODO: verify and maybe move
  free(read_buffer->data);
  free(read_buffer);
  
  // In this case a packet was available and parsed.
  return OK;
}

void delete_byte_list(plc4c_list_element *element) {
  char* item;
  item = element->value;
  // TODO: Fix this!
  //free(item);
}
void delete_s7_parameter_list_element(plc4c_list_element *element) {
  plc4c_s7_read_write_s7_parameter *item;
  item = element->value;
  free(item);
}
void delete_s7_read_response_payload_list_element(plc4c_list_element *element) {
  plc4c_s7_read_write_s7_var_payload_data_item *item;
  item = element->value;
  plc4c_utils_list_delete_elements(item->data, delete_byte_list);
  free(item->data);
  free(item);
}
void delete_s7_write_request_payload_list_element(plc4c_list_element *element) {
  plc4c_s7_read_write_s7_var_payload_data_item *item;
  item = element->value;
  free(item);
}
void delete_s7_write_response_payload_list_element(plc4c_list_element *element) {
  plc4c_s7_read_write_s7_var_payload_status_item *item;
  item = element->value;
  free(item);
}
void delete_mlfb_list(plc4c_list_element *element) {
  char* item;
  item = element->value;
  // TODO: Fix this!
  //free(item);
}
void delete_szl_list(plc4c_list_element *element){
  plc4c_s7_read_write_szl_data_tree_item *item;
  item = element->value;
  plc4c_utils_list_delete_elements(item->mlfb,delete_mlfb_list);
  free(item->mlfb);
  free(item);
}
void delete_s7_user_data_payload_list_element(plc4c_list_element *element) {
  plc4c_s7_read_write_s7_payload_user_data_item *item;
  item = element->value;

  plc4c_utils_list_delete_elements(
    item->s7_payload_user_data_item_cpu_function_read_szl_response_items, 
    delete_szl_list);
  free(item->s7_payload_user_data_item_cpu_function_read_szl_response_items);
  // TODO: Depending on the type of data-item, clean free the allocated memory.
  //free(item->szl_id);
  free(item);
}
void delete_s7_read_request_parameter_list_element(plc4c_list_element *element) {
  plc4c_s7_read_write_s7_var_request_parameter_item *item;
  item = element->value;
  free(item);
}
void delete_s7_write_request_parameter_list_element(plc4c_list_element *element) {
  plc4c_s7_read_write_s7_var_request_parameter_item *item;
  item = element->value;
  free(item);
}
void delete_s7_user_data_parameter_list_element(plc4c_list_element *element) {
  plc4c_s7_read_write_s7_parameter_user_data_item *item;
  item = element->value;
  // this is basically a copy of generated sources logic
  if (item->s7_parameter_user_data_item_cpu_functions_cpu_function_type == 8) {
    free(item->s7_parameter_user_data_item_cpu_functions_data_unit_reference_number);
    free(item->s7_parameter_user_data_item_cpu_functions_last_data_unit);
    free(item->s7_parameter_user_data_item_cpu_functions_error_code);
  }
  free(item);
}
void delete_copt_parameter_list_element(plc4c_list_element *element) {
  plc4c_s7_read_write_cotp_parameter *item;
  item = element->value;
  free(item);
}
void delete_payload_user_data_item_list_element(
    plc4c_list_element *element) {
  plc4c_s7_read_write_s7_payload_user_data_item *item;
  item = (plc4c_s7_read_write_s7_payload_user_data_item *)element->value;
  // TODO: Depending on the type of data-item, clean free the allocated memory.
  //free(item->szl_id);
  free(item);
}

void delete_parameter_user_data_item_list_element(
  plc4c_list_element *element) {
  plc4c_s7_read_write_s7_parameter_user_data_item *item;
  item = element->value;
  free(item);
}

void plc4c_driver_s7_destroy_receive_packet(
    plc4c_s7_read_write_tpkt_packet* packet) {

  plc4c_s7_read_write_cotp_packet *cotp_packet;
  plc4c_s7_read_write_s7_message *s7_message;
  plc4c_s7_read_write_s7_parameter *s7_param;
  plc4c_s7_read_write_s7_payload *s7_payload;

  cotp_packet = packet->payload;
  s7_message = packet->payload->payload;

  // Destroy s7 message with its parameters
  // and payload (note last 2 are optional)
  if (s7_message) {
    if ((s7_param = s7_message->parameter)) {
      switch (s7_param->_type) {
        case plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_read_var_request:
          plc4c_utils_list_delete_elements(s7_param->s7_parameter_read_var_request_items, 
              delete_s7_read_request_parameter_list_element);
          free(s7_param->s7_parameter_read_var_request_items);
          break;
        case plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_write_var_request:
          plc4c_utils_list_delete_elements(s7_param->s7_parameter_write_var_request_items, 
              delete_s7_write_request_parameter_list_element);
          free(s7_param->s7_parameter_write_var_request_items);
          break;
        case plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_user_data:
          plc4c_utils_list_delete_elements(s7_param->s7_parameter_user_data_items, 
              delete_s7_user_data_parameter_list_element);
          free(s7_param->s7_parameter_user_data_items);
          break;
        case plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_write_var_response:
          // TODO: something
          break;
        case plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_read_var_response:
          // TODO: something
          //just num of items so nothing to do
          break;
      }
      free(s7_param);
    }

    if ((s7_payload = s7_message->payload)) {
      switch (s7_payload->_type) {
        case plc4c_s7_read_write_s7_payload_type_plc4c_s7_read_write_s7_payload_read_var_response:
          plc4c_utils_list_delete_elements(s7_payload->s7_payload_read_var_response_items, 
              delete_s7_read_response_payload_list_element);
          free(s7_payload->s7_payload_read_var_response_items);
          break;
        case plc4c_s7_read_write_s7_payload_type_plc4c_s7_read_write_s7_payload_write_var_request:
          plc4c_utils_list_delete_elements(s7_payload->s7_payload_write_var_request_items, 
              delete_s7_write_request_payload_list_element);
          free(s7_payload->s7_payload_write_var_request_items);
          break;
        case plc4c_s7_read_write_s7_payload_type_plc4c_s7_read_write_s7_payload_write_var_response:
          plc4c_utils_list_delete_elements(s7_payload->s7_payload_write_var_response_items, 
              delete_s7_write_response_payload_list_element);
          free(s7_payload->s7_payload_write_var_response_items);
          break;
        case plc4c_s7_read_write_s7_payload_type_plc4c_s7_read_write_s7_payload_user_data:
          plc4c_utils_list_delete_elements(s7_payload->s7_payload_user_data_items, 
              delete_s7_user_data_payload_list_element);
          free(s7_payload->s7_payload_user_data_items);
          break;
      }
      free(s7_payload);
    }
    free(s7_message);
  }

  // Destroy copt message inc paramters list
  plc4c_utils_list_delete_elements(cotp_packet->parameters,
        delete_copt_parameter_list_element);
  free(cotp_packet->parameters);
  free(cotp_packet);

  // Destroy tpkt message
  free(packet);
}

/**
 * Create a COTP connection request packet.
 *
 * @param configuration configuration of the current connection.
 * @param plc4c_s7_read_write_tpkt_packet COTP connection-request (return)
 * @return OK, if the packet was correctly prepared, otherwise not-OK.
 */
plc4c_return_code plc4c_driver_s7_create_cotp_connection_request(
    plc4c_driver_s7_config* configuration,
    plc4c_s7_read_write_tpkt_packet** packet) {
  
  *packet = malloc(sizeof(plc4c_s7_read_write_tpkt_packet));
  if (*packet == NULL) {
    return NO_MEMORY;
  }

  (*packet)->payload = malloc(sizeof(plc4c_s7_read_write_cotp_packet));
  if ((*packet)->payload == NULL) {
    return NO_MEMORY;
  }

  (*packet)->payload->_type = 
      plc4c_s7_read_write_cotp_packet_type_plc4c_s7_read_write_cotp_packet_connection_request;
  (*packet)->payload->cotp_packet_connection_request_destination_reference = 0x0000;
  (*packet)->payload->cotp_packet_connection_request_source_reference = 0x000F;
  (*packet)->payload->cotp_packet_connection_request_protocol_class =
      plc4c_s7_read_write_cotp_protocol_class_CLASS_0;

  // Add the COTP parameters: Called TSAP, Calling TSAP and TPDU Size.
  plc4c_utils_list_create(&((*packet)->payload->parameters));
  plc4c_s7_read_write_cotp_parameter* called_tsap_param;
  plc4c_s7_read_write_cotp_parameter* calling_tsap_param;
  plc4c_s7_read_write_cotp_parameter* tpdu_size_param;
  
  called_tsap_param = malloc(sizeof(plc4c_s7_read_write_cotp_parameter));
  if (called_tsap_param == NULL) {
    return NO_MEMORY;
  }
  called_tsap_param->_type =
      plc4c_s7_read_write_cotp_parameter_type_plc4c_s7_read_write_cotp_parameter_called_tsap;
  called_tsap_param->cotp_parameter_called_tsap_tsap_id = configuration->called_tsap_id;

  plc4c_utils_list_insert_head_value((*packet)->payload->parameters, called_tsap_param);
  calling_tsap_param = malloc(sizeof(plc4c_s7_read_write_cotp_parameter));
  if (calling_tsap_param == NULL) {
    return NO_MEMORY;
  }
  calling_tsap_param->_type =
      plc4c_s7_read_write_cotp_parameter_type_plc4c_s7_read_write_cotp_parameter_calling_tsap;
  calling_tsap_param->cotp_parameter_calling_tsap_tsap_id =
      configuration->calling_tsap_id;

  plc4c_utils_list_insert_head_value((*packet)->payload->parameters, calling_tsap_param);
  tpdu_size_param = malloc(sizeof(plc4c_s7_read_write_cotp_parameter));
  if (tpdu_size_param == NULL) {
    return NO_MEMORY;
  }
  tpdu_size_param->_type =
      plc4c_s7_read_write_cotp_parameter_type_plc4c_s7_read_write_cotp_parameter_tpdu_size;
  tpdu_size_param->cotp_parameter_tpdu_size_tpdu_size =
      configuration->cotp_tpdu_size;

  plc4c_utils_list_insert_head_value(
      (*packet)->payload->parameters, tpdu_size_param);

  // For a COTP connection request, there is no payload.
  (*packet)->payload->payload = NULL;

  return OK;
}


void plc4c_driver_s7_destroy_cotp_connection_request(
    plc4c_s7_read_write_tpkt_packet *packet) {
    plc4c_utils_list_delete_elements(packet->payload->parameters,
        delete_copt_parameter_list_element);
    free(packet->payload->parameters);
    free(packet->payload);
    free(packet);
}

/**
 * Create a S7 connection request packet.
 *
 * @param configuration configuration of the current connection.
 * @param plc4c_s7_read_write_tpkt_packet S7 connection-request (return)
 * @return OK, if the packet was correctly prepared, otherwise not-OK.
 */
plc4c_return_code plc4c_driver_s7_create_s7_connection_request(
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
      calloc(1,sizeof(plc4c_s7_read_write_s7_message));
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

  (*s7_connect_request_packet)->payload->payload->payload = NULL;
/*      malloc(sizeof(plc4c_s7_read_write_s7_payload));
  if ((*s7_connect_request_packet)->payload->payload->payload == NULL) {
    return NO_MEMORY;
  }
  (*s7_connect_request_packet)->payload->payload->payload->_type =
      plc4c_s7_read_write_s7_payload_type_plc4c_s7_read_write_s7_payload_user_data;*/
  return OK;
}

void plc4c_driver_s7_destroy_s7_connection_request(
    plc4c_s7_read_write_tpkt_packet* packet) {

  free(packet->payload->payload->parameter);
  free(packet->payload->payload);
  free(packet->payload);
  free(packet);

}

/**
 * Create a S7 identify remote request packet
 *
 * @param configuration configuration of the current connection.
 * @param plc4c_s7_read_write_tpkt_packet S7 identify remote request (return)
 * @return OK, if the packet was correctly prepared, otherwise not-OK.
 */
plc4c_return_code plc4c_driver_s7_create_s7_identify_remote_request(
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
      calloc(1,sizeof(plc4c_s7_read_write_s7_message));
  if ((*s7_identify_remote_request_packet)->payload->payload == NULL) {
    return NO_MEMORY;
  }
  (*s7_identify_remote_request_packet)->payload->payload->_type =
      plc4c_s7_read_write_s7_message_type_plc4c_s7_read_write_s7_message_user_data;

  // Create a Parameter
  (*s7_identify_remote_request_packet)->payload->payload->parameter =
      malloc(sizeof(plc4c_s7_read_write_s7_parameter));
  if ((*s7_identify_remote_request_packet)->payload->payload->parameter == NULL) {
    return NO_MEMORY;
  }
  (*s7_identify_remote_request_packet)->payload->payload->parameter->_type =
      plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_user_data;
  plc4c_utils_list_create(
      &((*s7_identify_remote_request_packet)
            ->payload->payload->parameter->s7_parameter_user_data_items));
  // Create the Parameter Item
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

  // Create the Payload
  (*s7_identify_remote_request_packet)->payload->payload->payload = malloc(sizeof(plc4c_s7_read_write_s7_payload));
  if ((*s7_identify_remote_request_packet)->payload->payload->parameter == NULL) {
    return NO_MEMORY;
  }
  (*s7_identify_remote_request_packet)->payload->payload->payload->_type =
      plc4c_s7_read_write_s7_payload_type_plc4c_s7_read_write_s7_payload_user_data;
  plc4c_utils_list_create(
      &((*s7_identify_remote_request_packet)
            ->payload->payload->payload->s7_payload_user_data_items));
  // Create the Payload Item
  plc4c_s7_read_write_s7_payload_user_data_item* payload_item =
      malloc(sizeof(plc4c_s7_read_write_s7_payload_user_data_item));
  payload_item->_type =
      plc4c_s7_read_write_s7_payload_user_data_item_type_plc4c_s7_read_write_s7_payload_user_data_item_cpu_function_read_szl_request;
  payload_item->return_code = plc4c_s7_read_write_data_transport_error_code_OK;
  payload_item->transport_size =
      plc4c_s7_read_write_data_transport_size_OCTET_STRING;
  payload_item->s7_payload_user_data_item_cpu_function_read_szl_request_szl_index = 0x0000;
  payload_item->s7_payload_user_data_item_cpu_function_read_szl_request_szl_id = malloc(sizeof(plc4c_s7_read_write_szl_id));
  if (payload_item->s7_payload_user_data_item_cpu_function_read_szl_request_szl_id == NULL) {
    return NO_MEMORY;
  }
  payload_item->s7_payload_user_data_item_cpu_function_read_szl_request_szl_id->type_class =
      plc4c_s7_read_write_szl_module_type_class_CPU;
  payload_item->s7_payload_user_data_item_cpu_function_read_szl_request_szl_id->sublist_extract = 0x00;
  payload_item->s7_payload_user_data_item_cpu_function_read_szl_request_szl_id->sublist_list =
      plc4c_s7_read_write_szl_sublist_MODULE_IDENTIFICATION;
  plc4c_utils_list_insert_head_value(
      (*s7_identify_remote_request_packet)
          ->payload->payload->payload->s7_payload_user_data_items,
      payload_item);

  return OK;
}

void plc4c_driver_s7_destroy_s7_identify_remote_request(
    plc4c_s7_read_write_tpkt_packet* packet) {

  // Delete the s7 payload list
  plc4c_s7_read_write_s7_message *s7_message;
  s7_message = packet->payload->payload;
  plc4c_utils_list_delete_elements(
      s7_message->payload->s7_payload_user_data_items,
      delete_payload_user_data_item_list_element);
  free(s7_message->payload->s7_payload_user_data_items);
  free(s7_message->payload);

  // Delete the parameters list 
  plc4c_utils_list_delete_elements(
      s7_message->parameter->s7_parameter_user_data_items,
      delete_parameter_user_data_item_list_element);
  free(s7_message->parameter->s7_parameter_user_data_items);
  free(s7_message->parameter);

  // Free copt payload, tpkp payload and tpk packet
  free(packet->payload->payload);
  free(packet->payload);
  free(packet);

}
/**
 * Create a S7 read request packet
 *
 * @param configuration configuration of the current connection.
 * @param plc4c_s7_read_write_tpkt_packet S7 read request (return)
 * @return OK, if the packet was correctly prepared, otherwise not-OK.
 */
plc4c_return_code plc4c_driver_s7_create_s7_read_request(
    plc4c_read_request* read_request,
    plc4c_s7_read_write_tpkt_packet** packet) {

  plc4c_driver_s7_config* configuration;
  plc4c_s7_read_write_cotp_packet *cotp_packet;
  plc4c_s7_read_write_s7_message *s7_message;
  plc4c_list_element* cur_item;

  configuration = read_request->connection->configuration;
  
  *packet = malloc(sizeof(plc4c_s7_read_write_tpkt_packet));
  if (*packet == NULL) {
    return NO_MEMORY;
  }

  (*packet)->payload = malloc(sizeof(plc4c_s7_read_write_cotp_packet));
    
  // Terse local variable for clarity
  cotp_packet = (*packet)->payload;

  if (cotp_packet == NULL) {
    return NO_MEMORY;
  }
  cotp_packet->_type = plc4c_s7_read_write_cotp_packet_type_plc4c_s7_read_write_cotp_packet_data;
  cotp_packet->cotp_packet_data_tpdu_ref = configuration->pdu_id++;
  cotp_packet->cotp_packet_data_eot = true;
  cotp_packet->parameters = NULL;
  cotp_packet->payload = calloc(1,sizeof(plc4c_s7_read_write_s7_message));
  s7_message = cotp_packet->payload;
  if(s7_message == NULL) {
    return NO_MEMORY;
  }
  s7_message->_type = plc4c_s7_read_write_s7_message_type_plc4c_s7_read_write_s7_message_request;
  s7_message->parameter = malloc(sizeof(plc4c_s7_read_write_s7_parameter));
  if(s7_message->parameter == NULL) {
    return NO_MEMORY;
  }
  s7_message->parameter->_type = plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_read_var_request;
  plc4c_utils_list_create(&s7_message->parameter->s7_parameter_read_var_request_items);

  cur_item = read_request->items->tail;
  
  while (cur_item != NULL) {

    plc4c_item* item;
    plc4c_s7_read_write_s7_var_request_parameter_item* parsed_item_address;
    plc4c_s7_read_write_s7_var_request_parameter_item* updated_item_address;

    item = cur_item->value;
    plc4c_s7_read_write_s7_var_request_parameter_item_field* field = item->address;
    // Get the item address from the API request.
    parsed_item_address = field->parameter_item;

    // Create a copy of the request item...
    updated_item_address = malloc(sizeof(plc4c_s7_read_write_s7_var_request_parameter_item));
    if (updated_item_address == NULL) {
      return NO_MEMORY;
    }
    updated_item_address->_type = parsed_item_address->_type;
    updated_item_address->s7_var_request_parameter_item_address_address =
        malloc(sizeof(plc4c_s7_read_write_s7_address));
    if (updated_item_address->s7_var_request_parameter_item_address_address == NULL) {
      return NO_MEMORY;
    }
    // Memcpy inplace of fields assignment, as all fields where assigned
    memcpy(updated_item_address->s7_var_request_parameter_item_address_address,
      parsed_item_address->s7_var_request_parameter_item_address_address,
      sizeof(plc4c_s7_read_write_s7_address));

    // In case of TIME values, we read 4 bytes for each value instead.
    plc4c_driver_s7_time_transport_size(
      &updated_item_address->s7_var_request_parameter_item_address_address->s7_address_any_transport_size);
    
    // Add the new item to the request.
    plc4c_utils_list_insert_head_value(
        s7_message->parameter->s7_parameter_read_var_request_items,
        updated_item_address);

    cur_item = cur_item->next;
  }

  s7_message->payload = NULL;
  return OK;
}

void delete_s7_parameter_read_var_request_item_element(plc4c_list_element *element) {
  
  plc4c_s7_read_write_s7_var_request_parameter_item *item;
  item = element->value;
  free(item->s7_var_request_parameter_item_address_address);
  free(item);
}

void plc4c_driver_s7_destroy_s7_read_request(
    plc4c_s7_read_write_tpkt_packet* packet) {

  plc4c_s7_read_write_s7_message *s7_message;
  s7_message = packet->payload->payload;

  // Free the s7 message its parameter, and its parameter list
  plc4c_utils_list_delete_elements(
    s7_message->parameter->s7_parameter_read_var_request_items, 
    delete_s7_parameter_read_var_request_item_element);
  free(s7_message->parameter->s7_parameter_read_var_request_items);
  free(s7_message->parameter);
  free(s7_message);

  //s7_parameter_read_var_request_items
  free(packet->payload);
  free(packet);
}

/**
 * Create a S7 write request packet
 *
 * @param configuration configuration of the current connection.
 * @param plc4c_s7_read_write_tpkt_packet S7 write request (return)
 * @return OK, if the packet was correctly prepared, otherwise not-OK.
 */
plc4c_return_code plc4c_driver_s7_create_s7_write_request(
    plc4c_write_request* request,
    plc4c_s7_read_write_tpkt_packet** request_packet) {


  plc4c_driver_s7_config* configuration;
  plc4c_s7_read_write_cotp_packet *cotp_packet;
  plc4c_s7_read_write_s7_message *s7_packet;
  plc4c_list_element* list_item;
  plc4c_s7_read_write_s7_parameter* s7_parameters;
  plc4c_s7_read_write_s7_payload* s7_payload;

  configuration = request->connection->configuration;

  *request_packet = malloc(sizeof(plc4c_s7_read_write_tpkt_packet));
  if (*request_packet == NULL) 
    return NO_MEMORY;
  
  (*request_packet)->payload = malloc(sizeof(plc4c_s7_read_write_cotp_packet));
  if ((*request_packet)->payload == NULL) 
    return NO_MEMORY;
  
  // Terse local variable for clarity
  cotp_packet = (*request_packet)->payload;

  cotp_packet->_type = plc4c_s7_read_write_cotp_packet_type_plc4c_s7_read_write_cotp_packet_data;
  cotp_packet->cotp_packet_data_tpdu_ref = configuration->pdu_id++;
  cotp_packet->cotp_packet_data_eot = true;
  cotp_packet->parameters = NULL;

  // Allocate and initalise payload->payload 
  cotp_packet->payload = calloc(1,sizeof(plc4c_s7_read_write_s7_message));
  s7_packet = cotp_packet->payload;
  if (s7_packet == NULL)
    return NO_MEMORY;
  s7_packet->_type = plc4c_s7_read_write_s7_message_type_plc4c_s7_read_write_s7_message_request;

  // Allocate and initalise payload->payload->parameter
  s7_packet->parameter = malloc(sizeof(plc4c_s7_read_write_s7_parameter));
  s7_parameters = s7_packet->parameter;
  if (s7_parameters == NULL)
    return NO_MEMORY;
  s7_parameters->_type = plc4c_s7_read_write_s7_parameter_type_plc4c_s7_read_write_s7_parameter_write_var_request;
  plc4c_utils_list_create(&s7_parameters->s7_parameter_write_var_request_items);
  
  //  Allocate and initalise payload->payload->payload
  s7_packet->payload = malloc(sizeof(plc4c_s7_read_write_s7_payload));
  s7_payload = s7_packet->payload;
  if (s7_payload == NULL) 
    return NO_MEMORY;
  s7_payload->_type = plc4c_s7_read_write_s7_payload_type_plc4c_s7_read_write_s7_payload_write_var_request;
  plc4c_utils_list_create(&s7_payload->s7_payload_write_var_request_items);
  
  list_item = request->items->tail;
  
  while (list_item != NULL) {

    plc4c_request_value_item *item;
    plc4c_item *parsed_item;
    plc4c_s7_read_write_s7_var_request_parameter_item *parsed_param;
    plc4c_s7_read_write_s7_var_request_parameter_item *request_param;
    plc4c_data *parsed_value;
    plc4c_s7_read_write_s7_var_payload_data_item *request_value;
    plc4c_return_code return_code;
    int8_t* data_array;
    
    // Set things off to a good start 
    return_code = OK;

    // Get the item address from the API request.
    item = list_item->value; 
    parsed_item = item->item;

    plc4c_s7_read_write_s7_var_request_parameter_item_field* field = parsed_item->address;

    parsed_param = field->parameter_item;
    parsed_value = item->value;

    // Make a copy of the param
    request_param = malloc(sizeof(plc4c_s7_read_write_s7_var_request_parameter_item));
    if (!request_param) 
      return NO_MEMORY;
    
    request_param->_type = parsed_param->_type;
    request_param->s7_var_request_parameter_item_address_address = 
        malloc(sizeof(plc4c_s7_read_write_s7_address));
    if (!request_param->s7_var_request_parameter_item_address_address)
      return NO_MEMORY;
    
    memcpy(request_param->s7_var_request_parameter_item_address_address,
      parsed_param->s7_var_request_parameter_item_address_address, 
      sizeof(plc4c_s7_read_write_s7_address));

    if (return_code != OK) 
      return return_code;
    
    // Make a copy of the value 
    request_value = malloc(sizeof(plc4c_s7_read_write_s7_var_payload_data_item));
    
    // get transport data size (defined via transport size)
    request_value->transport_size = 
        plc4c_s7_read_write_transport_size_get_data_transport_size(
          parsed_param->s7_var_request_parameter_item_address_address->s7_address_any_transport_size);
    request_value->return_code = 0;
    
    plc4c_utils_list_create(&request_value->data);
    
    plc4c_add_data_to_request(parsed_value, request_value);

    // Add the new parameter to the request
    plc4c_utils_list_insert_head_value(s7_parameters->s7_parameter_write_var_request_items, request_param);
    
    // Add the new value to the request
    plc4c_utils_list_insert_head_value(s7_payload->s7_payload_write_var_request_items, request_value);

    list_item = list_item->next;
  }

  return OK;
}


void delete_s7_data_list_element(plc4c_list_element* element) {

  uint8_t *item;
  item = element->value;
#ifdef DEEP_WRITE_DATA_COPY
  free(item);
#endif
}

void delete_s7_payload_write_request_list_element(
    plc4c_list_element* element) {

  plc4c_s7_read_write_s7_var_payload_data_item *item;
  item = element->value;
  
  plc4c_utils_list_delete_elements(item->data,
      delete_s7_data_list_element);
  free(item->data);
  free(item);
}

void delete_s7_parameter_write_request_list_element(
    plc4c_list_element* element) {

  plc4c_s7_read_write_s7_var_request_parameter_item *item;
  item = element->value;

  free(item->s7_var_request_parameter_item_address_address);
  free(item);

}

void plc4c_driver_s7_destroy_s7_write_request(
    plc4c_s7_read_write_tpkt_packet* packet) {
  
  plc4c_s7_read_write_s7_message *s7_message;
  s7_message = packet->payload->payload;

  // Delete the s7 payload items
  plc4c_utils_list_delete_elements(
      s7_message->payload->s7_payload_write_var_request_items,
      delete_s7_payload_write_request_list_element);

  free(s7_message->payload->s7_payload_write_var_request_items);
  free(s7_message->payload);

  // Delete the s7 parameter items
  plc4c_utils_list_delete_elements(
      s7_message->parameter->s7_parameter_write_var_request_items,
      delete_s7_parameter_write_request_list_element);
  free(s7_message->parameter->s7_parameter_write_var_request_items);
  free(s7_message->parameter);

  free(s7_message);
  free(packet->payload);
  free(packet);
}

/**
 * Adjust transport sizes for time values
 *
 * @param transport_size current transport size of item
 */
void plc4c_driver_s7_time_transport_size(plc4c_s7_read_write_transport_size *transport_size) {
  // In case of TIME values, we read 4 bytes for each value instead.
  switch (*transport_size) {
    case plc4c_s7_read_write_transport_size_TIME:
      *transport_size = plc4c_s7_read_write_transport_size_DINT;
      break;
    case plc4c_s7_read_write_transport_size_DATE:
      *transport_size = plc4c_s7_read_write_transport_size_UINT;
      break;
    case plc4c_s7_read_write_transport_size_TIME_OF_DAY:
    case plc4c_s7_read_write_transport_size_TOD:
      *transport_size = plc4c_s7_read_write_transport_size_UDINT;
      break;
  }
}
