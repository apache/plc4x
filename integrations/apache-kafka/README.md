<!--

  Licensed to the Apache Software Foundation (ASF) under one or more
  contributor license agreements.  See the NOTICE file distributed with
  this work for additional information regarding copyright ownership.
  The ASF licenses this file to You under the Apache License, Version 2.0
  (the "License"); you may not use this file except in compliance with
  the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

-->
Welcome to your new Kafka Connect connector!

# Running in development


The [docker-compose.yml](docker-compose.yml) that is included in this repository is based on the Confluent Platform Docker
images. Take a look at the [quickstart](http://docs.confluent.io/3.0.1/cp-docker-images/docs/quickstart.html#getting-started-with-docker-client)
for the Docker images. 

The hostname `confluent` must be resolvable by your host. You will need to determine the ip address of your docker-machine using `docker-machine ip confluent` 
and add this to your `/etc/hosts` file. For example if `docker-machine ip confluent` returns `192.168.99.100` add this:

```
192.168.99.100  confluent
```


```
docker-compose up -d
```


Start the connector with debugging enabled.
 
```
./bin/debug.sh
```

Start the connector with debugging enabled. This will wait for a debugger to attach.

```
export SUSPEND='y'
./bin/debug.sh
```