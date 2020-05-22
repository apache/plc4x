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

package org.apache.plc4x.simulator.server.s7;

/**
 * Handler for PLC Server.
 *
 * @author julian
 * Created by julian on 22.05.20
 */
public interface S7PlcHandler {

    void onConnectionInitiated();

    void onConnectionEstablished();

    void onConnectionClosed();

    S7Int readIntFromDataBlock(int dbNumber, int byteAddress, byte bitAddress);

    class S7Int {
        private final Short signed;
        private final Integer unsigned;

        public S7Int(short signed) {
            this.signed = signed;
            this.unsigned = null;
        }

        public S7Int(int unsigned) {
            if (unsigned < 0) {
                throw new IllegalArgumentException("Signed value cannot be negative!");
            }
            this.unsigned = unsigned;
            this.signed = null;
        }

        public static S7Int _int(short signed) {
            return new S7Int(signed);
        }

        public static S7Int _uint(int unsigned) {
            return new S7Int(unsigned);
        }

        public boolean isSigned() {
            return signed != null;
        }

        public boolean isUnsigned() {
            return unsigned != null;
        }

        public short getSigned() {
            if (!isSigned()) {
                throw new UnsupportedOperationException("Cannot get signed on unsigned!");
            }
            return signed;
        }

        public Integer getUnsigned() {
            if (!isUnsigned()) {
                throw new UnsupportedOperationException("Cannot get unsigned on signed!");
            }
            return unsigned;
        }

        @Override
        public String toString() {
            if (isSigned()) {
                return "Signed(" + signed + ")";
            } else {
                return "Unsigned(" + unsigned + ")";
            }
        }
    }

}
