package manager

type EnvManager struct {
}

func NewEnvManager() *EnvManager {
	return &EnvManager{}
}

func (m *EnvManager) Init() error {
	return nil
}
