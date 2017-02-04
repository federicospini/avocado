pragma solidity ^0.4.8;


contract StateChannels {
    uint8 constant PHASE_OPEN = 0;
    uint8 constant PHASE_CHALLENGE = 1;
    uint8 constant PHASE_CLOSED = 2;
    
    event Error(string message);

    event LogAcceptedChannel(
        bytes32 channelId,
        address address0,
        address address1,
        uint8 phase,
        uint challengePeriod,
        uint closingBlock,
        bytes state,
        uint sequenceNumber
    );
    
    event LogAcceptedUpdate(
        bytes32 channelId,
        address address0,
        address address1,
        uint8 phase,
        uint challengePeriod,
        uint closingBlock,
        bytes state,
        uint sequenceNumber
    );

    mapping (bytes32 => Channel) channels;

    struct Channel {
        bytes32 channelId;
        address address0;
        address address1;
        uint8 phase;
        uint challengePeriod;
        uint closingBlock;
        bytes state;
        uint sequenceNumber;
    }

    function ecrecovery(bytes32 hash, bytes sig) constant returns (address) {
        bytes32 r;
        bytes32 s;
        bytes32 v_test;
        uint8 v_temp = 1;
        uint8 v;
        address a;
        
        if (sig.length != 65) {
            throw;
        }

        // The signature format is a compact form of:
        //   {bytes32 r}{bytes32 s}{uint8 v}
        // Compact means, uint8 is not padded to 32 bytes.
        assembly {
            r := mload(add(sig, 32))
            s := mload(add(sig, 64))
            v_test := mload(add(sig, 65))
            v_temp := mload(add(sig, 65))
            // v_temp := and(mload(add(sig, 65)), 255)
        }
        
        // // old geth sends a `v` value of [0,1], while the new, in line with the YP sends [27,28]
        if (v_temp < 27) {
            v = v_temp + 27;
        } else {
            v = v_temp;
        }
        
        a = ecrecover(hash, v, r, s);
        
        return a;
    }
    
    function ecverify(bytes32 hash, bytes sig, address signer) constant returns (bool b) {
        b = ecrecovery(hash, sig) == signer;
        return b;
    }

    function getChannel(bytes32 channelId) returns(
        address address0,
        address address1,
        uint8 phase,
        uint challengePeriod,
        uint closingBlock,
        bytes state,
        uint sequenceNumber
    ) {
        address0 = channels[channelId].address0;
        address1 = channels[channelId].address1;
        phase = channels[channelId].phase;
        challengePeriod = channels[channelId].challengePeriod;
        closingBlock = channels[channelId].closingBlock;
        state = channels[channelId].state;
        sequenceNumber = channels[channelId].sequenceNumber;
    }

    function newChannel(
        bytes32 channelId,
        address address0,
        address address1,
        bytes state,
        uint256 challengePeriod,
        bytes signature0,
        bytes signature1
    ) {
        if (channels[channelId].channelId == channelId) {
            Error("channel with that channelId already exists");
            return;
        }

        bytes32 fingerprint = sha3(
            'newChannel',
            channelId,
            address0,
            address1,
            state,
            challengePeriod
        );

        if (!ecverify(fingerprint, signature0, address0)) {
            Error("signature0 invalid");
            return;
        }

        if (!ecverify(fingerprint, signature1, address1)) {
            Error("signature1 invalid");
            return;
        }

        Channel memory channel = Channel(
            channelId,
            address0,
            address1,
            PHASE_OPEN,
            challengePeriod,
            0,
            state,
            0
        );

        channels[channelId] = channel;

        LogAcceptedChannel(
            channel.channelId,
            channel.address0,
            channel.address1,
            channel.phase,
            channel.challengePeriod,
            channel.closingBlock,
            channel.state,
            channel.sequenceNumber
        );
    }

    function updateState(
        bytes32 channelId,
        uint256 sequenceNumber,
        bytes state,
        bytes signature0,
        bytes signature1
    ) {
        tryClose(channelId);
        
        Channel memory channel = channels[channelId];

        if (channel.phase == PHASE_CLOSED) {
            Error("channel closed");
            return;
        }

        bytes32 fingerprint = sha3(
            'updateState',
            channelId,
            sequenceNumber,
            state
        );

        if (!ecverify(fingerprint, signature0, channel.address0)) {
            Error("signature0 invalid");
            return;
        }

        if (!ecverify(fingerprint, signature1, channel.address1)) {
            Error("signature1 invalid");
            return;
        }

        if (sequenceNumber <= channel.sequenceNumber) {
            Error("sequence number too low");
            return;
        }

        channel.state = state;
        channel.sequenceNumber = sequenceNumber;
        
        // SE NON FUNZIONE DECOMMENTARE O STUDIARE IL MODELLO DI COMPUTAZIONE
        // channels[channelId] = channel;
        
        LogAcceptedChannel(
            channel.channelId,
            channel.address0,
            channel.address1,
            channel.phase,
            channel.challengePeriod,
            channel.closingBlock,
            channel.state,
            channel.sequenceNumber
        );
    }

    function startChallengePeriod(
        bytes32 channelId,
        bytes signature,
        address signer
    ) {
        Channel memory channel = channels[channelId];
        
        if (channel.phase != PHASE_OPEN) {
            Error("channel not open");
            return;
        }

        bytes32 fingerprint = sha3(
            'startChallengePeriod',
            channelId
        );

        if (signer == channel.address0) {
            if (!ecverify(fingerprint, signature, channel.address0)) {
                Error("signature invalid");
                return;
            }
        } else if (signer == channel.address1) {
            if (!ecverify(fingerprint, signature, channel.address1)) {
                Error("signature invalid");
                return;
            }
        } else {
            Error("signer invalid");
            return;
        }

        channel.closingBlock = block.number + channel.challengePeriod;
        channel.phase = PHASE_CHALLENGE;
        
        // SE NON FUNZIONE DECOMMENTARE O STUDIARE IL MODELLO DI COMPUTAZIONE
        // channels[channelId] = channel;
    }

    function tryClose(
        bytes32 channelId
    ) {
        Channel memory channel = channels[channelId];
        
        if (
            channel.phase == PHASE_CHALLENGE &&
            block.number > channel.closingBlock
        ) {
            channel.phase = PHASE_CLOSED;
        }
    }

    /**
     * USEFUL ONLY FOR DEBUGGING PURPOSES
     * TO REMOVE IN PRODUCTION
     */

    event LogDeletedChannel(
        bytes32 channelId,
        address address0,
        address address1,
        uint8 phase,
        uint challengePeriod,
        uint closingBlock,
        bytes state,
        uint sequenceNumber
    );

    function deleteChannel(
        bytes32 channelId
    ) {
        Channel channel = channels[channelId];
        
        if (channel.channelId != channelId) {
            Error("channel with that channelId does not exist");
            return;
        }

        LogDeletedChannel(
            channel.channelId,
            channel.address0,
            channel.address1,
            channel.phase,
            channel.challengePeriod,
            channel.closingBlock,
            channel.state,
            channel.sequenceNumber
        );

        channel.channelId = 0x0;
        channels[channelId] = channel;
    }
}
