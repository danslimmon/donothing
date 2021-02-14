package donothing

type Procedure struct{}

func (pcd *Procedure) Short(s string) {}

func (pcd *Procedure) AddStep(func(*Step)) {}

func (pcd *Procedure) Check() error {
	return nil
}

func (pcd *Procedure) Execute() error {
	return nil
}

func NewProcedure() *Procedure {
	return nil
}
