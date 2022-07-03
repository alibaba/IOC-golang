package autowire

var structDescriptorsMap = make(map[string]*StructDescriptor)

func RegisterStructDescriptor(sdid string, descriptor *StructDescriptor) {
	structDescriptorsMap[sdid] = descriptor
}

func GetStructDescriptor(sdid string) *StructDescriptor {
	return structDescriptorsMap[sdid]
}

func GetStructDescriptorsMap() map[string]*StructDescriptor {
	return structDescriptorsMap
}
