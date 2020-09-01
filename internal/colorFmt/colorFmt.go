package colorFmt

import "github.com/fatih/color"

var (
	Infof     = color.New(color.FgHiBlue).PrintfFunc()
	Warnf     = color.New(color.FgHiMagenta).PrintfFunc()
	Fatalf    = color.New(color.FgHiRed).PrintfFunc()
	Successf  = color.New(color.FgHiGreen).PrintfFunc()
	Openf     = Successf
	Closedf   = Fatalf
	Filteredf = color.New(color.FgHiYellow).PrintfFunc()

	Sinfof     = color.New(color.FgHiBlue).SprintfFunc()
	Swarnf     = color.New(color.FgHiMagenta).SprintfFunc()
	Sfatalf    = color.New(color.FgHiRed).SprintfFunc()
	Ssuccessf  = color.New(color.FgHiGreen).SprintfFunc()
	Sopenf     = Ssuccessf
	Sclosedf   = Sfatalf
	Sfilteredf = color.New(color.FgHiYellow).SprintfFunc()
)
