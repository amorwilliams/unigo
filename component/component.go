package component

type Component interface {
	Awake() error
	Start() error
	Stop() error
}

type ComponentFactory interface {
}
