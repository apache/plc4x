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
# Welcome to Apache PLC4X!

Apache PLC4X is an effort to create a set of libraries for communicating with industrial grade programmable logic controllers (PLCs) in a uniform way.
We are planning on shipping libraries for usage in: 

1) Java
2) Scala
3) C/C++

As well as provide direct integration into other Apache projects, such as:

1) Apache Edgent
2) Apache Mynewt 

## Environment

Currently the project is configured to require the following software:

1) Java 8 JDK: For running Maven in general as well as compiling the Java and Scala modules `JAVA_HOME configured to point to that.
2) (Optional) Graphwiz: For generating the graphs in the documentation (http://www.graphviz.org/)

## Getting Started

As currently a lot of the work is being done in the documentation and this is handled by mavens site plugin, in order to get all parts built, it is important to also trigger Mavens site generation.
This is done by adding the `site:site` goal to the execution.

Unix/Linux/Mac:

    ./mvnw clean install site:site
    
Windows:

    mvnw.cmd clean install site:site
    
This will generate all artifacts as well as their documentation. 
In order to locally fully test the generated documentation site, it is advisable to add another goal `site:stage` to the build, which will copy the sub-sites of all modules into one directory `target/staging`. 
Without this the links between modules will not work correctly.

Unix/Linux/Mac:

    ./mvnw clean install site:site site:stage
    
Windows:

    mvnw.cmd clean install site:site site:stage

