<img src="https://camo.githubusercontent.com/fef71cfa1b21130afaf48a60c683ff11a24a6fe6/68747470733a2f2f75706c6f61642e77696b696d656469612e6f72672f77696b6970656469612f636f6d6d6f6e732f7468756d622f642f64622f4170616368655f536f6674776172655f466f756e646174696f6e5f4c6f676f5f253238323031362532392e7376672f3130323470782d4170616368655f536f6674776172655f466f756e646174696f6e5f4c6f676f5f253238323031362532392e7376672e706e67"  width="200" />

# <img src="https://camo.githubusercontent.com/86abd95b803d973f9dbda5ae4f46998971aa7296/68747470733a2f2f706c6334782e6170616368652e6f72672f696d616765732f6170616368655f706c6334785f6c6f676f2e706e67" width="100" />  Writing generated Driver


### Clone latest source and build

In order to make sure to be up to date, clone the latest sources from the official [git](https://github.com/apache/plc4x). If you want to submit your work later on, don't forget to check out the [contributing](https://plc4x.apache.org/developers/contributing.html) page. 

Once downloaded, follow the `README` instructions to build from sources. The following guide will use the 0.7.0-SNAPSHOT version, if your version differs, make sure to specify the correct version inside your `pom` files. We will also create an example driver called `Brol`

NOTE: make sure to add the [Apache Source Header](https://www.apache.org/legal/src-headers.html) in every file you create.

### Create modules

To create a new driver, you need to create following modules inside the project:

- `plc4x-protocols` : create a module called `plc4x-protocols-brol`
- `plc4j-drivers` : create a module called `plc4j-driver-brol`

The `plc4x-protocols` modules contain the `mspec` file which is used to describe the protocol used by the driver, independent from the language. The `plc4j-drivers` is a module of the Java implementation of PLC4X that will be used to implement the protocol into Java.



### plc4x-protocols

- `pom.xml`: make sure to have to specify the parent `pom` and add the following dependency

```xml
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 		http://maven.apache.org/xsd/maven-4.0.0.xsd">
  <parent>
    <artifactId>plc4x-protocols</artifactId>
    <groupId>org.apache.plc4x</groupId>
    <version>0.7.0-SNAPSHOT</version>
  </parent>
  <modelVersion>4.0.0</modelVersion>

  <artifactId>plc4x-protocols-your-driver</artifactId>

  <name>Your Driver</name>
  <description>Base protocol specifications for the brol protocol</description>

  <dependencies>
    <dependency>
      <groupId>org.apache.plc4x</groupId>
      <artifactId>plc4x-build-utils-protocol-base-mspec</artifactId>
      <version>0.7.0-SNAPSHOT</version>
    </dependency>
</dependencies>
</project>
```

- `/src/main/java` : create a new class called `BrolProtocol` in the package `org.apache.plc4x.protocol.Brol`

```java
package org.apache.plc4x.protocol.brol;

import org.apache.plc4x.plugins.codegenerator.language.mspec.parser.MessageFormatParser;
import org.apache.plc4x.plugins.codegenerator.protocol.Protocol;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.exceptions.GenerationException;

import java.io.InputStream;
import java.util.Map;

public class brol implements Protocol {

    @Override
    public String getName() {
        return "brol"; //this will be used by the PlcDriverManager to find the driver
    }

    @Override
    public Map<String, TypeDefinition> getTypeDefinitions() throws GenerationException {
        InputStream schemaInputStream = EipProtocol.class.getResourceAsStream("/protocols/brol/brol.mspec");
        if(schemaInputStream == null) {
            throw new GenerationException("Error loading message-format schema for protocol '" + getName() + "'");
        }
        return new MessageFormatParser().parse(schemaInputStream);
    }
}
```

- to export your protocol as service, you need to create a `org.apache.plc4x.plugins.codegenerator.protocol.Protocol` file under `/src/main/resources/META-INF/services` with following content:

```
#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#
org.apache.plc4x.protocol.brol.brolProtocol
```

- finally you can create the protocol description under `/src/main/resources/protocols/Brol/Brol.mspec`

For the content of the `mspec` pleaser refer to its dedicated [section](#mspec).

### plc4j-drivers

- `pom.xml`: make sure to specify the parent `pom` and add required dependencies. Depending on the features you want to implement, you can add some later on. Here is the minimum

```xml
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
  <modelVersion>4.0.0</modelVersion>

  <parent>
    <groupId>org.apache.plc4x</groupId>
    <artifactId>plc4j-drivers</artifactId>
    <version>0.7.0-SNAPSHOT</version>
  </parent>

  <artifactId>plc4j-driver-Brol</artifactId>
  <name>PLC4J: Driver: Brol</name>
  <description>Implementation of a PLC4X driver able to speak using the Brol protocol.</description>

  <build>
    <plugins>
      <plugin>
        <groupId>org.apache.plc4x.plugins</groupId>
        <artifactId>plc4x-maven-plugin</artifactId>
        <executions>
          <execution>
            <id>test</id>
            <phase>generate-sources</phase>
            <goals>
              <goal>generate-driver</goal>
            </goals>
            <configuration>
              <protocolName>Brol</protocolName>
              <languageName>java</languageName>
              <outputFlavor>read-write</outputFlavor>
            </configuration>
          </execution>
        </executions>
      </plugin>
      <plugin>
        <groupId>org.apache.maven.plugins</groupId>
        <artifactId>maven-dependency-plugin</artifactId>
        <configuration>
          <usedDependencies combine.children="append">
            <usedDependency>org.apache.plc4x:plc4x-build-utils-language-java</usedDependency>
            <usedDependency>org.apache.plc4x:plc4x-protocols-Brol</usedDependency>
          </usedDependencies>
        </configuration>
      </plugin>
    </plugins>
  </build>

  <dependencies>
    <dependency>
      <groupId>org.apache.plc4x</groupId>
      <artifactId>plc4j-api</artifactId>
      <version>0.7.0-SNAPSHOT</version>
    </dependency>
    <dependency>
      <groupId>org.apache.plc4x</groupId>
      <artifactId>plc4j-spi</artifactId>
      <version>0.7.0-SNAPSHOT</version>
    </dependency>

      <!--Make sure to add the transport used by your protocol-->
    <dependency>
      <groupId>org.apache.plc4x</groupId>
      <artifactId>plc4j-transport-tcp</artifactId>
      <version>0.7.0-SNAPSHOT</version>
    </dependency>

    <dependency>
      <groupId>io.netty</groupId>
      <artifactId>netty-buffer</artifactId>
    </dependency>

    <dependency>
      <groupId>org.apache.commons</groupId>
      <artifactId>commons-lang3</artifactId>
    </dependency>

    <dependency>
      <groupId>org.apache.plc4x</groupId>
      <artifactId>plc4x-build-utils-language-java</artifactId>
      <version>0.7.0-SNAPSHOT</version>
      <!-- Scope is 'provided' as this way it's not shipped with the driver -->
      <scope>provided</scope>
    </dependency>

    <dependency>
      <groupId>org.apache.plc4x</groupId>
      <artifactId>plc4x-protocols-eip</artifactId>
      <version>0.7.0-SNAPSHOT</version>
      <!-- Scope is 'provided' as this way it's not shipped with the driver -->
      <scope>provided</scope>
    </dependency>
    <dependency>
      <groupId>com.fasterxml.jackson.core</groupId>
      <artifactId>jackson-annotations</artifactId>
    </dependency>
      <!--Used to export the Driver as OSGi service-->
    <dependency>
      <groupId>org.osgi</groupId>
      <artifactId>osgi.cmpn</artifactId>
      <version>6.0.0</version>
    </dependency>
      <dependency>
          <groupId>org.apache.plc4x</groupId>
          <artifactId>plc4j-utils-test-utils</artifactId>
          <version>0.7.0-SNAPSHOT</version>
          <scope>test</scope>
      </dependency>

  </dependencies>

</project>

```

- Create following packages and classes ([here](https://github.com/etiennerobinet/plc4x/tree/develop/plc4j/drivers/eip/src/main/java/org/apache/plc4x/java/eip/readwrite) an example): 
  -  `org.apache.plc4x.java.Brol.readwrite` 
    - `BrolDriver.java` used to describe the port used, the protocol code to use when parsing the URI, the default transport and the ability to write/read/subscribe 
  -  `org.apache.plc4x.java.Brol.readwrite.configuration` 
    - `BrolConfiguration.java` here you can describe the parameters you can give through the URI and can be later on used in the logic of the Protocol
  -  `org.apache.plc4x.java.Brol.readwrite.field `
    - `BrolField.java` used to describe the a field in your implementation. This will be later used in the Protocol logic to encode/decode into/from the packet
    - `BrolFieldHandler.java`this will be used to create the fields matching the Pattern defined in `Brolfield.java` . Here you will also describe how the values from the read request will be handled
  -  `org.apache.plc4x.java.Brol.readwrite.protocol` 
    - `BrolProtocolLogic.java` describes the logic of the [protocol](#ProtocolLogic). Here you will implement the way to connect, disconnect, send and receive package using the plc4x interfaces like `PlcReadRequest`,`PlcWriteResponse,` etc.

### Mspec

### ProtocolLogic