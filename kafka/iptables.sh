#!/bin/bash

# Allow internal communication between Kafka Brokers
iptables -A INPUT -p tcp -m tcp --dport 9092 -j ACCEPT
iptables -A INPUT -p tcp -m tcp --dport 9093 -j ACCEPT

# Allow Zookeeper internal communication
iptables -A INPUT -p tcp -m tcp --dport 2181 -j ACCEPT

# Allow communication between Kafka and Zookeeper
iptables -A INPUT -p tcp -m tcp --dport 2181 -j ACCEPT
iptables -A OUTPUT -p tcp -m tcp --sport 9092 -j ACCEPT

# Allow operating Kafka using Kafka CLI tools
iptables -A INPUT -p tcp -m tcp --dport 9092 -j ACCEPT
iptables -A OUTPUT -p tcp -m tcp --sport 9092 -j ACCEPT

# Allow communication between network interfaces (if there are multiple network interfaces)
iptables -A INPUT -i lo -j ACCEPT
iptables -A INPUT -i enp0s31f6 -j ACCEPT

# Allow Ping to pass through
iptables -A INPUT -p icmp --icmp-type echo-request -j ACCEPT

# Allow SSH connections
iptables -A INPUT -p tcp -m tcp --dport 22 -j ACCEPT

# Allow response connections
iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT

# Protect Kafka security and block unwanted connections
iptables -P INPUT DROP

# Allow hosts to initiate connections
iptables -A OUTPUT -j ACCEPT
