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

#include "DefaultPlcReadRequest.h"

using namespace org::apache::plc4x::cpp::api::model;

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
						DefaultPlcReadRequest::DefaultPlcReadRequest(PlcReader* plcReader, std::map<std::string, PlcField*> mplcFields)
						{
                            _plcReader = plcReader;
                            _mplcFields = mplcFields;
						}

                        // ====================================================
                        PlcReader* DefaultPlcReadRequest::getReader()
                        {
                            return _plcReader;
                        }

						// ==========================================================
						int DefaultPlcReadRequest::getNumberOfFields()
						{
							return _mplcFields.size();
						}

						// =========================================================
						PlcField* DefaultPlcReadRequest::getField(std::string strName)
						{
                            PlcField* plcField = nullptr;
                            
                            std::map<std::string, PlcField*>::iterator iterator = _mplcFields.find(strName);

							if (iterator != _mplcFields.end())
							{
                                plcField = iterator->second;
							}

							return plcField;							
						}

						// ======================================================
						std::map<std::string, PlcField*> DefaultPlcReadRequest::getNamedFields()
						{
							return std::map<std::string, PlcField*>();
						}

						// =====================================================
						std::vector<PlcField*> DefaultPlcReadRequest::getFields()
						{
							std::pair<std::string, PlcField*> me;
							std::vector<PlcField*> vplcFields;

							BOOST_FOREACH(me, _mplcFields)
							{
                                vplcFields.push_back(me.second);
							}

							return vplcFields;
						}

						// =========================================================
						std::vector<std::string> DefaultPlcReadRequest::getFieldNames()
						{							
							// TODO: check if already exists...
							std::pair<std::string, PlcField*> me;
							std::vector<std::string> vstrNames;

							BOOST_FOREACH(me, _mplcFields)
							{
                                vstrNames.push_back(me.first);
							}

							return vstrNames;
						}

                        DefaultPlcReadRequest::Builder::Builder(PlcReader* plcReader, PlcFieldHandler* plcFieldHandler)
                        {
                            _plcReader = plcReader;
                            _plcFieldHandler = plcFieldHandler;
                        }

                        DefaultPlcReadRequest::Builder* DefaultPlcReadRequest::Builder::addItem(std::string strName, std::string strFieldQuery)
                        {
                            std::map<std::string, std::string>::iterator iterator = _mFields.find(strName);
                            if (iterator != _mFields.end())
                            {
                                BOOST_THROW_EXCEPTION(PlcRuntimeException("Duplicate field definition '" + strName + "'"));
                            }
                            _mFields.insert(std::pair<std::string,std::string>(strName, strFieldQuery));
                            return this;
                        }

                        PlcReadRequest* DefaultPlcReadRequest::Builder::build()
                        {
                            std::map<std::string, PlcField*> mParsedFields;
                            std::pair<std::string, std::string> me;

                            BOOST_FOREACH(me, _mFields)
                            {
                                PlcField* parsedField = _plcFieldHandler->createField(me.second);
                                mParsedFields.insert(std::pair<std::string, PlcField*>(me.first, parsedField));
                            }
                            return new DefaultPlcReadRequest(_plcReader, mParsedFields);
                        }
					}
				}
			}
		}
	}
}
