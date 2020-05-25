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
#ifndef PLC4C_SUBSCRIBE_H_
#define PLC4C_SUBSCRIBE_H_
#ifdef __cplusplus
extern "C" {
#endif

#include "plc4c/types.h"

/*
 *
 * Subscriptions
 *
 */

/**
 * Destroys a given subscription-request.
 *
 * @param subscription_request the subscription-request.
 */
void plc4c_subscription_request_destroy(
    plc4c_subscription_request *subscription_request);

/**
 * Returns the plc4c_connection for a give subscription request
 * @param subscription_request plc4c_subscription_request
 * @return plc4c_connection
 */
plc4c_connection *plc4c_subscription_request_get_connection(
    plc4c_subscription_request *subscription_request);

/**
 * Sets the plc4c_connection for a given subscription request
 * @param subscription_request plc4c_subscription_request
 * @param connection plc4c_connection
 */
void plc4c_subscription_request_set_connection(
    plc4c_subscription_request *subscription_request,
    plc4c_connection *connection);

/**
 * Adds a cyclic (polling) item to the subscription request.
 * @param subscription_request the subscription-request object.
 * @param address address for the resource.
 * @param polling_interval_in_millis interval in millis in which new data should
 * be mada available.
 * @return plc4c_return_code
 */
plc4c_return_code plc4c_subscription_request_add_cyclic_item(
    plc4c_subscription_request *subscription_request, char *address,
    long polling_interval_in_millis);

/**
 * Adds a change-of-state item to the subscription request.
 * @param subscription_request the subscription-request object.
 * @param address address for the resource.
 * @return plc4c_return_code
 */
plc4c_return_code plc4c_subscription_request_add_change_of_state_item(
    plc4c_subscription_request *subscription_request, char *address);

/**
 * Adds a event item to the subscription request.
 * @param subscription_request the subscription-request object.
 * @param address address for the resource.
 * @return plc4c_return_code
 */
plc4c_return_code plc4c_subscription_request_add_event_item(
    plc4c_subscription_request *subscription_request, char *address);

/**
 * Actually executes the subscription-request.
 * @param connection connection this subscription-request will be executed on.
 * @param subscription_request the subscription-request object.
 * @param subscription_request_execution pointer to a data-structure handling
 * one execution of the subscription-request.
 * @return plc4c_return_code
 */
plc4c_return_code plc4c_subscription_request_execute(
    plc4c_subscription_request *subscription_request,
    plc4c_subscription_request_execution **subscription_request_execution);

/**
 * Check if new data is available for a given subscription handle.
 * TODO: The subscription handle must contain a reference to the system-task
 * associated with it as well as to the subscription request.
 * @param subscription_handle the subscription handle
 * @return plc4c_return_code
 */
bool plc4c_subscription_check_data_available(void *subscription_handle);

/**
 * Destroys a given subscription-request execution.
 *
 * @param subscription_request_execution the subscription-request execution.
 */
void plc4c_subscription_request_execution_destroy(
    plc4c_subscription_request_execution *subscription_request_execution);

/**
 * Check if the subscription-request is completed successfully.
 *
 * @param subscription_request_execution the subscription-request execution.
 * @return true if the subscription-request is completed successfully.
 */
bool plc4c_subscription_request_execution_check_finished_successfully(
    plc4c_subscription_request_execution *subscription_request_execution);

/**
 * Check if the subscription-request is completed unsuccessfully.
 *
 * @param subscription_request_execution the subscription-request execution.
 * @return true if the subscription-request is completed with an error.
 */
bool plc4c_subscription_request_execution_check_finished_with_error(
    plc4c_subscription_request_execution *subscription_request_execution);

/**
 * Retrieve the subscription-response from a given subscription-request
 * execution.
 *
 * @param subscription_request_execution the subscription-request execution.
 * @return the subscription-response.
 */
plc4c_subscription_response *plc4c_subscription_request_execution_get_response(
    plc4c_subscription_request_execution *subscription_request_execution);

/**
 * Destroys a given subscription-response.
 *
 * @param subscription_response the subscription-response.
 */
void plc4c_subscription_response_destroy(
    plc4c_subscription_response *subscription_response);

/*
 *
 * Event handling
 *
 */

/**
 * Gets a list of events for a given subscription-handle from the system.
 * @param subscription_handle
 * @param events
 * @return
 */
plc4c_return_code plc4c_subscription_get_subscription_events(
    void *subscription_handle, plc4c_list **events);

/*
 *
 * Unsubscriptions
 *
 */

/**
 * Destroys a given unsubscription-request.
 *
 * @param unsubscription_request the unsubscription-request.
 */
void plc4c_unsubscription_request_destroy(
    plc4c_unsubscription_request *unsubscription_request);

/**
 * Returns the plc4c_connection for a give unsubscription request
 * @param unsubscription_request plc4c_unsubscription_request
 * @return plc4c_connection
 */
plc4c_connection *plc4c_unsubscription_request_get_connection(
    plc4c_unsubscription_request *unsubscription_request);

/**
 * Sets the plc4c_connection for a given unsubscription request
 * @param unsubscription_request plc4c_unsubscription_request
 * @param connection plc4c_connection
 */
void plc4c_unsubscription_request_set_connection(
    plc4c_unsubscription_request *unsubscription_request,
    plc4c_connection *connection);

/**
 * Actually executes the unsubscription-request.
 * @param connection connection this unsubscription-request will be executed on.
 * @param unsubscription_request the unsubscription-request object.
 * @param unsubscription_request_execution pointer to a data-structure handling
 * one execution of the unsubscription-request.
 * @return plc4c_return_code
 */
plc4c_return_code plc4c_unsubscription_request_execute(
    plc4c_unsubscription_request *unsubscription_request,
    plc4c_unsubscription_request_execution **unsubscription_request_execution);

/**
 * Destroys a given unsubscription-request execution.
 *
 * @param unsubscription_request_execution the unsubscription-request execution.
 */
void plc4c_unsubscription_request_execution_destroy(
    plc4c_unsubscription_request_execution *unsubscription_request_execution);

/**
 * Check if the unsubscription-request is completed successfully.
 *
 * @param unsubscription_request_execution the unsubscription-request execution.
 * @return true if the unsubscription-request is completed successfully.
 */
bool plc4c_unsubscription_request_execution_check_finished_successfully(
    plc4c_unsubscription_request_execution *unsubscription_request_execution);

/**
 * Check if the unsubscription-request is completed unsuccessfully.
 *
 * @param unsubscription_request_execution the unsubscription-request execution.
 * @return true if the unsubscription-request is completed with an error.
 */
bool plc4c_unsubscription_request_execution_check_finished_with_error(
    plc4c_unsubscription_request_execution *unsubscription_request_execution);

/**
 * Retrieve the unsubscription-response from a given unsubscription-request
 * execution.
 *
 * @param unsubscription_request_execution the unsubscription-request execution.
 * @return the unsubscription-response.
 */
plc4c_unsubscription_response *
plc4c_unsubscription_request_execution_get_response(
    plc4c_unsubscription_request_execution *unsubscription_request_execution);

/**
 * Destroys a given unsubscription-response.
 *
 * @param unsubscription_response the unsubscription-response.
 */
void plc4c_unsubscription_response_destroy(
    plc4c_unsubscription_response *unsubscription_response);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_SUBSCRIBE_H_