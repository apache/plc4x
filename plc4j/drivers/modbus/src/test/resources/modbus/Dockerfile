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

RUN pip install --no-cache-dir pymodbus==3.8.6

# Generate self-signed certificates for TLS testing at build time.
RUN apt-get update && apt-get install -y --no-install-recommends openssl && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /app/certs \
    && openssl req -x509 -newkey rsa:2048 -keyout /app/certs/ca.key -out /app/certs/ca.crt \
       -days 365 -nodes -subj "/CN=Test-CA" \
    && openssl req -newkey rsa:2048 -keyout /app/certs/server.key -out /app/certs/server.csr \
       -nodes -subj "/CN=localhost" \
    && printf "subjectAltName=DNS:localhost,IP:127.0.0.1" > /tmp/san.ext \
    && openssl x509 -req -in /app/certs/server.csr -CA /app/certs/ca.crt -CAkey /app/certs/ca.key \
       -CAcreateserial -out /app/certs/server.crt -days 365 -extfile /tmp/san.ext \
    && rm /app/certs/server.csr /tmp/san.ext

COPY modbus_server.py .

EXPOSE 5020 5021 5022 5023 5025/udp 5026/udp 5027/udp

CMD ["python", "-u", "modbus_server.py"]
