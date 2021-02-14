package donothing

type Step struct{}

func (step *Step) Name(s string) {}

func (step *Step) Short(s string) {}

func (step *Step) AddStep(func(*Step)) {}

func (step *Step) OutputString(name string, desc string) {}

func (step *Step) InputInt(name string, required bool) {}
