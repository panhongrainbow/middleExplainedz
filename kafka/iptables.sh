# 允許Kafka Broker內網通信
iptables -A INPUT -p tcp -m tcp --dport 9092 -j ACCEPT

# 允許Kafka Broker進行REPLICATION
iptables -A INPUT -p tcp -m tcp --dport 9093 -j ACCEPT

# 允許Zookeeper內部通信 
iptables -A INPUT -p tcp -m tcp --dport 2181 -j ACCEPT  

# 允許Kafka與Zookeeper通信
iptables -A INPUT -p tcp -m tcp --dport 2181 -j ACCEPT
iptables -A OUTPUT -p tcp -m tcp --sport 9092 -j ACCEPT

# 允許使用Kafka CLI工具進行操作  
iptables -A INPUT -p tcp -m tcp --dport 9092 -j ACCEPT 
iptables -A OUTPUT -p tcp -m tcp --sport 9092 -j ACCEPT

# 允許網路接口之間進行通信(如果有多個網路接口)
iptables -A INPUT -i lo -j ACCEPT 
iptables -A INPUT -i enp0s31f6 -j ACCEPT

# 允許Ping通過 
iptables -A INPUT -p icmp --icmp-type echo-request -j ACCEPT

# 允許SSH連接
iptables -A INPUT -p tcp -m tcp --dport 22 -j ACCEPT

# 允許回應連接
iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT

# 保護Kafka安全性,阻止不需要的連接  
iptables -P INPUT DROP  

# 允許主機發起連接
iptables -A OUTPUT -j ACCEPT
