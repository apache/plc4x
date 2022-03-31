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
#include <plc4c/spi/evaluation_helper.h>
#include "mqt_t__property.h"

// Code generated by code-generation. DO NOT EDIT.

// Array of discriminator values that match the enum type constants.
// (The order is identical to the enum constants so we can use the
// enum constant to directly access a given types discriminator values)
const plc4c_mqtt_read_write_mqt_t__property_discriminator plc4c_mqtt_read_write_mqt_t__property_discriminators[] = {
  {/* plc4c_mqtt_read_write_mqt_t__property__payloa_d__forma_t__indicator */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_PAYLOAD_FORMAT_INDICATOR },
  {/* plc4c_mqtt_read_write_mqt_t__property__messag_e__expir_y__interval */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_MESSAGE_EXPIRY_INTERVAL },
  {/* plc4c_mqtt_read_write_mqt_t__property__conten_t__type */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_CONTENT_TYPE },
  {/* plc4c_mqtt_read_write_mqt_t__property__respons_e__topic */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_RESPONSE_TOPIC },
  {/* plc4c_mqtt_read_write_mqt_t__property__correlatio_n__data */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_CORRELATION_DATA },
  {/* plc4c_mqtt_read_write_mqt_t__property__subscriptio_n__identifier */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_SUBSCRIPTION_IDENTIFIER },
  {/* plc4c_mqtt_read_write_mqt_t__property__expir_y__interval */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_SESSION_EXPIRY_INTERVAL },
  {/* plc4c_mqtt_read_write_mqt_t__property__assigne_d__clien_t__identifier */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_ASSIGNED_CLIENT_IDENTIFIER },
  {/* plc4c_mqtt_read_write_mqt_t__property__serve_r__kee_p__alive */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_SERVER_KEEP_ALIVE },
  {/* plc4c_mqtt_read_write_mqt_t__property__authenticatio_n__method */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_AUTHENTICATION_METHOD },
  {/* plc4c_mqtt_read_write_mqt_t__property__authenticatio_n__data */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_AUTHENTICATION_DATA },
  {/* plc4c_mqtt_read_write_mqt_t__property__reques_t__proble_m__information */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_REQUEST_PROBLEM_INFORMATION },
  {/* plc4c_mqtt_read_write_mqt_t__property__wil_l__dela_y__interval */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_WILL_DELAY_INTERVAL },
  {/* plc4c_mqtt_read_write_mqt_t__property__reques_t__respons_e__information */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_REQUEST_RESPONSE_INFORMATION },
  {/* plc4c_mqtt_read_write_mqt_t__property__respons_e__information */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_RESPONSE_INFORMATION },
  {/* plc4c_mqtt_read_write_mqt_t__property__serve_r__reference */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_SERVER_REFERENCE },
  {/* plc4c_mqtt_read_write_mqt_t__property__reaso_n__string */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_REASON_STRING },
  {/* plc4c_mqtt_read_write_mqt_t__property__receiv_e__maximum */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_RECEIVE_MAXIMUM },
  {/* plc4c_mqtt_read_write_mqt_t__property__topi_c__alia_s__maximum */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_TOPIC_ALIAS_MAXIMUM },
  {/* plc4c_mqtt_read_write_mqt_t__property__topi_c__alias */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_TOPIC_ALIAS },
  {/* plc4c_mqtt_read_write_mqt_t__property__maximu_m__qos */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_MAXIMUM_QOS },
  {/* plc4c_mqtt_read_write_mqt_t__property__retai_n__available */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_RETAIN_AVAILABLE },
  {/* plc4c_mqtt_read_write_mqt_t__property__use_r__property */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_USER_PROPERTY },
  {/* plc4c_mqtt_read_write_mqt_t__property__maximu_m__packe_t__size */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_MAXIMUM_PACKET_SIZE },
  {/* plc4c_mqtt_read_write_mqt_t__property__wildcar_d__subscriptio_n__available */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_WILDCARD_SUBSCRIPTION_AVAILABLE },
  {/* plc4c_mqtt_read_write_mqt_t__property__subscriptio_n__identifie_r__available */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_SUBSCRIPTION_IDENTIFIER_AVAILABLE },
  {/* plc4c_mqtt_read_write_mqt_t__property__share_d__subscriptio_n__available */
   .propertyType = plc4c_mqtt_read_write_mqt_t__property_type_SHARED_SUBSCRIPTION_AVAILABLE }

};

// Function returning the discriminator values for a given type constant.
plc4c_mqtt_read_write_mqt_t__property_discriminator plc4c_mqtt_read_write_mqt_t__property_get_discriminator(plc4c_mqtt_read_write_mqt_t__property_type type) {
  return plc4c_mqtt_read_write_mqt_t__property_discriminators[type];
}

// Create an empty NULL-struct
static const plc4c_mqtt_read_write_mqt_t__property plc4c_mqtt_read_write_mqt_t__property_null_const;

plc4c_mqtt_read_write_mqt_t__property plc4c_mqtt_read_write_mqt_t__property_null() {
  return plc4c_mqtt_read_write_mqt_t__property_null_const;
}


// Parse function.
plc4c_return_code plc4c_mqtt_read_write_mqt_t__property_parse(plc4c_spi_read_buffer* readBuffer, plc4c_mqtt_read_write_mqt_t__property** _message) {
  uint16_t startPos = plc4c_spi_read_get_pos(readBuffer);
  plc4c_return_code _res = OK;

  // Allocate enough memory to contain this data structure.
  (*_message) = malloc(sizeof(plc4c_mqtt_read_write_mqt_t__property));
  if(*_message == NULL) {
    return NO_MEMORY;
  }

  // Simple Field (propertyType)
  plc4c_mqtt_read_write_mqt_t__property_type* propertyType;
  _res = plc4c_mqtt_read_write_mqt_t__property_type_parse(readBuffer, (void*) &propertyType);
  if(_res != OK) {
    return _res;
  }
  (*_message)->property_type = *propertyType;

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_PAYLOAD_FORMAT_INDICATOR) { /* MQTT_Property_PAYLOAD_FORMAT_INDICATOR */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__payloa_d__forma_t__indicator;
                    
    // Simple Field (value)
    uint8_t value = 0;
    _res = plc4c_spi_read_unsigned_byte(readBuffer, 8, (uint8_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__payloa_d__forma_t__indicator_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_MESSAGE_EXPIRY_INTERVAL) { /* MQTT_Property_MESSAGE_EXPIRY_INTERVAL */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__messag_e__expir_y__interval;
                    
    // Simple Field (value)
    uint32_t value = 0;
    _res = plc4c_spi_read_unsigned_int(readBuffer, 32, (uint32_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__messag_e__expir_y__interval_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_CONTENT_TYPE) { /* MQTT_Property_CONTENT_TYPE */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__conten_t__type;
                    
    // Simple Field (value)
    plc4c_mqtt_read_write_mqt_t__string* value;
    _res = plc4c_mqtt_read_write_mqt_t__string_parse(readBuffer, (void*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__conten_t__type_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_RESPONSE_TOPIC) { /* MQTT_Property_RESPONSE_TOPIC */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__respons_e__topic;
                    
    // Simple Field (value)
    plc4c_mqtt_read_write_mqt_t__string* value;
    _res = plc4c_mqtt_read_write_mqt_t__string_parse(readBuffer, (void*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__respons_e__topic_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_CORRELATION_DATA) { /* MQTT_Property_CORRELATION_DATA */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__correlatio_n__data;
  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_SUBSCRIPTION_IDENTIFIER) { /* MQTT_Property_SUBSCRIPTION_IDENTIFIER */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__subscriptio_n__identifier;
                    
    // Simple Field (value)
    uint32_t value = 0;
    _res = plc4c_spi_read_unsigned_int(readBuffer, 32, (uint32_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__subscriptio_n__identifier_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_SESSION_EXPIRY_INTERVAL) { /* MQTT_Property_EXPIRY_INTERVAL */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__expir_y__interval;
                    
    // Simple Field (value)
    uint32_t value = 0;
    _res = plc4c_spi_read_unsigned_int(readBuffer, 32, (uint32_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__expir_y__interval_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_ASSIGNED_CLIENT_IDENTIFIER) { /* MQTT_Property_ASSIGNED_CLIENT_IDENTIFIER */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__assigne_d__clien_t__identifier;
                    
    // Simple Field (value)
    plc4c_mqtt_read_write_mqt_t__string* value;
    _res = plc4c_mqtt_read_write_mqt_t__string_parse(readBuffer, (void*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__assigne_d__clien_t__identifier_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_SERVER_KEEP_ALIVE) { /* MQTT_Property_SERVER_KEEP_ALIVE */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__serve_r__kee_p__alive;
                    
    // Simple Field (value)
    uint16_t value = 0;
    _res = plc4c_spi_read_unsigned_short(readBuffer, 16, (uint16_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__serve_r__kee_p__alive_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_AUTHENTICATION_METHOD) { /* MQTT_Property_AUTHENTICATION_METHOD */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__authenticatio_n__method;
                    
    // Simple Field (value)
    plc4c_mqtt_read_write_mqt_t__string* value;
    _res = plc4c_mqtt_read_write_mqt_t__string_parse(readBuffer, (void*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__authenticatio_n__method_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_AUTHENTICATION_DATA) { /* MQTT_Property_AUTHENTICATION_DATA */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__authenticatio_n__data;
  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_REQUEST_PROBLEM_INFORMATION) { /* MQTT_Property_REQUEST_PROBLEM_INFORMATION */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__reques_t__proble_m__information;
                    
    // Simple Field (value)
    uint8_t value = 0;
    _res = plc4c_spi_read_unsigned_byte(readBuffer, 8, (uint8_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__reques_t__proble_m__information_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_WILL_DELAY_INTERVAL) { /* MQTT_Property_WILL_DELAY_INTERVAL */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__wil_l__dela_y__interval;
                    
    // Simple Field (value)
    uint32_t value = 0;
    _res = plc4c_spi_read_unsigned_int(readBuffer, 32, (uint32_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__wil_l__dela_y__interval_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_REQUEST_RESPONSE_INFORMATION) { /* MQTT_Property_REQUEST_RESPONSE_INFORMATION */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__reques_t__respons_e__information;
                    
    // Simple Field (value)
    uint8_t value = 0;
    _res = plc4c_spi_read_unsigned_byte(readBuffer, 8, (uint8_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__reques_t__respons_e__information_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_RESPONSE_INFORMATION) { /* MQTT_Property_RESPONSE_INFORMATION */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__respons_e__information;
                    
    // Simple Field (value)
    plc4c_mqtt_read_write_mqt_t__string* value;
    _res = plc4c_mqtt_read_write_mqt_t__string_parse(readBuffer, (void*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__respons_e__information_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_SERVER_REFERENCE) { /* MQTT_Property_SERVER_REFERENCE */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__serve_r__reference;
                    
    // Simple Field (value)
    plc4c_mqtt_read_write_mqt_t__string* value;
    _res = plc4c_mqtt_read_write_mqt_t__string_parse(readBuffer, (void*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__serve_r__reference_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_REASON_STRING) { /* MQTT_Property_REASON_STRING */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__reaso_n__string;
                    
    // Simple Field (value)
    plc4c_mqtt_read_write_mqt_t__string* value;
    _res = plc4c_mqtt_read_write_mqt_t__string_parse(readBuffer, (void*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__reaso_n__string_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_RECEIVE_MAXIMUM) { /* MQTT_Property_RECEIVE_MAXIMUM */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__receiv_e__maximum;
                    
    // Simple Field (value)
    uint16_t value = 0;
    _res = plc4c_spi_read_unsigned_short(readBuffer, 16, (uint16_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__receiv_e__maximum_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_TOPIC_ALIAS_MAXIMUM) { /* MQTT_Property_TOPIC_ALIAS_MAXIMUM */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__topi_c__alia_s__maximum;
                    
    // Simple Field (value)
    uint16_t value = 0;
    _res = plc4c_spi_read_unsigned_short(readBuffer, 16, (uint16_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__topi_c__alia_s__maximum_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_TOPIC_ALIAS) { /* MQTT_Property_TOPIC_ALIAS */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__topi_c__alias;
                    
    // Simple Field (value)
    uint16_t value = 0;
    _res = plc4c_spi_read_unsigned_short(readBuffer, 16, (uint16_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__topi_c__alias_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_MAXIMUM_QOS) { /* MQTT_Property_MAXIMUM_QOS */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__maximu_m__qos;
                    
    // Simple Field (value)
    uint8_t value = 0;
    _res = plc4c_spi_read_unsigned_byte(readBuffer, 8, (uint8_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__maximu_m__qos_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_RETAIN_AVAILABLE) { /* MQTT_Property_RETAIN_AVAILABLE */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__retai_n__available;
                    
    // Simple Field (value)
    uint8_t value = 0;
    _res = plc4c_spi_read_unsigned_byte(readBuffer, 8, (uint8_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__retai_n__available_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_USER_PROPERTY) { /* MQTT_Property_USER_PROPERTY */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__use_r__property;
                    
    // Simple Field (name)
    plc4c_mqtt_read_write_mqt_t__string* name;
    _res = plc4c_mqtt_read_write_mqt_t__string_parse(readBuffer, (void*) &name);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__use_r__property_name = name;


                    
    // Simple Field (value)
    plc4c_mqtt_read_write_mqt_t__string* value;
    _res = plc4c_mqtt_read_write_mqt_t__string_parse(readBuffer, (void*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__use_r__property_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_MAXIMUM_PACKET_SIZE) { /* MQTT_Property_MAXIMUM_PACKET_SIZE */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__maximu_m__packe_t__size;
                    
    // Simple Field (value)
    uint32_t value = 0;
    _res = plc4c_spi_read_unsigned_int(readBuffer, 32, (uint32_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__maximu_m__packe_t__size_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_WILDCARD_SUBSCRIPTION_AVAILABLE) { /* MQTT_Property_WILDCARD_SUBSCRIPTION_AVAILABLE */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__wildcar_d__subscriptio_n__available;
                    
    // Simple Field (value)
    uint8_t value = 0;
    _res = plc4c_spi_read_unsigned_byte(readBuffer, 8, (uint8_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__wildcar_d__subscriptio_n__available_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_SUBSCRIPTION_IDENTIFIER_AVAILABLE) { /* MQTT_Property_SUBSCRIPTION_IDENTIFIER_AVAILABLE */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__subscriptio_n__identifie_r__available;
                    
    // Simple Field (value)
    uint8_t value = 0;
    _res = plc4c_spi_read_unsigned_byte(readBuffer, 8, (uint8_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__subscriptio_n__identifie_r__available_value = value;

  } else 
  if(propertyType == plc4c_mqtt_read_write_mqt_t__property_type_SHARED_SUBSCRIPTION_AVAILABLE) { /* MQTT_Property_SHARED_SUBSCRIPTION_AVAILABLE */
    (*_message)->_type = plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__share_d__subscriptio_n__available;
                    
    // Simple Field (value)
    uint8_t value = 0;
    _res = plc4c_spi_read_unsigned_byte(readBuffer, 8, (uint8_t*) &value);
    if(_res != OK) {
      return _res;
    }
    (*_message)->mqt_t__property__share_d__subscriptio_n__available_value = value;

  }

  return OK;
}

plc4c_return_code plc4c_mqtt_read_write_mqt_t__property_serialize(plc4c_spi_write_buffer* writeBuffer, plc4c_mqtt_read_write_mqt_t__property* _message) {
  plc4c_return_code _res = OK;

  // Simple Field (propertyType)
  _res = plc4c_mqtt_read_write_mqt_t__property_type_serialize(writeBuffer, &_message->property_type);
  if(_res != OK) {
    return _res;
  }

  // Switch Field (Depending of the current type, serialize the sub-type elements)
  switch(_message->_type) {
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__payloa_d__forma_t__indicator: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_byte(writeBuffer, 8, _message->mqt_t__property__payloa_d__forma_t__indicator_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__messag_e__expir_y__interval: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_int(writeBuffer, 32, _message->mqt_t__property__messag_e__expir_y__interval_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__conten_t__type: {

      // Simple Field (value)
      _res = plc4c_mqtt_read_write_mqt_t__string_serialize(writeBuffer, _message->mqt_t__property__conten_t__type_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__respons_e__topic: {

      // Simple Field (value)
      _res = plc4c_mqtt_read_write_mqt_t__string_serialize(writeBuffer, _message->mqt_t__property__respons_e__topic_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__correlatio_n__data: {

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__subscriptio_n__identifier: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_int(writeBuffer, 32, _message->mqt_t__property__subscriptio_n__identifier_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__expir_y__interval: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_int(writeBuffer, 32, _message->mqt_t__property__expir_y__interval_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__assigne_d__clien_t__identifier: {

      // Simple Field (value)
      _res = plc4c_mqtt_read_write_mqt_t__string_serialize(writeBuffer, _message->mqt_t__property__assigne_d__clien_t__identifier_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__serve_r__kee_p__alive: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_short(writeBuffer, 16, _message->mqt_t__property__serve_r__kee_p__alive_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__authenticatio_n__method: {

      // Simple Field (value)
      _res = plc4c_mqtt_read_write_mqt_t__string_serialize(writeBuffer, _message->mqt_t__property__authenticatio_n__method_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__authenticatio_n__data: {

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__reques_t__proble_m__information: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_byte(writeBuffer, 8, _message->mqt_t__property__reques_t__proble_m__information_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__wil_l__dela_y__interval: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_int(writeBuffer, 32, _message->mqt_t__property__wil_l__dela_y__interval_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__reques_t__respons_e__information: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_byte(writeBuffer, 8, _message->mqt_t__property__reques_t__respons_e__information_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__respons_e__information: {

      // Simple Field (value)
      _res = plc4c_mqtt_read_write_mqt_t__string_serialize(writeBuffer, _message->mqt_t__property__respons_e__information_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__serve_r__reference: {

      // Simple Field (value)
      _res = plc4c_mqtt_read_write_mqt_t__string_serialize(writeBuffer, _message->mqt_t__property__serve_r__reference_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__reaso_n__string: {

      // Simple Field (value)
      _res = plc4c_mqtt_read_write_mqt_t__string_serialize(writeBuffer, _message->mqt_t__property__reaso_n__string_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__receiv_e__maximum: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_short(writeBuffer, 16, _message->mqt_t__property__receiv_e__maximum_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__topi_c__alia_s__maximum: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_short(writeBuffer, 16, _message->mqt_t__property__topi_c__alia_s__maximum_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__topi_c__alias: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_short(writeBuffer, 16, _message->mqt_t__property__topi_c__alias_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__maximu_m__qos: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_byte(writeBuffer, 8, _message->mqt_t__property__maximu_m__qos_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__retai_n__available: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_byte(writeBuffer, 8, _message->mqt_t__property__retai_n__available_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__use_r__property: {

      // Simple Field (name)
      _res = plc4c_mqtt_read_write_mqt_t__string_serialize(writeBuffer, _message->mqt_t__property__use_r__property_name);
      if(_res != OK) {
        return _res;
      }

      // Simple Field (value)
      _res = plc4c_mqtt_read_write_mqt_t__string_serialize(writeBuffer, _message->mqt_t__property__use_r__property_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__maximu_m__packe_t__size: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_int(writeBuffer, 32, _message->mqt_t__property__maximu_m__packe_t__size_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__wildcar_d__subscriptio_n__available: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_byte(writeBuffer, 8, _message->mqt_t__property__wildcar_d__subscriptio_n__available_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__subscriptio_n__identifie_r__available: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_byte(writeBuffer, 8, _message->mqt_t__property__subscriptio_n__identifie_r__available_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__share_d__subscriptio_n__available: {

      // Simple Field (value)
      _res = plc4c_spi_write_unsigned_byte(writeBuffer, 8, _message->mqt_t__property__share_d__subscriptio_n__available_value);
      if(_res != OK) {
        return _res;
      }

      break;
    }
  }

  return OK;
}

uint16_t plc4c_mqtt_read_write_mqt_t__property_length_in_bytes(plc4c_mqtt_read_write_mqt_t__property* _message) {
  return plc4c_mqtt_read_write_mqt_t__property_length_in_bits(_message) / 8;
}

uint16_t plc4c_mqtt_read_write_mqt_t__property_length_in_bits(plc4c_mqtt_read_write_mqt_t__property* _message) {
  uint16_t lengthInBits = 0;

  // Simple field (propertyType)
  lengthInBits += plc4c_mqtt_read_write_mqt_t__property_type_length_in_bits(&_message->property_type);

  // Depending of the current type, add the length of sub-type elements ...
  switch(_message->_type) {
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__payloa_d__forma_t__indicator: {

      // Simple field (value)
      lengthInBits += 8;

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__messag_e__expir_y__interval: {

      // Simple field (value)
      lengthInBits += 32;

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__conten_t__type: {

      // Simple field (value)
      lengthInBits += plc4c_mqtt_read_write_mqt_t__string_length_in_bits(_message->mqt_t__property__conten_t__type_value);

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__respons_e__topic: {

      // Simple field (value)
      lengthInBits += plc4c_mqtt_read_write_mqt_t__string_length_in_bits(_message->mqt_t__property__respons_e__topic_value);

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__correlatio_n__data: {

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__subscriptio_n__identifier: {

      // Simple field (value)
      lengthInBits += 32;

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__expir_y__interval: {

      // Simple field (value)
      lengthInBits += 32;

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__assigne_d__clien_t__identifier: {

      // Simple field (value)
      lengthInBits += plc4c_mqtt_read_write_mqt_t__string_length_in_bits(_message->mqt_t__property__assigne_d__clien_t__identifier_value);

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__serve_r__kee_p__alive: {

      // Simple field (value)
      lengthInBits += 16;

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__authenticatio_n__method: {

      // Simple field (value)
      lengthInBits += plc4c_mqtt_read_write_mqt_t__string_length_in_bits(_message->mqt_t__property__authenticatio_n__method_value);

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__authenticatio_n__data: {

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__reques_t__proble_m__information: {

      // Simple field (value)
      lengthInBits += 8;

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__wil_l__dela_y__interval: {

      // Simple field (value)
      lengthInBits += 32;

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__reques_t__respons_e__information: {

      // Simple field (value)
      lengthInBits += 8;

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__respons_e__information: {

      // Simple field (value)
      lengthInBits += plc4c_mqtt_read_write_mqt_t__string_length_in_bits(_message->mqt_t__property__respons_e__information_value);

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__serve_r__reference: {

      // Simple field (value)
      lengthInBits += plc4c_mqtt_read_write_mqt_t__string_length_in_bits(_message->mqt_t__property__serve_r__reference_value);

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__reaso_n__string: {

      // Simple field (value)
      lengthInBits += plc4c_mqtt_read_write_mqt_t__string_length_in_bits(_message->mqt_t__property__reaso_n__string_value);

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__receiv_e__maximum: {

      // Simple field (value)
      lengthInBits += 16;

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__topi_c__alia_s__maximum: {

      // Simple field (value)
      lengthInBits += 16;

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__topi_c__alias: {

      // Simple field (value)
      lengthInBits += 16;

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__maximu_m__qos: {

      // Simple field (value)
      lengthInBits += 8;

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__retai_n__available: {

      // Simple field (value)
      lengthInBits += 8;

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__use_r__property: {

      // Simple field (name)
      lengthInBits += plc4c_mqtt_read_write_mqt_t__string_length_in_bits(_message->mqt_t__property__use_r__property_name);


      // Simple field (value)
      lengthInBits += plc4c_mqtt_read_write_mqt_t__string_length_in_bits(_message->mqt_t__property__use_r__property_value);

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__maximu_m__packe_t__size: {

      // Simple field (value)
      lengthInBits += 32;

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__wildcar_d__subscriptio_n__available: {

      // Simple field (value)
      lengthInBits += 8;

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__subscriptio_n__identifie_r__available: {

      // Simple field (value)
      lengthInBits += 8;

      break;
    }
    case plc4c_mqtt_read_write_mqt_t__property_type_plc4c_mqtt_read_write_mqt_t__property__share_d__subscriptio_n__available: {

      // Simple field (value)
      lengthInBits += 8;

      break;
    }
  }

  return lengthInBits;
}

