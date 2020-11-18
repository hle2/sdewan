# Steps for running v1 API microservices

### Steps to install packages
**1. Create namespace for SDEWAN Central Controller v1Microservices**

`$ kubectl create namespace sdewan-system`

**2. Create Databases used by SDEWAN Central Controller v1 Microservices for Etcd and Mongo**

`$ kubectl apply -f sccdb.yaml -n sdewan-system`

**3. create SDEWAN Central Controller v1 Microservices**

`$ kubectl apply -f scc.yaml -n sdewan-system`
`$ kubectl apply -f scc_rsync.yaml -n sdewan-system`

**4. install monitor resources**

`$ kubectl apply -f monitor-deploy.yaml`