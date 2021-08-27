# Fabric-2.0-Dummy-Network

---

**Note: This README presumes that you have installed all the prerequisites of Hyperledger Fabric.**

---

## Step by Step Guide

Follow the step by step guide to start the network and set up the server

#### CA-NETWORK

Go to the ca-network directory and follow the commands:

```sh
sudo ./clean.sh
```

This command will clean all the containers which is exited.

```sh
docker-compose -f docker-compose-ca.yaml up -d
```

This command will start the Certificate Authority containers

```sh
./network.sh
```

This command will start generating the certificates of the above mentioned organizations.

#### FABRIC-NETWORK

Go to the network directory and follow the commands:

```sh
mkdir channel-artifacts && mkdir crypto-config
```

This command will create the channel-artifcats and crypto-config directory.

```sh
cp -rvf ../ca-network/organizations/fabric-ca ./crypto-config
```

This command will copy fabric-ca directory from ca-network to crypto-config directory.

```sh
cp -rvf ../ca-network/organizations/ordererOrganizations ./crypto-config
```

This command will copy ordererOrganizations directory from ca-network to crypto-config directory.

```sh
cp -rvf ../ca-network/organizations/peerOrganizations ./crypto-config
```

This command will copy peerOrganizations directory from ca-network to crypto-config directory.

```sh
./artifact.sh
```

This command will generate all the channel related artifacts for the Fabric Network.

```sh
docker-compose -f docker-compose-cli.yaml -f docker-compose-etcdraft2.yaml -f docker-compose-couchdb.yaml up -d
```

This command will start all the peers, orderers and couchdb containers for the network.

```sh
docker exec -it cli bash
```

This command will open the cli for executing the peers command

Once you are in the container you will be able to see the output like this

```sh
bash-5.0#
```

Now, following are the commands which needs to be executed from the cli

```sh
./scripts/channel.sh
```

## Explorer

```sh
 - Update the Private Keys Files in the following path in connection-profile/org-network.json file

 -- network/crypto-config/peerOrganizations/admin.org.com/users/Admin@superadmin.org.com/msp/keystore/YOUR-PRIVATE-KEY_sk

 - Once Done run command
 docker-compose up -d
```
