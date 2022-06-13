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

#ifndef _PLC_WRITER
#define _PLC_WRITER

namespace org
{
	namespace apache
	{
		namespace plc4x
		{
			namespace cpp
			{
				namespace base
				{
					namespace messages
					{
						
						class PlcWriter
						{
							/**
							* Writes a given value to a PLC.
							*
							* @param writeRequest object describing the type, location and value that whould be written.
							* @return a {@link CompletableFuture} giving async access to the response of the write operation.
							*/
							// Todo: CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest);
						public:
							
						};
					}
				}
			}
		}
	}
}

#endif