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

#include "S7PlcConnection.h"
#include <boost/algorithm/string.hpp>
#include <boost/lexical_cast.hpp>


namespace org
{
	namespace apache
	{
		namespace plc4x
		{
			namespace cpp
			{
				namespace s7
				{
                    S7PlcConnection::S7PlcConnection() 
                    {
                        _lCloseDeviceTimeoutMS = CONF.getLong("plc4x.s7connection.close.device,timeout", 1000);
                    };

                    S7PlcConnection::~S7PlcConnection() 
                    {
                    
                    };

                    S7PlcConnection::S7PlcConnection(std::string strHost, int iRack, int iSlot, string strParams)
                    {
                        setIPAddress(strHost);
                        setPort(_ISO_ON_TCP_PORT);
   
                        _iRack = iRack;
                        _iSlot = iSlot;

                        short sCurParamPduSize = 1024;
                        short sCurParamMaxAmqCaller = 8;
                        short sCurParamMaxAmqCallee = 8;
                        S7ControllerType curParamControllerType = S7ControllerType::ANY;

                        if (strParams.length() > 0)
                        {
                            vector<string> vecParams;
                            boost::split(vecParams, strParams, boost::is_any_of("&"), boost::token_compress_on);
                            for (vector<string>::iterator itParams= vecParams.begin(); itParams < vecParams.end(); itParams++)
                            {
                                vector<string> vecParamElements;
                                boost::split(vecParamElements, *itParams, boost::is_any_of("="), boost::token_compress_on);
                                string strParamName = vecParamElements[0];
                                if (vecParamElements.size() == 2)
                                {
                                    string strParamValue = vecParamElements[1];
                                    boost::algorithm::to_lower(strParamValue);
                                    
                                    if (strParamValue == "pdu-size")
                                    {
                                        sCurParamPduSize = boost::lexical_cast<short>(strParamValue);
                                    }
                                    else if (strParamValue == "max-amq-caller")
                                    {
                                        sCurParamMaxAmqCaller = boost::lexical_cast<short>(strParamValue);
                                    }
                                    else if (strParamValue == "max-amq-callee")
                                    {
                                        sCurParamMaxAmqCallee = boost::lexical_cast<short>(strParamValue);
                                    }
                                    else if (strParamValue == "controller-type")
                                    {
                                        S7ControllerTypeMap s7ControllerTypeMap;
                                        curParamControllerType = s7ControllerTypeMap[strParamValue];
                                    }
                                    else
                                    {
                                        string strMessage = "Unknown parameter " + strParamName + " with value " + strParamValue;
                                        LOG_DEBUG(strMessage);
                                    }
                                }
                                else {
                                    string strMessage = "Unknown no-value parameter " + strParamName;
                                    LOG_DEBUG(strMessage);
                                }
                            }
                        }

                        // It seems that the LOGO devices are a little picky about the pdu-size.
                        // Instead of handling this out, they just hang up without any error message.
                        // So in case of a LOGO controller, set this to a known working value.
                        if (curParamControllerType == S7ControllerType::LOGO && sCurParamPduSize == 1024) 
                        {
                            sCurParamPduSize = 480;
                        }

                        // IsoTP uses pre defined sizes. Find the smallest box,
                        // that would be able to contain the requested pdu size.
                        _sParamPduSize = sCurParamPduSize;
                        _sParamMaxAmqCaller = sCurParamMaxAmqCaller;
                        _sParamMaxAmqCallee = sCurParamMaxAmqCallee;
                        paramControllerType = curParamControllerType;

                        string strMessage = "Setting up S7cConnection with: host-name " + strHost +
                            " rack " + std::to_string(_iRack) + " slot " + std::to_string(_iSlot) +
                            " pdu-size " + std::to_string(_sParamPduSize) + " slot " + std::to_string(_iSlot) +
                            " max-amq-caller " + std::to_string(_sParamMaxAmqCaller) +
                            " max-amq-callee " + std::to_string(_sParamMaxAmqCallee);
                        LOG_INFO(strMessage);

                        return;
                    }

                   

				}
			}
		}
	}
}
