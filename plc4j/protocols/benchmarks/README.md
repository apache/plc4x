The benchmarks module contains micro-benchmarks for multiple protocols.

It uses JMH for benchmarking (http://openjdk.java.net/projects/code-tools/jmh/)

The recommendation from JMH id to run bechmarks in a separate module like this one:

_The recommended way to run a JMH benchmark is to use Maven to setup a standalone project that depends on the jar files of your application. This approach is preferred to ensure that the benchmarks are correctly initialized and produce reliable results. It is possible to run benchmarks from within an existing project, and even from within an IDE, however setup is more complex and the results are less reliable._
Source: http://openjdk.java.net/projects/code-tools/jmh/

To run the test you can use a plugin for your IDE. If you want to use maven you can use the profile `run-benchmark`:

`mvn -Prun-benchmark verify`