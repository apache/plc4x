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
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

"""Regression tests for the protocol model CI check, without invoking dotnet."""

import contextlib
import io
import tempfile
import unittest
from pathlib import Path
from unittest import mock

import generate_protocol_models as models


class ProtocolModelCheckTests(unittest.TestCase):
    def run_check(self, checked_in, generated):
        with tempfile.TemporaryDirectory() as temporary_directory:
            driver = Path(temporary_directory) / "driver"
            model = driver / "model"
            model.mkdir(parents=True)
            (model / "Example.cs").write_bytes(checked_in)

            def generate(protocol, mspec, output, extra_arguments):
                output_model = output / "model"
                output_model.mkdir(parents=True)
                (output_model / "Example.cs").write_bytes(generated)

            diagnostics = io.StringIO()
            with (
                mock.patch.object(models, "MODELS", (("example", None, driver, ()),)),
                mock.patch.object(models, "generate", side_effect=generate),
                mock.patch("sys.argv", ["generate_protocol_models.py", "--check"]),
                contextlib.redirect_stderr(diagnostics),
            ):
                status = models.main()
            return status, diagnostics.getvalue()

    def test_identical_files_without_final_lf_fail(self):
        status, diagnostics = self.run_check(b"class Example {}", b"class Example {}")
        self.assertEqual(status, 1)
        self.assertEqual(diagnostics.count("missing final LF:"), 2)
        self.assertIn("Example.cs", diagnostics)

    def test_checked_in_file_without_final_lf_fails(self):
        status, diagnostics = self.run_check(b"class Example {}", b"class Example {}\n")
        self.assertEqual(status, 1)
        self.assertEqual(diagnostics.count("missing final LF:"), 1)

    def test_generated_file_without_final_lf_fails(self):
        status, diagnostics = self.run_check(b"class Example {}\n", b"class Example {}")
        self.assertEqual(status, 1)
        self.assertEqual(diagnostics.count("missing final LF:"), 1)

    def test_matching_files_with_final_lf_pass(self):
        status, diagnostics = self.run_check(b"class Example {}\n", b"class Example {}\n")
        self.assertEqual(status, 0)
        self.assertEqual(diagnostics, "")

    def test_crlf_checkout_matches_lf_output(self):
        status, diagnostics = self.run_check(b"class Example {}\r\n", b"class Example {}\n")
        self.assertEqual(status, 0)
        self.assertEqual(diagnostics, "")

    def test_content_drift_with_final_lf_still_fails(self):
        status, diagnostics = self.run_check(b"class Old {}\n", b"class New {}\n")
        self.assertEqual(status, 1)
        self.assertIn("changed generated file: Example.cs", diagnostics)


if __name__ == "__main__":
    unittest.main()
