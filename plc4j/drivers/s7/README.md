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
# S7 Driver

## Siemens LOGO

The S7 driver supports communication with Siemens LOGO devices, if the S7 communication is enabled.

However the implementation seems to be not as advanced as the normal S7 devices. 
In case of an unexpected message, the PLC doesn't resond with an error message, but it just closes the connection without and error message.
Therefore we need to provide the type of controller as part of the connection string.

A valid Siemens LOGO connection string looks as follows:

   s7://10.10.64.21/0/0?controller-type=LOGO
   
By specifying the controller-type, the driver doesn't try to read the CPU Services Data. 