// Copyright (c) 2017 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vppcalls

import (
	"fmt"
	govppapi "git.fd.io/govpp.git/api"
	log "github.com/ligato/cn-infra/logging/logrus"
	l2ba "github.com/ligato/vpp-agent/plugins/defaultplugins/l2plugin/bin_api/l2"
	"github.com/ligato/vpp-agent/plugins/defaultplugins/l2plugin/model/l2"
)

// VppAddBridgeDomain adds new bridge domain
func VppAddBridgeDomain(bdIdx uint32, bridgeDomain *l2.BridgeDomains_BridgeDomain, vppChan *govppapi.Channel) error {
	log.DefaultLogger().Debug("Adding VPP bridge domain ", bridgeDomain.Name)
	req := &l2ba.BridgeDomainAddDel{}
	req.BdID = bdIdx
	req.IsAdd = 1

	// Set bridge domain params
	req.Learn = boolToUint(bridgeDomain.Learn)
	req.ArpTerm = boolToUint(bridgeDomain.ArpTermination)
	req.Flood = boolToUint(bridgeDomain.Flood)
	req.UuFlood = boolToUint(bridgeDomain.UnknownUnicastFlood)
	req.Forward = boolToUint(bridgeDomain.Forward)
	req.MacAge = uint8(bridgeDomain.MacAge)

	reply := &l2ba.BridgeDomainAddDelReply{}
	err := vppChan.SendRequest(req).ReceiveReply(reply)
	if err != nil {
		return fmt.Errorf("Adding bridge domain failed with error %v", err)
	}
	if 0 != reply.Retval {
		return fmt.Errorf("Adding bridge domain returned %d", reply.Retval)
	}

	log.DefaultLogger().WithFields(log.Fields{"Name": bridgeDomain.Name, "Index": bdIdx}).Print("Bridge domain added.")
	return nil
}

// VppUpdateBridgeDomain updates bridge domain parameters
func VppUpdateBridgeDomain(oldBdIdx uint32, newBdIdx uint32, newBridgeDomain *l2.BridgeDomains_BridgeDomain, vppChan *govppapi.Channel) error {
	log.DefaultLogger().Debug("Updating VPP bridge domain parameters ", newBridgeDomain.Name)

	if oldBdIdx != 0 {
		VppDeleteBridgeDomain(oldBdIdx, vppChan)
	}

	req := &l2ba.BridgeDomainAddDel{}
	req.BdID = newBdIdx
	req.IsAdd = 1

	// Set bridge domain params
	req.Learn = boolToUint(newBridgeDomain.Learn)
	req.ArpTerm = boolToUint(newBridgeDomain.ArpTermination)
	req.Flood = boolToUint(newBridgeDomain.Flood)
	req.UuFlood = boolToUint(newBridgeDomain.UnknownUnicastFlood)
	req.Forward = boolToUint(newBridgeDomain.Forward)
	req.MacAge = uint8(newBridgeDomain.MacAge)

	reply := &l2ba.BridgeDomainAddDelReply{}
	err := vppChan.SendRequest(req).ReceiveReply(reply)
	if err != nil {
		return fmt.Errorf("Updating bridge domain failed with error %v", err)
	}
	if 0 != reply.Retval {
		return fmt.Errorf("Updating bridge domain returned %d", reply.Retval)
	}

	log.DefaultLogger().WithFields(log.Fields{"Name": newBridgeDomain.Name, "Index": newBdIdx}).Debug("Bridge domain Updated.")
	return nil
}

// VppDeleteBridgeDomain removes existing bridge domain
func VppDeleteBridgeDomain(bdIdx uint32, vppChan *govppapi.Channel) error {
	req := &l2ba.BridgeDomainAddDel{}
	req.BdID = bdIdx
	req.IsAdd = 0

	reply := &l2ba.BridgeDomainAddDelReply{}
	err := vppChan.SendRequest(req).ReceiveReply(reply)
	if err != nil {
		log.DefaultLogger().WithFields(log.Fields{"Error": err}).Error("Error while removing bridge domain")
		return err
	}
	if 0 != reply.Retval {
		log.DefaultLogger().WithFields(log.Fields{"Return value": reply.Retval}).Error("Unexpected return value")
	}

	return nil
}

func boolToUint(value bool) uint8 {
	if value {
		return 1
	}
	return 0
}
