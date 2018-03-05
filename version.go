package p9p

// Version returns the version of the 9p protocol used the established connection. If
// the connection is dead, or has never been established, or there is a version
// disagreement between the client and server, an error is returned
func (c *Conn) Version() (string, error) {
	if c.state != StEstablished {
		return "", ErrNoConn
	}
	if c.version == "" {
		return "", ErrBadVersion
	}
	return c.version, nil
}
