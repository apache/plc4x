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

#ifndef _DEFAULT_PLC_READ_REQUEST
#define _DEFAULT_PLC_READ_REQUEST

#include "InternalPlcReadRequest.h"
#include "InternalPlcFieldRequest.h"
#include "PlcReader.h"
#include <org/apache/plc4x/cpp/api/model/PlcField.h>
#include "../connection/PlcFieldHandler.h"

#include <map>
#include <boost/foreach.hpp>




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
                        using namespace org::apache::plc4x::cpp::api::messages;
                        using namespace org::apache::plc4x::cpp::base::connection;

                        class DefaultPlcReadRequest : public InternalPlcReadRequest
						{
						public:
							// Todo:
							/*CompletableFuture<PlcReadResponse> execute() 
							{
								return reader.read(this);
							}*/

							// Todo: implement java Builder-pattern for C++							
							
							DefaultPlcReadRequest(PlcReader* plcReader, std::map<std::string, PlcField*> mplcFields);
							int getNumberOfFields();						
							std::vector<std::string> getFieldNames();			
							PlcField* getField(std::string name);
							std::vector<PlcField*> getFields();
							std::map<std::string, PlcField*> getNamedFields();

                            class Builder
                            {
                            public:
                                Builder(PlcReader* plcReader, PlcFieldHandler* plcFieldHandler);
                                Builder* addItem(std::string strName, std::string strFieldQuery);
                                PlcReadRequest* build();

                            private:
                                PlcReader* _plcReader;
                                PlcFieldHandler* _plcFieldHandler;
                                std::map<std::string, std::string> _mFields;
                            };


						protected:

							PlcReader* getReader();
							

						private:					

							DefaultPlcReadRequest();
							PlcReader* _plcReader;
							std::map<std::string, PlcField*> _mplcFields;
						};
					}
				}
			}
		}
	}
}

#endif