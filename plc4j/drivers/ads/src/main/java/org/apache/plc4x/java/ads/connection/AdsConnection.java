package org.apache.plc4x.java.ads.connection;

import org.apache.plc4x.java.api.messages.PlcBrowseRequest;
import org.apache.plc4x.java.api.messages.PlcBrowseResponse;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.*;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;

import java.util.concurrent.CompletableFuture;

public class AdsConnection  extends DefaultNettyPlcConnection {

    public AdsConnection(boolean canRead, boolean canWrite, boolean canSubscribe, boolean canBrowse, PlcFieldHandler fieldHandler, PlcValueHandler valueHandler, Configuration configuration, ChannelFactory channelFactory, boolean awaitSessionSetupComplete, boolean awaitSessionDisconnectComplete, boolean awaitSessionDiscoverComplete, ProtocolStackConfigurer stackConfigurer, BaseOptimizer optimizer) {
        super(canRead, canWrite, canSubscribe, canBrowse, fieldHandler, valueHandler, configuration, channelFactory, awaitSessionSetupComplete, awaitSessionDisconnectComplete, awaitSessionDiscoverComplete, stackConfigurer, optimizer);
    }

    @Override
    public CompletableFuture<PlcBrowseResponse> browse(PlcBrowseRequest browseRequest) {
        return super.browse(browseRequest);
    }

}
