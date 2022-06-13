<!--
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
# Plc to Google IoT Core Adapter

This Adapter opens a MQTT Bridge to Google Cloud IoT Core.

To run this code, you need an account to the google cloud and a project.
You need to create a device registry and add a device to it, as described here:
https://cloud.google.com/iot/docs/how-tos/devices

Then, you can run the PlcToGoogleIoTCoreSample to connect tho google and send some values into the cloud.
Some sample arguments:

    -project-id=myprojectname
    -registry-id=plc4x-test
    -cloud-region=europe-west1
    -device-id=plc4x-test-device
    -private-key-file=../../../rsa_private_pkcs8
    -algorithm=RS256

Some documentation can be found here:
https://github.com/GoogleCloudPlatform/java-docs-samples/tree/master/iot/api-client/manager#cloud-iot-core-java-mqtt-example
This code was adapted from:
https://github.com/GoogleCloudPlatform/java-docs-samples/tree/master/iot/api-client/manager#cloud-iot-core-java-mqtt-example

To retrieve the messages, you can subscribe to Cloud Pub/Sub.
For example, you could install the google cloud sdk (gcloud) and run:

    gcloud auth login
    gcloud config set myprojectname
    gcloud iot devices configs list --project=myprojectname \\
        --region=europe-west1 \\
        --registry=myplc4x-test-registry \\
        --device=myplc4x-test-device \\
        --limit=5
    gcloud pubsub subscriptions create --topic plc4x-test-events plc4x-test-subscription
    gcloud pubsub subscriptions pull --auto-ack plc4x-test-subscription


To pull more than one message, use the option --limit [number]

For further reference to the Cloud Pub/Sub, see:
https://cloud.google.com/pubsub/docs/quickstart-cli
https://cloud.google.com/sdk/gcloud/reference/alpha/pubsub/subscriptions/pull

