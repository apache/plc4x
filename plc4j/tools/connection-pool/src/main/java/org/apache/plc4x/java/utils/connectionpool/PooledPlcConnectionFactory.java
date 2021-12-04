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
package org.apache.plc4x.java.utils.connectionpool;

import org.apache.commons.pool2.BaseKeyedPooledObjectFactory;
import org.apache.commons.pool2.PooledObject;
import org.apache.commons.pool2.impl.DefaultPooledObject;
import org.apache.plc4x.java.api.PlcConnection;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public abstract class PooledPlcConnectionFactory extends BaseKeyedPooledObjectFactory<PoolKey, PlcConnection> {

    private static final Logger LOGGER = LoggerFactory.getLogger(PooledPlcConnectionFactory.class);

    @Override
    public PooledObject<PlcConnection> wrap(PlcConnection plcConnection) {
        LOGGER.debug("Wrapping connection {}", plcConnection);
        return new DefaultPooledObject<>(plcConnection);
    }

    @Override
    public void destroyObject(PoolKey key, PooledObject<PlcConnection> p) throws Exception {
        p.getObject().close();
    }

    @Override
    public boolean validateObject(PoolKey key, PooledObject<PlcConnection> p) {
        return p.getObject().isConnected();
    }
}
