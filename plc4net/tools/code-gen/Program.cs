//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using org.apache.plc4net.tools.codegen;
using org.apache.plc4net.tools.codegen.output;

// The pure-.NET mspec -> C# generator. Replaces the `update-generated-code`
// Maven profile: no JDK, no freemarker.
//
//   dotnet run --project plc4net/tools/code-gen -- \
//     <protocol> <mspec-source[;mspec-source...]> <output-dir> [namespace]
//
// <output-dir> receives model/*.cs. It is wiped of *.cs under model/ first so
// a removed mspec type does not leave a stale file behind.

if (args.Length < 3)
{
    Console.Error.WriteLine(
        "usage: plc4net-code-gen <protocol> <mspec-source[;mspec-source...]> <output-dir> [namespace]");
    return 2;
}

var protocol = args[0];
var mspecPath = args[1];
var outDir = args[2];
var ns = args.Length > 3
    ? args[3]
    : $"org.apache.plc4net.drivers.{protocol.Replace("-", "")}.readwrite";

// A protocol's types can be split across files and directories (knxnetip
// spans knxnetip.mspec, device-info.mspec and the generated
// knx-master-data.mspec). The source argument is one or more ';'-separated
// paths, each a .mspec file or a directory of them; every *.mspec found is
// compiled as one unit, the way the Java plugin does.
var mspecFiles = new List<string>();
foreach (var part in mspecPath.Split(';', StringSplitOptions.RemoveEmptyEntries))
{
    if (Directory.Exists(part))
    {
        mspecFiles.AddRange(Directory.GetFiles(part, "*.mspec"));
    }
    else if (File.Exists(part))
    {
        // An explicit file means exactly that file - the siblings are
        // independent protocols that merely share a directory
        // (ads.mspec / ads-discovery.mspec both define AmsNetId).
        mspecFiles.Add(Path.GetFullPath(part));
    }
    else
    {
        Console.Error.WriteLine($"mspec source not found: {part}");
        return 2;
    }
}

var sorted = mspecFiles.Distinct().OrderBy(p => p, StringComparer.Ordinal).ToArray();
if (sorted.Length == 0)
{
    Console.Error.WriteLine($"no .mspec files found for '{mspecPath}'");
    return 2;
}
var model = MspecModelBuilder.BuildFiles(sorted);
foreach (var diagnostic in model.Diagnostics)
{
    Console.Error.WriteLine($"warning: {diagnostic}");
}
var files = new CSharpGenerator(model, protocol, ns).Generate();

// Stage every source file before replacing the current model directory. This
// keeps a previously generated driver intact if a disk/write failure occurs.
Directory.CreateDirectory(outDir);
var stageRoot = Path.Combine(outDir, $".plc4net-code-gen-{Guid.NewGuid():N}");
var stageModelDir = Path.Combine(stageRoot, "model");
var modelDir = Path.Combine(outDir, "model");
var backupDir = Path.Combine(outDir, $".plc4net-code-gen-backup-{Guid.NewGuid():N}");

try
{
    foreach (var (relativePath, source) in files)
    {
        var target = Path.Combine(stageRoot, relativePath);
        Directory.CreateDirectory(Path.GetDirectoryName(target)!);
        File.WriteAllText(target, source);
    }

    if (Directory.Exists(modelDir))
    {
        Directory.Move(modelDir, backupDir);
    }
    try
    {
        Directory.Move(stageModelDir, modelDir);
    }
    catch
    {
        if (Directory.Exists(backupDir) && !Directory.Exists(modelDir))
        {
            Directory.Move(backupDir, modelDir);
        }
        throw;
    }

    if (Directory.Exists(backupDir))
    {
        Directory.Delete(backupDir, recursive: true);
    }
}
finally
{
    if (Directory.Exists(stageRoot))
    {
        Directory.Delete(stageRoot, recursive: true);
    }
}

Console.WriteLine(
    $"{protocol}: {model.Types.Count} types, {model.Enums.Count} enums, " +
    $"{model.DataIos.Count} dataIo, {files.Count} files -> {outDir}");
return 0;
