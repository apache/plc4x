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
[![Maven central](https://img.shields.io/maven-central/v/org.apache.plc4x/plc4j-api.svg)](https://img.shields.io/maven-central/v/org.apache.plc4x/plc4j-api.svg)
[![License](https://img.shields.io/github/license/apache/plc4x.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Last commit](https://img.shields.io/github/last-commit/apache/plc4x.svg)]()
[![Platform compatibility](https://img.shields.io/github/workflow/status/apache/plc4x/Platform%20compatibility?label=Platform%20compatibility)](https://github.com/apache/plc4x/actions/workflows/ensure-platforms.yml)
[![Twitter](https://img.shields.io/twitter/follow/ApachePLC4X.svg?label=Follow&style=social)](https://twitter.com/ApachePLC4X)


<h1 align="center">
  <br>
   <img src="https://plc4x.apache.org/images/apache_plc4x_logo.png" 
   alt="Apache PLC4X Logo" title="Apache PLC4X Logo"/>
  <br>
</h1>
<h3 align="center">工业物联网转接器</h3>
<h4 align="center">PLC4X的最终目标是创建一组库，可以提供统一接口访问任何类型的PLC</h4>

***

# 目录

  * [关于](#关于)
  * [入门](#入门)
  * [开发者](#开发者)
  * [社区](#社区)
  * [贡献](#贡献)
  * [许可证](#许可证)

***

## 关于

Apache PLC4X 致力于创建一组用于以统一方式与工业级可编程逻辑控制器 (PLC) 通信的库。
我们计划将库用于以下用途：

1. Java
2. Go
3. C (未准备好)
4. Python (未准备好)
5. C# (.Net) (未准备好)

PLC4X同时与其他Apache项目集成，例如：

* [Apache Calcite](https://calcite.apache.org/)
* [Apache Camel](https://camel.apache.org/)
* [Apache Kafka-Connect](https://kafka.apache.org)
* [Apache Karaf](https://karaf.apache.org/)
* [Apache NiFi](https://nifi.apache.org/)

并且提供了独立的（Java）应用程序，例如：

* OPC-UA 服务器：使您能够使用带有 OPC-UA 的 PLC4X 与传统设备进行通信。
* PLC4X 服务器：使您能够与中央 PLC4X 服务器通信，然后通过 PLC4X 与设备通信。

它还提供（Java）工具以供在应用程序内部使用：

* 连接缓存：用于重用和共享 PLC 连接的框架的新实现
* 连接池：用于重用和共享 PLC 连接的框架的旧实现
* OPM：Object-Plc-Mapping：允许将 PLC 字段绑定到类似于JPA的java POJOs中的属性
* Scraper：用于定期和重复数据收集的实用程序。

## 入门

根据编程语言的不同，用法会有所不同，因此请到 PLC4X 网站上的 [Getting Started](https://plc4x.apache.org/users/gettingstarted.html) 查找选择的语言。

### Java

注意：目前支持构建 Apache PLC4X 所有部分的 Java 版本至少为 Java 11

请参阅网站上的 PLC4J 用户指南，开始在您的 Java 应用程序中使用 PLC4X：
[https://plc4x.apache.org/users/getting-started/plc4j.html](https://plc4x.apache.org/users/getting-started/plc4j.html)

## 开发者

### 环境

目前，该项目配置需要以下软件：

1. Java 11 JDK：一般用来运行 Maven 以及编译配置为指向它的 Java 和 Scala 模块“JAVA_HOME”。
2. Git（在处理源代码分发时）
3. （可选，用于运行所有测试）libpcap/Npcap 用于 Java 中的原始套接字测试或使用“被动模式”驱动程序
4. （可选，用于构建网站）[Graphviz](https://www.graphviz.org/)：用于在文档中生成图表

警告：代码生成需要使用一些额外的VM设置。当从根目录中开始编译，在`.mvn/jvm.config`中的设置将会被自动应用。当仅需要编译一个子模块，这样设置vm的参数就显得很重要：`--add-exports jdk.compiler/com.sun.tools.javac.api=ALL-UNNAMED --add-exports jdk.compiler/com.sun.tools.javac.file=ALL-UNNAMED --add-exports jdk.compiler/com.sun.tools.javac.parser=ALL-UNNAMED --add-exports jdk.compiler/com.sun.tools.javac.tree=ALL-UNNAMED --add-exports jdk.compiler/com.sun.tools.javac.util=ALL-UNNAMED`。例如，在Intellij中，在IDE中这样就这样设置：Preferences | Build, Execution, Deployment | Build Tools | Maven | Runner: JVM Options。

更多详细的描述可以在我们的官网找到：

https://plc4x.apache.org/developers/preparing/index.html

#### 为了编译`PLC4C`我们另外需要：

所有的要求都由构建本身检索。

#### 为了编译`PLC4Go`我们另外需要：

All requirements are retrieved by the build itself所有的要求都有构建本身检索。

#### 为了编译`PLC4Py`我们另外需要：

1. Python 3.7 或者更高
2. Python pyenv

#### 为了编译`PLC4Net`我们另外需要:

1. DotNet SDK 6.0

通过此设置，你将能够构建PLC4X的Java部分。

在进行完整构建时，如果不满足所有要求，我们会自动运行先决条件检查并通过解释使构建失败。

### 通过Docker进行构建

如果您不想在普通系统上设置环境并且安装了 Docker，您还可以在 Docker 容器中构建所有内容：

```
   docker build -t plc4x .
```

### 入门

如果您必须在系统上至少安装 Java 11 并连接到 Maven Central（用于下载外部第三方依赖项）。 构建需要使用 Maven 3.6，因此请确保它已安装并在你的系统上可用。

注意：在 repo 中安装了一个方便的 Maven-Wrapper，使用时，它会自动下载并安装 Maven。 如果你想使用它，请使用 `./mvnw` 或 `mvnw` 而不是普通的 `mvn` 命令。

注意：从sources-zip 运行时，`mvnw` 可能无法在`Mac` 或`Linux` 上执行。 这可以通过在目录中运行以下命令轻松修复。

```
$ chmod +x mvnw
```

注意：如果您在 `Windows` 系统上工作，请在以下构建命令中使用 `mvnw.cmd` 而不是 `./mvnw`。

构建 PLC4X Java jar 并将它们安装在本地 maven 存储库中

```
./mvnw install
```

您现在可以构建使用 PLC4X 的 Java 应用程序。 PLC4X 示例是一个很好的起点，可在`plc4j/examples` 目录中找到。

可以通过启用 `with-go` 配置文件来构建 `Go` 驱动程序：

```
./mvnw -P with-go install 
```

`C#/.Net` 实现当前处于`work in progress` 状态。为了能够构建`C#/.Net` 模块，您当前需要激活：`with-dotnet` 配置文件。

```
./mvnw -P with-dotnet install
```

Python 实现目前处于某种不干净的状态，仍需要重构。
为了能够构建 Python 模块，您当前需要激活：
`with-sandbox` 和 `with-python` 配置文件。

```
./mvnw -P with-sandbox,with-python install
```

为了全部构建，可以使用以下命令：

```
./mvnw -P with-c,with-dotnet,with-go,with-python,with-sandbox install
```

## 社区

使用以下渠道之一加入 PLC4X 社区。 我们很乐意提供帮助！

### 邮件列表

订阅我们的邮件列表：
* Apache PLC4X开发者列表: [dev-subscribe@plc4x.apache.org](mailto:dev-subscribe@plc4x.apache.org)
* Apache PLC4X提交列表: [commits-subscribe@plc4x.apache.org](mailto:commits-subscribe@plc4x.apache.org)
* Apache PLC4XJira通知列表: [issues-subscribe@plc4x.apache.org](mailto:issues-subscribe@plc4x.apache.org)

同时可见: [https://plc4x.apache.org/mailing-lists.html](https://plc4x.apache.org/mailing-lists.html)

### Twitter

获取最新的PLC4X消息: [https://twitter.com/ApachePlc4x](https://twitter.com/ApachePlc4x)

## 贡献

您可以通过多种形式参与 PLC4X 项目。

有，但不局限于：

* 提供信息和见解
* 测试PLC4X并提供反馈
* 提交Pull Requests
* 提交错误报告
* 在我们的邮件列表上积极沟通
* 推广项目（文章、博客以及访谈）
* 文档

我们是一群非常友好的人，所以不要害怕向前迈进。如果您想为 PLC4X 做出贡献，请查看我们的[贡献指南](https://plc4x.apache.org/developers/contributing.html)！

## 许可证

Apache PLC4X是基于Apache 2.0版本许可证下进行发布的。
