package builder

const (
	Success = 0

	ErrCodeNameRequired = -10001
)

var errMap = map[int16]string{
	Success: "ok",
	ErrCodeNameRequired: "ErrCodeNameRequired"
}

type resourcePoolConfig struct {
	name string
}

type Builder interface {
	Build() int16
}

type ResourcPoolBuild struct {
	resourcePoolConfig
}

func (build *ResourcPoolBuild) Build() int16 {
	if len(build.name) == 0 {
		return ErrCodeNameRequired
	}

	return Success
}

func (build *ResourcPoolBuild) SetName(name string) int16 {
	if len(name) == 0 {
		return ErrCodeNameRequired
	}
	build.name = name
	return Success
}
