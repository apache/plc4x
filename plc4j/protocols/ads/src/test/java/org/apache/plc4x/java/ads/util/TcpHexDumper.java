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
package org.apache.plc4x.java.ads.util;

import org.apache.commons.io.HexDump;
import org.apache.commons.io.IOUtils;
import org.junit.jupiter.api.extension.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.net.ServerSocket;
import java.net.Socket;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

public class TcpHexDumper implements BeforeEachCallback, AfterEachCallback, ParameterResolver {

    private static final Logger logger = LoggerFactory.getLogger(TcpHexDumper.class);

    private ExecutorService pool = Executors.newCachedThreadPool();

    private Integer portToUse;

    private ServerSocket serverSocket;

    int shutdownTimeout = 10;

    public void init(int port) throws IOException, InterruptedException {
        if (serverSocket != null) {
            stop();
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
                throw new RuntimeException(e);
            }
            pool.submit(() -> {
                InputStream inputStream;
                try {
                    inputStream = accept.getInputStream();
                    logger.info("Starting to read now");
                } catch (IOException e) {
                    throw new RuntimeException(e);
                }
                byte[] buffer = new byte[4096];
                try {
                    while (IOUtils.read(inputStream, buffer) > 0) {
                        ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
                        HexDump.dump(buffer, 0, byteArrayOutputStream, 0);
                        logger.info("Dump:\n{}", byteArrayOutputStream);
                    }
                } catch (IOException e) {
                    throw new RuntimeException(e);
                }
            });
        });
        logger.info("Started pool");
    }

    public void stop() throws IOException, InterruptedException {
        serverSocket.close();
        pool.awaitTermination(shutdownTimeout, TimeUnit.SECONDS);
        logger.info("Stopped");
    }

    @Override
    public void afterEach(ExtensionContext context) throws Exception {
        stop();
    }

    @Override
    public void beforeEach(ExtensionContext context) throws Exception {
        init(initPortToUse(context));
        shutdownTimeout = initShutdownTimeout(context);
    }

    @Override
    public boolean supportsParameter(ParameterContext parameterContext, ExtensionContext extensionContext) throws ParameterResolutionException {
        return parameterContext.getParameter().getType()
            .equals(int.class);
    }

    @Override
    public Object resolveParameter(ParameterContext parameterContext, ExtensionContext extensionContext) throws ParameterResolutionException {
        try {
            return initPortToUse(extensionContext);
        } catch (IOException e) {
            throw new ParameterResolutionException("Could not find a free port", e);
        }
    }

    private int initShutdownTimeout(ExtensionContext context) {
        ExtendWithTcpHexDumper annotation = context.getRequiredTestClass().getAnnotation(ExtendWithTcpHexDumper.class);
        return annotation.shutdownTimeout();
    }

    private int initPortToUse(ExtensionContext context) throws IOException {
        if (portToUse == null) {
            ExtendWithTcpHexDumper annotation = context.getRequiredTestClass().getAnnotation(ExtendWithTcpHexDumper.class);
            int port = annotation.port();
            portToUse = port != 0 ? port : findFreePort();
        }
        return portToUse;
    }

    private static int findFreePort() throws IOException {
        try (ServerSocket socket = new ServerSocket(0)) {
            socket.setReuseAddress(true);
            return socket.getLocalPort();
        }
    }
}
