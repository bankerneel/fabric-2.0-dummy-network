#!/bin/bash

export CORE_PEER_LOCALMSPID=superadminMSP
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/superadmin.sosh.com/peers/peer0.superadmin.sosh.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/superadmin.sosh.com/users/Admin@superadmin.sosh.com/msp
export CORE_PEER_ADDRESS=peer0.superadmin.sosh.com:7051
export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/sosh.com/orderers/orderer.sosh.com/msp/tlscacerts/tlsca.sosh.com-cert.pem

peer channel create -f channel-artifacts/soshchannel.tx -c soshchannel -o orderer.sosh.com:7050 --tls --cafile $ORDERER_CA
peer channel join -b soshchannel.block
peer channel update -o orderer.sosh.com:7050 -c soshchannel -f channel-artifacts/superadminAnchor.tx --tls --cafile $ORDERER_CA

export CORE_PEER_LOCALMSPID=usersMSP
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/users.sosh.com/peers/peer0.users.sosh.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/users.sosh.com/users/Admin@users.sosh.com/msp
export CORE_PEER_ADDRESS=peer0.users.sosh.com:9051
export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/sosh.com/orderers/orderer.sosh.com/msp/tlscacerts/tlsca.sosh.com-cert.pem

peer channel join -b soshchannel.block
peer channel update -o orderer.sosh.com:7050 -c soshchannel -f channel-artifacts/usersAnchor.tx --tls --cafile $ORDERER_CA