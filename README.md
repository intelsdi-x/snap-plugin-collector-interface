# snap collector plugin - interface
This plugin collects metrics from /proc/interface kernel interface about the traffic (octets per second), packets per second and errors of interfaces.  

It's used in the [snap framework](http://github.com:intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license-and-authors)
6. [Acknowledgements](#acknowledgements)

## Getting Started
### System Requirements
* [golang 1.4+](https://golang.org/dl/)

### Operating systems
All OSs currently supported by snap:
* Linux/amd64

### Installation
#### Download interface plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [GitHub Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-interface  
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-interface.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `/build/rootfs/`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Ensure `$SNAP_PATH` is exported  
`export SNAP_PATH=$GOPATH/src/github.com/intelsdi-x/snap/build`

## Documentation

### Collected Metrics
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

### Examples
Example running interface, passthru processor, and writing data to a file.

This is done from the snap directory.

In one terminal window, open the snap daemon (in this case with logging set to 1 and trust disabled):
```
$ $SNAP_PATH/bin/snapd -l 1 -t 0
```

In another terminal window:
Load interface plugin
```
$ $SNAP_PATH/bin/snapctl plugin load snap-plugin-collector-interface
```
See available metrics for your system
```
$ $SNAP_PATH/bin/snapctl metric list
```

Create a task manifest file (e.g. `iface-file.json`):

Put your desired interface name instead of "\<interface_name\>"    
    
```json
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "1s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/intel/procfs/iface/<interface_name>/bytes_recv": {},
                "/intel/procfs/iface/<interface_name>/bytes_sent": {}, 
                "/intel/procfs/iface/<interface_name>/err_recv": {},
                "/intel/procfs/iface/<interface_name>/fifo_recv": {} 
            },
            "config": {
                "/intel/mock": {
                    "password": "secret",
                    "user": "root"
                }
            },
            "process": [
                {
                    "plugin_name": "passthru",
                    "process": null,
                    "publish": [
                        {                         
                            "plugin_name": "file",
                            "config": {
                                "file": "/tmp/published_interface"
                            }
                        }
                    ],
                    "config": null
                }
            ],
            "publish": null
        }
    }
}
```

Load passthru plugin for processing:
```
$ $SNAP_PATH/bin/snapctl plugin load build/plugin/snap-processor-passthru
Plugin loaded
Name: passthru
Version: 1
Type: processor
Signed: false
Loaded Time: Fri, 20 Nov 2015 11:44:03 PST
```

Load file plugin for publishing:
```
$ $SNAP_PATH/bin/snapctl plugin load build/plugin/snap-publisher-file
Plugin loaded
Name: file
Version: 3
Type: publisher
Signed: false
Loaded Time: Fri, 20 Nov 2015 11:41:39 PST
```

Create task:
```
$ $SNAP_PATH/bin/snapctl task create -t examples/tasks/mem-file.json
Using task manifest to create task
Task created
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
Name: Task-02dd7ff4-8106-47e9-8b86-70067cd0a850
State: Running
```

Stop task:
```
$ $SNAP_PATH/bin/snapctl task stop 02dd7ff4-8106-47e9-8b86-70067cd0a850
Task stopped:
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-interface/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-interface/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@MarcinKrolik](https://github.com/marcin-krolik/)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.