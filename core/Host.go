package core

import (
	"net"
	"strings"
)

/**
 * Stores hostname.
 */

const (
	HOSTNAME = iota
	IPV4ADDRESS
	IPV6ADDRESS
)

type Host struct {
	/** hostName field
	 */
	hostname string

	/** address field
	 */
	addressType int

	inetAddress net.IP
}

/** Constructor given host name or IP address.
 */
func NewHost(hname string) *Host {
	if hname == "" {
		return nil
	}

	host := &Host{}

	host.hostname = hname
	if host.isIPv6Address(hname) {
		host.addressType = IPV6ADDRESS
	}
	host.addressType = IPV4ADDRESS

	return host
}

/** constructor
 * @param name String to set
 * @param addrType int to set
 */
/*func NewHost2(name string, addrType int) *Host {
    host := &Host{};
    host.addressType = addrType
    host.hostname = strings.ToLower(strings.TrimSpace(name))
    return host;
}*/

/**
 * Return the host name in encoded form.
 * @return String
 */
func (host *Host) String() string {
	if host.addressType == IPV6ADDRESS && !host.isIPv6Reference(host.hostname) {
		return "[" + host.hostname + "]"
	}
	return host.hostname
}

/**
 * Compare for equality of hosts.
 * Host names are compared by textual equality. No dns lookup
 * is performed.
 * @param obj Object to set
 * @return boolean
 */
/*func (host *Host) equals(Object obj) bool{
    if (!host.getClass().equals(obj.getClass())) {
        return false;
    }
    Host otherHost = (Host) obj;
    return otherHost.hostname.equals(hostname);

}*/

/** get the HostName field
 * @return String
 */
func (host *Host) GetHostName() string {
	return host.hostname
}

/** get the Address field
 * @return String
 */
func (host *Host) GetAddress() string {
	return host.hostname
}

/**
 * Convenience function to get the raw IP destination address
 * of a SIP message as a String.
 * @return String
 */
func (host *Host) GetIpAddress() string {
	var rawIpAddress string
	if host.hostname == "" {
		return ""
	}

	if host.addressType == HOSTNAME {
		//try {
		if host.inetAddress == nil {
			host.inetAddress = net.ParseIP(host.hostname)
		}
		rawIpAddress = host.inetAddress.String() //getHostAddress();
		//} catch (UnknownHostException ex) {
		//    dbgPrint("Could not resolve hostname " + ex);
		//}
	} else {
		rawIpAddress = host.hostname
	}
	return rawIpAddress
}

/**
 * Set the hostname member.
 * @param h String to set
 */
func (host *Host) SetHostName(hname string) {
	host.inetAddress = nil
	if host.isIPv6Address(hname) {
		host.addressType = IPV6ADDRESS
	} else {
		host.addressType = HOSTNAME
	}
	// Null check bug fix sent in by jpaulo@ipb.pt
	if hname != "" {
		host.hostname = strings.ToLower(strings.TrimSpace(hname))
	}
}

/** Set the IP Address.
 *@param address is the address string to set.
 */
func (host *Host) SetHostAddress(address string) {
	host.inetAddress = nil
	if host.isIPv6Address(address) {
		host.addressType = IPV6ADDRESS
	} else {
		host.addressType = IPV4ADDRESS
	}
	if address != "" {
		host.hostname = strings.TrimSpace(address)
	}
}

/**
 * Set the address member
 * @param address address String to set
 */
func (host *Host) SetAddress(address string) {
	host.SetHostAddress(address)
}

/** Return true if the address is a DNS host name
 *  (and not an IPV4 address)
 *@return true if the hostname is a DNS name
 */
func (host *Host) IsHostName() bool {
	return host.addressType == HOSTNAME
}

/** Return true if the address is a DNS host name
 *  (and not an IPV4 address)
 *@return true if the hostname is host address.
 */
func (host *Host) IsIPAddress() bool {
	return host.addressType != HOSTNAME
}

/** Get the inet address from host host.
 * Caches the inet address returned from dns lookup to avoid
 * lookup delays.
 *
 *@throws UnkownHostexception when the host name cannot be resolved.
 */
func (host *Host) GetInetAddress() net.IP {
	if host.hostname == "" {
		return nil
	}
	if host.inetAddress != nil {
		return host.inetAddress
	}
	host.inetAddress = net.ParseIP(host.hostname)
	return host.inetAddress

}

//----- IPv6
/**
 * Verifies whether the <code>address</code> could
 * be an IPv6 address
 */
func (host *Host) isIPv6Address(address string) bool {
	return address != "" && strings.Index(address, ":") != -1
}

/**
 * Verifies whether the ipv6reference, i.e. whether it enclosed in
 * square brackets
 */
func (host *Host) isIPv6Reference(address string) bool {
	return address[0] == '[' && address[len(address)-1] == ']'
}

func (host *Host) Clone() interface{} {
	retval := &Host{}
	retval.addressType = host.addressType
	retval.hostname = host.hostname
	return retval
}
