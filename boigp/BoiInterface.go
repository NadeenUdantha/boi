package boigp

type BoiInterface interface {
	Init(n int)
	Update(i int)
	Score(i int) float64
	Swap(i, j int)
	Mutate(i int)
	Crossover(i, a, b int)
}
