/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

#include "DefaultPlcReadRequest.h"

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
						// =========================================================
						DefaultPlcReadRequest::DefaultPlcReadRequest()
						{
						}

						// =========================================================
						DefaultPlcReadRequest::DefaultPlcReadRequest(PlcReader* reader, std::map<std::string, PlcField> fields)
						{
							this->_reader = reader;
							this->_fields = fields;
						}

						// ==========================================================
						int DefaultPlcReadRequest::getNumberOfFields()
						{
							return _fields.size();
						}

						// =========================================================
						PlcField* DefaultPlcReadRequest::getField(std::string name)
						{
							// sgl: check 1st, whether the key is contained in vector 
							// (otherwise an exception is thrown)...
							std::map<std::string, PlcField>::iterator iterator = _fields.find(name);

							if (iterator == _fields.end())
							{
								return nullptr;
							}

							return &(iterator->second);
							
						}

						// ======================================================
						std::map<std::string, PlcField> DefaultPlcReadRequest::getNamedFields()
						{
							// Todo:
							/*return fields.entrySet()
								.stream()
								.map(stringPlcFieldEntry->Pair.of(stringPlcFieldEntry.getKey(), stringPlcFieldEntry.getValue()))
								.collect(Collectors.toCollection(LinkedList::new));
*/
							return std::map<std::string, PlcField>();
						}

						// ====================================================
						PlcReader* DefaultPlcReadRequest::getReader()
						{
							return _reader;
						}

						// =====================================================
						std::vector<PlcField> DefaultPlcReadRequest::getFields()
						{
							// TODO: check if already exists...
							std::pair<std::string, PlcField> me;
							std::vector<PlcField> vNames;

							BOOST_FOREACH(me, _fields)
							{
								vNames.push_back(me.second);
							}

							return vNames;
						}

						// =========================================================
						std::vector<std::string> DefaultPlcReadRequest::getFieldNames()
						{							
							// TODO: check if already exists...
							std::pair<std::string, PlcField> me;
							std::vector<std::string> vNames;

							BOOST_FOREACH(me, _fields)
							{
								vNames.push_back(me.first);
							}

							return vNames;
						}						
					}
				}
			}
		}
	}
}
