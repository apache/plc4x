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

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#ifndef _WIN32
#include <arpa/inet.h>
#include <netdb.h>
#include <sys/socket.h>
#else
#include <winsock.h>
#define bzero(b, len) (memset((b), '\0', (len)), (void)0)
#define MSG_DONTWAIT 0
#endif

#include "plc4x_constants.h"

#define MAX 80

void connection_func(unsigned int connFd) {
  char buff[MAX];
  int n;
  // infinite loop for chat
  for (;;) {
    bzero(buff, MAX);

    // read the message from client and copy it in buffer
    int received_bytes = recv(connFd, (char*)buff, sizeof(buff), 0);
    if (received_bytes < 0) {
      printf("Server error reading ...\n");
    }

    // print buffer which contains the client contents
    printf("From client: %s\t To client : ", buff);
    bzero(buff, MAX);
    n = 0;
    // copy server message in the buffer
    while ((buff[n++] = getchar()) != '\n')
      ;

    // and send that buffer to client
    size_t bytes_sent = send(connFd, buff, sizeof(buff), MSG_DONTWAIT);
    if (bytes_sent < 0) {
      printf("Server error sending ...\n");
    }

    // if msg contains "Exit" then server exit and chat ended.
    if (strncmp("exit", buff, 4) == 0) {
      printf("Server Exit...\n");
      break;
    }
  }
}

int main(int argc, char** argv) {
  int len;
  unsigned int sockFd, connFd;
  struct sockaddr_in serverAddress, cli;

  // socket create and verification
  sockFd = socket(AF_INET, SOCK_STREAM, 0);
  if (sockFd == -1) {
    printf("socket creation failed...\n");
    exit(0);
  } else {
    printf("Socket successfully created. Listening to TCP port \n");
  }
  bzero(&serverAddress, sizeof(serverAddress));

  // assign IP, PORT
  serverAddress.sin_family = AF_INET;
  serverAddress.sin_addr.s_addr = htonl(INADDR_ANY);
  serverAddress.sin_port =
      htons(PLC4C_PLC4X_READ_WRITE_PLC4X_CONSTANTS_PLC4X_TCP_DEFAULT_PORT());

  // Binding newly created socket to given IP and verification
  if ((bind(sockFd, (struct sockaddr*)&serverAddress, sizeof(serverAddress))) !=
      0) {
    printf("socket bind failed...\n");
    exit(0);
  } else {
    printf("Socket successfully bound..\n");
  }

  // Now server is ready to listen and verification
  if ((listen(sockFd, 5)) != 0) {
    printf("Listen failed...\n");
    exit(0);
  } else {
    printf("Server listening..\n");
  }
  len = sizeof(cli);

  // Accept the data packet from client and verification
  connFd = accept(sockFd, (struct sockaddr*)&cli, &len);
  if (connFd < 0) {
    printf("server accept failed...\n");
    exit(0);
  } else {
    printf("server accept the client...\n");
  }

  // Function for chatting between client and server
  connection_func(connFd);

  // After chatting close the socket
  shutdown(sockFd, 2);
}