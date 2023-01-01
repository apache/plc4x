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

## Setting up in CLion

Per default CLion will not be able to run the tests. 
However, it's pretty simple to set up.

In order to make this happen, you should add another `Profile` to your `CMake` configuration.

You can do this by going to: 

`Preferences...` / `Build, Execution, Deployment` / `CMake`

In the `Profiles` list, click on the `+` button to add a new profile.

I gave this profile the name `Test`.

In the settings select the `Build type` = `Debug`, `Toolchain` = `Use Default`.
`CMake options` = `-DUNITY_VERSION:STRING=2.5.2 -DBUILD_PHASE=test-compile`.

Leave the rest unchanged (which is actually empty).

After saving you can select the `Test` profile in the `Configurations` Drop-Down.
After that is selected, the tests should be available. 
