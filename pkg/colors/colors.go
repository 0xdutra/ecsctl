package colors

import (
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

var Green = color.New(color.FgHiGreen).SprintFunc()
var Red = color.New(color.FgHiRed).SprintFunc()
var Out = colorable.NewColorableStdout()
