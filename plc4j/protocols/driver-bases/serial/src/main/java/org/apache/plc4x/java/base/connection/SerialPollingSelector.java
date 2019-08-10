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

package org.apache.plc4x.java.base.connection;

import org.apache.commons.lang3.NotImplementedException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.nio.channels.SelectionKey;
import java.nio.channels.Selector;
import java.nio.channels.spi.AbstractSelectableChannel;
import java.nio.channels.spi.AbstractSelector;
import java.nio.channels.spi.SelectorProvider;
import java.util.HashSet;
import java.util.Set;

/**
 * TODO write comment
 *
 * @author julian
 * Created by julian on 2019-08-10
 */
class SerialPollingSelector extends AbstractSelector {

    private static final Logger logger = LoggerFactory.getLogger(SerialPollingSelector.class);

    public SerialPollingSelector(SelectorProvider selectorProvider) {
        super(selectorProvider);
    }

    @Override
    public Set<SelectionKey> keys() {
        return new HashSet<>();
    }

    @Override
    public Set<SelectionKey> selectedKeys() {
        return new HashSet<>();
    }

    @Override
    public int selectNow() throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public int select(long timeout) throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public int select() throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public Selector wakeup() {
        // TODO do we have to do something here?
        return this;
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
        throw new NotImplementedException("");
    }

}
