#!/bin/bash
rm -rf channel-artifacts/*
export FABRIC_CFG_PATH=$PWD
configtxgen -outputBlock channel-artifacts/genesis.block -channelID ordererchannel -profile OrdererChannel
configtxgen -outputCreateChannelTx channel-artifacts/soshchannel.tx -channelID soshchannel -profile soshchannel
configtxgen --outputAnchorPeersUpdate channel-artifacts/superadminAnchor.tx -channelID soshchannel -profile soshchannel -asOrg superadminMSP
configtxgen --outputAnchorPeersUpdate channel-artifacts/usersAnchor.tx -channelID soshchannel -profile soshchannel -asOrg usersMSP
