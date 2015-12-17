// +build linux

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Corporation

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

	"fmt"
	"os"
	"strings"
	"time"
	"io/ioutil"
	"path/filepath"

	"github.com/vektra/errors"

	str "github.com/intelsdi-x/snap-plugin-utilities/strings"
	"github.com/intelsdi-x/snap-plugin-utilities/ns"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"

)

const (
	ifaceinfo string = "/proc/net/dev"
	VENDOR          = "intel"
	OS              = "linux"
	PLUGIN          = "iface"
	VERSION         = 1
)

type ifacePlugin struct {
	stats map[string]interface{}
	host  string
}

func New() *ifacePlugin {
	fh, err := os.Open(ifaceinfo)

	if err != nil {
		return nil
	}
	defer fh.Close()

	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}

	iface := &ifacePlugin{stats: map[string]interface{}{}, host: host}

	return iface
}

func parseHeader(line string) ([]string, error ){

	l := strings.Split(line, "|")

	if len(l) < 3 {
		return nil, errors.New(fmt.Sprintf("Wrong header format {%s}", line))
	}

	header := strings.Fields(l[1])

	if len(header) < 8 {
		return nil, errors.New(fmt.Sprintf("Wrong header length. Expected 8 is {%d}", len(header)))
	}

	recv := make([]string, len(header))
	sent := make([]string, len(header))
	copy(recv, header)
	copy(sent, header)

	str.ForEach(
		recv,
		func (s string) string {
			return s + "_recv"
		})

	str.ForEach(
		sent,
		func (s string) string {
			return s + "_sent"
		})

	return append(recv, sent...), nil
}

func getStats(stats map[string]interface{}) error {

	content, err := ioutil.ReadFile(ifaceinfo)

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
			return errors.New(fmt.Sprintf("Wrong interface line format {%v}", len(ifdata)))
		}

		iname := strings.TrimSpace(ifdata[0])
		ivals := strings.Fields(ifdata[1])

		if len(ivals) != 16 {
			return errors.New(fmt.Sprintf("Wrong data length. Expected 16 is {%d}", len(ivals)))
		}

		istats := map[string]interface{}{}
		for i := 0; i < 16; i++ {
			stat := header[i]
			val := ivals[i]
			istats[stat] = val
		}

		stats[iname] = istats
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


func (iface *ifacePlugin) GetMetricTypes(_ plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	metricTypes := []plugin.PluginMetricType{}
	if err := getStats(iface.stats); err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	namespaces := []string{}

	err := ns.FromMap(iface.stats, filepath.Join(VENDOR, OS, PLUGIN), &namespaces)

	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaces {
		metricType := plugin.PluginMetricType{Namespace_: strings.Split(namespace, string(os.PathSeparator))}
		metricTypes = append(metricTypes, metricType)
	}
	return metricTypes, nil
}

func (iface *ifacePlugin) CollectMetrics(metricTypes []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {
	metrics := []plugin.PluginMetricType{}
	getStats(iface.stats)
	for _, metricType := range metricTypes {
		ns := metricType.Namespace()
		if len(ns) < 5 {
			return nil, errors.New(fmt.Sprintf("Namespace length is too short (len = %d)", len(ns)))
		}

		val := getMapValueByNamespace(iface.stats, ns[3:])

		metric := plugin.PluginMetricType{
			Namespace_: ns,
			Data_:      val,
			Source_:    iface.host,
			Timestamp_: time.Now(),
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (iface *ifacePlugin) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	return cpolicy.New(), nil
}
