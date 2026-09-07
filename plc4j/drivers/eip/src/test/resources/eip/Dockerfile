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

FROM python:3.12-slim

WORKDIR /app

# cpppo is the de-facto open-source EtherNet/IP / CIP simulator. It exposes
# named tags over CIP just like a Logix controller does, so it covers the
# read/write paths our driver actually exercises.
RUN pip install --no-cache-dir cpppo==5.2.5

# The cpppo --print flag echoes a line to stdout each time a request is
# served, which we use as a readiness signal. The "TAG=TYPE[N]" syntax
# defines tags with a given CIP data type and array length.
EXPOSE 44818

CMD ["python", "-u", "-m", "cpppo.server.enip", \
     "--print", \
     "--address", "0.0.0.0:44818", \
     "hurz_BOOL=BOOL", \
     "hurz_SINT=SINT", \
     "hurz_INT=INT", \
     "hurz_DINT=DINT", \
     "hurz_REAL=REAL", \
     "hurz_DINT_ARR=DINT[10]"]
