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
  <protocol> <mspec-source[;mspec-source...]> <output-dir> [namespace]
```

Each source can be one `.mspec` file or a directory containing `.mspec`
files. Separate multiple sources with semicolons.

Generation is staged and then replaces `<output-dir>/model`, so a failed
write leaves the previous generated model intact. Regenerating these protocol model files is an explicit
driver-maintenance step; the driver slices that consume them commit the
result. This is separate from the checked-in ANTLR parser artifacts below,
whose reproducibility is enforced by CI.

## Supported subset and fail-fast behavior

The generator only emits a protocol when every field has a known C# wire
mapping. It deliberately fails before writing output for unsupported field
keywords, terminated arrays, unsupported data-IO shapes, unsupported temporal
primitives, non-big-endian byte order, and unimplemented `STATIC_CALL`
targets. This prevents a successful generation from silently changing a wire
layout or deferring a missing implementation to production.

Counted and byte-length arrays are supported; byte arrays with a byte-length
use the buffer bulk-read path. Modbus `rtuCrcCheck` and `asciiLrcCheck` are
generated as concrete helpers. Other protocol-specific static helpers remain
an explicit generator gap and cause a diagnostic failure until implemented.

## Grammars and the checked-in parsers

`src/generated/` holds the ANTLR 4.13.2 output (lexer, parser, listener and
`.tokens` side-cars) for the two mspec grammars in
`code-generation/protocol-base-mspec/src/main/antlr4/org/apache/plc4x/plugins/codegenerator/language/mspec/`:
`MSpec.g4` and `expression/Expression.g4`. plc4net keeps no copy of them; they
are used exactly as they are upstream, the same files the Java toolchain
compiles. The output is **checked in**, not generated at build time —
day-to-day work needs no JDK and no ANTLR tool.

ANTLR pastes the text of a semantic predicate into the generated code
unchanged, and the lexer's `EmptyLine` rule spells its predicate the way the
Java runtime does, `{getCharPositionInLine() == 0}?`. The generated
`MSpecLexer` is a `partial` class, so `src/MSpecLexer.Predicates.cs` supplies
that method (it returns the C# runtime's `Column`) and the grammar needs no
C# port; without that file the generated lexer does not compile.
`MSpecLexerTests` pins the behaviour.

### Regenerating the parsers

This requires Python 3 and a Java runtime. The script downloads the pinned
ANTLR archive; it does not require a manually installed ANTLR command.

```bash
python plc4net/tools/code-gen/generate_parsers.py
python plc4net/tools/code-gen/generate_parsers.py --check
```

The script downloads the pinned ANTLR archive into the user cache, verifies
its SHA-256 checksum, applies the ASF header and repository line-ending rules,
and omits ANTLR's debug-only `.interp` files. The `--check` mode regenerates in
a temporary directory and fails when the checked-in artifacts drift. The .NET
workflow runs that check whenever PLC4NET or either shared grammar changes.

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
