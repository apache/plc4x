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

#ifndef _DEFAULT_PLC_WRITE_REQUEST
#define _DEFAULT_PLC_WRITE_REQUEST

#include "InternalPlcWriteRequest.h"
#include "InternalPlcFieldRequest.h"
#include "PlcWriter.h"

#include <map>
#include <boost/foreach.hpp>

using namespace org::apache::plc4x::cpp::api::messages;


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
						
						class DefaultPlcWriteRequest : public InternalPlcWriteRequest 
						{
						public:
							// Todo:
							/*CompletableFuture<PlcReadResponse> execute() 
							{
								return reader.read(this);
							}*/

							// Todo: implement java Builder-pattern for C++							
							
							DefaultPlcWriteRequest(PlcWriter* writer, std::map<std::string, PlcField*> fields);
							int getNumberOfFields();						
							std::vector<std::string> getFieldNames();			
							PlcField* getField(std::string name);
							std::vector<PlcField*> getFields();
							std::map<std::string, PlcField*> getNamedFields();
							// Todo: implement callback to: execute() { _writer.write(this); }

						protected:
							PlcWriter* getWriter();
							
						private:					
							DefaultPlcWriteRequest();
							PlcWriter* _plcWriter;
							std::map<std::string, PlcField*> _mplcFields;
						};
					}
				}
			}
		}
	}
}

#endif