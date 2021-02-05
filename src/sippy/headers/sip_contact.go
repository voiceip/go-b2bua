// Copyright (c) 2003-2005 Maxim Sobolev. All rights reserved.
// Copyright (c) 2006-2015 Sippy Software, Inc. All rights reserved.
// Copyright (c) 2015 Andrii Pylypenko. All rights reserved.
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
// list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation and/or
// other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
// ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package sippy_header

import (
    "sippy/conf"
    "sippy/net"
)

type SipContact struct {
    compactName
    *sipAddressHF
    Asterisk bool
}

var _sip_contact_name compactName = newCompactName("Contact", "m")

func NewSipContact(config sippy_conf.Config) *SipContact {
    return &SipContact{
        compactName  : _sip_contact_name,
        Asterisk     : false,
        sipAddressHF : newSipAddressHF(
                            NewSipAddress("Anonymous",
                                NewSipURL("", config.GetMyAddress(), config.GetMyPort(), false))),
    }
}

func NewSipContactFromAddress(addr *SipAddress) *SipContact {
    return &SipContact{
        compactName  : _sip_contact_name,
        Asterisk : false,
        sipAddressHF : newSipAddressHF(addr),
    }
}

func (self *SipContact) GetCopy() *SipContact {
    return &SipContact{
        compactName  : _sip_contact_name,
        sipAddressHF : self.sipAddressHF.getCopy(),
        Asterisk     : self.Asterisk,
    }
}

func (self *SipContact) GetCopyAsIface() SipHeader {
    return self.GetCopy()
}

func CreateSipContact(body string) []SipHeader {
    rval := []SipHeader{}
    if body == "*" {
        rval = append(rval, &SipContact{
            Asterisk     : true,
            compactName  : _sip_contact_name,
        })
    } else {
        addresses := CreateSipAddressHFs(body)
        for _, addr := range addresses {
            rval = append(rval, &SipContact{
                            sipAddressHF : addr,
                            Asterisk : false,
                            compactName  : _sip_contact_name,
                        })
        }
    }
    return rval
}

func (self *SipContact) StringBody() string {
    return self.LocalStringBody(nil)
}

func (self *SipContact) LocalStringBody(hostport *sippy_net.HostPort) string {
    if self.Asterisk {
        return "*"
    }
    return self.sipAddressHF.LocalStringBody(hostport)
}

func (self *SipContact) String() string {
    return self.LocalStr(nil, false)
}

func (self *SipContact) LocalStr(hostport *sippy_net.HostPort, compact bool) string {
    hname := self.Name()
    if compact {
        hname = self.CompactName()
    }
    return hname + ": " + self.LocalStringBody(hostport)
}
