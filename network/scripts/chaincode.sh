peer lifecycle chaincode package wallet.tar.gz --path github.com/hyperledger/fabric-samples/chaincode/wallet --lang golang --label $1

export CORE_PEER_LOCALMSPID=superadminMSP
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/superadmin.sosh.com/peers/peer0.superadmin.sosh.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/superadmin.sosh.com/users/Admin@superadmin.sosh.com/msp
export CORE_PEER_ADDRESS=peer0.superadmin.sosh.com:7051
export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/sosh.com/orderers/orderer.sosh.com/msp/tlscacerts/tlsca.sosh.com-cert.pem

peer lifecycle chaincode install wallet.tar.gz
export CCID=$(peer lifecycle chaincode queryinstalled | cut -d ' ' -f 3 | sed s/.$// | grep $1)
peer lifecycle chaincode approveformyorg --package-id $CCID --channelID soshchannel --name wallet --version 1 --sequence $2 --waitForEvent --tls --cafile $ORDERER_CA
peer lifecycle chaincode checkcommitreadiness --channelID soshchannel --name wallet --version 1  --sequence $2 --tls --cafile $ORDERER_CA

export CORE_PEER_LOCALMSPID=usersMSP
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/users.sosh.com/peers/peer0.users.sosh.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/users.sosh.com/users/Admin@users.sosh.com/msp
export CORE_PEER_ADDRESS=peer0.users.sosh.com:9051
export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/sosh.com/orderers/orderer.sosh.com/msp/tlscacerts/tlsca.sosh.com-cert.pem

peer lifecycle chaincode install wallet.tar.gz
export CCID=$(peer lifecycle chaincode queryinstalled | cut -d ' ' -f 3 | sed s/.$// | grep $1)
peer lifecycle chaincode approveformyorg --package-id $CCID --channelID soshchannel --name wallet --version 1 --sequence $2 --waitForEvent --tls --cafile $ORDERER_CA
peer lifecycle chaincode checkcommitreadiness --channelID soshchannel --name wallet --version 1  --sequence $2 --tls --cafile $ORDERER_CA

export CORE_PEER_LOCALMSPID=superadminMSP
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/superadmin.sosh.com/peers/peer0.superadmin.sosh.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/superadmin.sosh.com/users/Admin@superadmin.sosh.com/msp
export CORE_PEER_ADDRESS=peer0.superadmin.sosh.com:7051
export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/sosh.com/orderers/orderer.sosh.com/msp/tlscacerts/tlsca.sosh.com-cert.pem

peer lifecycle chaincode commit -o orderer.sosh.com:7050 --channelID soshchannel --name wallet --version 1 --sequence $2 --tls true --cafile $ORDERER_CA --peerAddresses peer0.superadmin.sosh.com:7051 --peerAddresses peer0.users.sosh.com:9051  --tlsRootCertFiles ./crypto/peerOrganizations/superadmin.sosh.com/peers/peer0.superadmin.sosh.com/tls/ca.crt --tlsRootCertFiles ./crypto/peerOrganizations/users.sosh.com/peers/peer0.users.sosh.com/tls/ca.crt
