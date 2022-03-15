#! /usr/bin/python3

from __future__ import print_function

from socket import inet_ntop, AF_INET, AF_INET6
from socket import ntohs
from struct import pack

from bcc import BPF
from bcc.utils import printb

debug = 0

# define BPF program
bpf_text = open("raw.c").read()

struct_init = {'ipv4':
                   {'count':
                        """
                        struct ipv4_flow_key_t flow_key = {};
                        flow_key.saddr = skp->__sk_common.skc_rcv_saddr;
                        flow_key.daddr = skp->__sk_common.skc_daddr;
                        flow_key.dport = ntohs(dport);
                        ipv4_count.increment(flow_key);""",
                    'trace':
                        """
                        struct ipv4_data_t data4 = {.pid = pid, .ip = ipver};
                        data4.uid = bpf_get_current_uid_gid();
                        data4.ts_us = bpf_ktime_get_ns() / 1000;
                        data4.saddr = skp->__sk_common.skc_rcv_saddr;
                        data4.daddr = skp->__sk_common.skc_daddr;
                        data4.dport = ntohs(dport);
                        bpf_get_current_comm(&data4.task, sizeof(data4.task));
                        ipv4_events.perf_submit(ctx, &data4, sizeof(data4));"""
                    },
               'ipv6':
                   {'count':
                        """
                        struct ipv6_flow_key_t flow_key = {};
                        bpf_probe_read(&flow_key.saddr, sizeof(flow_key.saddr),
                            skp->__sk_common.skc_v6_rcv_saddr.in6_u.u6_addr32);
                        bpf_probe_read(&flow_key.daddr, sizeof(flow_key.daddr),
                            skp->__sk_common.skc_v6_daddr.in6_u.u6_addr32);
                        flow_key.dport = ntohs(dport);
                        ipv6_count.increment(flow_key);""",
                    'trace':
                        """
                        struct ipv6_data_t data6 = {.pid = pid, .ip = ipver};
                        data6.uid = bpf_get_current_uid_gid();
                        data6.ts_us = bpf_ktime_get_ns() / 1000;
                        bpf_probe_read(&data6.saddr, sizeof(data6.saddr),
                            skp->__sk_common.skc_v6_rcv_saddr.in6_u.u6_addr32);
                        bpf_probe_read(&data6.daddr, sizeof(data6.daddr),
                            skp->__sk_common.skc_v6_daddr.in6_u.u6_addr32);
                        data6.dport = ntohs(dport);
                        bpf_get_current_comm(&data6.task, sizeof(data6.task));
                        ipv6_events.perf_submit(ctx, &data6, sizeof(data6));"""
                    }
               }

bpf_text = bpf_text.replace("IPV4_CODE", struct_init['ipv4']['trace'])
bpf_text = bpf_text.replace("IPV6_CODE", struct_init['ipv6']['trace'])

dports = [80, 443]
dports_if = ' && '.join(['dport != %d' % ntohs(dport) for dport in dports])
bpf_text = bpf_text.replace('FILTER_PORT', 'if (%s) { currsock.delete(&pid); return 0; }' % dports_if)

bpf_text = bpf_text.replace('FILTER_PID', '')
bpf_text = bpf_text.replace('FILTER_PORT', '')
bpf_text = bpf_text.replace('FILTER_UID', '')

# initialize BPF
b = BPF(text=bpf_text)
b.attach_kprobe(event="tcp_v4_connect", fn_name="trace_connect_entry")
b.attach_kprobe(event="tcp_v6_connect", fn_name="trace_connect_entry")
b.attach_kretprobe(event="tcp_v4_connect", fn_name="trace_connect_v4_return")
b.attach_kretprobe(event="tcp_v6_connect", fn_name="trace_connect_v6_return")

def print_ipv4_event(cpu, data, size):
    event = b["ipv4_events"].event(data)
    global start_ts
    # if args.timestamp:
    if start_ts == 0:
        start_ts = event.ts_us
    printb(b"%-9.3f" % ((float(event.ts_us) - start_ts) / 1000000), nl="")
    # if args.print_uid:
    #     printb(b"%-6d" % event.uid, nl="")
    printb(b"%-6d %-12.12s %-2d %-16s %-16s %-4d" % (event.pid,
        event.task, event.ip,
        inet_ntop(AF_INET, pack("I", event.saddr)).encode(),
        inet_ntop(AF_INET, pack("I", event.daddr)).encode(), event.dport))


def print_ipv6_event(cpu, data, size):
    event = b["ipv6_events"].event(data)
    global start_ts
    # if args.timestamp:
    if start_ts == 0:
        start_ts = event.ts_us
    printb(b"%-9.3f" % ((float(event.ts_us) - start_ts) / 1000000), nl="")
    # if args.print_uid:
    #     printb(b"%-6d" % event.uid, nl="")
    printb(b"%-6d %-12.12s %-2d %-16s %-16s %-4d" % (event.pid,
        event.task, event.ip,
        inet_ntop(AF_INET6, event.saddr).encode(), inet_ntop(AF_INET6, event.daddr).encode(),
        event.dport))


print("Tracing connect ... Hit Ctrl-C to end")
print("%-9s" % ("TIME(s)"), end="")
print("%-6s %-12s %-2s %-16s %-16s %-4s" % ("PID", "COMM", "IP", "SADDR",
                                            "DADDR", "DPORT"))

start_ts = 0

# read events
b["ipv4_events"].open_perf_buffer(print_ipv4_event)
b["ipv6_events"].open_perf_buffer(print_ipv6_event)
while 1:
    try:
        b.perf_buffer_poll()
    except KeyboardInterrupt:
        exit()
