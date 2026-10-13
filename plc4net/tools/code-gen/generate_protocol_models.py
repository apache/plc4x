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
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

"""Regenerate PLC4Net protocol models or verify that they are up to date."""

import argparse
import subprocess
import sys
import tempfile
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[3]
PLC4NET_ROOT = REPOSITORY_ROOT / "plc4net"
CODE_GENERATOR = PLC4NET_ROOT / "tools" / "code-gen"
MODELS = (
    (
        "modbus",
        REPOSITORY_ROOT
        / "protocols"
        / "modbus"
        / "src"
        / "main"
        / "resources"
        / "protocols"
        / "modbus"
        / "modbus.mspec",
        PLC4NET_ROOT / "drivers" / "modbus",
        (),
    ),
    (
        "s7",
        REPOSITORY_ROOT
        / "protocols"
        / "s7"
        / "src"
        / "main"
        / "resources"
        / "protocols"
        / "s7"
        / "s7.mspec",
        PLC4NET_ROOT / "drivers" / "s7",
        ("--manual-static-calls",),
    ),
)


def generate(protocol, mspec, output, extra_arguments):
    subprocess.run(
        [
            "dotnet",
            "run",
            "--project",
            str(CODE_GENERATOR),
            "--configuration",
            "Release",
            "--",
            protocol,
            str(mspec),
            str(output),
            *extra_arguments,
        ],
        check=True,
        cwd=REPOSITORY_ROOT,
    )


def differs(expected, actual):
    expected_files = {
        path.relative_to(expected)
        for path in expected.rglob("*")
        if path.is_file()
    }
    actual_files = {
        path.relative_to(actual)
        for path in actual.rglob("*")
        if path.is_file()
    }
    missing_files = expected_files - actual_files
    unexpected_files = actual_files - expected_files
    changed_files = [
        path
        for path in expected_files & actual_files
        if (expected / path).read_bytes().replace(b"\r\n", b"\n")
        != (actual / path).read_bytes().replace(b"\r\n", b"\n")
    ]
    if missing_files or unexpected_files or changed_files:
        print(f"Generated model differs: {actual}", file=sys.stderr)
        for path in sorted(missing_files):
            print(f"  missing from generated output: {path}", file=sys.stderr)
        for path in sorted(unexpected_files):
            print(f"  unexpected generated output: {path}", file=sys.stderr)
        for path in sorted(changed_files):
            print(f"  changed generated file: {path}", file=sys.stderr)
        return True
    return False


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--check",
        action="store_true",
        help="regenerate into a temporary directory and fail on checked-in model drift",
    )
    arguments = parser.parse_args()

    if not arguments.check:
        for protocol, mspec, driver_directory, extra_arguments in MODELS:
            generate(protocol, mspec, driver_directory, extra_arguments)
        return 0

    with tempfile.TemporaryDirectory(prefix="plc4net-models-") as temporary_directory:
        temporary_root = Path(temporary_directory)
        has_drift = False
        for protocol, mspec, driver_directory, extra_arguments in MODELS:
            output = temporary_root / protocol
            generate(protocol, mspec, output, extra_arguments)
            has_drift = differs(output / "model", driver_directory / "model") or has_drift
        return 1 if has_drift else 0


if __name__ == "__main__":
    sys.exit(main())
