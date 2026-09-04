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
# plc4net-code-gen

The pure-.NET `.mspec` → C# generator. Parses `.mspec` with the checked-in
ANTLR-for-C# lexer/parser, walks the tree into a type-model IR, and emits the
model class, `StaticParse`, `Serialize` and `GetLengthInBits` per type — the
replacement for the Java freemarker `model-template` / `io-template` (whose
C# `io-template` was never migrated, so it only ever produced data classes).

## Run

```bash
dotnet run --project plc4net/tools/code-gen -c Release -- \
  <protocol> <mspec-path[;mspec-path...]> <output-dir> [namespace]
```

`<output-dir>/model/*.cs` is wiped and rewritten, so a removed mspec type
leaves no stale file. The CI job `generated-code-is-current` and the Maven
`update-generated-code` profile both invoke this for Modbus and S7.

## Grammar (`MSpec.g4`, `Expression.g4`) and the checked-in parsers

`src/generated/` holds the ANTLR 4.13.2 output (lexer, parser, listener,
`.interp` / `.tokens` side-cars). It is **checked in**, not generated at
build time — day-to-day work needs no JDK and no ANTLR tool.

`MSpec.g4` / `Expression.g4` are copied from
`code-generation/protocol-base-mspec/src/main/antlr4/...` with **one edit**:
the lexer's column-zero predicate is ported from the ANTLR Java runtime API
to the C# one —

| upstream (`code-generation`)          | here                 |
|--------------------------------------|----------------------|
| `{getCharPositionInLine() == 0}?`    | `{Column == 0}?`     |

So the grammar is not byte-shared with the other languages; an upstream
grammar change has to be re-applied here, the predicate re-ported, the ANTLR
tool re-run, and ASF headers re-prepended to the generated `.cs`.

### Regenerating the parsers (needs the ANTLR tool + a JRE)

```bash
# antlr-4.13.2-complete.jar on the classpath
cd plc4net/tools/code-gen
java -jar antlr-4.13.2-complete.jar -Dlanguage=CSharp -o src/generated -package org.apache.plc4net.tools.codegen.grammar MSpec.g4 Expression.g4
# then re-prepend the ASF header to each new *.cs (the license-headers CI job checks this)
```

The `.interp` / `.tokens` side-cars have no comment syntax and carry no
header (RAT-excluded in `plc4net/pom.xml`).

## Layout

| path                         | role                                             |
|------------------------------|--------------------------------------------------|
| `MspecReader.cs`             | text → ANTLR parse tree                           |
| `MspecModelBuilder.cs`       | parse tree → `model/` IR                          |
| `MspecExpressionParser.cs`   | a quoted mspec expression → `model/terms/` tree   |
| `model/`                     | the type-model IR                                 |
| `output/CSharpGenerator.cs`  | IR → C# (class + parse / serialize / length)      |
| `output/CSharpExpressionRenderer.cs` | `Term` → a C# expression                  |
| `Program.cs`                 | the CLI                                           |
