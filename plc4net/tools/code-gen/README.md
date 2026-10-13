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
  <protocol> <mspec-source[;mspec-source...]> <output-dir> [namespace] [--manual-static-calls]
```

Each source can be one `.mspec` file or a directory containing `.mspec`
files. Separate multiple sources with semicolons.

Generation is staged and then replaces `<output-dir>/model`, so a failed
write leaves the previous generated model intact. The files follow the
repository's `.editorconfig`: LF line endings and no final newline. Regenerating these protocol model files is an explicit
driver-maintenance step; the driver slices that consume them commit the
result. This is separate from the checked-in ANTLR parser artifacts below,
whose reproducibility is enforced by CI.

### Regenerating protocol models

From the repository root, use the Maven wrapper profiles to regenerate the
checked-in Modbus and S7 protocol models or to verify that they have not
drifted from their `.mspec` sources:

```bash
./mvnw -f plc4net/pom.xml -Pupdate-generated-code generate-sources
./mvnw -f plc4net/pom.xml -Pcheck-generated-code generate-sources
```

The first command updates the checked-in model directories. The second writes
only to a temporary directory and fails when generated output differs from the
checked-in files; CI runs that check.

The profiles use `python` by default. On a system where Python 3 is exposed
only as `python3`, add `-Dpython.executable=python3` to either command.

## Supported subset and fail-fast behavior

The generator only emits a protocol when every field has a known C# wire
mapping. It deliberately fails before writing output for unsupported field
keywords, terminated arrays, unsupported data-IO shapes, unsupported temporal
primitives, non-big-endian byte order, and unimplemented `STATIC_CALL`
targets. This prevents a successful generation from silently changing a wire
layout or deferring a missing implementation to production.

Counted and byte-length arrays are supported; byte arrays with a byte-length
use the buffer bulk-read path. Modbus `rtuCrcCheck` and `asciiLrcCheck` are
generated as concrete helpers. Any other `STATIC_CALL` fails the generation,
unless `--manual-static-calls` is given.

With `--manual-static-calls`, every `STATIC_CALL` the generator has no
implementation for is left to the hand-written part of the protocol's partial
`<Protocol>StaticHelper` class, the way plc4j and plc4go keep their `StaticHelper`
by hand. The generated code calls `<Protocol>StaticHelper.<Function>(...)` as the
mspec writes it, and the generated half of the class does not declare the
function, so a function that nobody wrote fails the compile of the driver
instead of the generation. S7 is generated this way: its string, time and date
codings live in `S7StaticHelper.Manual.cs` next to the generated model.

A case of a discriminated type that does not pin a discriminator read from the
wire keeps the value it was parsed with, as plc4j does: Modbus'
`ModbusPDUError` pins `errorFlag` but not `functionFlag`, so it takes the
function code as an extra trailing constructor and `StaticParse` parameter and
returns it from the discriminator accessor. The message then serializes back to
the bytes it was parsed from, which a checksum over the serialized form relies
on. A discriminator that is a parser argument rather than a wire field is not
kept; it takes the type's default.

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
