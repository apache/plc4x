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

#ifndef _DBG_TRACE_H_
#define _DBG_TRACE_H_

#define TRACE_BUFFERSIZE	0x1000000

// TRACE MACROS to main header file
#if defined(_DEBUG)

#define DEBUG_TRACE_0(f) TracerW(f)
#define DEBUG_TRACE_1(f,a) TracerW(f,a)
#define DEBUG_TRACE_2(f,a,b) TracerW(f,a,b)
#define DEBUG_TRACE_3(f,a,b,c) TracerW(f,a,b,c)
void TracerW(const wchar_t* formatStr, ...);

#define DEBUG_TRACE0(f) Tracer(f)
#define DEBUG_TRACE1(f,a) Tracer(f,a)
#define DEBUG_TRACE2(f,a,b) Tracer(f,a,b)
#define DEBUG_TRACE3(f,a,b,c) Tracer(f,a,b,c)
void Tracer(const char* formatStr, ...);

#else

#define DEBUG_TRACE_0(f) ((void)NULL)
#define DEBUG_TRACE_1(f,a) ((void)NULL)
#define DEBUG_TRACE_2(f,a,b) ((void)NULL)
#define DEBUG_TRACE_3(f,a,b,c) ((void)NULL)

#define DEBUG_TRACE0(f) ((void)NULL)
#define DEBUG_TRACE1(f,a) ((void)NULL)
#define DEBUG_TRACE2(f,a,b) ((void)NULL)
#define DEBUG_TRACE3(f,a,b,c) ((void)NULL)

#endif

#endif
