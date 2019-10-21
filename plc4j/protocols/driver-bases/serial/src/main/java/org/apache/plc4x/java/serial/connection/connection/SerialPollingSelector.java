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

package org.apache.plc4x.java.serial.connection.connection;

import io.netty.util.concurrent.DefaultEventExecutor;
import io.netty.util.concurrent.DefaultPromise;
import org.apache.commons.lang3.NotImplementedException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.nio.channels.SelectionKey;
import java.nio.channels.Selector;
import java.nio.channels.spi.AbstractSelectableChannel;
import java.nio.channels.spi.AbstractSelector;
import java.nio.channels.spi.SelectorProvider;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;
import java.util.stream.Collectors;

/**
 * TODO write comment
 *
 * @author julian
 * Created by julian on 2019-08-10
 */
class SerialPollingSelector extends AbstractSelector {

    private static final Logger logger = LoggerFactory.getLogger(SerialPollingSelector.class);

    private final List<SelectionKey> registeredChannels;
    private final Set<SelectorEvent> events = ConcurrentHashMap.newKeySet();

    // Use a Netty Promise
    private final DefaultEventExecutor executor = new DefaultEventExecutor();
    private DefaultPromise<Void> selectPromise;

    public static class SelectorEvent {

        private final SelectionKey key;

        private final int event;
        public SelectorEvent(SelectionKey key, int event) {
            this.key = key;
            this.event = event;
        }

        public SelectionKey getKey() {
            return this.key;
        }

        public int getEvent() {
            return event;
        }

    }
    public SerialPollingSelector(SelectorProvider selectorProvider) {
        super(selectorProvider);
        registeredChannels = new ArrayList<>();
    }

    @Override
    public Set<SelectionKey> keys() {
        return new HashSet<>(registeredChannels);
    }

    /**
     * Returns all keys that are in the events queue
     * @return
     */
    @Override
    public Set<SelectionKey> selectedKeys() {
        return events.stream().map(SelectorEvent::getKey).collect(Collectors.toSet());
    }

    @Override
    public int selectNow() throws IOException {
        // throw new NotImplementedException("");
        logger.debug("selectNow()");
        // check if one channel is active
        return events.size();
    }

    @Override
    public int select(long timeout) throws IOException {
        logger.debug("select({})", timeout);
        if (events.size() > 0) {
            return events.size();
        }
        this.selectPromise = new DefaultPromise<>(executor);
        try {
            if (selectPromise.await(timeout)) {
                logger.debug("Promise was cancelled, new Events should be there.");
            } else {
                logger.debug("Promise timed out, expecting no new events.");
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new RuntimeException("Was interrupted", e);
        }
        return events.size();
    }

    @Override
    public int select() throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public Selector wakeup() {
        // throw new NotImplementedException("Not implemented for this selector, should not be needed.");
        // NOOP
        return this;
    }

    public void addEvent(SelectorEvent event) {
        logger.debug("Adding Event to Selector, canceling Promise...");
        this.events.add(event);
        // Add the OP to the SelectionKey
        ((SerialSelectionKey) event.key).addReadyOp(event.event);
        // Close the future so that the select is fired imediatly
        if (!selectPromise.isDone()) {
            selectPromise.setSuccess(null);
        } else {
            logger.debug("Promise is already cancelled, skipping that.");
        }
    }

    public void removeEvent(SerialSelectionKey serialSelectionKey) {
        events.removeIf(event -> event.key.equals(serialSelectionKey));
    }

    @Override
    protected void implCloseSelector() throws IOException {
        // TODO should we do something here?
    }

    @Override
    protected SelectionKey register(AbstractSelectableChannel ch, int ops, Object att) {
        logger.debug("Registering Channel for selector {} with operations {}", ch, ops);
        if (!(ch instanceof SerialSocketChannel)) {
            throw new IllegalArgumentException("Given channel has to be of type " + SerialSocketChannel.class);
        }
        final SerialSelectionKey key = new SerialSelectionKey(ch, this, ops);
        // Attach attr
        key.attach(att);
        synchronized (this) {
            // TODO is this always the case??
            final int index = registeredChannels.size();
            registeredChannels.add(key);
            key.setIndex(index);
        }
        return key;
    }

}
