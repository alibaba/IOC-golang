package protocol_impl

type Param struct {
	Address    string
	Timeout    string
	ExportPort string
}

func (p *Param) Init(iocProtocol *IOCProtocol) (*IOCProtocol, error) {
	iocProtocol.address = p.Address
	iocProtocol.exportPort = p.ExportPort
	iocProtocol.timeout = p.Timeout
	return iocProtocol, nil
}
