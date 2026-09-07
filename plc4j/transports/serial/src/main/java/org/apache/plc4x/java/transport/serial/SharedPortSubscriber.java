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
package org.apache.plc4x.java.transport.serial;

/**
 * Receives the broadcast of everything a shared serial port reads. Every
 * subscriber sees the FULL wire traffic (protocol codecs pick out their own
 * frames), mirroring the plc4go shared-port design.
 */
interface SharedPortSubscriber {

    /** A chunk read from the physical port. Called from the reader thread. */
    void onData(byte[] data, int offset, int length);

    /** The physical port failed fatally; the shared port has been evicted. */
    void onFailure(Throwable cause);
}
