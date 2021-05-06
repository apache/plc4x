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

#include <plc4c/spi/types_private.h>
#include <plc4c/transport_tcp.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#ifndef _WIN32
#include <netdb.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#else
#include <winsock.h>
#define strtok_r strtok_s
#define bzero(b,len) (memset((b), '\0', (len)), (void) 0)
#define MSG_DONTWAIT 0
#endif

extern int errno;

plc4c_return_code plc4c_transport_tcp_configure_function(
    char* transport_connect_information, plc4c_list* parameters, void** configuration) {
  plc4c_transport_tcp_config* tcp_configuration = malloc(sizeof(plc4c_transport_tcp_config));
  if (tcp_configuration == NULL) {
    return NO_MEMORY;
  }

  char *port;
  char *host = strtok_r(transport_connect_information, ":", &port);
  tcp_configuration->address = host;
  // If no port was specified, generally use the default port for this driver
  if (strlen(port) == 0) {
    // TODO: Currently return an error.
    return INTERNAL_ERROR;
  } else {
    tcp_configuration->port = atoi(port);
  }

  *configuration = tcp_configuration;
  return OK;
}

plc4c_return_code plc4c_transport_tcp_open_function(void* config) {
  int sockfd;
  int connfd;
  struct sockaddr_in servaddr;

  plc4c_transport_tcp_config* tcp_config = config;

#ifdef _WIN32
  WSADATA wsa;
  int wsa_res = WSAStartup(MAKEWORD(2,2), &wsa);
  // Something happened when initializing the WinSock API usage
  if (wsa_res != 0) {
    return INTERNAL_ERROR;
  }
#endif
  tcp_config->sockfd = socket(AF_INET, SOCK_STREAM, 0);
  if (tcp_config->sockfd < 0) {
    return CONNECTION_ERROR;
  }

  // Configure where to connect to.
  bzero(&servaddr, sizeof(struct sockaddr_in));
  servaddr.sin_family = AF_INET;
  servaddr.sin_addr.s_addr = inet_addr(tcp_config->address);
  servaddr.sin_port = htons(tcp_config->port);

  int result = connect(tcp_config->sockfd, (struct sockaddr*) &servaddr,
                       sizeof(servaddr));
  if(result != 0) {
    char* error_msg = strerror(errno);
    printf("%s\n", error_msg);
    return CONNECTION_ERROR;
  }
  return OK;
}

plc4c_return_code plc4c_transport_tcp_close_function(void* config) {
  plc4c_transport_tcp_config* tcp_config = config;

  // If the sockfd is zero we're not really connected.
  if(tcp_config->sockfd == 0) {
    return INTERNAL_ERROR;
  }

  // Stop receiving as well as sending.
  shutdown(tcp_config->sockfd, 2);
  return OK;
}

plc4c_return_code plc4c_transport_tcp_send_message_function(
    void* transport_configuration, plc4c_spi_write_buffer* message) {
  plc4c_transport_tcp_config* tcp_config = transport_configuration;

  size_t bytes_sent = send(tcp_config->sockfd, message->data, message->length, MSG_DONTWAIT);
  if(bytes_sent < 0) {
    return CONNECTION_ERROR;
  }

  return OK;
}

plc4c_return_code plc4c_transport_tcp_select_message_function(
    void* transport_configuration, uint8_t min_size,
    accept_message_function accept_message, plc4c_spi_read_buffer** message) {
      
  plc4c_transport_tcp_config* tcp_config = transport_configuration;

  // First try to read the minimum number of bytes the driver needs to know
  // how big a packet is.
  uint8_t* size_buffer = malloc(sizeof(uint8_t) * min_size);
  if(size_buffer == NULL) {
    return NO_MEMORY;
  }
  int received_bytes = recv(tcp_config->sockfd, size_buffer, min_size, 0);
  // TODO: if the value is negative, it's more a "please remove this much of corrupt data" ...
  if(received_bytes < 0) {
    return CONNECTION_ERROR;
  }

  // Now that we have enough data to find out how many bytes we need, have
  // the accept_message function find out how many
  int16_t message_size = accept_message(size_buffer, min_size);
  if(message_size < 0) {
    return INTERNAL_ERROR;
  }
  uint8_t* message_buffer = malloc(sizeof(uint8_t) * message_size);
  if(message_size < 0) {
    return NO_MEMORY;
  }

  // Copy the size_buffer to the start of the new buffer.
  memcpy(message_buffer, size_buffer, min_size);
  free(size_buffer);

  // Read the rest of the packet.
  received_bytes = recv(tcp_config->sockfd, message_buffer + min_size, message_size - min_size, 0);
  if(received_bytes != message_size - min_size) {
    return CONNECTION_ERROR;
  }

  // Create a new read-buffer with the read message data.
  plc4c_spi_read_buffer_create(message_buffer, message_size, message);
  // TODO: leaks: message_buffer
  return OK;
}

plc4c_transport *plc4c_transport_tcp_create() {
  plc4c_transport *transport =
      (plc4c_transport *)malloc(sizeof(plc4c_transport));
  transport->transport_code = "tcp";
  transport->transport_name = "TCP transport";
  transport->configure = &plc4c_transport_tcp_configure_function;
  transport->open = &plc4c_transport_tcp_open_function;
  transport->close = &plc4c_transport_tcp_close_function;
  transport->send_message = &plc4c_transport_tcp_send_message_function;
  transport->select_message = &plc4c_transport_tcp_select_message_function;
  return transport;
}
