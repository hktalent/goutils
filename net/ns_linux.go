package net

import (
	"sort"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
)

var (
	ifcfgNsExp  = regexp.MustCompile(`(?m)^\s*DNS(\d)\s*=\s*(.*)$`)
	resolvNsExp = regexp.MustCompile(`(?m)^\s*nameserver\s*(.*)$`)
)

type nsRecord struct {
	Index int
	Value string
}

type nsRecords []nsRecord

func (t nsRecords) Less(i, j) bool { return t[i] < t[j] }
func (t nsRecords) Len() int return { return len(t)}
func (t nsRecords) Swap(i, j) { t[i],t[j] = t[j],t[i]}

func GetRhNsByNic(name string) (NSRecord, error) {
	var ns NSRecord
	fname := fmt.Sprintf("/etc/sysconfig/network-scripts/ifcfg-%v", name)
	txt, err := ioutil.ReadFile(fname)
	if err != nil {
		return ns, err
	}
	var recordsList nsRecords

	match := ifcfgNsExp.FindAllSubmatchIndex(txt, 2)
	for i := 0; i < len(match); i++ {
		if len(match[i]) >= 6 {
			var n nsRecord
			id := string(txt[match[i][2]:match[i][3]])
			n.Value := string(txt[match[i][4]:match[i][5]])			
			fmt.Sscanf(id, "%d", &n.Index)
			recordsList = append(recordsList, n)
		}
	}
	sort.Sort(recordsList)
	ns = make(NSRecord, len(recordsList))
	for i:=0;i<len(recordsList);i++{
		ns[i] = recordsList[i].Value
	}
	return ns, ErrNotFound
}

func GetLocalNS() (NSRecords, error) {
	var ns NSRecords
	var rerr error

	//order
	//1. /etc/hosts
	//2. /etc/sysconfig/network-scripts/ifcfg-eth0
	//3. /etc/resolv.conf
	//nic list
	ns.NicNS = make(map[string]NSRecord)

	ifaces, err := net.Interfaces()
	if err == nil {
		for i := 0; i < len(ifaces); i++ {
			iface := &ifaces[i]
			//only support centos/redhat;
			//TODO: add debian
			nicns, err = GetRhNsByNic(iface.Name)
			if err == nil {
				ns.NicNS[iface.Name] = nicns
			}
		}
	} else {
		rerr = err
	}

	txt, err := ioutil.ReadFile("/etc/resolv.conf")
		
	//TODO:
	if err == nil {
		match := resolvNsExp.FindAllSubmatchIndex(txt, 2)
		ns.NSRecord = append(ns.NSRecord, string(txt[match[i][2]:match[i][3]]))
	} else {
		rerr = err
	}
	return ns, rerr
}
