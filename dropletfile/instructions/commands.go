package instructions

import (
	"strings"

	"github.com/getopendroplet/droplet/dropletfile/parser"

	"github.com/pkg/errors"
)

// KeyValuePair represent an arbitrary named value (useful in slice instead of map[string] string to preserve ordering)
type KeyValuePair struct {
	Key   string
	Value string
}

func (kvp *KeyValuePair) String() string {
	return kvp.Key + "=" + kvp.Value
}

// KeyValuePairOptional is the same as KeyValuePair but Value is optional
type KeyValuePairOptional struct {
	Key     string
	Value   *string
	Comment string
}

func (kvpo *KeyValuePairOptional) ValueString() string {
	v := ""
	if kvpo.Value != nil {
		v = *kvpo.Value
	}
	return v
}

// Command is implemented by every command present in a dropletfile
type Command interface {
	Name() string
	Location() []parser.Range
}

// KeyValuePairs is a slice of KeyValuePair
type KeyValuePairs []KeyValuePair

// withNameAndCode is the base of every command in a Dropletfile (String() returns its source code)
type withNameAndCode struct {
	code     string
	name     string
	location []parser.Range
}

func (c *withNameAndCode) String() string {
	return c.code
}

// Name of the command
func (c *withNameAndCode) Name() string {
	return c.name
}

// Location of the command in source
func (c *withNameAndCode) Location() []parser.Range {
	return c.location
}

func newWithNameAndCode(req parseRequest) withNameAndCode {
	return withNameAndCode{code: strings.TrimSpace(req.original), name: req.command, location: req.location}
}

// SourcesAndDest represent a list of source files and a destination
type SourcesAndDest []string

// Sources list the source paths
func (s SourcesAndDest) Sources() []string {
	res := make([]string, len(s)-1)
	copy(res, s[:len(s)-1])
	return res
}

// Dest path of the operation
func (s SourcesAndDest) Dest() string {
	return s[len(s)-1]
}

// SingleWordExpander is a provider for variable expansion where 1 word => 1 output
type SingleWordExpander func(word string) (string, error)

// SupportsSingleWordExpansion interface marks a command as supporting variable expansion
type SupportsSingleWordExpansion interface {
	Expand(expander SingleWordExpander) error
}

func expandKvp(kvp KeyValuePair, expander SingleWordExpander) (KeyValuePair, error) {
	key, err := expander(kvp.Key)
	if err != nil {
		return KeyValuePair{}, err
	}
	value, err := expander(kvp.Value)
	if err != nil {
		return KeyValuePair{}, err
	}
	return KeyValuePair{Key: key, Value: value}, nil
}

func expandKvpsInPlace(kvps KeyValuePairs, expander SingleWordExpander) error {
	for i, kvp := range kvps {
		newKvp, err := expandKvp(kvp, expander)
		if err != nil {
			return err
		}
		kvps[i] = newKvp
	}
	return nil
}

func expandSliceInPlace(values []string, expander SingleWordExpander) error {
	for i, v := range values {
		newValue, err := expander(v)
		if err != nil {
			return err
		}
		values[i] = newValue
	}
	return nil
}

// ArgCommand : ARG name[=value]
type ArgCommand struct {
	withNameAndCode
	Args []KeyValuePairOptional
}

// Expand variables
func (c *ArgCommand) Expand(expander SingleWordExpander) error {
	for i, v := range c.Args {
		p, err := expander(v.Key)
		if err != nil {
			return err
		}
		v.Key = p
		if v.Value != nil {
			p, err = expander(*v.Value)
			if err != nil {
				return err
			}
			v.Value = &p
		}
		c.Args[i] = v
	}
	return nil
}

// ConfigCommand : CONFIG /etc/nginx /etc/hosts
type ConfigCommand struct {
	withNameAndCode
	Configs []string
}

// CopyCommand : COPY foo /path
type CopyCommand struct {
	withNameAndCode
	SourcesAndDest
	Chown string
	Chmod string
}

// CronCommand : CRON * * * * * df -h
type CronCommand struct {
	withNameAndCode
	Minute        string
	Hour          string
	DayOfTheMonth string
	Month         string
	DayOfTheWeek  string
	ShellDependantCmdLine
}

// Expand variables
func (c *CopyCommand) Expand(expander SingleWordExpander) error {
	expandedChown, err := expander(c.Chown)
	if err != nil {
		return err
	}
	c.Chown = expandedChown
	return expandSliceInPlace(c.SourcesAndDest, expander)
}

// DeleteCommand : DELETE /path
type DeleteCommand struct {
	withNameAndCode
	SourcesAndDest
}

// EnvCommand : ENV key1 value1 [keyN valueN...]
type EnvCommand struct {
	withNameAndCode
	Env KeyValuePairs // kvp slice instead of map to preserve ordering
}

// Expand variables
func (c *EnvCommand) Expand(expander SingleWordExpander) error {
	return expandKvpsInPlace(c.Env, expander)
}

// ExposeCommand : EXPOSE 6667/tcp 7000/tcp
type ExposeCommand struct {
	withNameAndCode
	Ports []string
}

// LabelCommand : LABEL some json data describing the image
type LabelCommand struct {
	withNameAndCode
	Labels   KeyValuePairs // kvp slice instead of map to preserve ordering
	noExpand bool
}

// Expand variables
func (c *LabelCommand) Expand(expander SingleWordExpander) error {
	if c.noExpand {
		return nil
	}
	return expandKvpsInPlace(c.Labels, expander)
}

// PackageCommand : PACKAGE nginx
type PackageCommand struct {
	withNameAndCode
	Action   string
	Packages []string
}

// Expand variables
func (c *PackageCommand) Expand(expander SingleWordExpander) error {
	action, err := expander(c.Action)
	if err != nil {
		return err
	}
	c.Action = action
	return expandSliceInPlace(c.Packages, expander)
}

// ShellDependantCmdLine represents a cmdline optionally prepended with the shell
type ShellDependantCmdLine struct {
	CmdLine      StrSlice
	PrependShell bool
}

// RunCommand : RUN some command yo
type RunCommand struct {
	withNameAndCode
	withExternalData
	ShellDependantCmdLine
}

// Stage represents a single stage in a multi-stage build
type Stage struct {
	Name       string
	Commands   []Command
	SourceCode string
	Location   []parser.Range
	Comment    string
}

// AddCommand to the stage
func (s *Stage) AddCommand(cmd Command) {
	// todo: validate cmd type
	s.Commands = append(s.Commands, cmd)
}

// IsCurrentStage check if the stage name is the current stage
func IsCurrentStage(s []Stage, name string) bool {
	if len(s) == 0 {
		return false
	}
	return s[len(s)-1].Name == name
}

// CurrentStage return the last stage in a slice
func CurrentStage(s []Stage) (*Stage, error) {
	if len(s) == 0 {
		return nil, errors.New("no build stage in current context")
	}
	return &s[len(s)-1], nil
}

// HasStage looks for the presence of a given stage name
func HasStage(s []Stage, name string) (int, bool) {
	for i, stage := range s {
		// Stage name is case-insensitive by design
		if strings.EqualFold(stage.Name, name) {
			return i, true
		}
	}
	return -1, false
}

// UserCommand : USER foo
type UserCommand struct {
	withNameAndCode
	User string
}

// Expand variables
func (c *UserCommand) Expand(expander SingleWordExpander) error {
	p, err := expander(c.User)
	if err != nil {
		return err
	}
	c.User = p
	return nil
}

// WorkdirCommand : WORKDIR /tmp
type WorkdirCommand struct {
	withNameAndCode
	Path string
}

// Expand variables
func (c *WorkdirCommand) Expand(expander SingleWordExpander) error {
	p, err := expander(c.Path)
	if err != nil {
		return err
	}
	c.Path = p
	return nil
}

type withExternalData struct {
	m map[interface{}]interface{}
}

func (c *withExternalData) getExternalValue(k interface{}) interface{} {
	return c.m[k]
}

func (c *withExternalData) setExternalValue(k, v interface{}) {
	if c.m == nil {
		c.m = map[interface{}]interface{}{}
	}
	c.m[k] = v
}
