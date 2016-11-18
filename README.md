# Snap collector plugin - interface
This plugin collects metrics from /proc/interface kernel interface about the traffic (octets per second), packets per second and errors of interfaces.  

It's used in the [Snap framework](https://github.com/intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
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
* [golang 1.5+](https://golang.org/dl/)

### Operating systems
All OSs currently supported by Snap:
* Linux/amd64

### Installation
#### Download the plugin binary:

You can get the pre-built binaries for your OS and architecture from the plugin's [GitHub Releases](https://github.com/intelsdi-x/snap-plugin-collector-interface/releases) page. Download the plugin from the latest release and load it into `snapd` (`/opt/snap/plugins` is the default location for Snap packages).

#### To build the plugin binary:

Fork https://github.com/intelsdi-x/snap-plugin-collector-interface
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-interface.git
```

Build the Snap interface plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `./build/`

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap#getting-started)
* Load the plugin and create a task, see example in [Examples](#examples).

Configuration parameters:
- `proc_path` path to '1/dev/net' file (helpful for running plugin in Docker container)
In fact this file should NOT be <proc_path>/dev/net because this is a symlink and therefore
will resolve in the exact same way both within or outside of a container. Therefore PID 1
has to be used

## Documentation

### Collected Metrics
List of collected metrics is described in [METRICS.md](METRICS.md).

Plugin reads metrics from `<proc_path>/1/dev/net` file. (see comment above)
Path to above file can be provided in configuration in task manifest as `proc_path`. If configuration is not provided, plugin will try
to read from default location which is `/proc/1/dev/net`.

### Examples
Example of running Snap interface collector and writing data to file.

Ensure [Snap daemon is running](https://github.com/intelsdi-x/snap#running-snap):
* initd: `service snap-telemetry start`
* systemd: `systemctl start snap-telemetry`
* command line: `snapd -l 1 -t 0 &`

Download and load Snap plugins:
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-interface/latest/linux/x86_64/snap-plugin-collector-interface
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
$ chmod 755 snap-plugin-*
$ snapctl plugin load snap-plugin-collector-interface
$ snapctl plugin load snap-plugin-publisher-file
```

See all available metrics:

```
$ snapctl metric list
```

Download an [example task file](examples/task/iface-file.json) and load it:
```
$ curl -sfLO https://raw.githubusercontent.com/intelsdi-x/snap-plugin-collector-interface/master/examples/task/iface-file.json
$ snapctl task create -t iface-file.json
Using task manifest to create task
Task created
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
Name: Task-02dd7ff4-8106-47e9-8b86-70067cd0a850
State: Running
```

See realtime output from `snapctl task watch <task_id>` (CTRL+C to exit)

This data is published to a file `/tmp/published_interface` per task specification

Stop task:
```
$ $SNAP_PATH/bin/snapctl task stop 02dd7ff4-8106-47e9-8b86-70067cd0a850
Task stopped:
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-interface/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-interface/pulls).

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap.

To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[Snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [Marcin Krolik](https://github.com/marcin-krolik/)
