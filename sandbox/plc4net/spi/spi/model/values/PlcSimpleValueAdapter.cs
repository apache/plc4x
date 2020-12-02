namespace org.apache.plc4net.spi.model.values
{
    public abstract class PlcSimpleValueAdapter : PlcValueAdapter
    {
       
        public  bool IsSimple()
        {
            return true;
        }

        public int GetLength()
        {
            return 1;
        }
        
    }
}