#!/usr/bin/env python3
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

"""Regenerate the checked-in C# parsers from the shared MSpec grammars."""

from __future__ import annotations

import argparse
import hashlib
import shutil
import subprocess
import sys
import tempfile
import urllib.request
from pathlib import Path


ANTLR_VERSION = "4.13.2"
ANTLR_SHA256 = "eae2dfa119a64327444672aff63e9ec35a20180dc5b8090b7a6ab85125df4d76"
ANTLR_URL = (
    "https://repo1.maven.org/maven2/org/antlr/antlr4/"
    f"{ANTLR_VERSION}/antlr4-{ANTLR_VERSION}-complete.jar"
)
NAMESPACE = "org.apache.plc4net.tools.codegen.grammar"

ASF_HEADER = """//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// \"License\"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// \"AS IS\" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

"""


def sha256(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as source:
        for block in iter(lambda: source.read(1024 * 1024), b""):
            digest.update(block)
    return digest.hexdigest()


def obtain_antlr(cache_dir: Path) -> Path:
    cache_dir.mkdir(parents=True, exist_ok=True)
    jar = cache_dir / f"antlr4-{ANTLR_VERSION}-complete.jar"
    if not jar.exists() or sha256(jar) != ANTLR_SHA256:
        with urllib.request.urlopen(ANTLR_URL) as response, jar.open("wb") as output:
            shutil.copyfileobj(response, output)
    actual_hash = sha256(jar)
    if actual_hash != ANTLR_SHA256:
        jar.unlink(missing_ok=True)
        raise RuntimeError(
            f"ANTLR archive checksum mismatch: expected {ANTLR_SHA256}, got {actual_hash}"
        )
    return jar


def normalize(text: str) -> str:
    return text.replace("\r\n", "\n").replace("\r", "\n")


def generate(repository_root: Path, output_dir: Path, antlr_jar: Path) -> None:
    grammar_root = Path(
        "code-generation/protocol-base-mspec/src/main/antlr4/"
        "org/apache/plc4x/plugins/codegenerator/language/mspec"
    )
    grammars = (grammar_root / "MSpec.g4", grammar_root / "expression/Expression.g4")

    for grammar in grammars:
        subprocess.run(
            [
                "java",
                "-jar",
                str(antlr_jar),
                "-Dlanguage=CSharp",
                "-listener",
                "-no-visitor",
                "-package",
                NAMESPACE,
                "-Xexact-output-dir",
                "-o",
                str(output_dir),
                str(grammar).replace("\\", "/"),
            ],
            cwd=repository_root,
            check=True,
        )

    for generated_file in output_dir.glob("*.cs"):
        generated = normalize(generated_file.read_text(encoding="utf-8"))
        generated_file.write_text(
            (ASF_HEADER + generated).rstrip("\n"),
            encoding="utf-8",
            newline="\n",
        )

    for generated_file in output_dir.glob("*.tokens"):
        generated = normalize(generated_file.read_text(encoding="utf-8"))
        generated_file.write_text(
            ASF_HEADER.rstrip("\n") + "\n" + generated,
            encoding="utf-8",
            newline="\n",
        )

    for interpreter_data in output_dir.glob("*.interp"):
        interpreter_data.unlink()


def compare_or_replace(staged: Path, checked_in: Path, check: bool) -> int:
    expected_names = sorted(path.name for path in staged.iterdir())
    actual_names = sorted(path.name for path in checked_in.iterdir())
    differences = []

    if expected_names != actual_names:
        differences.append("generated file list")

    for name in sorted(set(expected_names) & set(actual_names)):
        expected = normalize((staged / name).read_text(encoding="utf-8"))
        actual = normalize((checked_in / name).read_text(encoding="utf-8"))
        if expected != actual:
            differences.append(name)

    if check:
        if differences:
            print("Checked-in parser artifacts differ: " + ", ".join(differences), file=sys.stderr)
            return 1
        print("Checked-in parser artifacts are reproducible.")
        return 0

    shutil.rmtree(checked_in)
    shutil.copytree(staged, checked_in)
    print(f"Regenerated {len(expected_names)} parser artifacts in {checked_in}")
    return 0


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--check",
        action="store_true",
        help="verify that regeneration produces the checked-in artifacts",
    )
    parser.add_argument(
        "--cache-dir",
        type=Path,
        default=Path.home() / ".cache" / "plc4x" / "antlr",
        help="directory used for the verified ANTLR archive",
    )
    args = parser.parse_args()

    repository_root = Path(__file__).resolve().parents[3]
    checked_in = Path(__file__).resolve().parent / "src" / "generated"
    antlr_jar = obtain_antlr(args.cache_dir)

    with tempfile.TemporaryDirectory(prefix="plc4net-antlr-") as temporary_dir:
        staged = Path(temporary_dir)
        generate(repository_root, staged, antlr_jar)
        return compare_or_replace(staged, checked_in, args.check)


if __name__ == "__main__":
    raise SystemExit(main())
