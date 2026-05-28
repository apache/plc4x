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

bacpypes3 0.0.102 constructs the Application via Application.from_args; the
SimpleArgumentParser exposes --address / --instance / --vendoridentifier etc.
We assemble argv manually with values pulled from the BACNET_LOCAL_ADDRESS env
var (set by docker-compose to the container's static IP) and then attach a
handful of object instances for the tests to read/write/subscribe against.

Exposes:
    DEVICE:1234              objectName="plc4x-it", vendorIdentifier=0x4D4D
    ANALOG_VALUE:0..4        writable Real PresentValue
    BINARY_VALUE:0..1        writable Enumerated PresentValue
    ANALOG_INPUT:0           read-only Real PresentValue with a 2-second
                             sawtooth so SubscribeCOV consumers see traffic
    MULTI_STATE_VALUE:0      writable Unsigned PresentValue
"""

import asyncio
import os
import sys

from bacpypes3.apdu import SimpleAckPDU
from bacpypes3.app import Application
from bacpypes3.argparse import SimpleArgumentParser
from bacpypes3.constructeddata import Array
from bacpypes3.errors import ExecutionError
from bacpypes3.local.analog import AnalogInputObject, AnalogValueObject
from bacpypes3.local.binary import BinaryValueObject
from bacpypes3.local.multistate import MultiStateValueObject
from bacpypes3.primitivedata import Real, Unsigned
from bacpypes3.service.object import ReadWritePropertyMultipleServices


# bacpypes3 0.0.102 ships a stub WritePropertyMultiple handler that raises
# UnrecognizedService(). Replace it with a working dispatch so the plc4go
# integration tests can verify WPM wire-format end-to-end. Mirrors the
# per-property cast_out + obj.write_property flow that do_WritePropertyRequest
# uses, applied to each (object, property) pair in the request.
async def _do_write_property_multiple(self, apdu):
    for spec in apdu.listOfWriteAccessSpecs:
        obj = self.get_object_id(spec.objectIdentifier)
        if not obj:
            raise ExecutionError(errorClass="object", errorCode="unknownObject")
        for prop in spec.listOfProperties:
            property_type = obj.get_property_type(prop.propertyIdentifier)
            array_index = prop.propertyArrayIndex
            priority = prop.priority
            if issubclass(property_type, Array):
                if array_index is None:
                    pass
                elif array_index == 0:
                    property_type = Unsigned
                else:
                    property_type = property_type._subtype
            value = prop.value.cast_out(property_type, null=(priority is not None))
            await obj.write_property(prop.propertyIdentifier, value, array_index, priority)
    await self.response(SimpleAckPDU(context=apdu))


ReadWritePropertyMultipleServices.do_WritePropertyMultipleRequest = _do_write_property_multiple


# In docker-compose the test runs on bridge 172.30.0.10/24; the host can
# override this for a non-docker invocation.
DEFAULT_ADDRESS = os.environ.get("BACNET_LOCAL_ADDRESS", "172.30.0.10/24:47808")
DEFAULT_INSTANCE = os.environ.get("BACNET_INSTANCE", "1234")
DEFAULT_NAME = os.environ.get("BACNET_NAME", "plc4x-it")
# bacpypes3 ships a bundled vendor-id → vendor-info map and refuses to start
# with an unknown vendor id. The library pre-registers two ids: 0 (ASHRAE,
# spec maintainer) and 999 (intended for tests / unregistered vendors).
DEFAULT_VENDOR = os.environ.get("BACNET_VENDOR", "999")


async def sawtooth(obj: AnalogInputObject) -> None:
    """Background coroutine bumping AnalogInput.0 by 1.0 every 2 seconds so
    SubscribeCOV consumers see a stream of notifications."""
    value = 0.0
    while True:
        await asyncio.sleep(2.0)
        value = (value + 1.0) % 100.0
        obj.presentValue = Real(value)


async def main() -> None:
    # Build the argv that SimpleArgumentParser expects.
    parser = SimpleArgumentParser()
    args = parser.parse_args([
        "--address", DEFAULT_ADDRESS,
        "--instance", DEFAULT_INSTANCE,
        "--name", DEFAULT_NAME,
        "--vendoridentifier", DEFAULT_VENDOR,
        "--debug",
        "bacpypes3.ipv4.IPv4DatagramServer",
        "bacpypes3.ipv4.bvll",
        "bacpypes3.ipv4.link",
        "bacpypes3.ipv4.service",
        "bacpypes3.npdu",
        "bacpypes3.apdu",
        "bacpypes3.app.Application",
    ])
    print(f"starting bacpypes3 device {args.instance} @ {args.address}", flush=True)

    app = Application.from_args(args)

    # Writable scratchpad objects for Read/Write integration tests.
    for i in range(5):
        app.add_object(
            AnalogValueObject(
                objectIdentifier=("analog-value", i),
                objectName=f"AV-{i}",
                presentValue=Real(0.0),
                outOfService=False,
            )
        )
    for i in range(2):
        app.add_object(
            BinaryValueObject(
                objectIdentifier=("binary-value", i),
                objectName=f"BV-{i}",
                presentValue="inactive",
            )
        )

    # Read-only analog input with a slow generator — exercises SubscribeCOV.
    # covIncrement must be non-None or bacpypes3's present_value_filter
    # crashes inside property_change with `Real - NoneType`. 0.5 means every
    # sawtooth tick (1.0) triggers a COV notification.
    ai0 = AnalogInputObject(
        objectIdentifier=("analog-input", 0),
        objectName="AI-0",
        presentValue=Real(0.0),
        outOfService=False,
        covIncrement=Real(0.5),
    )
    app.add_object(ai0)

    app.add_object(
        MultiStateValueObject(
            objectIdentifier=("multi-state-value", 0),
            objectName="MSV-0",
            presentValue=1,
            numberOfStates=4,
        )
    )

    print("device ready", flush=True)
    asyncio.create_task(sawtooth(ai0))

    # Idle forever; bacpypes3's event loop keeps the network stack alive.
    await asyncio.Future()


if __name__ == "__main__":
    # -u (PYTHONUNBUFFERED) on the CMD makes stdout flush per line for compose logs.
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        sys.exit(0)
