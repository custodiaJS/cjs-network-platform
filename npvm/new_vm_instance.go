package npvm

func NewVmInstance() (*NpVm, error) {
	imageres, err := _LoadVmImageResources()
	if err != nil {
		return nil, err
	}
	vmobj := &NpVm{
		imageResources: imageres,
	}
	return vmobj, nil
}
