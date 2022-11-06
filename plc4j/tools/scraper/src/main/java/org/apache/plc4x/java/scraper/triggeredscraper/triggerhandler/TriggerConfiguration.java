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
package org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.scraper.exception.ScraperConfigurationException;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScrapeJobImpl;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * basic configuration for all available triggers and handling of regarding condition
 */
public class TriggerConfiguration{
    private static final Logger logger = LoggerFactory.getLogger(TriggerConfiguration.class);

    private static final String TRIGGER = "TRIGGER_VAR";
    private static final String SCHEDULED       = "SCHEDULED";
    private static final String PREVIOUS_DEF    = "PREV";

    private static final double TOLERANCE_FLOATING_EQUALITY = 1e-6;

    private static final Pattern TRIGGER_STRATEGY_PATTERN =
        Pattern.compile("\\((?<strategy>[A-Z_0-9]+),(?<scheduledInterval>\\d+)(,(\\((?<triggerVar>[^!=<>()]+)\\))((?<comp>[!=<>]{1,2}))(\\((?<compVar>[PREVa-z0-9.\\-]+)\\))((?<concatConn>[ANDOR]{2,3})(\\((?<triggerVar2>[^!=<>()]+)\\))((?<comp2>[!=<>]{1,2}))(\\((?<compVar2>[PREVa-z0-9.\\-]+)\\)))?)?\\)");

    private final TriggerType triggerType;
    private final Long scrapeInterval;
    private TriggeredScrapeJobImpl triggeredScrapeJobImpl;
    private List<TriggerElement> triggerElementList;

    /**
     * default constructor when an Field should be used for triggering
     * @param triggerType type of trigger from enum
     * @param scrapeInterval scrape interval of triggered variable
     * @param triggerElementList list of triggerElemts with concat that combined is used as triger
     * @param triggeredScrapeJobImpl the job which is valid for the configuration
     * @throws ScraperConfigurationException when something goes wrong with configuration
     */
    public TriggerConfiguration(TriggerType triggerType,
                                String scrapeInterval,
                                List<TriggerElement> triggerElementList,
                                TriggeredScrapeJobImpl triggeredScrapeJobImpl)
                                throws ScraperConfigurationException {
        this.triggerElementList = triggerElementList;
        this.triggerType = triggerType;
        this.triggeredScrapeJobImpl = triggeredScrapeJobImpl;
        this.scrapeInterval = parseScrapeInterval(scrapeInterval);

        String exceptionMessage;

        if(this.triggerType.equals(TriggerType.TRIGGER_VAR) ) {
            //test for valid field-connection string, on exception quit job and return message to user
            if(this.triggerElementList.isEmpty()){
                exceptionMessage = String.format("No items in trigger List for trigger-type TRIGGER_VAR for Job %s!", triggeredScrapeJobImpl.getJobName());
                throw new ScraperConfigurationException(exceptionMessage);
            }
            checkTriggerVarList();

            //ToDo add more and other trigger
        }
        else{
            exceptionMessage = String.format("TriggerType %s is not yet implemented", this.triggerType);
            throw new ScraperConfigurationException(exceptionMessage);
        }

    }

    /**
     * default constructor when scheduled trigger shall be performed
     * @param triggerType type of trigger from enum
     * @param scrapeInterval scrape interval of data from block
     * @throws ScraperConfigurationException when something goes wrong with configuration
     */
    public TriggerConfiguration(TriggerType triggerType, String scrapeInterval) throws ScraperConfigurationException {
        this.triggerType = triggerType;
        this.scrapeInterval = parseScrapeInterval(scrapeInterval);
        this.triggerElementList = new ArrayList<>();
    }

    /**
     * checks the trigger list for correct syntax
     * @throws ScraperConfigurationException if syntax isn't correct an exception is thrown
     */
    private void checkTriggerVarList() throws ScraperConfigurationException {
        boolean first = true;
        for(TriggerElement triggerElement:this.triggerElementList){
            if(!first && triggerElement.getConcatType()==null){
                throw new ScraperConfigurationException("A concat for the second and following trigger must be given!");
            }
            first = false;
        }
    }

    /**
     * parses String of scrape interval
     * @param scrapeInterval string extracted from RegEx
     * @return converted value
     * @throws ScraperConfigurationException if parsing could not be performed
     */
    private long parseScrapeInterval(String scrapeInterval) throws ScraperConfigurationException {
        try {
            return Long.parseLong(scrapeInterval);
        }
        catch (Exception e){
            handleException(e);
            String exceptionMessage = String.format("No valid numeric for scrapeInterval for Job %s: %s",triggeredScrapeJobImpl.getJobName(),scrapeInterval);
            throw new ScraperConfigurationException(exceptionMessage);
        }
    }

    /**
     * evaluates the trigger dependent of base type and converts acquired respectively ref-value to the needed datatype
     * @param acquiredValues acquired value
     * @return true when condition is matched, false otherwise
     * @throws ScraperException when something goes wrong
     */
    boolean evaluateTrigger(List<Object> acquiredValues) throws ScraperException {
        TriggerEvaluation triggerEvaluation = new TriggerEvaluation(acquiredValues,triggerElementList);
        return triggerEvaluation.evaluateTrigger();
    }


    /**
     * defines the used base type for comparison
     * @return the detected base type
     * @throws ScraperException when an unsupported Type is chosen,which is not (yet) implemented for comparison
     * ToDo check how to handle time-variables if needed
     */
    private static Class<?> validateDataType(PlcField plcField) throws ScraperConfigurationException {
        if(plcField!=null){
            Class<?> javaDataType = plcField.getPlcValueType().getDefaultJavaType();
            if(!javaDataType.equals(Boolean.class)
                && !javaDataType.equals(Integer.class)
                && !javaDataType.equals(Long.class)
                && !javaDataType.equals(Double.class)
            ){
                String exceptionMessage = String.format("Unsupported plc-trigger variable %s with converted data-type %s used",plcField,plcField.getPlcValueType().getDefaultJavaType());
                throw new ScraperConfigurationException(exceptionMessage);
            }
            return javaDataType;
        }
        else{
            String exceptionMessage = "Null plc-trigger variable used";
            throw new ScraperConfigurationException(exceptionMessage);
        }

    }

    /**
     * nested class performing the trigger evaluation
     */
    class TriggerEvaluation{
        private List<Object> acquiredValuesList;
        private List<TriggerElement> triggerElementList;

        TriggerEvaluation(List<Object> acquiredValuesList, List<TriggerElement> triggerElementList) {
            this.acquiredValuesList = acquiredValuesList;
            this.triggerElementList = triggerElementList;
        }

        /**
         * does the evaluation of the trigger conditions are met
         * //ToDo refactor this to improve readability
         * @return true if trigger conditions are met, false otherwise
         * @throws ScraperException if something went wrong
         */
        boolean evaluateTrigger() throws ScraperException {
            List<Boolean> triggerResultList = new ArrayList<>();
            if(logger.isTraceEnabled()){
                String connString = "empty";
                if(!triggerElementList.isEmpty()) {
                    connString = triggerElementList.get(0).getPlcConnectionString();
                }
                logger.trace("eval values for job {} and {}: {}",triggeredScrapeJobImpl.getJobName(),connString,acquiredValuesList);
            }
            //iterate through all items of acquirement-list
            for(int countElements=0; countElements<acquiredValuesList.size();countElements++){
                TriggerElement triggerElement = triggerElementList.get(countElements);
                Object acquiredObject = acquiredValuesList.get(countElements);
                if(validateDataType(triggerElement.getPlcField()).equals(Boolean.class)){
                    //if given type is Boolean
                    boolean currentValue;
                    boolean refValue;
                    try{
                        currentValue = (boolean) acquiredObject;
                        refValue = (boolean) triggerElement.getCompareValue();
                    }
                    catch (Exception e){
                        handleException(e);
                        return false;
                    }
                    if(triggerElement.getComparatorType().equals(Comparator.EQUAL)){
                        triggerResultList.add(currentValue == refValue);
                    }
                    else {
                        triggerResultList.add(currentValue != refValue);
                    }
                }
                if(validateDataType(triggerElement.getPlcField()).equals(Double.class)
                    || validateDataType(triggerElement.getPlcField()).equals(Integer.class)
                    || validateDataType(triggerElement.getPlcField()).equals(Long.class)) {
                    //if given type is numerical
                    boolean skipComparison = false; //comparison shall be skipped if previous values was null
                    double currentValue;
                    double refValue = 0;
                    try{

                        if(acquiredObject instanceof Short){
                            currentValue = ((Short) acquiredObject).doubleValue();
                        }
                        else {
                            if (acquiredObject instanceof Integer) {
                                currentValue = ((Integer) acquiredObject).doubleValue();
                            }
                            else {
                                if (acquiredObject instanceof Long) {
                                    currentValue = ((Long) acquiredObject).doubleValue();
                                }
                                else{
                                    if (acquiredObject instanceof Double) {
                                        currentValue = (Double) acquiredObject;
                                    }else {
                                        currentValue = (double) acquiredObject;
                                    }
                                }
                            }

                        }
                        if(triggerElement.getPreviousMode()){
                            if(triggerElement.getCompareValue()==null){
                                triggerElement.setCompareValue(currentValue);
                                triggerElement.setReservedCompareValue(currentValue);
                                triggerResultList.add(true);
                                if(logger.isTraceEnabled()) {
                                    logger.trace("Initially set compare value to {}", currentValue);
                                }
                                skipComparison=true;
                            }
                            else{
                                refValue = (double) triggerElement.getCompareValue();
                            }
                        }
                        else {
                            refValue = (double) triggerElement.getCompareValue();
                        }

                    }
                    catch (Exception e){
                        handleException(e);
                        return false;
                    }

                    boolean triggerResult = false;
                    if(!skipComparison) {
                        switch (triggerElement.getComparatorType()) {
                            case EQUAL:
                                triggerResult = isApproximately(currentValue, refValue, TOLERANCE_FLOATING_EQUALITY);
                                break;
                            case UNEQUAL:
                                triggerResult = !isApproximately(currentValue, refValue, TOLERANCE_FLOATING_EQUALITY);
                                break;
                            case SMALLER:
                                triggerResult = currentValue < refValue;
                                break;
                            case SMALLER_EQUAL:
                                triggerResult = currentValue <= refValue;
                                break;
                            case GREATER:
                                triggerResult = currentValue > refValue;
                                break;
                            case GREATER_EQUAL:
                                triggerResult = currentValue >= refValue;
                                break;
                            default:
                                triggerResult = false;
                        }
                    }

                    if(triggerResult && triggerElement.getPreviousMode()){
                        triggerElement.setReservedCompareValue(currentValue);
                        if(logger.isTraceEnabled()) {
                            logger.trace("Subcondition matched. Previous value: {}, current compare value {} for Job {}",
                                triggerElement.getReservedCompareValue(),
                                triggerElement.getCompareValue(),
                                triggerElement.getTriggerJob());
                        }
                    }
                    triggerResultList.add(triggerResult);

                }

            }
            if(triggerResultList.isEmpty()){
                if(logger.isDebugEnabled()) {
                    logger.debug("No results could be acquired - setting trigger to false");
                }
                return false;
            }
            //check if there is more then one condition for trigger
            if(triggerResultList.size()>1) {
                if(logger.isTraceEnabled()) {
                    logger.trace("{}", triggerResultList);
                }
                boolean combinedResult=triggerResultList.get(0);
                for (int countElements = 1; countElements < acquiredValuesList.size(); countElements++) {
                    switch (triggerElementList.get(countElements).getConcatType()){
                        case AND:
                            combinedResult = combinedResult && triggerResultList.get(countElements);
                            break;
                        case OR:
                            combinedResult = combinedResult || triggerResultList.get(countElements);
                            break;
                        default:
                            //should not happen
                            combinedResult = false;
                    }
                }
                if(combinedResult) {
                    triggerElementList.forEach(TriggerElement::overrideCompareValue);
                }
                return combinedResult;
            }
            else{
                if(triggerResultList.get(0)) {
                    triggerElementList.forEach(TriggerElement::overrideCompareValue);
                }
                //return first result because its the only one
                return triggerResultList.get(0);
            }
        }
    }

    /**
     * creates the TriggerConfiguration for a given ScrapeJob from triggerConfig-String
     * @param jobTriggerStrategy config-string from file
     * @param triggeredScrapeJob job belonging to the config
     * @return created TriggerConfiguration
     * @throws ScraperConfigurationException when something goes wrong
     */
    public static TriggerConfiguration createConfiguration(String jobTriggerStrategy,TriggeredScrapeJobImpl triggeredScrapeJob) throws ScraperConfigurationException {
        Matcher matcher = TRIGGER_STRATEGY_PATTERN.matcher(jobTriggerStrategy);

        if(matcher.matches()){
            String triggerStrategy = matcher.group("strategy");
            String scheduledMs = matcher.group("scheduledInterval");
            if(logger.isDebugEnabled()) {
                logger.debug("Strategy: {}, scheduled ms: {}", triggerStrategy, scheduledMs);
            }

            String triggerVar = matcher.group("triggerVar");
            String comparatorString = matcher.group("comp");
            String comparatorVariable = matcher.group("compVar");

            switch (triggerStrategy){
                case TRIGGER:

                    if(triggerVar ==null || comparatorString==null || comparatorVariable==null){
                        throw new ScraperConfigurationException("TRIGGER_VAR trigger strategy needs the trigger-condition - information missing! given configString: "+jobTriggerStrategy);
                    }

                    List<TriggerElement> triggerElements = new ArrayList<>();

                    //TODO Change this (probably only 1 source to get the connection directly)
                    String connectionString = triggeredScrapeJob.getSourceConnections().get(triggeredScrapeJob.getSourceConnections().keySet().iterator().next());
                    TriggerElement triggerElement = new TriggerElement(
                        comparatorString,
                        null,
                        comparatorVariable,
                        triggerVar,
                        triggerStrategy,
                        connectionString);

                    triggerElement.setTriggerJob(triggeredScrapeJob.getJobName());
                    triggerElements.add(triggerElement);

                    String concatConn = matcher.group("concatConn");
                    String triggerVar2 = matcher.group("triggerVar2");
                    String comparatorString2 = matcher.group("comp2");
                    String comparatorVariable2 = matcher.group("compVar2");

                    if(triggerVar2 != null && comparatorString2 != null && comparatorVariable2 != null && concatConn != null){
                        TriggerElement triggerElement2 = new TriggerElement(
                            comparatorString2,
                            concatConn,
                            comparatorVariable2,
                            triggerVar2,
                            triggerStrategy,
                            connectionString);


                        triggerElement2.setTriggerJob(triggeredScrapeJob.getJobName());
                        triggerElements.add(triggerElement2);

                    }

                    //ToDo add clever Strategy to concat more than two conditions if needed
                    return new TriggerConfiguration(TriggerType.TRIGGER_VAR,scheduledMs,triggerElements,triggeredScrapeJob);
                case SCHEDULED:
                    if(triggerVar !=null || comparatorString!=null || comparatorVariable!=null){
                        throw new ScraperConfigurationException("SCHEDULED trigger strategy must only be used with scheduled interval - nothing more!  given configString: "+jobTriggerStrategy);
                    }
                    return new TriggerConfiguration(TriggerType.SCHEDULED,scheduledMs);
                default:
                    throw new ScraperConfigurationException("Unknown Trigger Strategy "+triggerStrategy);
            }


        }
        throw new ScraperConfigurationException("Invalid trigger strategy string description: "+jobTriggerStrategy);
    }

    private void handleException(Exception e){
        //push up if needed
        if(logger.isDebugEnabled()) {
            logger.debug("Exception: ", e);
        }
    }

    TriggerType getTriggerType() {
        return triggerType;
    }

    public long getScrapeInterval() {
        return scrapeInterval;
    }

    public List<TriggerElement> getTriggerElementList() {
        return triggerElementList;
    }

    /**
     * check for approximate equality to avoid "Floating-point expressions shall not be tested for equality or inequality." Sonar-Bug
     * @param self current value
     * @param other reference value
     * @param within tolerance band
     * @return if approximate equal, false otherwise
     */
    private static boolean isApproximately(double self, double other, double within)
    {
        return Math.abs(self - other) <= within;
    }

    public enum Comparator {
        EQUAL,
        UNEQUAL,
        GREATER,
        GREATER_EQUAL,
        SMALLER,
        SMALLER_EQUAL
    }

    //ToDo replace constant TriggerType by more generic ones --> PLC4X-89
    public enum TriggerType {
        SCHEDULED,
        S7_TRIGGER_VAR,
        TRIGGER_VAR
    }

    public enum ConcatType {
        AND,
        OR
    }

    public static class TriggerElement{
        private Comparator comparatorType;
        private ConcatType concatType;
        //if trigger should be compared to previous value
        private Boolean previousMode;
        private Object compareValue;
        private PlcField plcField;
        private String plcFieldString;

        private String plcConnectionString;

        private String uuid;

        private String triggerJob;

        //storage for overwrite if condition matched
        private Object reservedCompareValue;

        public TriggerElement() {
            this.comparatorType = null;
            this.concatType = null;
            this.previousMode = false;
            this.compareValue = null;
            this.plcField = null;
            this.plcFieldString = null;
            this.reservedCompareValue = null;
            this.plcConnectionString="not defined";
            this.triggerJob = "Not yet defined";
            this.uuid = "";
        }

        public TriggerElement(Comparator comparatorType, ConcatType concatType, Boolean previousMode, Object compareValue, PlcField plcField, String plcFieldString) {
            this.comparatorType = comparatorType;
            this.concatType = concatType;
            this.previousMode = previousMode;
            this.compareValue = compareValue;
            this.plcField = plcField;
            this.plcFieldString = plcFieldString;
        }

        public TriggerElement(Comparator comparatorType, Object compareValue, PlcField plcField) {
            this();
            this.comparatorType = comparatorType;
            this.compareValue = compareValue;
            this.plcField = plcField;
        }

        TriggerElement(String comparator, String concatType, String compareValue, String plcField, String triggerStrategy, String plcConnectionString) throws ScraperConfigurationException {
            this();
            this.plcFieldString = plcField;
            this.plcConnectionString = plcConnectionString;
            if(triggerStrategy.equals(TRIGGER)){
                try {
                    this.plcField = prepareField(plcFieldString);
                }
                catch (Exception e){
                    if(logger.isDebugEnabled()) {
                        logger.debug("Exception occurred parsing a PlcField");
                    }
                    throw new ScraperConfigurationException("Exception on parsing S7Field (" + plcField + "): " + e.getMessage());
                }
                this.compareValue = convertCompareValue(compareValue,this.plcField);
                this.comparatorType = detectComparatorType(comparator);
                matchTypeAndComparator();
            }

            this.concatType = detectConcatType(concatType);

        }

        //I used this because the prepareField method is deprecated with generated drivers
        //So I need to create the field using the connection string here
        private PlcField prepareField(String fieldQuery) throws PlcConnectionException {
            PlcDriverManager driverManager = new PlcDriverManager();
            PlcDriver driver = driverManager.getDriverForUrl(plcConnectionString);
            return driver.prepareField(fieldQuery);
        }

        /**
         * parses the ref-value to a given value, as well as checking if ref-value matches to the given data-type
         * @param compareValue compare-value extracted by regex
         * @return converted object to needed data-type
         * @throws ScraperException when something does not match or parsing fails
         */
        private Object convertCompareValue(String compareValue, PlcField plcField) throws ScraperConfigurationException {
            Class<?> javaDataType = validateDataType(plcField);
            if(javaDataType.equals(Boolean.class)){
                switch (compareValue){
                    case "1":
                    case "true":
                        return true;
                    case "0":
                    case "false":
                        return false;
                    default:
                        String exceptionMessage = String.format("No valid compare Value at DataType Boolean for trigger: %s",compareValue);
                        throw new ScraperConfigurationException(exceptionMessage);
                }
            }
            if(javaDataType.equals(Double.class)
                || javaDataType.equals(Integer.class)
                || javaDataType.equals(Long.class)){
                try {
                    //everything fits to Double for conversion ... so for first step use only double
                    //ToDo if different handling dependent on specific datatype is needed then differ
                    if(PREVIOUS_DEF.equals(compareValue)){
                        this.previousMode=true;
                        return null;
                    }
                    return Double.parseDouble(compareValue);
                }
                catch (Exception e){
                    logger.debug(e.getMessage(), e);
                    String exceptionMessage = String.format("No valid compare Value at DataType Numeric for trigger: %s",compareValue);
                    throw new ScraperConfigurationException(exceptionMessage);
                }
            }
            String exceptionMessage = "Invalid Datatype detected ... should not happen and be catcht earlier - please report";
            throw new ScraperConfigurationException(exceptionMessage);
        }

        /**
         * converts parsed comparator from regex to ComparatorType
         * @throws ScraperException when no valid comparator has been used
         */
        private Comparator detectComparatorType(String comparator) throws ScraperConfigurationException {
            switch (comparator){
                case "==":
                    return Comparator.EQUAL;
                case "!=":
                    return Comparator.UNEQUAL;
                case "<=":
                    return Comparator.SMALLER_EQUAL;
                case "<":
                    return Comparator.SMALLER;
                case ">=":
                    return Comparator.GREATER_EQUAL;
                case ">":
                    return Comparator.GREATER;
                default:
                    throw new ScraperConfigurationException("Invalid comparator detected!");
            }
        }

        /**
         * convertes parsed comparator from regex to ComparatorType
         * @throws ScraperException when no valid comparator has been used
         */
        private ConcatType detectConcatType(String concat) throws ScraperConfigurationException {
            //concat is not necessary in every case, correct usage is checked later on
            if(concat==null){
                return null;
            }
            switch (concat){
                case "AND":
                    return ConcatType.AND;
                case "OR":
                    return ConcatType.OR;
                default:
                    throw new ScraperConfigurationException("Invalid concat between triggerVars detected: "+concat);
            }
        }

        /**
         * matches data-type and comparator for a valid combination
         * @throws ScraperException when invalid combination is detected
         */
        private void matchTypeAndComparator() throws ScraperConfigurationException {
            if(validateDataType(this.plcField).equals(Boolean.class)
                && !(this.comparatorType.equals(Comparator.EQUAL) || this.comparatorType.equals(Comparator.UNEQUAL))){
                String exceptionMessage = String.format("Trigger-Data-Type (%s) and Comparator (%s) do not match",this.plcField.getPlcValueType().getDefaultJavaType(),this.comparatorType);
                throw new ScraperConfigurationException(exceptionMessage);
            }
            //all other combinations are valid
        }

        Comparator getComparatorType() {
            return comparatorType;
        }

        ConcatType getConcatType() {
            return concatType;
        }

        Boolean getPreviousMode() {
            return previousMode;
        }

        Object getCompareValue() {
            return compareValue;
        }

        PlcField getPlcField() {
            return plcField;
        }

        String getPlcFieldString() {
            return plcFieldString;
        }

        void setCompareValue(Object compareValue) {
            this.compareValue = compareValue;
        }

        Object getReservedCompareValue() {
            return reservedCompareValue;
        }

        void setReservedCompareValue(Object reservedCompareValue) {
            this.reservedCompareValue = reservedCompareValue;
        }

        String getTriggerJob() {
            return triggerJob;
        }

        void setTriggerJob(String triggerJob) {
            this.triggerJob = triggerJob;
        }

        void overrideCompareValue(){
            if(this.previousMode && this.reservedCompareValue!=null){
                if(logger.isDebugEnabled()) {
                    logger.debug("Compare value overridden, before: {}, now: {}; for Trigger {}", this.compareValue, this.reservedCompareValue, this.triggerJob);
                }
                this.compareValue = this.reservedCompareValue;
            }
        }

        public String getPlcConnectionString() {
            return plcConnectionString;
        }

        public void setPlcConnectionString(String plcConnectionString) {
            this.plcConnectionString = plcConnectionString;
        }

        public String getUuid() {
            return uuid;
        }

        public void setUuid(String uuid) {
            this.uuid = uuid;
        }
    }
}