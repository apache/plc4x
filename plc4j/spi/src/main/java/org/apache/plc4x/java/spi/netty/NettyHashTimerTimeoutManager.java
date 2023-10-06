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
package org.apache.plc4x.java.spi.netty;

import io.netty.util.HashedWheelTimer;
import io.netty.util.Timeout;
import io.netty.util.Timer;
import java.util.Set;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import org.apache.plc4x.java.spi.TimedOperation;
import org.apache.plc4x.java.spi.TimeoutManager;

public class NettyHashTimerTimeoutManager implements TimeoutManager {

    private final Timer timer;

    public NettyHashTimerTimeoutManager() {
        this(100L);
    }

    /**
     * Creates a new NettyHashTimerTimeoutManager that checks for timeouts every `tickInMilliseconds` milliseconds.
     * @param tickInMilliseconds milliseconds between timeout checks.
     */
    public NettyHashTimerTimeoutManager(long tickInMilliseconds) {
        HashedWheelTimer wheelTimer = new HashedWheelTimer(tickInMilliseconds, TimeUnit.MILLISECONDS);
        timer = wheelTimer;
        wheelTimer.start();
    }

    @Override
    public CompletionCallback<?> register(TimedOperation operation) {
        Timeout newTimeout = timer.newTimeout(timeout -> {
            if (timeout.isCancelled()) {
                return;
            }
            TimeoutException exception = new TimeoutException();
            operation.getOnTimeoutConsumer().accept(exception);
        }, operation.getTimeout().toMillis(), TimeUnit.MILLISECONDS);

        return new TimeoutCompletionCallback<>(newTimeout);
    }

    @Override
    public void stop() {
        Set<Timeout> timeouts = timer.stop();
        timeouts.forEach(Timeout::cancel);
    }

    static class TimeoutCompletionCallback<T> implements CompletionCallback<T> {

        private final Timeout timeout;

        TimeoutCompletionCallback(Timeout timeout) {
            this.timeout = timeout;
        }

        @Override
        public void complete() {
            timeout.cancel();
        }
    }
}
