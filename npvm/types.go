package npvm

type _NpVmImageResources struct {
	kernel   string
	firmware string
}

type NpVm struct {
	imageResources *_NpVmImageResources
}
