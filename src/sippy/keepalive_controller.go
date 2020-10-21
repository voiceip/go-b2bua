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
package sippy

import (
    "sippy/log"
    "sippy/headers"
    "sippy/types"
)

type keepaliveController struct {
    ua          sippy_types.UA
    triedauth   bool
    ka_tr       sippy_types.ClientTransaction
    keepalives  int
    logger      sippy_log.ErrorLogger
}

func newKeepaliveController(ua sippy_types.UA, logger sippy_log.ErrorLogger) *keepaliveController {
    if ua.GetKaInterval() <= 0 {
        return nil
    }
    self := &keepaliveController{
        ua          : ua,
        triedauth   : false,
        keepalives  : 0,
        logger      : logger,
    }
    return self
}

func (self *keepaliveController) Start() {
    StartTimeout(self.keepAlive, self.ua.GetSessionLock(), self.ua.GetKaInterval(), 1, self.logger)
}

func (self *keepaliveController) RecvResponse(resp sippy_types.SipResponse, tr sippy_types.ClientTransaction) {
    var err error
    var challenge *sippy_header.SipWWWAuthenticateBody
    var req sippy_types.SipRequest
    var new_auth_fn sippy_header.NewSipXXXAuthorizationFunc

    if self.ua.GetState() != sippy_types.UA_STATE_CONNECTED {
        return
    }
    code, _ := resp.GetSCode()
    if self.ua.GetUsername() != "" && self.ua.GetPassword() != "" && ! self.triedauth {
        if code == 401 && resp.GetSipWWWAuthenticate() != nil {
            challenge, err = resp.GetSipWWWAuthenticate().GetBody()
            if err != nil {
                self.logger.Error("error parsing 401 auth: " + err.Error())
                return
            }
            new_auth_fn = func(realm, nonce, method, uri, username, password string) sippy_header.SipHeader {
                return sippy_header.NewSipAuthorization(realm, nonce, method, uri, username, password)
            }
        } else if code == 407 && resp.GetSipProxyAuthenticate() != nil {
            challenge, err = resp.GetSipProxyAuthenticate().GetBody()
            if err != nil {
                self.logger.Error("error parsing 407 auth: " + err.Error())
                return
            }
            new_auth_fn = func(realm, nonce, method, uri, username, password string) sippy_header.SipHeader {
                return sippy_header.NewSipProxyAuthorization(realm, nonce, method, uri, username, password)
            }
        }
        if challenge != nil {
            req, err = self.ua.GenRequest("INVITE", self.ua.GetLSDP(), challenge.GetNonce(), challenge.GetRealm(), new_auth_fn)
            if err != nil {
                self.logger.Error("Cannot create INVITE: " + err.Error())
                return
            }
            self.ka_tr, err = self.ua.PrepTr(req)
            if err == nil {
                self.triedauth = true
            }
            self.ua.SipTM().BeginClientTransaction(req, self.ka_tr)
            return
        }
    }
    if code < 200 {
        return
    }
    self.ka_tr = nil
    self.keepalives += 1
    if code == 408 || code == 481 || code == 486 {
        if self.keepalives == 1 {
            //print "%s: Remote UAS at %s:%d does not support re-INVITES, disabling keep alives" % (self.ua.cId, self.ua.rAddr[0], self.ua.rAddr[1])
            StartTimeout(func() { self.ua.Disconnect(nil, "") }, self.ua.GetSessionLock(), 600, 1, self.logger)
            return
        }
        //print "%s: Received %d response to keep alive from %s:%d, disconnecting the call" % (self.ua.cId, code, self.ua.rAddr[0], self.ua.rAddr[1])
        self.ua.Disconnect(nil, "")
        return
    }
    StartTimeout(self.keepAlive, self.ua.GetSessionLock(), self.ua.GetKaInterval(), 1, self.logger)
}

func (self *keepaliveController) keepAlive() {
    var err error
    var req sippy_types.SipRequest

    if self.ua.GetState() != sippy_types.UA_STATE_CONNECTED {
        return
    }
    req, err = self.ua.GenRequest("INVITE", self.ua.GetLSDP(), "", "", nil)
    if err != nil {
        self.logger.Error("Cannot create INVITE: " + err.Error())
        return
    }
    self.triedauth = false
    self.ka_tr, err = self.ua.PrepTr(req)
    if err == nil {
        self.ua.SipTM().BeginClientTransaction(req, self.ka_tr)
    }
}

func (self *keepaliveController) Stop() {
    if ka_tr := self.ka_tr; ka_tr != nil {
        ka_tr.Cancel()
        self.ka_tr = nil
    }
}
