package opts

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"time"
)

//Parse with os.Args
func (o *node) Parse() ParsedOpts {
	return o.ParseArgs(os.Args[1:])
}

//ParseArgs with the provided arguments
func (o *node) ParseArgs(args []string) ParsedOpts {
	if err := o.process(args); err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
	return o
}

//process is the same as ParseArgs except
//it returns an error on failure
func (o *node) process(args []string) error {
	//cannot be processed - already encountered error - programmer error
	if o.erred != nil {
		return fmt.Errorf("[opts] process error: %s", o.erred)
	}
	//1. set config via JSON file
	if o.cfgPath != "" {
		b, err := ioutil.ReadFile(o.cfgPath)
		if err == nil {
			v := o.val.Interface() //*struct
			err = json.Unmarshal(b, v)
			if err != nil {
				o.erred = fmt.Errorf("Invalid config file: %s", err)
				return errors.New(o.Help())
			}
		}
	}
	flagset := flag.NewFlagSet(o.name, flag.ContinueOnError)
	flagset.SetOutput(ioutil.Discard)
	//pre-loop through the options and
	//add shortnames and env names where possible
	for _, opt := range o.flags {
		//should generate shortname?
		if len(opt.name) >= 3 && opt.shortName == "" {
			//not already taken?
			if s := opt.name[0:1]; !o.optnames[s] {
				opt.shortName = s
				o.optnames[s] = true
			}
		}
		env := camel2const(opt.name)
		if o.useEnv && (opt.envName == "" || opt.envName == "!") &&
			opt.name != "help" && opt.name != "version" &&
			!o.envnames[env] {
			opt.envName = env
		}
	}
	for _, opt := range o.flags {
		//2. set config via environment
		envVal := ""
		if opt.useEnv || o.useEnv {
			envVal = os.Getenv(opt.envName)
		}
		//3. set config via Go's pkg/flags
		addr := opt.val.Addr().Interface()
		switch addr := addr.(type) {
		case flag.Value:
			flagset.Var(addr, opt.name, "")
			if opt.shortName != "" {
				flagset.Var(addr, opt.shortName, "")
			}
		case *[]string:
			sep := ""
			switch opt.typeName {
			case "commalist":
				sep = ","
			case "spacelist":
				sep = " "
			}
			fv := &sepList{sep: sep, strs: addr}
			flagset.Var(fv, opt.name, "")
			if opt.shortName != "" {
				flagset.Var(fv, opt.shortName, "")
			}
		case *bool:
			str2bool(envVal, addr)
			flagset.BoolVar(addr, opt.name, *addr, "")
			if opt.shortName != "" {
				flagset.BoolVar(addr, opt.shortName, *addr, "")
			}
		case *string:
			str2str(envVal, addr)
			flagset.StringVar(addr, opt.name, *addr, "")
			if opt.shortName != "" {
				flagset.StringVar(addr, opt.shortName, *addr, "")
			}
		case *int:
			str2int(envVal, addr)
			flagset.IntVar(addr, opt.name, *addr, "")
			if opt.shortName != "" {
				flagset.IntVar(addr, opt.shortName, *addr, "")
			}
		case *time.Duration:
			flagset.DurationVar(addr, opt.name, *addr, "")
			if opt.shortName != "" {
				flagset.DurationVar(addr, opt.shortName, *addr, "")
			}
		default:
			return fmt.Errorf("[opts] Option '%s' has unsupported type", opt.name)
		}
	}
	//set user config
	err := flagset.Parse(args)
	if err != nil {
		//insert flag errors into help text
		o.erred = err
		o.internalOpts.Help = true
	}
	//internal opts (--help and --version)
	if o.internalOpts.Help {
		return errors.New(o.Help())
	} else if o.internalOpts.Version {
		fmt.Println(o.version)
		os.Exit(0)
	}
	//fill each individual arg
	args = flagset.Args()
	for i, argument := range o.args {
		if len(args) > 0 {
			str := args[0]
			args = args[1:]
			argument.val.SetString(str)
		} else if argument.defstr == "" {
			//not-set and no default!
			o.erred = fmt.Errorf("Argument #%d '%s' has no default value", i+1, argument.name)
			return errors.New(o.Help())
		}
	}
	//use command? peek at args
	if len(o.cmds) > 0 && len(args) > 0 {
		a := args[0]
		//matching command, use it
		if sub, exists := o.cmds[a]; exists {
			//user wants name to be set
			if o.cmdname != nil {
				o.cmdname.SetString(a)
			}
			return sub.process(args[1:])
		}
	}
	//fill arglist? assign remaining as slice
	if o.arglist != nil {
		if len(args) < o.arglist.min {
			o.erred = fmt.Errorf("Too few arguments (expected %d, got %d)", o.arglist.min, len(args))
			return errors.New(o.Help())
		}
		o.arglist.val.Set(reflect.ValueOf(args))
		args = nil
	}
	//we *should* have consumed all args at this point.
	//this prevents:  ./foo --bar 42 -z 21 ping --pong 7
	//where --pong 7 is ignored
	if len(args) != 0 {
		o.erred = fmt.Errorf("Unexpected arguments: %+v", args)
		return errors.New(o.Help())
	}
	return nil
}
