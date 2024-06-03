package generator

type GeneratorInterface interface {
	Generate() error
}

type Vars interface {
	Validate() error
}
