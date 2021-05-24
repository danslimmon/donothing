okay, so we have questions about how the `Run` method and inputs and outputs should relate to each
other. the way i initially wvrote the README, `Run` takes a function with a `donothing.Inputs`
struct as its only parameter, and a `donothing.Outputs` struct as its only output.

let's first describe the _behavior_ we want, and then design the interface to suit that behavior.

# the behavior we want

a given step's `Run` function may need to refer to a value determined by a previous step. another
thing it may need to do is prompt the user for a value.

at a given point in the development of a procedure, a step may need to prompt the user for a value
because the step where that value is determined is not yet automated. if that prior step is later
automated , we want the subsequent step to automatically see that there's a new output by the
specified name, and use it

so maybe what we need is:

- a `fact` is a uniquely named, typed variable.
- the facts taken as input by a step are defined with `step.Inputs()`
- if we reach a step that requires an input, but that input's value is not yet determined, we prompt
    the user for a value (using the prompt defined for that input, as passed to `step.Inputs()`)
- if a fact is created during a step (either because the step's `Run()` function returned it or
    because the user was prompted for its value), that fact is added to the procedure's "fact sheet"

```
    pcd.AddStep(func(step *donothing.Step) {
        step.Name("blahBlah")
        step.Short("Blah the blah blah")
        step.Run(func(facts *donothing.Facts) error {
            x := facts.GetInt("factX", "The value of X")
            // ... do something ...
            facts.SetInt("y", "factY")
            return nil
        })
    })
```
