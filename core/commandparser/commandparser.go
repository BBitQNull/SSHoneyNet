package commandparser

import (
	"context"
)

type CmdParserService interface {
	CommandParser(ctx context.Context, cmd string) (*Script, error)
}

type Script struct {
	Lines []*CommandLine `@@*`
}

type CommandLine struct {
	Pipeline []*Command   `@@ ("|" @@)*`
	Redir    *Redirection `@@?`
}

type Command struct {
	Name  string          `@Ident`
	Flags []FlagWithValue `@@*`
	Args  []Argument      `@@*`
}

type FlagWithValue struct {
	Name  string  `("-" @Ident | "--" @Ident | @("--" @Ident "=" @Ident))`
	Value *string `(@Ident | @String)?`
}

type Argument struct {
	Value *string `(@String | @Ident)`
}
type Redirection struct {
	File string `(">" @String | ">>" @String | "<" @String | "<<" @String)`
}
