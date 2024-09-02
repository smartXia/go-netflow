package constants

import (
	"github.com/fatih/color"
	"github.com/rfyiamcool/go-netflow"
)

const Thold = 1024 * 1024 // 1mb

var (
	Nf      netflow.Interface
	Yellow  = color.New(color.FgYellow).SprintFunc()
	Red     = color.New(color.FgRed).SprintFunc()
	Info    = color.New(color.FgGreen).SprintFunc()
	Blue    = color.New(color.FgBlue).SprintFunc()
	Magenta = color.New(color.FgHiMagenta).SprintFunc()
)
