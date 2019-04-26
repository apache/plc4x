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
# Python binding for the interop server

This module provides the (experimental) Python support for the interop server.
Or, simpler, a python 3 binding for PLC4X.

The only thing which needs to be done as _installation_ is to run the `initialize_interop_server.sh` skript to build the interop server and copy it to `lib/`  forder.

Then, you are good to go.

Some tests can be found in `test/test_PlcDriverManager.py`.

Here is some example code:

```python
try:
    manager = PlcDriverManager()

    connection = None
    try:
        connection = manager.get_connection("s7://192.168.167.210/0/1")
        for _ in range(100):
            result = connection.execute(Request(fields={"field1": "%M0:USINT"}))
            print("Response Code is " + str(result.get_field("field1").get_response_code()))
            # We now that we want to get an int...
            print("Response Value is " + str(result.get_field("field1").get_int_value()))

    except PlcException as e:
        raise Exception(str(e.url))
    finally:
        if connection is not None:
            connection.close()
finally:
    manager.close()
```

the `PlcDriverManager` spawns an interop server in the background, thus it is important to close it afterwards.
Otherwise this process keeps alive and you have to kill by yourself.

All generated files (from thrift) are in `org.apache.plc4x.interop`.
I built a very simple Python API in `org.apache.plc4x`.