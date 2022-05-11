# sysperf - USE method visualization tool

## What is the USE method?
The USE Method can be summarized as:

For every resource, check utilization, saturation, and errors.
It's intended to be used early in a performance investigation, to identify systemic bottlenecks.

Terminology definitions:

resource: all physical server functional components (CPUs, disks, busses, ...) [1]
utilization: the average time that the resource was busy servicing work [2]
saturation: the degree to which the resource has extra work which it can't service, often queued
errors: the count of error events

The metrics are usually expressed in the following terms:

utilization: as a percent over a time interval. eg, "one disk is running at 90% utilization".
saturation: as a queue length. eg, "the CPUs have an average run queue length of four".
errors: scalar counts. eg, "this network interface has had fifty late collisions".

## Requirements:

- Must be root to run perf tool
- perf (linux-tools-common linux-tools-generic linux-tools-`uname -r`)
- bcc tools (bpfcc-tools)

## Currently, working on CPU 

- [x] Raw metric for CPU utilization: system-wide average
- [x] Raw metric for CPU saturation: run-queue length or scheduler latency
- [ ] show saturation in frontend
- [ ] show latency in frontend