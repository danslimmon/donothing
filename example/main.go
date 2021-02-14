// The example package contains an example implementation of a donothing script.
//
// The procedure implemented here is the following little arithmetic trick:
//
//     Multiply your phone number -- treating it as a seven-digit
//	   number (without its area code) -- by 8.  Then write down the
//     following three numbers:
//
//         - your phone number,
//         - 8, and
//         - the product of your phone number and 8.
//
//     Add up all the individual digits in those three numbers.  If
//     the sum is more than one digit, take that sum and add up its
//     digits.  Continue adding up digits until only one digit is left.
//
// This trick was found at [this website](https://www.pleacher.com/handley/puzzles/mathmagc.html).
//
// To run through the example code, execute `go build && ./example`. To print the procedure's
// markdown documentation, run `./example --print`.
package example

import (
	"os"

	"github.com/danslimmon/donothing"
)

// manual returns the manual implementation of the example procedure.
//
// In this implementation, the user will be prompted to execute each step.
func manual() *donothing.Procedure {
	pcd := donothing.NewProcedure()
	pcd.Short("The magic of 8")

	pcd.AddStep(func(step *donothing.Step) {
		step.Name("inputPhoneNumber")
		step.Short("Enter your phone number")
		step.OutputString(
			// The name of this output, by which other steps will refer to it
			"phoneNumber",
			// A description for this output, which we'll use to prompt the user
			"Your phone number",
		)
	})

	pcd.AddStep(func(step *donothing.Step) {
		step.Name("multiplyPhoneNumber")
		step.Short("Multiply your phone number by 8")
		step.OutputString("phoneNumberTimesEight", "Your phone number times 8")
	})

	pcd.AddStep(func(step *donothing.Step) {
		step.Name("addDigits")
		step.Short("Add up the digits")
		step.InputInt("phoneNumber", true)
		step.InputInt("phoneNumberTimesEight", true)
	})

	return pcd
}

// automated returns the automated implementation of the example procedure.
//
// In this implementation, the user will be prompted only for their phone number.
func automated() *donothing.Procedure {
	return nil
}

func main() {
	// Switch these comments around to use the automated version of the procedure instead of the
	// manual one.
	pcd := manual()
	//pcd := automated()

	if err := pcd.Check(); err != nil {
		panic(err)
	}
	if err := pcd.Render(os.Stdout); err != nil {
		panic(err)
	}
	if err := pcd.Execute(); err != nil {
		panic(err)
	}
}
