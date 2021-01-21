package dagger

import (
	"context"
	"strings"
	"testing"
)

func TestValidateEmptyValue(t *testing.T) {
	cc := &Compiler{}
	v, err := cc.Compile("", "#dagger: compute: _")
	if err != nil {
		t.Fatal(err)
	}
	if err := v.Get("#dagger.compute").Validate("#Script"); err != nil {
		t.Fatal(err)
	}
}

func TestLocalScript(t *testing.T) {
	ctx := context.TODO()

	cc := &Compiler{}
	src := `[{do: "local", dir: "foo"}]`
	v, err := cc.Compile("", src)
	if err != nil {
		t.Fatal(err)
	}
	s, err := v.Script()
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	err = s.Walk(ctx, func(op *Op) error {
		n++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatal(n)
	}
}

func TestWalkBootScript(t *testing.T) {
	ctx := context.TODO()

	cc := &Compiler{}
	cfg, err := cc.Compile("clientconfig.cue", defaultBootScript)
	if err != nil {
		t.Fatal(err)
	}
	script, err := cfg.Get("bootscript").Script()
	if err != nil {
		t.Fatal(err)
	}
	dirs, err := script.LocalDirs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) != 1 {
		t.Fatal(dirs)
	}
	if dirs[0] != "." {
		t.Fatal(dirs)
	}
}

func TestWalkBiggerScript(t *testing.T) {
	t.Skip("FIXME")

	ctx := context.TODO()
	cc := &Compiler{}
	script, err := cc.CompileScript("boot.cue", `
[
//	{
//		do: "load"
//		from: {
//			do: "local"
//			dir: "ga"
//		}
//	},
	{
		do: "local"
		dir: "bu"
	},
	{
		do: "copy"
		from: [
			{
				do: "local"
				dir: "zo"
			}
		]
	},
	{
		do: "exec"
		args: ["ls"]
		mount: "/mnt": input: [
			{
				do: "local"
				dir: "meu"
			}
		]
	}
]
`)
	if err != nil {
		t.Fatal(err)
	}
	dirs, err := script.LocalDirs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) != 4 {
		t.Fatal(dirs)
	}
	wanted := "ga bu zo meu"
	got := strings.Join(dirs, " ")
	if wanted != got {
		t.Fatal(got)
	}
}