# Steps for running v1 API microservices

### Steps to install packages
**1. Install cert-manager**

`$ kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v0.16.1/cert-manager.yaml`

**2. Create namespace for SDEWAN Central Controller v1Microservices**

`$ kubectl create namespace sdewan-system`

**3. Create Databases used by SDEWAN Central Controller v1 Microservices for Etcd and Mongo**

`$ kubectl apply -f scc_db.yaml -n sdewan-system`

**4. create SDEWAN Central Controller v1 Microservices**

`$ kubectl apply -f scc.yaml -n sdewan-system`
`$ kubectl apply -f scc_rsync.yaml -n sdewan-system`

**5. install monitor resources**

`$ ./monitor-deploy.sh