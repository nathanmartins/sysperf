#! /usr/bin/python3

from __future__ import print_function

import json
from socket import inet_ntop, AF_INET, AF_INET6
from socket import ntohs
from struct import pack

from bcc import BPF

pretty_print = False

# load BPF program source code
bpf_text = open("ebpf.c").read()

# replace ports to be traced...
traced_port = [80, 443]
traced_port_if = ' && '.join(['dport != %d' % ntohs(port) for port in traced_port])
bpf_text = bpf_text.replace('FILTER_PORT', 'if (%s) { currsock.delete(&pid); return 0; }' % traced_port_if)

# initialize bcc's BPF object
bpf_obj = BPF(text=bpf_text)
bpf_obj.attach_kprobe(event="tcp_v4_connect", fn_name="trace_connect_entry")
bpf_obj.attach_kprobe(event="tcp_v6_connect", fn_name="trace_connect_entry")
bpf_obj.attach_kretprobe(event="tcp_v4_connect", fn_name="trace_connect_v4_return")
bpf_obj.attach_kretprobe(event="tcp_v6_connect", fn_name="trace_connect_v6_return")


def print_ipv4_event(cpu, data, size):
    event = bpf_obj["ipv4_events"].event(data)
    global start_ts
    if start_ts == 0:
        start_ts = event.ts_us

    parsed = {
        "took": (float(event.ts_us) - start_ts) / 1000000,
        "pid": event.pid,
        "ip version":  event.ip,
        "origin": inet_ntop(AF_INET, pack("I", event.saddr)),
        "destination": inet_ntop(AF_INET, pack("I", event.daddr)),
        "port": event.dport,
    }

    if pretty_print:
        print(json.dumps(parsed, indent=4, sort_keys=True))
    else:
        print(json.dumps(parsed))


def print_ipv6_event(cpu, data, size):
    event = bpf_obj["ipv6_events"].event(data)
    global start_ts
    if start_ts == 0:
        start_ts = event.ts_us
    parsed = {
        "took": (float(event.ts_us) - start_ts) / 1000000,
        "pid": event.pid,
        "ip version":  event.ip,
        "origin": inet_ntop(AF_INET6, pack("I", event.saddr)),
        "destination": inet_ntop(AF_INET6, pack("I", event.daddr)),
        "port": event.dport,
    }
    if pretty_print:
        print(json.dumps(parsed, indent=4, sort_keys=True))
    else:
        print(json.dumps(parsed))


if pretty_print:
    print(json.dumps({"message": "tracing tcp connections"}, indent=4, sort_keys=True))
else:
    print(json.dumps({"message": "tracing tcp connections"}))

start_ts = 0

# read events
bpf_obj["ipv4_events"].open_perf_buffer(print_ipv4_event)
bpf_obj["ipv6_events"].open_perf_buffer(print_ipv6_event)
while True:
    try:
        bpf_obj.perf_buffer_poll()
    except KeyboardInterrupt:
        exit()
