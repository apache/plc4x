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
package org.apache.plc4x.java.mock.connection.tcp;

import org.apache.commons.io.HexDump;
import org.apache.commons.io.IOUtils;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayOutputStream;
import java.io.Closeable;
import java.io.IOException;
import java.io.InputStream;
import java.net.ServerSocket;
import java.net.Socket;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

public class TcpHexDumper implements Closeable {

    private static final Logger logger = LoggerFactory.getLogger(TcpHexDumper.class);

    private final ExecutorService pool;

    private final Integer portToUse;

    private final Integer shutdownTimeout;

    private ServerSocket serverSocket;

    private TcpHexDumper(Integer portToUse) {
        this(portToUse, 10);
    }

    public TcpHexDumper(Integer portToUse, Integer shutdownTimeout) {
        pool = Executors.newCachedThreadPool();
        this.portToUse = portToUse;
        this.shutdownTimeout = shutdownTimeout;
    }

    private void init(int port) throws IOException, InterruptedException {
        if (serverSocket != null) {
            stop(true);
        }
        serverSocket = new ServerSocket(port);
        logger.info("Starting pool");
        pool.submit(() -> {
            Socket accept;
            try {
                logger.info("Waiting for an incoming connection");
                accept = serverSocket.accept();
                logger.info("Accepted {} and starting listener", accept);
            } catch (IOException e) {
                throw new PlcRuntimeException(e);
            }
            pool.submit(() -> {
                InputStream inputStream;
                try {
                    inputStream = accept.getInputStream();
                    logger.info("Starting to read now");
                } catch (IOException e) {
                    throw new PlcRuntimeException(e);
                }
                byte[] buffer = new byte[4096];
                try {
                    while (IOUtils.read(inputStream, buffer) > 0) {
                        ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
                        HexDump.dump(buffer, 0, byteArrayOutputStream, 0);
                        logger.info("Dump:\n{}", byteArrayOutputStream);
                    }
                } catch (IOException e) {
                    throw new PlcRuntimeException(e);
                }
            });
        });
        logger.info("Started pool");
    }

    private void stop() throws IOException, InterruptedException {
        stop(false);
    }

    private void stop(boolean await) throws IOException, InterruptedException {
        serverSocket.close();
        if (await) {
            pool.awaitTermination(shutdownTimeout, TimeUnit.SECONDS);
        } else {
            pool.shutdownNow();
        }
        logger.info("Stopped");
    }

    public void after() {
        try {
            stop(true);
        } catch (InterruptedException e) {
            logger.info("Shutdown error", e);
            Thread.currentThread().interrupt();
        } catch (IOException e) {
            logger.info("Shutdown error", e);
        }
    }

    public void before() throws Throwable {
        init(portToUse);
    }

    public static TcpHexDumper runOn(int port) throws IOException, InterruptedException {
        TcpHexDumper tcpHexDumper = new TcpHexDumper(port);
        tcpHexDumper.init(port);
        return tcpHexDumper;
    }

    public Integer getPort() {
        return serverSocket.getLocalPort();
    }

    @Override
    public void close() throws IOException {
        try {
            stop();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IOException(e);
        }
    }
}
