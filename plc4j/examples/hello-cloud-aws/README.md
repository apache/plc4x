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
# AWS IoT Core Example

This example requires AWS IoT Core.

To run this code, you need to set up the AWS IoT Core Device.

For how to configure AWS IoT Core, described here: https://docs.aws.amazon.com/iot/latest/developerguide/creating-a-virtual-thing.html 
or https://aws.amazon.com/premiumsupport/knowledge-center/iot-core-publish-mqtt-messages-python/?nc1=h_ls 

Then, you can run the AWSIoTCoreSample,can send and receive PLC data to the cloud.
Some sample arguments:
    
    --cert DeviceCertPath.pem.crt 
    --key DevicePrivateKey.pem.key 
    --rootca AWSRootCACert.pem 
    --endpoint [yourendpoint].iot.ap-[region].amazonaws.com
    --username AWS_ID
    --password AWS_PW 
    --clientId DeviceID
    --connection-string PlcConnectionString
    --field-addresses address

An AWS Iot Core Java example is here: https://github.com/aws/aws-iot-device-sdk-java-v2/blob/main/samples/RawPubSub/src/main/java/rawpubsub/RawPubSub.java


