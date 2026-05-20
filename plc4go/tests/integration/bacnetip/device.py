#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

"""bacpypes3-based BACnet/IP virtual device used by plc4go integration tests.

Exposes:
    DEVICE:1234              objectName="plc4x-it", vendor 0x4D4D (Apache PLC4X)
    ANALOG_VALUE:0..4        writable Real PresentValue
    BINARY_VALUE:0..1        writable Enumerated PresentValue
    ANALOG_INPUT:0           read-only Real PresentValue with a slow sawtooth
                             generator so SubscribeCOV tests see notifications
    MULTI_STATE_VALUE:0      writable Unsigned PresentValue

The device binds to 0.0.0.0:47808 inside the container; docker-compose maps
that to the host's 47808/udp.
"""

import asyncio

from bacpypes3.app import Application
from bacpypes3.local.analog import AnalogInputObject, AnalogValueObject
from bacpypes3.local.binary import BinaryValueObject
from bacpypes3.local.device import DeviceObject
from bacpypes3.local.multistate import MultiStateValueObject
from bacpypes3.primitivedata import Real


async def sawtooth(obj: AnalogInputObject) -> None:
    """Background coroutine that bumps AnalogInput.0 by 1.0 every 2 seconds so
    SubscribeCOV consumers see a stream of notifications."""
    value = 0.0
    while True:
        await asyncio.sleep(2.0)
        value = (value + 1.0) % 100.0
        obj.presentValue = Real(value)


def build_app() -> Application:
    device = DeviceObject(
        objectIdentifier=("device", 1234),
        objectName="plc4x-it",
        vendorIdentifier=0x4D4D,
        modelName="plc4x-it-simulator",
        maxApduLengthAccepted=1476,
        segmentationSupported="segmentedBoth",
        maxSegmentsAccepted=16,
        protocolVersion=1,
        protocolRevision=14,
    )

    app = Application(device, ("0.0.0.0/24", 47808))

    for i in range(5):
        app.add_object(
            AnalogValueObject(
                objectIdentifier=("analogValue", i),
                objectName=f"AV-{i}",
                presentValue=Real(0.0),
                outOfService=False,
            )
        )

    for i in range(2):
        app.add_object(
            BinaryValueObject(
                objectIdentifier=("binaryValue", i),
                objectName=f"BV-{i}",
                presentValue="inactive",
            )
        )

    ai0 = AnalogInputObject(
        objectIdentifier=("analogInput", 0),
        objectName="AI-0",
        presentValue=Real(0.0),
        outOfService=False,
    )
    app.add_object(ai0)

    app.add_object(
        MultiStateValueObject(
            objectIdentifier=("multiStateValue", 0),
            objectName="MSV-0",
            presentValue=1,
            numberOfStates=4,
        )
    )

    asyncio.create_task(sawtooth(ai0))
    return app


async def main() -> None:
    build_app()
    # Application.run() blocks on the BACnet event loop forever; this
    # container only exists to host that loop, so just await the
    # never-completing future.
    await asyncio.Future()


if __name__ == "__main__":
    asyncio.run(main())
