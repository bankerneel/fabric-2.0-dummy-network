#!/bin/bash

if [ -d "organizations/peerOrganizations" ]; then
  rm -Rf organizations/peerOrganizations && rm -Rf organizations/ordererOrganizations
fi

echo
echo "##########################################################"
echo "##### Generate certificates using Fabric CA's ############"
echo "##########################################################"

docker-compose -f docker-compose-ca.yaml up -d

. organizations/fabric-ca/registerEnroll.sh

sleep 10

echo "##########################################################"
echo "############ Create superadmin Identities ######################"
echo "##########################################################"

createsuperadmin

echo "##########################################################"
echo "############ Create users Identities ######################"
echo "##########################################################"

createusers

echo "##########################################################"
echo "############ Create Orderer Org Identities ###############"
echo "##########################################################"

createOrderer

echo
echo "Generate CCP files for superadmin and users"
./organizations/ccp-generate.sh