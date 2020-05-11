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

#include <plc4c/plc4c.h>
#include <plc4c/subscribe.h>
#include <plc4c/spi/types_private.h>
#include <stdlib.h>

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
    plc4c_subscription_request *subscription_request) {
  free(subscription_request);
}

/**
 * Returns the plc4c_connection for a give subscription request
 * @param subscription_request plc4c_subscription_request
 * @return plc4c_connection
 */
plc4c_connection *plc4c_subscription_request_get_connection(
    plc4c_subscription_request *subscription_request) {
  return subscription_request->connection;
}

/**
 * Sets the plc4c_connection for a given subscription request
 * @param subscription_request plc4c_subscription_request
 * @param connection plc4c_connection
 */
void plc4c_subscription_request_set_connection(
    plc4c_subscription_request *subscription_request,
    plc4c_connection *connection) {
  subscription_request->connection = connection;
}

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
    plc4c_subscription_request_execution **subscription_request_execution) {
  // Inject the default subscription context into the system task.
  plc4c_subscription_request_execution *new_subscription_request_execution =
      malloc(sizeof(plc4c_subscription_request_execution));
  new_subscription_request_execution->subscription_request = subscription_request;
  new_subscription_request_execution->subscription_response = NULL;
  new_subscription_request_execution->system_task = NULL;

  plc4c_system_task *system_task;
  plc4c_connection_get_driver(plc4c_subscription_request_get_connection(subscription_request))
      ->subscribe_function(new_subscription_request_execution, &system_task);

  // Increment the number of running tasks for this connection.
  plc4c_connection_task_added(subscription_request->connection);
  // Add the new task to the task-list.
  plc4c_utils_list_insert_tail_value(
      plc4c_system_get_task_list(plc4c_connection_get_system(
          plc4c_subscription_request_get_connection(subscription_request))),
      system_task);

  *subscription_request_execution = new_subscription_request_execution;
  return OK;
}

/**
 * Destroys a given subscription-request execution.
 *
 * @param subscription_request_execution the subscription-request execution.
 */
void plc4c_subscription_request_execution_destroy(
    plc4c_subscription_request_execution *subscription_request_execution) {
  free(subscription_request_execution);
}

/**
 * Check if the subscription-request is completed successfully.
 *
 * @param subscription_request_execution the subscription-request execution.
 * @return true if the subscription-request is completed successfully.
 */
bool plc4c_subscription_request_execution_check_finished_successfully(
    plc4c_subscription_request_execution *subscription_request_execution) {
  if (subscription_request_execution == NULL) {
    return true;
  }
  if (subscription_request_execution->system_task == NULL) {
    return true;
  }
  return subscription_request_execution->system_task->completed;
}

/**
 * Check if the subscription-request is completed unsuccessfully.
 *
 * @param subscription_request_execution the subscription-request execution.
 * @return true if the subscription-request is completed with an error.
 */
bool plc4c_subscription_request_execution_check_finished_with_error(
    plc4c_subscription_request_execution *subscription_request_execution) {
  // TODO: Implement this sensibly ...
  return false;
}

/**
 * Retrieve the subscription-response from a given subscription-request
 * execution.
 *
 * @param subscription_request_execution the subscription-request execution.
 * @return the subscription-response.
 */
plc4c_subscription_response *plc4c_subscription_request_execution_get_response(
    plc4c_subscription_request_execution *subscription_request_execution) {
  if (subscription_request_execution == NULL) {
    return NULL;
  }
  return subscription_request_execution->subscription_response;
}

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
    plc4c_unsubscription_request *unsubscription_request) {
  free(unsubscription_request);
}

/**
 * Returns the plc4c_connection for a give unsubscription request
 * @param unsubscription_request plc4c_unsubscription_request
 * @return plc4c_connection
 */
plc4c_connection *plc4c_unsubscription_request_get_connection(
    plc4c_unsubscription_request *unsubscription_request) {
  return unsubscription_request->connection;
}

/**
 * Sets the plc4c_connection for a given unsubscription request
 * @param unsubscription_request plc4c_unsubscription_request
 * @param connection plc4c_connection
 */
void plc4c_unsubscription_request_set_connection(
    plc4c_unsubscription_request *unsubscription_request,
    plc4c_connection *connection) {
  unsubscription_request->connection = connection;
}

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
    plc4c_unsubscription_request_execution **unsubscription_request_execution) {
  // Inject the default unsubscription context into the system task.
  plc4c_unsubscription_request_execution *new_unsubscription_request_execution =
      malloc(sizeof(plc4c_unsubscription_request_execution));
  new_unsubscription_request_execution->unsubscription_request = unsubscription_request;
  new_unsubscription_request_execution->unsubscription_response = NULL;
  new_unsubscription_request_execution->system_task = NULL;

  plc4c_system_task *system_task;
  plc4c_connection_get_driver(plc4c_unsubscription_request_get_connection(unsubscription_request))
      ->unsubscribe_function(new_unsubscription_request_execution, &system_task);

  // Increment the number of running tasks for this connection.
  plc4c_connection_task_added(unsubscription_request->connection);
  // Add the new task to the task-list.
  plc4c_utils_list_insert_tail_value(
      plc4c_system_get_task_list(plc4c_connection_get_system(
          plc4c_unsubscription_request_get_connection(unsubscription_request))),
      system_task);

  *unsubscription_request_execution = new_unsubscription_request_execution;
  return OK;
}

/**
 * Destroys a given unsubscription-request execution.
 *
 * @param unsubscription_request_execution the unsubscription-request execution.
 */
void plc4c_unsubscription_request_execution_destroy(
    plc4c_unsubscription_request_execution *unsubscription_request_execution) {
  free(unsubscription_request_execution);
}

/**
 * Check if the unsubscription-request is completed successfully.
 *
 * @param unsubscription_request_execution the unsubscription-request execution.
 * @return true if the unsubscription-request is completed successfully.
 */
bool plc4c_unsubscription_request_execution_check_finished_successfully(
    plc4c_unsubscription_request_execution *unsubscription_request_execution) {
  if (unsubscription_request_execution == NULL) {
    return true;
  }
  if (unsubscription_request_execution->system_task == NULL) {
    return true;
  }
  return unsubscription_request_execution->system_task->completed;
}

/**
 * Check if the unsubscription-request is completed unsuccessfully.
 *
 * @param unsubscription_request_execution the unsubscription-request execution.
 * @return true if the unsubscription-request is completed with an error.
 */
bool plc4c_unsubscription_request_execution_check_finished_with_error(
    plc4c_unsubscription_request_execution *unsubscription_request_execution) {
  // TODO: Implement this sensibly ...
  return false;
}

/**
 * Retrieve the unsubscription-response from a given unsubscription-request
 * execution.
 *
 * @param unsubscription_request_execution the unsubscription-request execution.
 * @return the unsubscription-response.
 */
plc4c_unsubscription_response *
plc4c_unsubscription_request_execution_get_response(
    plc4c_unsubscription_request_execution *unsubscription_request_execution) {
  if (unsubscription_request_execution == NULL) {
    return NULL;
  }
  return unsubscription_request_execution->unsubscription_response;
}
