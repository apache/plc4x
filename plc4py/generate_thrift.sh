#!/usr/bin/env bash

# Requires thrift 0.12.0
# Generates the "generated" files in the right dir structure

# Please be cautious with this, as there was some manual fixing necessary to get everything running (for me at least)

thrift --gen py:python3 -out ./ ../plc4j/utils/interop/src/main/thrift/interop.thrift