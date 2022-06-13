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

#include "DefaultPlcWriteRequest.h"

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
                        DefaultPlcWriteRequest::DefaultPlcWriteRequest()
						{
						}

						// =========================================================
						DefaultPlcWriteRequest::DefaultPlcWriteRequest(PlcWriter* plcWriter, std::map<std::string, PlcField*> mplcFields)
						{
							_plcWriter = plcWriter;
                            _mplcFields = mplcFields;
						}

                        // ====================================================
                        PlcWriter* DefaultPlcWriteRequest::getWriter()
                        {
                            return _plcWriter;
                        }

                        // ==========================================================
                        int DefaultPlcWriteRequest::getNumberOfFields()
                        {
                            return _mplcFields.size();
                        }

                        // =========================================================
                        PlcField* DefaultPlcWriteRequest::getField(std::string strName)
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
                        std::map<std::string, PlcField*> DefaultPlcWriteRequest::getNamedFields()
                        {
                            return std::map<std::string, PlcField*>();
                        }

                        // =====================================================
                        std::vector<PlcField*> DefaultPlcWriteRequest::getFields()
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
                        std::vector<std::string> DefaultPlcWriteRequest::getFieldNames()
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

						

											
					}
				}
			}
		}
	}
}
