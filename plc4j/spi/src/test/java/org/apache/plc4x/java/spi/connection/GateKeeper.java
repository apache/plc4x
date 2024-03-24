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

package org.apache.plc4x.java.spi.connection;

import java.util.concurrent.CountDownLatch;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Utility class to control entry and exit points.
 */
class GateKeeper {

    private final Logger logger = LoggerFactory.getLogger(GateKeeper.class);
    private final String gate;
    private final CountDownLatch in = new CountDownLatch(1);
    private final CountDownLatch out = new CountDownLatch(1);

    GateKeeper(String gate) {
        this.gate = gate;
    }

    boolean awaitIn() throws InterruptedException {
        logger.debug("Awaiting entry permit for {}", gate);
        in.await();
        return true;
    }

    boolean awaitOut() throws InterruptedException {
        logger.debug("Awaiting exit permit for {}", gate);
        out.await();
        return true;
    }

    boolean entered() {
        return in.getCount() == 0;
    }

    boolean exited() {
        return out.getCount() == 0;
    }

    void permitIn() {
        logger.info("Allowing permit for {}", gate);
        in.countDown();
    }

    public void permitOut() {
        logger.info("Allowing exit for {}", gate);
        out.countDown();
    }

    public String gate() {
        return gate;
    }

    @Override
    public String toString() {
        return "GateKeeper [" + gate + ", entered=" + entered() + ", exited=" + exited() + "]";
    }
}