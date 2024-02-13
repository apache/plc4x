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

package org.apache.plc4x.java.tools.ui.handler;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.plc4x.java.tools.ui.event.UiApplicationEvent;
import org.springframework.context.event.EventListener;
import org.springframework.stereotype.Component;
import org.springframework.web.socket.CloseStatus;
import org.springframework.web.socket.TextMessage;
import org.springframework.web.socket.WebSocketMessage;
import org.springframework.web.socket.WebSocketSession;
import org.springframework.lang.Nullable;

import java.io.IOException;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

@Component
public class WebSocketHandler implements org.springframework.web.socket.WebSocketHandler {

    private final Map<String, WebSocketSession> openSessions = new ConcurrentHashMap<>();
    private final ObjectMapper objectMapper = new ObjectMapper();

    /**
     * ApplicationEvent listener forwarding all ApplicationEvents from this application to all connected clients.
     * @param uiApplicationEvent event that we want to be forwarded.
     */
    @EventListener
    public void onApplicationEvent(UiApplicationEvent<?> uiApplicationEvent) {
        openSessions.forEach((s, webSocketSession) -> {
            try {
                webSocketSession.sendMessage(new TextMessage(objectMapper.writeValueAsString(uiApplicationEvent)));
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        });
    }

    /**
     * Register a new web-socket session.
     * @param session the new web-socket session.
     */
    @Override
    public void afterConnectionEstablished(@Nullable WebSocketSession session) {
        if(session != null) {
            openSessions.put(session.getId(), session);
        }
    }

    /**
     * Remove a web-socket session from the list.
     * @param session the web-socket session we want to remove
     * @param closeStatus the status of the closed session
     */
    @Override
    public void afterConnectionClosed(@Nullable WebSocketSession session, @Nullable CloseStatus closeStatus) {
        if((session != null) && (closeStatus != null)) {
            //noinspection resource
            openSessions.remove(session.getId());
        }
    }

    /**
     * Handle an incoming message on the web-socket session
     * (this should actually not happen, as we only use the web-sockets for server-to-client communication)
     * @param session the web-socket session the message is coming in on
     * @param message the message
     */
    @Override
    public void handleMessage(@Nullable WebSocketSession session, @Nullable WebSocketMessage<?> message) {
        System.out.println("handleMessage");
    }

    /**
     * Callback for handling transport errors.
     * @param session the web-socket session causing the error
     * @param exception the error that happened
     */
    @Override
    public void handleTransportError(@Nullable WebSocketSession session, @Nullable Throwable exception) {
        System.out.println("handleTransportError");
    }

    @Override
    public boolean supportsPartialMessages() {
        return false;
    }

}
