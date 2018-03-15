package reporter

import (
	"github.com/apex/up/reporter/discard"
	"github.com/apex/up/reporter/plain"
	"github.com/apex/up/reporter/text"
)

var (
	// Discard reporter.
	Discard = discard.Report

	// Plain reporter.
	Plain = plain.Report

	// Text reporter.
	Text = text.Report
)
