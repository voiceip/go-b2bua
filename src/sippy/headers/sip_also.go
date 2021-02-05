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
    "sippy/net"
)

type SipAlso struct {
    normalName
    *sipAddressHF
}

var _sip_also_name normalName = newNormalName("Also")

func CreateSipAlso(body string) []SipHeader {
    addresses := CreateSipAddressHFs(body)
    rval := make([]SipHeader, len(addresses))
    for i, addr := range addresses {
        rval[i] = &SipAlso{
            normalName   : _sip_also_name,
            sipAddressHF : addr,
        }
    }
    return rval
}

func NewSipAlso(addr *SipAddress) *SipAlso {
    return &SipAlso{
        normalName   : _sip_also_name,
        sipAddressHF : newSipAddressHF(addr),
    }
}

func (self *SipAlso) String() string {
    return self.LocalStr(nil, false)
}

func (self *SipAlso) LocalStr(hostport *sippy_net.HostPort, compact bool) string {
    return self.Name() + ": " + self.LocalStringBody(hostport)
}

func (self *SipAlso) GetCopy() *SipAlso {
    if self == nil {
        return nil
    }
    return &SipAlso{
        normalName   : _sip_also_name,
        sipAddressHF : self.sipAddressHF.getCopy(),
    }
}

func (self *SipAlso) GetCopyAsIface() SipHeader {
    return self.GetCopy()
}
