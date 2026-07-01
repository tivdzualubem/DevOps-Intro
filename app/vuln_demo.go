package main

import "golang.org/x/text/language"

// runVulnerableDependencyDemo deliberately calls a vulnerable dependency
// version for the Lab 9 govulncheck red/green CI demonstration.
func runVulnerableDependencyDemo() {
	_, _ = language.Parse("en-US")
}
