# snap collector plugin - interface

## Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Description (optional)
----------|-----------------------
/intel/procfs/iface/\<interface_name\>/bytes_recv | The total number of bytes of data received by the interface
/intel/procfs/iface/\<interface_name\>/bytes_sent | The total number of bytes of data transmitted by the interface
/intel/procfs/iface/\<interface_name\>/compressed_recv | The number of compressed packets received by the device driver
/intel/procfs/iface/\<interface_name\>/compressed_sent | The number of compressed packets transmitted by the device driver
/intel/procfs/iface/\<interface_name\>/drop_recv | The total number of packets dropped by the device driver while receiving
/intel/procfs/iface/\<interface_name\>/drop_sent | The total number of packets dropped by the device driver while transmitting
/intel/procfs/iface/\<interface_name\>/errs_recv | The total number of receive errors detected by the device driver
/intel/procfs/iface/\<interface_name\>/errs_sent | The total number of transmit errors detected by the device driver
/intel/procfs/iface/\<interface_name\>/fifo_recv | The number of FIFO buffer errors while receiving
/intel/procfs/iface/\<interface_name\>/fifo_sent | The number of FIFO buffer errors  while transmitting
/intel/procfs/iface/\<interface_name\>/frame_recv | The number of packet framing errors while receiving
/intel/procfs/iface/\<interface_name\>/frame_sent | The number of packet framing errors while transmitting
/intel/procfs/iface/\<interface_name\>/multicast_recv | The number of multicast frames received by the device driver
/intel/procfs/iface/\<interface_name\>/multicast_sent | The number of multicast frames transmitted by the device driver
/intel/procfs/iface/\<interface_name\>/packets_recv | The total number of packets of data received by the interface
/intel/procfs/iface/\<interface_name\>/packets_sent | The total number of packets of data transmitted by the interface