package exercises

import "fmt"

// Store provides access to exercise definitions. Currently backed by an
// in-memory registry; will be replaced with database-backed storage when
// the content authoring workflow is defined.
type Store struct {
	exercises map[string]Exercise
}

// NewStore creates a Store pre-loaded with the built-in exercises.
func NewStore() *Store {
	s := &Store{
		exercises: make(map[string]Exercise),
	}
	s.seed()
	return s
}

// Get returns an exercise by ID, or an error if not found.
func (s *Store) Get(id string) (Exercise, error) {
	ex, ok := s.exercises[id]
	if !ok {
		return Exercise{}, fmt.Errorf("exercise not found: %s", id)
	}
	return ex, nil
}

// List returns all exercises.
func (s *Store) List() []Exercise {
	out := make([]Exercise, 0, len(s.exercises))
	for _, ex := range s.exercises {
		out = append(out, ex)
	}
	return out
}

// seed populates the store with the initial set of exercises.
func (s *Store) seed() {
	for _, ex := range builtinExercises {
		s.exercises[ex.ID] = ex
	}
}

var builtinExercises = []Exercise{
	{
		ID:    "m01-hello-world",
		Title: "Hello World",
		Instructions: `Write a program that prints exactly:

` + "```" + `
Hello, World!
` + "```" + `

Use the ` + "`fmt`" + ` package to print to stdout.`,
		StarterCode: `package main

import "fmt"

func main() {
	// Your code here
}
`,
		Hint:           "Use fmt.Println() to print the text.",
		ExpectedOutput: "Hello, World!\n",
		ValidationMode: ValidationModeExact,
	},
	{
		ID:    "m01-variables",
		Title: "Variable Declaration",
		Instructions: "Declare a variable `name` with the value `\"Gopher\"` and print a greeting:\n\n```\nHello, Gopher!\n```",
		StarterCode: `package main

import "fmt"

func main() {
	// Declare a variable "name" and use it in a greeting
}
`,
		Hint:           "Use := for short variable declaration, then fmt.Printf or fmt.Println.",
		ExpectedOutput: "Hello, Gopher!\n",
		ValidationMode: ValidationModeExact,
	},
	{
		ID:    "m01-loop",
		Title: "For Loop",
		Instructions: `Write a program that prints the numbers 1 through 5, one per line:

` + "```" + `
1
2
3
4
5
` + "```",
		StarterCode: `package main

import "fmt"

func main() {
	// Use a for loop to print numbers 1-5
}
`,
		Hint:           "Use a for loop: for i := 1; i <= 5; i++ { ... }",
		ExpectedOutput: "1\n2\n3\n4\n5\n",
		ValidationMode: ValidationModeExact,
	},
	{
		ID:    "m01-fizzbuzz",
		Title: "FizzBuzz",
		Instructions: `Write a FizzBuzz program for numbers 1 to 15:
- Print "Fizz" for multiples of 3
- Print "Buzz" for multiples of 5
- Print "FizzBuzz" for multiples of both
- Print the number otherwise

Each output on its own line.`,
		StarterCode: `package main

import "fmt"

func main() {
	// Implement FizzBuzz for 1-15
}
`,
		Hint:           "Use the modulo operator (%). Check divisibility by 15 first, then 3, then 5.",
		ExpectedOutput: "1\n2\nFizz\n4\nBuzz\nFizz\n7\n8\nFizz\nBuzz\n11\nFizz\n13\n14\nFizzBuzz\n",
		ValidationMode: ValidationModeExact,
	},
}
