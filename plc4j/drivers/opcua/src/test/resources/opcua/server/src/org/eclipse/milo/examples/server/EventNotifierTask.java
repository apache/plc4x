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

package org.eclipse.milo.examples.server;

import static org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.Unsigned.ubyte;
import static org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.Unsigned.ushort;

import java.util.UUID;
import org.eclipse.milo.opcua.sdk.server.OpcUaServer;
import org.eclipse.milo.opcua.sdk.server.model.nodes.objects.BaseEventTypeNode;
import org.eclipse.milo.opcua.sdk.server.model.nodes.objects.ServerTypeNode;
import org.eclipse.milo.opcua.sdk.server.nodes.UaNode;
import org.eclipse.milo.opcua.stack.core.Identifiers;
import org.eclipse.milo.opcua.stack.core.types.builtin.ByteString;
import org.eclipse.milo.opcua.stack.core.types.builtin.DateTime;
import org.eclipse.milo.opcua.stack.core.types.builtin.LocalizedText;
import org.eclipse.milo.opcua.stack.core.types.builtin.NodeId;
import org.eclipse.milo.opcua.stack.core.types.builtin.QualifiedName;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class EventNotifierTask implements Runnable {

  private final Logger logger = LoggerFactory.getLogger(EventNotifierTask.class);
  private final OpcUaServer server;
  private Thread eventThread;

  public EventNotifierTask(OpcUaServer server) {
    this.server = server;
  }

  @Override
  public void run() {
    UaNode serverNode = server
      .getAddressSpaceManager()
      .getManagedNode(Identifiers.Server)
      .orElse(null);

    if (serverNode instanceof ServerTypeNode) {
      ((ServerTypeNode) serverNode).setEventNotifier(ubyte(1));

      // Post a bogus Event every couple seconds
      eventThread = new Thread(() -> {
        while (!eventThread.isInterrupted()) {
          try {
            BaseEventTypeNode eventNode = server.getEventFactory().createEvent(
                new NodeId(ushort(1), UUID.randomUUID()),
                Identifiers.BaseEventType
            );

            eventNode.setBrowseName(new QualifiedName(1, "foo"));
            eventNode.setDisplayName(LocalizedText.english("foo"));
            eventNode.setEventId(ByteString.of(new byte[]{0, 1, 2, 3}));
            eventNode.setEventType(Identifiers.BaseEventType);
            eventNode.setSourceNode(serverNode.getNodeId());
            eventNode.setSourceName(serverNode.getDisplayName().getText());
            eventNode.setTime(DateTime.now());
            eventNode.setReceiveTime(DateTime.NULL_VALUE);
            eventNode.setMessage(LocalizedText.english("event message!"));
            eventNode.setSeverity(ushort(2));

            logger.info("Posting event {}", eventNode.getEventId());
            server.getEventBus().post(eventNode);

            eventNode.delete();
          } catch (Throwable e) {
            logger.error("Error creating EventNode: {}", e.getMessage(), e);
          }

          try {
            Thread.sleep(500);
          } catch (InterruptedException ignored) {
            // ignored
          }
        }
      }, "bogus-event-poster");

      eventThread.setDaemon(true);
      eventThread.start();
    }
  }
}
