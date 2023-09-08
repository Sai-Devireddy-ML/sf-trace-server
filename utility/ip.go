// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package utility

import (
	"net"
	"net/http"
)

// ExtractIP calls ExtractIPFromHeader(r) to extract a valid IP address. If no valid IP can be extracted from headers,
// ParseIP(r.RemoteAddr) is called.
func ExtractIP(r *http.Request) net.IP {
	if ip := ExtractIPFromHeader(r.Header); ip != nil {
		return ip
	}
	return ParseIP(r.RemoteAddr)
}

func ExtractIPRUM(r *http.Request) net.IP {
	if ip := ExtractIPFromHeaderRum(r.Header); ip != nil {
		return ip
	}

	return ParseIP((r.RemoteAddr))
}

// ExtractIPFromHeader extracts host information from `Forwarded`, `X-Real-IP`, `X-Forwarded-For` headers,
// in this order. The first valid IP address extracted is returned.
func ExtractIPFromHeader(header http.Header) net.IP {
	for _, parseFn := range parseHeadersInOrder {
		if ip := ParseIP(parseFn(header)); ip != nil {
			return ip
		}
	}
	return nil
}

// Identical to ExtractIPFromHeader, but also extracts host information from `X-Original-Forwarded-For` header,
// with that header taking priority over all others present in the request body.
func ExtractIPFromHeaderRum(header http.Header) net.IP {
	for _, parseFn := range parseRumHeadersInOrder {
		if ip := ParseIP(parseFn(header)); ip != nil {
			return ip
		}
	}

	return nil
}

// ParseIP returns the IP address parsed from a given input if a valid IP can be extracted. Otherwise returns nil.
func ParseIP(inp string) net.IP {
	if inp == "" {
		return nil
	}
	host, _ := splitHost(inp)
	if ip := net.ParseIP(host); ip != nil {
		return ip
	}
	return nil
}
