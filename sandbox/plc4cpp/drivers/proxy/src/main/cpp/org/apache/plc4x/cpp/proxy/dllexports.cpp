/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

//-------- DllMain.cpp --------//
#include "dllexports.h"
#include "ProxyDriver.h"
#if defined (_WIN32)
    #include <windows.h>
#endif

using namespace org::apache::plc4x::cpp::proxy;
using namespace org::apache::plc4x::cpp::spi;

#if defined (_WIN32)
int WINAPI DllEntryPoint(HINSTANCE hinst, unsigned long reason, void*)
{
    return 1;
}
#endif

extern "C" PlcDriver* _CreatePlcDriverInstance()
{
    return new ProxyDriver;
}