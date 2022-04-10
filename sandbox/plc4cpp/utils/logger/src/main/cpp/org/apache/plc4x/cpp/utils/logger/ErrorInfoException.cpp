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

#include <boost/exception/all.hpp>
#include <exception>
#include <boost/lexical_cast.hpp>
#include "ErrorInfoException.h"
#include "ExLog.h"
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

                    ErrorInfoException::ErrorInfoException() :
                        _uiUserMessageID(-1),
                        pCategory(NULL),
                        _iEv(-1)
                    {
                        std::string strErrorMessage = " ErrorInfoException no Message";
                        LOG_ERROR(strErrorMessage);
                        DEBUG_TRACE1("%s\n", strErrorMessage.c_str());

                        return;
                    }

                    ErrorInfoException::ErrorInfoException(unsigned int uiUserMessageID, std::string strErrorInfo) :
                        _uiUserMessageID(uiUserMessageID),
                        _strWhat(strErrorInfo),
                        pCategory(NULL),
                        _iEv(-1)
                    {
                        LOG_ERROR(strErrorInfo);
                        DEBUG_TRACE1("%s\n", strErrorInfo.c_str());

                        return;
                    }

                    ErrorInfoException::ErrorInfoException(unsigned int uiUserMessageID, std::exception& ex) :
                        _uiUserMessageID(uiUserMessageID),
                        _strWhat(ex.what()),
                        pCategory(NULL),
                        _iEv(-1)
                    {
                        LOG_ERROR(ex.what());
                        DEBUG_TRACE1("%s\n", ex.what());

                        return;
                    }

                    ErrorInfoException::ErrorInfoException(unsigned int uiUserMessageID, boost::exception& ex) :
                        _uiUserMessageID(uiUserMessageID),
                        _strWhat(*boost::get_error_info<errmsg_info>(ex)),
                        pCategory(NULL),
                        _iEv(-1)
                    {
                        LOG_ERROR(*boost::get_error_info<errmsg_info>(ex));
                        DEBUG_TRACE1("%s\n", _strWhat.c_str());

                        return;
                    }

                    ErrorInfoException::ErrorInfoException(unsigned int uiUserMessageID, Error_Category::Category cat, int iErrorValue) :
                        _uiUserMessageID(uiUserMessageID),
                        pCategory(CreateErrorCategory(cat)),
                        _iEv(iErrorValue)
                    {
                        std::string strMessage = boost::lexical_cast<std::string>(iErrorValue);
                        LOG_ERROR(strMessage.c_str());
                        DEBUG_TRACE1("ErrorValue:%i\n", iErrorValue);

                        return;
                    }

                    ErrorInfoException::ErrorInfoException(unsigned int uiUserMessageID, Error_Category::Category cat, int iErrorValue, std::string strErrorInfo) :
                        _uiUserMessageID(uiUserMessageID),
                        _strWhat(strErrorInfo),
                        pCategory(CreateErrorCategory(cat)),
                        _iEv(iErrorValue)
                    {
                        LOG_ERROR(strErrorInfo);
                        DEBUG_TRACE1("%s\n", strErrorInfo.c_str());

                        return;
                    }


                    ErrorInfoException::ErrorInfoException(unsigned int uiUserMessageID, Error_Category::Category cat, int iErrorValue, std::exception& ex) :
                        _uiUserMessageID(uiUserMessageID),
                        _strWhat(ex.what()),
                        pCategory(CreateErrorCategory(cat)),
                        _iEv(iErrorValue)
                    {
                        LOG_ERROR(ex.what());
                        DEBUG_TRACE1("%s\n", ex.what());

                        return;
                    }

                    ErrorInfoException::ErrorInfoException(unsigned int uiUserMessageID, Error_Category::Category cat, int iErrorValue, boost::exception& ex) :
                        _uiUserMessageID(uiUserMessageID),
                        _strWhat(*boost::get_error_info<errmsg_info>(ex)),
                        pCategory(CreateErrorCategory(cat)),
                        _iEv(iErrorValue)
                    {
                        LOG_ERROR(*boost::get_error_info<errmsg_info>(ex));
                        DEBUG_TRACE1("%s\n", _strWhat.c_str());

                        return;
                    }

                    ErrorInfoException::~ErrorInfoException()
                    {
                        if (pCategory != NULL) { delete pCategory; };

                        return;
                    }

                    Error_Category* ErrorInfoException::CreateErrorCategory(Error_Category::Category cat) throw()
                    {
                        Error_Category* pCat = NULL;
                        try
                        {
                            switch (cat)
                            {
                            case Error_Category::Category::System:
                            default:
                                pCat = new Error_Category();
                            break;

                            case Error_Category::Category::Application:
                                pCat = new Error_Category_Application();
                            break;
                            }
                        }
                        catch (...)
                        {
                            pCat = NULL;
                        }

                        return pCat;
                    }
                }
            }
        }
    }
}
