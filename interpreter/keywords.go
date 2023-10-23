package interpreter

import "fmt"

type dest int

const (
	destNormal dest = iota
	destSkip
)

type kwd int

const (
	kwdChar kwd = iota
	kwdDest
	kwdProp
	kwdSpec
)

func (kwd kwd) String() string {
	switch kwd {
	case kwdChar:
		return "kwdChar"
	case kwdDest:
		return "kwdDest"
	case kwdProp:
		return "kwdProp"
	case kwdSpec:
		return "kwdSpec"
	default:
		return fmt.Sprintf("kwdUnknown(%d)", kwd)
	}
}

type ipfn int

const (
	ipfnBin ipfn = iota
	ipfnUnicode
	ipfnHex
)

func (ipfn ipfn) String() string {
	switch ipfn {
	case ipfnBin:
		return "ipfnBin"
	case ipfnUnicode:
		return "ipfnUnicode"
	case ipfnHex:
		return "ipfnHex"
	default:
		return fmt.Sprintf("ipfnUnknown(%d)", ipfn)
	}
}

type keyword struct {
	kwd kwd
	idx int
}

var keywords = map[string]keyword{
	"author":             {kwdDest, int(destSkip)},
	"colorschememapping": {kwdDest, int(destSkip)},
	"colortbl":           {kwdDest, int(destSkip)},
	"company":            {kwdDest, int(destSkip)},
	"datastore":          {kwdDest, int(destSkip)},
	"expandedcolortbl":   {kwdDest, int(destSkip)},
	"fldinst":            {kwdDest, int(destSkip)},
	"fonttbl":            {kwdDest, int(destSkip)},
	"generator":          {kwdDest, int(destSkip)},
	"header":             {kwdDest, int(destSkip)},
	"headerf":            {kwdDest, int(destSkip)},
	"headerl":            {kwdDest, int(destSkip)},
	"headerr":            {kwdDest, int(destSkip)},
	"list":               {kwdDest, int(destSkip)},
	"listlevel":          {kwdDest, int(destSkip)},
	"listname":           {kwdDest, int(destSkip)},
	"lsdlocked":          {kwdDest, int(destSkip)},
	"operator":           {kwdDest, int(destSkip)},
	"panose":             {kwdDest, int(destSkip)},
	"par":                {kwdChar, '\n'},
	"pgdscnxt":           {kwdDest, int(destSkip)},
	"pntxta":             {kwdDest, int(destSkip)},
	"pntxtb":             {kwdDest, int(destSkip)},
	"rquote":             {kwdChar, '\''},
	"shp":                {kwdDest, int(destSkip)},
	"sn":                 {kwdDest, int(destSkip)},
	"stylesheet":         {kwdDest, int(destSkip)},
	"themedata":          {kwdDest, int(destSkip)},
	"title":              {kwdDest, int(destSkip)},
	"u":                  {kwdSpec, int(ipfnUnicode)},
	"wgrffmtfilter":      {kwdDest, int(destSkip)},
	"xmlnstbl":           {kwdDest, int(destSkip)},
}
