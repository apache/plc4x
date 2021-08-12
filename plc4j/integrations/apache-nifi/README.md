# PLC4X Apache NiFi Integration

## Plc4xSinkProcessor

## Plc4xSourceProcessor

## Plc4xSourceRecordProcessor

This processor is <ins>record oriented</ins>, formatting output flowfile content using a Record Writer (for further information see [NiFi Documentation](https://nifi.apache.org/docs/nifi-docs/html/record-path-guide.html#overview)). An example of operation is included on the following path:
*./test-nifi-template/NIFI-PLC4XIntegration-record-example-1.12.xml*. This file is a Nifi Template that could be directly imported from the NiFi UI to test the operation.

The Plc4xSourceRecord Processor can be configured using the following **properties**:

- *PLC connection String:* PLC4X connection string used to connect to a given PLC device.
- *PLC resource address String:* PLC4X address string used identify the resource to read/write on a given PLC device (Multiple values supported). The expected  format is: {name}={address}(;{name}={address}*)
- *Record Writer:* Specifies the Controller Service to use for writing results to a FlowFile. The Record Writer may use Inherit Schema to emulate the inferred schema behavior, i.e. an explicit schema need not be defined in the writer, and will be supplied by the same logic used to infer the schema from the column types.
- *Force Reconnect every request:* Specifies if the connection to the PLC will be recreated on every trigger event

An *example* of these properties for reading values from a S7-1200:

- *PLC connection String:* *s7://10.105.143.7:102?remote-rack=0&remote-slot=1&controller-type=S7_1200*
- *PLC resource address String:* *var1=%DB1:DBX0.0:BOOL;var2=%DB1:DBX0.1:BOOL;var3=%DB1:DBB01:BYTE;var4=%DB1:DBW02:WORD;var5=%DB1:DBW04:INT*
- *Record Writer:* *PLC4x Embedded - AvroRecordSetWriter*
- *Force Reconnect every request:* *false*

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


The output flowfile will contain the PLC read values. This information is included in the flowfile content, following the Record Oriented presentation using an **schema** and the configuration specified in the Record Writer (format, schema inclusion, etc). In the schema, one field will be included for each of the variables defined in the Processor's  *PLC resource address String:* property, taking into account the specified datatype. Also, a *ts* (timestamp) field is additionally included containing the read date. An example of the content of a flowfile for the previously defined properties:

```
[ {
  "var1_boolean" : true,
  "var2_boolean" : false,
  "var3_byte" : [ false, false, false, false, false, false, true, true ],
  "var5_int" : 1992,
  "var4_word" : [ false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, true ],
  "ts" : 1628783058433
} ]
```

Also, it is important to keep in mind the Processor Scheduling Configuration. Using the parameter **Run Schedule** (for example to *1 sec*), the reading frequency can be set. Note that by default, this value is defined to 0 sec (as fast as possible).

Table of data mapping between plc data and avro

| PLC type | Avro Type |
|----------|-----------|
| PlcBigDecimal | floatType |
| PlcBigInteger | longType |
| PlcBitString | stringType |
| PlcBOOL | booleanType |
| PlcBYTE | stringType |
| PlcCHAR | stringType |
| PlcDATE_AND_TIME | stringType |
| PlcDATE | stringType |
| PlcDINT | stringType |
| PlcDWORD | stringType |
| PlcINT | intType |
| PlcLINT | stringType |
| PlcList | stringType |
| PlcLREAL | stringType |
| PlcLTIME | stringType |
| PlcLWORD | stringType |
| PlcNull | stringType |
| PlcREAL | doubleType |
| PlcSINT | intType |
| PlcSTRING | stringType |
| PlcStruct | stringType |
| PlcTIME_OF_DAY | stringType |
| PlcTIME | stringType |
| PlcUDINT | stringType |
| PlcUINT | stringType |
| PlcULINT | stringType |
| PlcUSINT | stringType |
| PlcWCHAR | stringType |
| PlcWORD | stringType |
| ELSE | stringType |
