// +build linux

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015-2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package iface

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/intelsdi-x/snap-plugin-utilities/ns"
	"github.com/intelsdi-x/snap-plugin-utilities/str"
)

const (
	// Name of plugin
	Name = "iface"
	// Version of plugin
	Version = 6

	nsVendor = "intel"
	nsClass  = "procfs"
	nsType   = "iface"
)

// prefix in metric namespace
var prefix = []string{nsVendor, nsClass, nsType}

// added init PID 1 so that we take the global namespace
// to be sure to gather all possible interfaces wherever
// the plugin gets executed
var ifaceInfo = "/proc/1/net/dev"

// New creates instance of interface info plugin
func New() *IfaceCollector {
	logger := log.New()
	return &IfaceCollector{
		logger:    logger,
		proc_path: ifaceInfo,
		stats:     map[string]interface{}{},
	}
}

type IfaceCollector struct {
	stats     map[string]interface{}
	logger    *log.Logger
	proc_path string
}

// GetConfigPolicy returns config policy
func (iface *IfaceCollector) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	cpolicy := plugin.NewConfigPolicy()
	cpolicy.AddNewStringRule([]string{nsVendor, nsClass, Name},
		"proc_path", false, plugin.SetDefaultString("/proc"))
	return *cpolicy, nil
}

// GetMetricTypes returns list of available metric types
// It returns error in case retrieval was not successful
func (iface *IfaceCollector) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {
	err := iface.setProcPath(cfg)
	if err != nil {
		return nil, err
	}

	mts := []plugin.Metric{}

	if err = iface.getStats(); err != nil {
		return nil, err
	}

	namespaces := []string{}

	err = ns.FromMap(iface.stats, filepath.Join(nsVendor, nsClass, nsType), &namespaces)

	if err != nil {
		return nil, err
	}

	// List of terminal metric names
	mList := make(map[string]bool)
	for _, namespace := range namespaces {
		metric := plugin.Metric{Namespace: plugin.NewNamespace(strings.Split(namespace, "/")...)}
		ns := metric.Namespace
		// Interface metric (aka last element in namespace)
		mItem := ns[len(ns)-1]
		// Keep it if not already seen before
		if !mList[mItem.Value] {
			mList[mItem.Value] = true
			mts = append(mts, plugin.Metric{
				Namespace: plugin.NewNamespace(prefix...).
					AddDynamicElement("interface_name", "name of interface").
					AddStaticElement(mItem.Value),
				Description: "dynamic interface metric: " + mItem.Value,
			})
		}
	}
	return mts, nil
}

// CollectMetrics returns list of requested metric values
// It returns error in case retrieval was not successful
func (iface *IfaceCollector) CollectMetrics(mts []plugin.Metric) ([]plugin.Metric, error) {
	err := iface.setProcPath(mts[0].Config)
	if err != nil {
		return nil, err
	}

	metrics := []plugin.Metric{}

	if err := iface.getStats(); err != nil {
		return nil, err
	}
	curTime := time.Now()
	for _, mt := range mts {
		ns := mt.Namespace
		if len(ns) < 5 {
			return nil, fmt.Errorf("Namespace length is too short (len = %d)", len(ns))
		}
		if ns[len(ns)-2].Value == "*" {
			for itf, istats := range iface.stats {
				val := getMapValueByNamespace(istats.(map[string]interface{}), ns.Strings()[4:])
				if val != nil {
					ns1 := plugin.NewNamespace(createNamespace(itf, ns[len(ns)-1].Value)...)
					ns1[len(ns1)-2].Name = itf
					interfaceConfiguration, err := getInterfaceConfiguration(itf)
					if err != nil {
						return nil, err
					}

					metric := plugin.Metric{
						Namespace: ns1,
						Data:      val,
						Tags:      interfaceConfiguration,
						Timestamp: curTime,
					}
					metrics = append(metrics, metric)
				}
			}
		} else {
			val := getMapValueByNamespace(iface.stats, ns.Strings()[3:])
			metric := plugin.Metric{
				Namespace: ns,
				Data:      val,
				Timestamp: curTime,
			}
			metrics = append(metrics, metric)
		}
	}
	return metrics, nil
}

// Function to check properness of configuration parameter
// and set plugin attribute accordingly
func (iface *IfaceCollector) setProcPath(cfg plugin.Config) error {

	procPath, err := cfg.GetString("proc_path")
	if err == nil && len(procPath) > 0 {
		procPathStats, err := os.Stat(procPath)
		if err != nil {
			return err
		}
		if !procPathStats.IsDir() {
			return errors.New(fmt.Sprintf("%s is not a directory", procPath))
		}
		iface.proc_path = procPath + "/1/net/dev"
	}
	return nil
}

// createNamespace returns namespace slice of strings composed from: vendor, class, type and components of metric name
func createNamespace(itf string, name string) []string {
	var suffix = []string{itf, name}
	return append(prefix, suffix...)
}

func parseHeader(line string) ([]string, error) {

	l := strings.Split(line, "|")

	if len(l) < 3 {
		return nil, fmt.Errorf("Wrong header format {%s}", line)
	}

	header := strings.Fields(l[1])

	if len(header) < 8 {
		return nil, fmt.Errorf("Wrong header length. Expected 8 is {%d}", len(header))
	}

	recv := make([]string, len(header))
	sent := make([]string, len(header))
	copy(recv, header)
	copy(sent, header)

	str.ForEach(
		recv,
		func(s string) string {
			return s + "_recv"
		})

	str.ForEach(
		sent,
		func(s string) string {
			return s + "_sent"
		})

	return append(recv, sent...), nil
}

func (iface *IfaceCollector) getStats() error {
	path := iface.proc_path
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")

	header, err := parseHeader(lines[1])

	if err != nil {
		return err
	}

	for _, line := range lines[2:] {

		if line == "" {
			continue
		}

		ifdata := strings.Split(line, ":")

		if len(ifdata) != 2 {
			return fmt.Errorf("Wrong interface line format {%v}", len(ifdata))
		}

		iname := strings.TrimSpace(ifdata[0])
		ivals := strings.Fields(ifdata[1])

		if len(ivals) != len(header) {
			return fmt.Errorf("Wrong data length. Expected {%d} is {%d}", len(header), len(ivals))
		}

		istats := map[string]interface{}{}
		for i := 0; i < 16; i++ {
			stat := header[i]
			val, err := strconv.ParseInt(ivals[i], 10, 64)
			if err != nil {
				log.WithFields(log.Fields{
					"iname":  iname,
					"stat":   stat,
					"strVal": ivals[i],
					"val":    val,
				}).Warn("Cannot parse metric value to number, metric value saved as -1, ", err)
				val = -1
			}
			istats[stat] = val
		}

		iface.stats[iname] = istats
	}

	return nil
}

func getMapValueByNamespace(map_ map[string]interface{}, ns []string) interface{} {
	if len(ns) == 0 {
		fmt.Println("Namespace length equal to zero!")
		return nil
	}

	current := ns[0]

	if len(ns) == 1 {
		if val, ok := map_[current]; ok {
			return val
		}
		return nil
	}

	if v, ok := map_[current].(map[string]interface{}); ok {
		return getMapValueByNamespace(v, ns[1:])
	}

	return nil
}

func getInterfaceConfiguration(ifaceName string) (map[string]string, error) {
	interfaceConfig, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, err
	}
	tags := make(map[string]string)
	tags["hardware_addr"] = string(interfaceConfig.HardwareAddr)
	tags["mtu"] = strconv.Itoa(interfaceConfig.MTU)
	return tags, nil
}
