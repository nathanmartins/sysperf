# -*- mode: ruby -*-
# vi: set ft=ruby :


Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/impish64"
  config.vm.provision "shell", inline: <<-SHELL
    apt-get update
    apt-get install -y bpfcc-tools  linux-headers-$(uname -r) clang build-essential libbpf-dev libbpf0  linux-tools-common linux-tools-generic golang
    curl https://get.docker.com | bash
  SHELL
end