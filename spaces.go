package main

import "strings"

// TabsSpaces contains all info needed about tabs and spaces for a file
type TabsSpaces struct {
	perTab int  // number of spaces per tab/indentation
	spaces bool // use spaces, or tabs?
}

var defaultTabsSpaces = TabsSpaces{4, false}

// modeBlank
var languageIndentation = map[TabsSpaces][]Mode{
	{4, false}: {modeC, modeGo, modeHIDL, modeLisp, modeMakefile, modeNroff, modeOCaml, modeRust, modeStandardML}, // Tabs
	{2, true}:  {modeAssembly, modeConfig, modeHTML, modeHaskell, modeJSON, modeLua, modeObjectPascal, modeOdin, modePolicyLanguage, modeShell, modeVim, modeVim, modeXML},
	{3, true}:  {modeAda}, // Ada is special
	{4, true}:  {modeBat, modeBattlestar, modeCMake, modeCS, modeCpp, modeCrystal, modeGit, modeJSON, modeJava, modeJavaScript, modeKotlin, modeLua, modeMakefile, modeMarkdown, modeNim, modeOak, modePython, modeSQL, modeScala, modeText, modeTypeScript, modeZig},
}

// Spaces checks if the given mode should use tabs or spaces.
// Returns true for spaces.
func Spaces(mode Mode) bool {
	for k, vs := range languageIndentation {
		for _, v := range vs {
			if v == mode {
				return k.spaces
			}
		}
	}
	return defaultTabsSpaces.spaces
}

// TabsSpacesFromMode takes a mode, like modeJava, and tries to return the appropriate
// settings for tabs and spaces, as a TabsSpaces struct.
func TabsSpacesFromMode(mode Mode) TabsSpaces {
	// Given e.mode, find the matching TabsSpaces struct and set that to e.tabs
	for k, vs := range languageIndentation {
		for _, v := range vs {
			if v == mode {
				return k
			}
		}
	}
	return defaultTabsSpaces
}

// String returns the string for one indentation
func (ts TabsSpaces) String() string {
	if !ts.spaces {
		return "\t"
	}
	return strings.Repeat(" ", ts.perTab)
}