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

# Logstash Java Plugin

[![Travis Build Status](https://travis-ci.org/logstash-plugins/logstash-filter-java_filter_example.svg)](https://travis-ci.org/logstash-plugins/logstash-filter-java_filter_example)

This is a Java plugin for [Logstash](https://github.com/elastic/logstash).

It is fully free and fully open source. The license is Apache 2.0, meaning you are free to use it however you want.

The documentation for Logstash Java plugins is available [here](https://www.elastic.co/guide/en/logstash/6.7/contributing-java-plugin.html).

Example input pipeline for logstash:
```
## logstash pipeline config - input
input {
	## use plc4x plugin (logstash-input-plc4x)
	plc4x {
		## define sources (opc-ua examples)
		sources => {
			source1 => "opcua:tcp://opcua-server:4840/"
			source2 => "opcua:tcp://opcua-server1:4840/"
		}
		## define jobs
		jobs => {
			job1 => {
				# pull rate in milliseconds
				rate => 1000
				# sources queried by job1
				sources => ["source1"]
				# defined queries [logstash_internal_fieldname => "IIoT query"]
				queries =>  {
					PreStage => "ns=2;i=3"
					MidStage => "ns=2;i=4"
					PostStage => "ns=2;i=5"
					Motor => "ns=2;i=6"
					ConvoyerBeltTimestamp => "ns=2;i=7"
					RobotArmTimestamp => "ns=2;i=8"
				}
			}
		}
	}
}
```

# Build the plugin:

    ./mvnw clean package -P with-java,with-logstash

# To install the plugin:
1) Copy the build plugin *logstash-input-plc4x-0.5.0.gem* from your target/gem directory to the Logstash plugin directory *logstash/plugin*
2) Follow installation instructions from Logstash documentation: https://www.elastic.co/guide/en/logstash/current/working-with-plugins.html
3) Create pipeline file (see example pipeline above)