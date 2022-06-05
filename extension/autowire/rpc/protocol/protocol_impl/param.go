package protocol_impl

type Param struct {
	Address    string
	ExportPort string
}

func (p *Param) Init(iocProtocol *IOCProtocol) (*IOCProtocol, error) {
	iocProtocol.address = p.Address
	iocProtocol.exportPort = p.ExportPort
	return iocProtocol, nil
}
