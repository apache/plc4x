package org.apache.plc4x.java.profinet.device;

import org.pcap4j.core.PcapHandle;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.function.Function;

public class ProfinetRunnable implements Runnable {

        private final Logger logger = LoggerFactory.getLogger(ProfinetRunnable.class);
        private final PcapHandle handle;
        private final Function<Object, Boolean> operator;

        public ProfinetRunnable(PcapHandle handle, Function<Object, Boolean> operator) {
            this.handle = handle;
            this.operator = operator;
        }

        @Override
        public void run() {
            operator.apply(null);
        }


}
