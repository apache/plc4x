/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package utils

import (
	"bytes"
	"context"
	"net"
	"sync"
	"time"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/pcap"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func GetIPAddresses(ctx context.Context, netInterface net.Interface, useArpBasedScan bool) (foundIps chan net.IP, err error) {
	foundIps = make(chan net.IP, 65536)
	addrs, err := netInterface.Addrs()
	if err != nil {
		return nil, errors.Wrap(err, "Error getting addresses")
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error().Msgf("panic-ed %v", err)
			}
		}()
		wg := &sync.WaitGroup{}
		for _, address := range addrs {
			// Check if context has been cancelled before continuing
			select {
			case <-ctx.Done():
				return
			default:
			}

			var ipnet *net.IPNet
			switch v := address.(type) {
			case *net.IPAddr:
				ipnet = &net.IPNet{IP: v.IP, Mask: v.IP.DefaultMask()}
			case *net.IPNet:
				ipnet = v
			default:
				continue
			}

			// Skip loop-back and IPv6
			if ipnet.IP.IsLoopback() || ipnet.IP.To4() == nil {
				continue
			}

			log.Debug().Stringer("IP", ipnet.IP).Stringer("Mask", ipnet.Mask).Msg("Expanding local subnet")
			if useArpBasedScan {
				if err := lockupIpsUsingArp(ctx, netInterface, ipnet, foundIps, wg); err != nil {
					log.Error().Err(err).Msg("failing to resolve using arp scan. Falling back to ip based scan")
					useArpBasedScan = false
				}
			}
			if !useArpBasedScan {
				if err := lookupIps(ctx, ipnet, foundIps, wg); err != nil {
					log.Error().Err(err).Msg("error looking up ips")
				}
			}
		}
		wg.Wait()
		log.Trace().Msg("Closing found ips channel")
		close(foundIps)
	}()
	return foundIps, nil
}

// As PING operations might be blocked by a firewall, responding to ARP packets is mandatory for IP based
// systems. So we are using an ARP scan to resolve the ethernet hardware addresses of each possible ip in range
// Only for devices that respond will we schedule a discovery.
func lockupIpsUsingArp(ctx context.Context, netInterface net.Interface, ipNet *net.IPNet, foundIps chan net.IP, wg *sync.WaitGroup) error {
	// We add on signal for error handling
	wg.Add(1)
	go func() { wg.Done() }()
	log.Debug().Msgf("Scanning for alive IP addresses for interface '%s' and net: %s", netInterface.Name, ipNet)
	// First find the pcap device name for the given interface.
	allDevs, _ := pcap.FindAllDevs()
	var devName string
	for _, dev := range allDevs {
		for _, devAddress := range dev.Addresses {
			if devAddress.IP.Equal(ipNet.IP) {
				devName = dev.Name
				break
			}
		}
	}
	if len(devName) == 0 {
		log.Error().Interface("allDevs", allDevs).Str("ip", ipNet.IP.String()).Msg("Device for discovery not found")
		return errors.New("Device for discovery not found")
	}

	// Open up a pcap handle for packet reads/writes.
	handle, err := pcap.OpenLive(devName, 65536, true, pcap.BlockForever)
	if err != nil {
		return errors.Wrap(err, "Error opening network interface")
	}

	// Start up a goroutine to read in packet data.
	stop := make(chan struct{}, 1)
	// As we don't know how much the handler will find we use a value of 1 and set that to done after the 10 sec in the cleanup function directly after
	wg.Add(1)
	// Handler for processing incoming ARP responses.
	go func(handle *pcap.Handle, iface net.Interface, stop chan struct{}) {
		defer func() {
			if err := recover(); err != nil {
				log.Error().Msgf("panic-ed %v", err)
			}
		}()
		src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
		in := src.Packets()
		for {
			var packet gopacket.Packet
			select {
			case <-stop:
				return
			case packet = <-in:
				if packet == nil {
					continue
				}
				arpLayer := packet.Layer(layers.LayerTypeARP)
				if arpLayer == nil {
					continue
				}
				arp := arpLayer.(*layers.ARP)
				// Filter our messages originating from us.
				if arp.Operation != layers.ARPReply || bytes.Equal(iface.HardwareAddr, arp.SourceHwAddress) {
					continue
				}
				// Schedule a discovery operation for this ip.
				ip := net.IP(arp.SourceProtAddress)
				log.Trace().Msgf("Scheduling discovery for IP %s", ip)
				timeout := time.NewTimer(2 * time.Second)
				go func(ip net.IP) {
					defer CleanupTimer(timeout)
					select {
					case <-ctx.Done():
					case foundIps <- DuplicateIP(ip):
					case <-timeout.C:
					}
				}(DuplicateIP(ip))
			}
		}
	}(handle, netInterface, stop)
	// Make sure we clean up after 10 seconds.
	defer func() {
		go func() {
			wg.Done()
			time.Sleep(10 * time.Second)
			handle.Close()
			close(stop)
		}()
	}()
	writeArp := func(handle *pcap.Handle, iface net.Interface, addr net.IPNet) error {
		// Set up all the layers' tags we can.
		eth := layers.Ethernet{
			SrcMAC:       iface.HardwareAddr,
			DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			EthernetType: layers.EthernetTypeARP,
		}
		arp := layers.ARP{
			AddrType:          layers.LinkTypeEthernet,
			Protocol:          layers.EthernetTypeIPv4,
			HwAddressSize:     6,
			ProtAddressSize:   4,
			Operation:         layers.ARPRequest,
			SourceHwAddress:   []byte(iface.HardwareAddr),
			SourceProtAddress: []byte(addr.IP.To4()),
			DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
		}
		// Set up buffer and options for serialization.
		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}
		log.Debug().Msgf("Sending ARP requests to all devices in network: %s", addr.String())
		// Send one ARP packet for every possible address.
		for ip := IncrementIP(addr.IP.Mask(ipNet.Mask)); addr.Contains(ip) && addr.Contains(IncrementIP(DuplicateIP(ip))); ip = IncrementIP(ip) {
			// Check if context has been cancelled before continuing
			select {
			case <-ctx.Done():
				return nil
			default:
			}
			arp.DstProtAddress = ip
			if err := gopacket.SerializeLayers(buf, opts, &eth, &arp); err != nil {
				return err
			}
			if err := handle.WritePacketData(buf.Bytes()); err != nil {
				return err
			}
		}
		return nil
	}
	// Write our scan packets out to the handle.
	if err := writeArp(handle, netInterface, *ipNet); err != nil {
		log.Printf("error writing packets on %v: %v", netInterface.Name, err)
		return err
	}
	return nil
}

// Simply takes the IP address and the netmask and schedules one discovery task for every possible IP
func lookupIps(ctx context.Context, ipnet *net.IPNet, foundIps chan net.IP, wg *sync.WaitGroup) error {
	log.Debug().Msgf("Scanning all IP addresses for network: %s", ipnet)
	// expand CIDR-block into one target for each IP
	// Remark: The last IP address a network contains is a special broadcast address. We don't want to check that one.
	for ip := IncrementIP(ipnet.IP.Mask(ipnet.Mask)); ipnet.Contains(ip) && ipnet.Contains(IncrementIP(DuplicateIP(ip))); ip = IncrementIP(ip) {
		// Check if context has been cancelled before continuing
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		wg.Add(1)
		timeout := time.NewTimer(2 * time.Second)
		go func(ip net.IP) {
			defer func() { wg.Done() }()
			defer CleanupTimer(timeout)
			select {
			case <-ctx.Done():
			case foundIps <- ip:
			case <-timeout.C:
			}
		}(DuplicateIP(ip))
		log.Trace().Stringer("IP", ip).Msg("Expanded CIDR")
	}

	log.Debug().Stringer("net", ipnet).Msg("Done expanding CIDR")

	return nil
}

func IncrementIP(ip net.IP) net.IP {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}

	return ip
}

func DuplicateIP(ip net.IP) net.IP {
	dup := make(net.IP, len(ip))
	copy(dup, ip)
	return dup
}
