package core

import (
	"bytes"
	"net"
	"strconv"
)

/**
* Holds the hostname:port.
 */

type HostPort struct {
	// host / ipv4/ ipv6/
	/** host field
	 */
	host *Host

	/** port field
	 *
	 */
	port int
}

/** Default constructor
 */
func NewHostPort() *HostPort {
	hostPort := &HostPort{} // marker for not set.
	hostPort.port = -1
	return hostPort
}

/**
 * Encode hostPort hostport into its string representation.
 * Note that hostPort could be different from the string that has
 * been parsed if something has been edited.
 * @return String
 */
func (hostPort *HostPort) String() string {
	var retval bytes.Buffer //= new StringBuffer();
	if hostPort.host != nil {
		retval.WriteString(hostPort.host.String())
		if hostPort.port != -1 {
			retval.WriteString(SIPSeparatorNames_COLON + strconv.Itoa(hostPort.port))
		}
	}
	return retval.String()
}

/** returns true if the two objects are equals, false otherwise.
 * @param other Object to set
 * @return boolean
 */
/*public boolean equals(Object other) {
            if (! hostPort.getClass().equals(other.getClass())) {
                return false;
            }
            HostPort that = (HostPort) other;
	    if ( (hostPort.port == null && that.port != null) ||
		 (hostPort.port != null && that.port == null) ) return false;
	    else if (hostPort.port == that.port && hostPort.host.equals(that.host))
		return true;
	    else
              return hostPort.host.equals(that.host) && hostPort.port.equals(that.port);
        }*/

/** get the Host field
 * @return host field
 */
func (hostPort *HostPort) GetHost() *Host {
	return hostPort.host
}

/** get the port field
 * @return int
 */
func (hostPort *HostPort) GetPort() int {
	return hostPort.port
}

/**
 * Returns boolean value indicating if Header has port
 * @return boolean value indicating if Header has port
 */
func (hostPort *HostPort) HasPort() bool {
	return hostPort.port != -1
}

/** remove port.
 */
func (hostPort *HostPort) RemovePort() {
	hostPort.port = -1
}

/**
 * Set the host member
 * @param h Host to set
 */
func (hostPort *HostPort) SetHost(h *Host) {
	hostPort.host = h
}

/**
 * Set the port member
 * @param p int to set
 */
func (hostPort *HostPort) SetPort(p int) {
	// -1 is same as remove port.
	hostPort.port = p
}

/** Return the internet address corresponding to the host.
 *@throws java.net.UnkownHostException if host name cannot be resolved.
 *@return the inet address for the host.
 */
func (hostPort *HostPort) GetInetAddress() net.IP {
	if hostPort.host == nil {
		return nil
	}
	return net.ParseIP(hostPort.host.GetHostName())
}

func (hostPort *HostPort) Clone() interface{} {
	retval := &HostPort{}
	if hostPort.host != nil {
		retval.host = hostPort.host.Clone().(*Host)
	} else {
		retval.host = nil
	}
	retval.port = hostPort.port
	return retval
}
