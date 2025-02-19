package domain

type CodeGenerator interface {
	Generate(length int) string
}
