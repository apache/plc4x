!--
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
-->
# PLC4X Apache NiFi Integration

## Plc4xSinkProcessor

## Plc4xSourceProcessor

## Plc4xSourceRecordProcessor

This processor is <ins>record oriented</ins>, formatting output flowfile content using a Record Writer (for further information see [NiFi Documentation](https://nifi.apache.org/docs/nifi-docs/html/record-path-guide.html#overview)). 

The Plc4xSourceRecord Processor can be configured using the following **properties**:

- *PLC connection String:* PLC4X connection string used to connect to a given PLC device.
- *Record Writer:* Specifies the Controller Service to use for writing results to a FlowFile. The Record Writer may use Inherit Schema to emulate the inferred schema behavior, i.e. an explicit schema need not be defined in the writer, and will be supplied by the same logic used to infer the schema from the column types.
- *Read timeout (miliseconds):* Specifies the time in milliseconds for the connection to return a timeout

Then, the PLC variables to be accessed are specificied using Nifi processor **Dynamic Properties**. For each variable, add a new property to the processor where the property name matches the variable name, and the variable value corresponds to the address tag. 

An *example* of these properties for reading values from a S7-1200:

- *PLC connection String:* *s7://10.105.143.7:102?remote-rack=0&remote-slot=1&controller-type=S7_1200*
- *Record Writer:* *PLC4x Embedded - AvroRecordSetWriter*
- *Read timeout (miliseconds):* *10000*
- *var1:* *%DB1:DBX0.0:BOOL*
- *var2:* *%DB1:DBX0.1:BOOL*
- *var3:* *%DB1:DBB01:BYTE*
- *var4:* *%DB1:DBW02:WORD*
- *var5:* *%DB1:DBW04:INT*

Another *example* of these properties for reading values using OPCUA:
- *PLC connection String:* *opcua:tcp://10.105.143.6:4840?discovery=false*
- *Record Writer:* *PLC4x Embedded - AvroRecordSetWriter*
- *Read timeout (miliseconds):* *10000*
- *AcyclicReceiveBit00:* *ns=2;i=11*
- *MaxCurrentI_max:* *ns=2;i=33*

For the **Record Writer** property, any writer included in NiFi could be used, such as JSON, CSV, etc (also custom writers can be created). In this example, an Avro Writer is supplied, configured as follows:

- *Schema Write Strategy:* Embed Avro Schema
- *Schema Cache:* No value set
- *Schema Protocol Version:* 1
- *Schema Access Strategy:* Inherit Record Schema
- *Schema Registry:* No value set
- *Schema Name:* ${schema.name}
- *Schema Version:* No value set
- *Schema Branch:* No value set
- *Schema Text:* ${avro.schema}
- *Compression Format:* NONE
- *Cache Size:* 1000
- *Encoder Pool Size:* 32


The output flowfile will contain the PLC read values. This information is included in the flowfile content, following the Record Oriented presentation using a **schema** and the configuration specified in the Record Writer (format, schema inclusion, etc). In the schema, one tag will be included for each of the variables defined taking into account the specified datatype. Also, a *ts* (timestamp) field is additionally included containing the read date. An example of the content of a flowfile for the previously defined properties:

```
[ {
  "var1" : true,
  "var2" : false,
  "var3" : [ false, false, false, false, false, false, true, true ],
  "var5" : 1992,
  "var4" : [ false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, true ],
  "ts" : 1628783058433
} ]
```

Also, it is important to keep in mind the Processor Scheduling Configuration. Using the parameter **Run Schedule** (for example to *1 sec*), the reading frequency can be set. Note that by default, this value is defined to 0 sec (as fast as possible).

Table of data mapping between plc data and Avro


| PLC type | Avro Type |
|----------:|-----------|
| PlcBOOL | boolean |
| PlcBYTE | bytes |
| PlcSINT | int |
| PlcINT | int |
| PlcLINT | long |
| PlcREAL | float |
| PlcLREAL | double |
| PlcCHAR | string |
| PlcDATE_AND_TIME | string |
| PlcDATE | string |
| PlcDINT | string |
| PlcDWORD | string |
| PlcLTIME | string |
| PlcLWORD | string |
| PlcNull | string |
| PlcSTRING | string |
| PlcTIME_OF_DAY | string |
| PlcTIME | string |
| PlcUDINT | string |
| PlcUINT | string |
| PlcULINT | string |
| PlcUSINT | string |
| PlcWCHAR | string |
| PlcWORD | string |
| ELSE | string |
