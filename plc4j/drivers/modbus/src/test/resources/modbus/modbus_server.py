#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

"""
Pymodbus TCP server for Modbus driver integration tests.

Pre-populates all 4 Modbus address spaces (coils, discrete inputs,
input registers, holding registers) with known test values that match
the register map defined in specs/047-pymodbus-integration-tests/data-model.md.

The server listens on port 5020 and prints "Modbus server started" when ready.
"""

import asyncio
import logging

from pymodbus.datastore import (
    ModbusSequentialDataBlock,
    ModbusSparseDataBlock,
    ModbusSlaveContext,
    ModbusServerContext,
)
import ssl

from pymodbus.framer import FramerType
from pymodbus.server import StartAsyncTcpServer, StartAsyncTlsServer, StartAsyncUdpServer

logging.basicConfig(level=logging.INFO)
log = logging.getLogger(__name__)

# Server ports — one per framing/transport mode
TCP_PORT = 5020
RTU_PORT = 5021
ASCII_PORT = 5022
TLS_PORT = 5023
UDP_PORT = 5025
UDP_RTU_PORT = 5026
UDP_ASCII_PORT = 5027

# Certificate paths (generated during Docker build)
CERT_DIR = "/app/certs"
SERVER_CERT = f"{CERT_DIR}/server.crt"
SERVER_KEY = f"{CERT_DIR}/server.key"
CA_CERT = f"{CERT_DIR}/ca.crt"

# PSK identity/key for TLS-PSK testing
PSK_IDENTITY = "modbus-test"
PSK_KEY = "0123456789abcdef0123456789abcdef"  # 32 hex chars = 16 bytes


def build_holding_registers():
    """
    Build holding register data block with pre-populated test values.

    Pymodbus 3.x applies a +1 offset to incoming addresses (no zero_mode),
    so all values are shifted right by 1 position. The driver sends wire
    address N-1 for tag address N; pymodbus then accesses index N in the array.

    Addresses 1-31: Read fixtures (matching existing test-simulation.xml values).
    Addresses 100-200: Write target area (initially zero).

    All values are unsigned 16-bit integers. Multi-register types are stored
    as consecutive big-endian UINT16 values.
    """
    # Start with 201 zero-initialized registers
    # Pymodbus applies a +1 offset: wire address N maps to data block index N+1.
    # So tag address T (wire addr T-1) is stored at index T.
    values = [0] * 201

    # BOOL at index 1 (tag: holding-register:1:BOOL) -> true (1)
    values[1] = 1

    # BYTE at index 2 (tag: holding-register:2:BYTE) -> 42
    values[2] = 42

    # WORD at index 3 (tag: holding-register:3:WORD) -> 42424
    values[3] = 42424

    # DWORD at indices 4-5 (tag: holding-register:4:DWORD) -> 4242442424
    values[4] = 64734  # 0xFCDE
    values[5] = 35000  # 0x88B8

    # LWORD at indices 6-9 (tag: holding-register:6:LWORD) -> 4242442424242424242
    values[6] = 15072   # 0x3AE0
    values[7] = 12008   # 0x2EE8
    values[8] = 19716   # 0x4D04
    values[9] = 18866   # 0x49B2

    # SINT at index 10 (tag: holding-register:10:SINT) -> -42
    values[10] = 65494  # 0xFFD6

    # USINT at index 11 (tag: holding-register:11:USINT) -> 42
    values[11] = 42

    # INT at index 12 (tag: holding-register:12:INT) -> -2424
    values[12] = 63112  # 0xF688

    # UINT at index 13 (tag: holding-register:13:UINT) -> 42424
    values[13] = 42424  # 0xA5B8

    # DINT at indices 14-15 (tag: holding-register:14:DINT) -> -242442424
    values[14] = 61836  # 0xF18C
    values[15] = 40776  # 0x9F48

    # UDINT at indices 16-17 (tag: holding-register:16:UDINT) -> 4242442424
    values[16] = 64734  # 0xFCDE
    values[17] = 35000  # 0x88B8

    # LINT at indices 18-21 (tag: holding-register:18:LINT) -> -4242442424242424242
    values[18] = 50463  # 0xC51F
    values[19] = 53527  # 0xD117
    values[20] = 45819  # 0xB2FB
    values[21] = 46670  # 0xB64E

    # ULINT at indices 22-25 (tag: holding-register:22:ULINT) -> 4242442424242424242
    values[22] = 15072  # 0x3AE0
    values[23] = 12008  # 0x2EE8
    values[24] = 19716  # 0x4D04
    values[25] = 18866  # 0x49B2

    # REAL at indices 26-27 (tag: holding-register:26:REAL) -> 3.141593
    values[26] = 16457  # 0x4049
    values[27] = 4060   # 0x0FDC

    # LREAL at indices 28-31 (tag: holding-register:28:LREAL) -> 2.71828182846
    values[28] = 16389  # 0x4005
    values[29] = 48906  # 0xBF0A
    values[30] = 35604  # 0x8B14
    values[31] = 24527  # 0x5FCF

    return ModbusSequentialDataBlock(0, values)


def build_coils():
    """
    Build coil data block with pre-populated test values.

    Pymodbus 3.x applies a +1 offset, so values are shifted right by 1.
    Indices 1-10: Read fixtures with pattern [T, F, T, T, F, T, F, F, T, F].
    Indices 100-120: Write target area (initially FALSE).
    """
    # Pymodbus accesses internal address 1 for wire address 0 (+1 offset).
    # Start the block at address 1 so it covers addresses 1 through N.
    values = [0] * 2000
    # Read fixture pattern (index 0 = internal address 1 = coil:1)
    values[0] = 1    # coil:1 -> TRUE
    values[1] = 0    # coil:2 -> FALSE
    values[2] = 1    # coil:3 -> TRUE
    values[3] = 1    # coil:4 -> TRUE
    values[4] = 0    # coil:5 -> FALSE
    values[5] = 1    # coil:6 -> TRUE
    values[6] = 0    # coil:7 -> FALSE
    values[7] = 0    # coil:8 -> FALSE
    values[8] = 1    # coil:9 -> TRUE
    values[9] = 0    # coil:10 -> FALSE
    return ModbusSequentialDataBlock(1, values)


def build_discrete_inputs():
    """
    Build discrete input data block with pre-populated test values.

    Pymodbus 3.x applies a +1 offset, so values are shifted right by 1.
    Indices 1-10: Read-only fixtures with pattern [F, T, F, T, T, F, T, T, F, T].
    This is the inverse of the coils pattern.
    """
    values = [0] * 2000
    # Index 0 = internal address 1 = discrete-input:1
    values[0] = 0    # discrete-input:1 -> FALSE
    values[1] = 1    # discrete-input:2 -> TRUE
    values[2] = 0    # discrete-input:3 -> FALSE
    values[3] = 1    # discrete-input:4 -> TRUE
    values[4] = 1    # discrete-input:5 -> TRUE
    values[5] = 0    # discrete-input:6 -> FALSE
    values[6] = 1    # discrete-input:7 -> TRUE
    values[7] = 1    # discrete-input:8 -> TRUE
    values[8] = 0    # discrete-input:9 -> FALSE
    values[9] = 1    # discrete-input:10 -> TRUE
    return ModbusSequentialDataBlock(1, values)


def build_input_registers():
    """
    Build input register data block with pre-populated test values.

    Pymodbus 3.x applies a +1 offset, so values are shifted right by 1.
    Indices 1-7: Read-only fixtures for representative data types.
    """
    values = [0] * 21

    # Same +1 offset as holding registers: tag address T stored at index T.
    # BOOL at index 1 (tag: input-register:1:BOOL) -> true (1)
    values[1] = 1

    # INT at index 2 (tag: input-register:2:INT) -> -2424
    values[2] = 63112  # 0xF688

    # UINT at index 3 (tag: input-register:3:UINT) -> 42424
    values[3] = 42424  # 0xA5B8

    # DINT at indices 4-5 (tag: input-register:4:DINT) -> -242442424
    values[4] = 61836  # 0xF18C
    values[5] = 40776  # 0x9F48

    # REAL at indices 6-7 (tag: input-register:6:REAL) -> 3.141593
    values[6] = 16457  # 0x4049
    values[7] = 4060   # 0x0FDC

    return ModbusSequentialDataBlock(0, values)


def build_server_context():
    """Build a shared server context with all 4 register types pre-populated."""
    slave_context = ModbusSlaveContext(
        di=build_discrete_inputs(),  # Discrete Inputs (FC2)
        co=build_coils(),            # Coils (FC1, FC5, FC15)
        hr=build_holding_registers(),  # Holding Registers (FC3, FC6, FC16)
        ir=build_input_registers(),  # Input Registers (FC4)
    )
    return ModbusServerContext(slaves=slave_context, single=True)


def build_tls_context():
    """Build an SSL context for the TLS server using the generated certificates."""
    ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
    ctx.load_cert_chain(certfile=SERVER_CERT, keyfile=SERVER_KEY)
    ctx.load_verify_locations(cafile=CA_CERT)
    # Don't require client certificates for testing
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE
    return ctx


async def run_server():
    """Start Modbus servers for all framing and transport modes.

    Each mode runs on a separate TCP port with the same pre-populated data.
    Pymodbus shares a single data store across all servers, so writes on one
    port are visible on the others.
    """
    context = build_server_context()

    log.info("Starting Modbus TCP   server on port %d", TCP_PORT)
    log.info("Starting Modbus RTU   server on port %d", RTU_PORT)
    log.info("Starting Modbus ASCII server on port %d", ASCII_PORT)
    log.info("Starting Modbus TLS   server on port %d", TLS_PORT)
    log.info("Starting Modbus UDP   server on port %d", UDP_PORT)
    log.info("Starting Modbus UDP-RTU   server on port %d", UDP_RTU_PORT)
    log.info("Starting Modbus UDP-ASCII server on port %d", UDP_ASCII_PORT)

    tasks = []

    # Plain TCP with MBAP framing
    tasks.append(asyncio.create_task(StartAsyncTcpServer(
        context=context,
        address=("0.0.0.0", TCP_PORT),
        framer=FramerType.SOCKET,
    )))

    # RTU over TCP
    tasks.append(asyncio.create_task(StartAsyncTcpServer(
        context=context,
        address=("0.0.0.0", RTU_PORT),
        framer=FramerType.RTU,
    )))

    # ASCII over TCP
    tasks.append(asyncio.create_task(StartAsyncTcpServer(
        context=context,
        address=("0.0.0.0", ASCII_PORT),
        framer=FramerType.ASCII,
    )))

    # TLS with MBAP framing (certificate-based)
    tls_ctx = build_tls_context()
    tasks.append(asyncio.create_task(StartAsyncTlsServer(
        context=context,
        address=("0.0.0.0", TLS_PORT),
        framer=FramerType.SOCKET,
        sslctx=tls_ctx,
    )))

    # UDP with MBAP framing
    tasks.append(asyncio.create_task(StartAsyncUdpServer(
        context=context,
        address=("0.0.0.0", UDP_PORT),
        framer=FramerType.SOCKET,
    )))

    # RTU over UDP
    tasks.append(asyncio.create_task(StartAsyncUdpServer(
        context=context,
        address=("0.0.0.0", UDP_RTU_PORT),
        framer=FramerType.RTU,
    )))

    # ASCII over UDP
    tasks.append(asyncio.create_task(StartAsyncUdpServer(
        context=context,
        address=("0.0.0.0", UDP_ASCII_PORT),
        framer=FramerType.ASCII,
    )))

    # Ready signal for Testcontainers wait strategy
    print("Modbus server started", flush=True)

    await asyncio.gather(*tasks)


if __name__ == "__main__":
    asyncio.run(run_server())
