package instructions

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/getopendroplet/droplet/dropletfile/command"
	"github.com/getopendroplet/droplet/dropletfile/parser"

	"github.com/pkg/errors"
)

type parseRequest struct {
	command    string
	args       []string
	attributes map[string]bool
	flags      *BFlags
	original   string
	location   []parser.Range
	comments   []string
}

var parseRunPreHooks []func(*RunCommand, parseRequest) error
var parseRunPostHooks []func(*RunCommand, parseRequest) error

func nodeArgs(node *parser.Node) []string {
	result := []string{}
	for ; node.Next != nil; node = node.Next {
		arg := node.Next
		if len(arg.Children) == 0 {
			result = append(result, arg.Value)
		} else if len(arg.Children) == 1 {
			//sub command
			result = append(result, arg.Children[0].Value)
			result = append(result, nodeArgs(arg.Children[0])...)
		}
	}

	return result
}

func newParseRequestFromNode(node *parser.Node) parseRequest {
	return parseRequest{
		command:    node.Value,
		args:       nodeArgs(node),
		attributes: node.Attributes,
		original:   node.Original,
		flags:      NewBFlagsWithArgs(node.Flags),
		location:   node.Location(),
		comments:   node.PrevComment,
	}
}

// ParseInstruction converts an AST to a typed instruction (either a command or a build stage beginning when encountering a `STAGE` statement)
func ParseInstruction(node *parser.Node) (v interface{}, err error) {
	defer func() {
		err = parser.WithLocation(err, node.Location())
	}()
	req := newParseRequestFromNode(node)
	switch node.Value {
	case command.Arg:
		return parseArg(req)
	case command.Config:
		return parseConfig(req)
	case command.Copy:
		return parseCopy(req)
	case command.Cron:
		return parseCron(req)
	case command.Delete:
		return parseDelete(req)
	case command.Env:
		return parseEnv(req)
	case command.Expose:
		return parseExpose(req)
	case command.Label:
		return parseLabel(req)
	case command.Run:
		return parseRun(req)
	case command.Stage:
		return parseStage(req)
	case command.User:
		return parseUser(req)
	case command.Package:
		return parsePackage(req)
	case command.Workdir:
		return parseWorkdir(req)
	}

	return nil, &UnknownInstruction{Instruction: node.Value, Line: node.StartLine}
}

// ParseCommand converts an AST to a typed Command
func ParseCommand(node *parser.Node) (Command, error) {
	s, err := ParseInstruction(node)

	if err != nil {
		return nil, err
	}

	if c, ok := s.(Command); ok {
		return c, nil
	}

	return nil, parser.WithLocation(errors.Errorf("%T is not a command type", s), node.Location())
}

// UnknownInstruction represents an error occurring when a command is unresolvable
type UnknownInstruction struct {
	Line        int
	Instruction string
}

func (e *UnknownInstruction) Error() string {
	return fmt.Sprintf("unknown instruction: %s", strings.ToUpper(e.Instruction))
}

type parseError struct {
	inner error
	node  *parser.Node
}

func (e *parseError) Error() string {
	return fmt.Sprintf("dropletfile parse error line %d: %v", e.node.StartLine, e.inner.Error())
}

func (e *parseError) Unwrap() error {
	return e.inner
}

// Parse a Dropletfile into a collection of buildable stages.
func Parse(ast *parser.Node) (stages []Stage, metaArgs []ArgCommand, err error) {
	for _, n := range ast.Children {
		cmd, err := ParseInstruction(n)

		if err != nil {
			return nil, nil, &parseError{inner: err, node: n}
		}

		if len(stages) == 0 {
			// meta arg case
			if a, isArg := cmd.(*ArgCommand); isArg {
				metaArgs = append(metaArgs, *a)
				continue
			}
		}

		switch c := cmd.(type) {
		case *Stage:
			stages = append(stages, *c)
		case Command:
			stage, err := CurrentStage(stages)
			if err != nil {
				return nil, nil, parser.WithLocation(err, n.Location())
			}
			stage.AddCommand(c)
		default:
			return nil, nil, parser.WithLocation(errors.Errorf("%T is not a command type", cmd), n.Location())
		}

	}

	return stages, metaArgs, nil
}

func parseKvps(args []string, cmdName string) (KeyValuePairs, error) {
	if len(args) == 0 {
		return nil, errAtLeastOneArgument(cmdName)
	}

	if len(args)%2 != 0 {
		// should never get here, but just in case
		return nil, errTooManyArguments(cmdName)
	}

	var res KeyValuePairs
	for j := 0; j < len(args); j += 2 {
		if len(args[j]) == 0 {
			return nil, errBlankCommandNames(cmdName)
		}

		name := args[j]
		value := args[j+1]
		res = append(res, KeyValuePair{Key: name, Value: value})
	}

	return res, nil
}

func parseArg(req parseRequest) (*ArgCommand, error) {
	if len(req.args) == 0 {
		return nil, errAtLeastOneArgument("ARG")
	}

	pairs := make([]KeyValuePairOptional, len(req.args))

	for i, arg := range req.args {
		kvpo := KeyValuePairOptional{}

		// 'arg' can just be a name or name-value pair. Note that this is different
		// from 'env' that handles the split of name and value at the parser level.
		// The reason for doing it differently for 'arg' is that we support just
		// defining an arg and not assign it a value (while 'env' always expects a
		// name-value pair). If possible, it will be good to harmonize the two.
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts[0]) == 0 {
				return nil, errBlankCommandNames("ARG")
			}

			kvpo.Key = parts[0]
			kvpo.Value = &parts[1]
		} else {
			kvpo.Key = arg
		}
		kvpo.Comment = getComment(req.comments, kvpo.Key)
		pairs[i] = kvpo
	}

	return &ArgCommand{
		Args:            pairs,
		withNameAndCode: newWithNameAndCode(req),
	}, nil
}

func parseConfig(req parseRequest) (*ConfigCommand, error) {
	if len(req.args) == 0 {
		return nil, errAtLeastOneArgument("CONFIG")
	}

	configsTab := req.args
	sort.Strings(configsTab)

	return &ConfigCommand{
		Configs:         configsTab,
		withNameAndCode: newWithNameAndCode(req),
	}, nil
}

func parseCopy(req parseRequest) (*CopyCommand, error) {
	if len(req.args) < 2 {
		return nil, errNoDestinationArgument("COPY")
	}

	flChown := req.flags.AddString("chown", "")
	flChmod := req.flags.AddString("chmod", "")

	if err := req.flags.Parse(); err != nil {
		return nil, err
	}

	return &CopyCommand{
		SourcesAndDest:  SourcesAndDest(req.args),
		withNameAndCode: newWithNameAndCode(req),
		Chown:           flChown.Value,
		Chmod:           flChmod.Value,
	}, nil
}

func parseCron(req parseRequest) (*CronCommand, error) {
	if len(req.args) < 5 {
		return nil, errAtLeastSixArgument("CRON")
	}

	schedule := make([]string, 5)
	copy(schedule, req.args)

	req.args = append(req.args[:0], req.args[5:]...)

	return &CronCommand{
		withNameAndCode:       newWithNameAndCode(req),
		Minute:                schedule[0],
		Hour:                  schedule[1],
		DayOfTheMonth:         schedule[2],
		Month:                 schedule[3],
		DayOfTheWeek:          schedule[4],
		ShellDependantCmdLine: parseShellDependentCommand(req, false),
	}, nil
}

func parseDelete(req parseRequest) (*DeleteCommand, error) {
	if len(req.args) == 0 {
		return nil, errNoDestinationArgument("DELETE")
	}

	return &DeleteCommand{
		SourcesAndDest:  SourcesAndDest(req.args),
		withNameAndCode: newWithNameAndCode(req),
	}, nil
}

func parseEnv(req parseRequest) (*EnvCommand, error) {
	if len(req.args) == 0 {
		return nil, errAtLeastOneArgument("ENV")
	}

	envs, err := parseKvps(req.args, "ENV")

	if err != nil {
		return nil, err
	}

	return &EnvCommand{
		Env:             envs,
		withNameAndCode: newWithNameAndCode(req),
	}, nil
}

func parseExpose(req parseRequest) (*ExposeCommand, error) {
	if len(req.args) == 0 {
		return nil, errAtLeastOneArgument("EXPOSE")
	}

	portsTab := req.args
	sort.Strings(portsTab)

	return &ExposeCommand{
		Ports:           portsTab,
		withNameAndCode: newWithNameAndCode(req),
	}, nil
}

func parseLabel(req parseRequest) (*LabelCommand, error) {
	if len(req.args) == 0 {
		return nil, errAtLeastOneArgument("LABEL")
	}

	labels, err := parseKvps(req.args, "LABEL")

	if err != nil {
		return nil, err
	}

	return &LabelCommand{
		Labels:          labels,
		withNameAndCode: newWithNameAndCode(req),
	}, nil
}

func parsePackage(req parseRequest) (*PackageCommand, error) {
	flAction := req.flags.AddString("action", "")

	if err := req.flags.Parse(); err != nil {
		return nil, err
	}

	return &PackageCommand{
		Packages:        []string(req.args),
		Action:          flAction.Value,
		withNameAndCode: newWithNameAndCode(req),
	}, nil
}

func parseRun(req parseRequest) (*RunCommand, error) {
	if len(req.args) == 0 {
		return nil, errAtLeastOneArgument("RUN")
	}

	cmd := &RunCommand{}

	for _, fn := range parseRunPreHooks {
		if err := fn(cmd, req); err != nil {
			return nil, err
		}
	}

	cmd.ShellDependantCmdLine = parseShellDependentCommand(req, false)
	cmd.withNameAndCode = newWithNameAndCode(req)

	for _, fn := range parseRunPostHooks {
		if err := fn(cmd, req); err != nil {
			return nil, err
		}
	}

	return cmd, nil
}

func parseStage(req parseRequest) (*Stage, error) {
	if len(req.args) != 1 {
		return nil, errExactlyOneArgument("STAGE")
	}

	stageName := strings.ToLower(req.args[0])
	if ok, _ := regexp.MatchString("^[a-z][a-z0-9-_\\.]*$", stageName); !ok {
		return nil, errors.Errorf("invalid name for build stage: %q, name can't start with a number or contain symbols", req.args[0])
	}

	return &Stage{
		Name:       stageName,
		SourceCode: strings.TrimSpace(req.original),
		Commands:   []Command{},
		Location:   req.location,
		Comment:    getComment(req.comments, stageName),
	}, nil
}

func parseUser(req parseRequest) (*UserCommand, error) {
	if len(req.args) != 1 {
		return nil, errExactlyOneArgument("USER")
	}

	return &UserCommand{
		User:            req.args[0],
		withNameAndCode: newWithNameAndCode(req),
	}, nil
}

func parseWorkdir(req parseRequest) (*WorkdirCommand, error) {
	if len(req.args) != 1 {
		return nil, errExactlyOneArgument("WORKDIR")
	}

	return &WorkdirCommand{
		Path:            req.args[0],
		withNameAndCode: newWithNameAndCode(req),
	}, nil
}

func parseShellDependentCommand(req parseRequest, emptyAsNil bool) ShellDependantCmdLine {
	args := handleJSONArgs(req.args, req.attributes)
	cmd := StrSlice(args)

	if emptyAsNil && len(cmd) == 0 {
		cmd = nil
	}

	return ShellDependantCmdLine{
		CmdLine:      cmd,
		PrependShell: !req.attributes["json"],
	}
}

func errAtLeastOneArgument(command string) error {
	return errors.Errorf("%s requires at least one argument", command)
}

func errAtLeastSixArgument(command string) error {
	return errors.Errorf("%s requires at least six argument", command)
}

func errExactlyOneArgument(command string) error {
	return errors.Errorf("%s requires exactly one argument", command)
}

func errNoDestinationArgument(command string) error {
	return errors.Errorf("%s requires at least two arguments, but only one was provided. Destination could not be determined.", command)
}

func errBlankCommandNames(command string) error {
	return errors.Errorf("%s names can not be blank", command)
}

func errTooManyArguments(command string) error {
	return errors.Errorf("Bad input to %s, too many arguments", command)
}

func getComment(comments []string, name string) string {
	if name == "" {
		return ""
	}
	for _, line := range comments {
		if strings.HasPrefix(line, name+" ") {
			return strings.TrimPrefix(line, name+" ")
		}
	}
	return ""
}
