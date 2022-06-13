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

#include <boost/exception/all.hpp>
#include <exception>
#include <boost/lexical_cast.hpp>
#include <typeinfo>
#include "ErrorInfoException.h"

#include "DbgTrace.h"
#include "BLogger.h"

namespace org
{
    namespace apache
    {
        namespace plc4x
        {
            namespace cpp
            {
                namespace utils
                {
                    using namespace std;
                    using namespace boost;

                    ExLog::ExLog()
                    {

                    }

                    ExLog::ExLog(unsigned int uiUserMessageID, std::exception& ex, const char* pError) :
                        _strWhat(ex.what())
                    {
                        string strExceptionClassName = typeid(ex).name();

                        if (strExceptionClassName.find(typeid(ExLog).name()) != std::string::npos)
                        {
                            _uiUserMessageID = (reinterpret_cast<ExLog&>(ex)).getUserMessageID();
                        }
                        else if (strExceptionClassName.find(typeid(ErrorInfoException).name()) != std::string::npos)
                        {
                            _uiUserMessageID = (reinterpret_cast<ErrorInfoException&>(ex)).getUserMessageID();
                        }
                        else
                        {
                            _uiUserMessageID = uiUserMessageID;
                        }

                        std::stringstream strStream;

                        _strLastErrorMessage = std::string(pError) + '\n';

                        _strLastErrorMessage = _strLastErrorMessage + std::string("MessageID: ") + lexical_cast<std::string>(_uiUserMessageID) + '\n';

                        strStream << diagnostic_information(ex) << std::endl;
                        _strLastErrorMessage = _strLastErrorMessage + strStream.str();

                        _strLastErrorMessage = _strLastErrorMessage + printStackTrace();

                        DEBUG_TRACE1("%s\n", _strLastErrorMessage.c_str());
                        LOG_ERROR(_strLastErrorMessage);

                        return;
                    }

                    ExLog::ExLog(unsigned int uiUserMessageID, boost::exception& ex, const char* pError) :
                        _strWhat(*get_error_info<errmsg_info>(ex))
                    {
                        string strExceptionClassName = typeid(ex).name();

                        if (strExceptionClassName.find(typeid(ExLog).name()) != std::string::npos)
                        {
                            _uiUserMessageID = (reinterpret_cast<ExLog&>(ex)).getUserMessageID();
                        }
                        else if (strExceptionClassName.find(typeid(ErrorInfoException).name()) != std::string::npos)
                        {
                            _uiUserMessageID = (reinterpret_cast<ErrorInfoException&>(ex)).getUserMessageID();
                        }
                        else
                        {
                            _uiUserMessageID = uiUserMessageID;
                        }

                        std::stringstream strStream;

                        _strLastErrorMessage = std::string(pError) + '\n';

                        _strLastErrorMessage = _strLastErrorMessage + std::string("MessageID: ") + lexical_cast<std::string>(_uiUserMessageID) + '\n';

                        strStream << diagnostic_information(ex) << std::endl;
                        _strLastErrorMessage = _strLastErrorMessage + strStream.str();

                        _strLastErrorMessage = _strLastErrorMessage + printStackTrace();

                        DEBUG_TRACE1("%s\n", _strLastErrorMessage.c_str());
                        LOG_ERROR(_strLastErrorMessage);

                        return;
                    }

                    ExLog::~ExLog()
                    {
                    }

                    void ExLog::logInfo(std::string strInfo)
                    {
                        LOG_INFO(strInfo);
                        DEBUG_TRACE1("%s\n", strInfo.c_str());

                        return;
                    }

                    std::string ExLog::printStackTrace(unsigned int uiMaxFrame /* = 63 */)
                    {
                        /* TODO: Implement Stackwalker */

                        std::string strTrace;

                        strTrace = "Stack Trace\n";

                        return strTrace;
                    }
                }
            }
        }
    }
}
