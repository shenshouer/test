# Use goreman to run `go get github.com/mattn/goreman` 
etcd1: etcd -name infra1 -listen-client-urls http://localhost:4001 -advertise-client-urls http://localhost:4001 -listen-peer-urls http://localhost:7001 -initial-advertise-peer-urls http://localhost:7001 -initial-cluster-token etcd-cluster-1 -initial-cluster 'infra1=http://localhost:7001,infra2=http://localhost:7002,infra3=http://localhost:7003' -initial-cluster-state new 
etcd2: etcd -name infra2 -listen-client-urls http://localhost:4002 -advertise-client-urls http://localhost:4002 -listen-peer-urls http://localhost:7002 -initial-advertise-peer-urls http://localhost:7002 -initial-cluster-token etcd-cluster-1 -initial-cluster 'infra1=http://localhost:7001,infra2=http://localhost:7002,infra3=http://localhost:7003' -initial-cluster-state new 
etcd3: etcd -name infra3 -listen-client-urls http://localhost:4003 -advertise-client-urls http://localhost:4003 -listen-peer-urls http://localhost:7003 -initial-advertise-peer-urls http://localhost:7003 -initial-cluster-token etcd-cluster-1 -initial-cluster 'infra1=http://localhost:7001,infra2=http://localhost:7002,infra3=http://localhost:7003' -initial-cluster-state new 
proxy: etcd -name proxy1 -proxy=on -listen-client-urls http://localhost:8080 -initial-cluster 'infra1=http://localhost:7001,infra2=http://localhost:7002,infra3=http://localhost:7003' 

