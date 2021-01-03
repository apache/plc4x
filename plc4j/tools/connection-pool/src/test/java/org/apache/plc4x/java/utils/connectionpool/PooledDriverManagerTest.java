package org.apache.plc4x.java.utils.connectionpool;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.utils.connectionpool.CachedDriverManager;
import org.apache.plc4x.java.utils.connectionpool.PooledDriverManager;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;

/**
 *
 * @author julian
 * Created by julian on 24.02.20
 */
class PooledDriverManagerTest implements WithAssertions {

    @Test
    void getCachedDriverManager() throws PlcConnectionException {
        CachedDriverManager mock = Mockito.mock(CachedDriverManager.class, Mockito.RETURNS_DEEP_STUBS);
        PooledDriverManager driverManager = new PooledDriverManager(key -> mock);

        assertThat(driverManager.getCachedManagers().size()).isEqualTo(0);
        PlcConnection connection = driverManager.getConnection("abc");

        assertThat(driverManager.getCachedManagers())
            .containsValue(mock)
            .containsKey("abc")
            .hasSize(1);

        verify(mock, times(1)).getConnection("abc");
    }
}