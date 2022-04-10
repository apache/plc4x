/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

#define TRACE_BUFFERSIZE	0x1000000

#if defined(_DEBUG)

#include <stdio.h>
#include <stdarg.h>
#include <wtypes.h>
#include <tchar.h>
#include <windows.h>


void TracerW(const wchar_t* formatStr, ...)
{
	static wchar_t arcBuffer[TRACE_BUFFERSIZE];
	va_list args;
	va_start(args, formatStr);
    swprintf_s(arcBuffer, TRACE_BUFFERSIZE, formatStr, args);
	va_end(args);
	::OutputDebugStringW(arcBuffer);
}


void Tracer(const char* formatStr, ...)
{
	static char arcBuffer[TRACE_BUFFERSIZE];
	va_list args;
	va_start(args, formatStr);
	vsprintf_s(arcBuffer, TRACE_BUFFERSIZE, formatStr, args);
	va_end(args);
	::OutputDebugStringA(arcBuffer);
}

#endif
