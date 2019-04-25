//utility types

package opts

import (
	"fmt"
	"os"

	"github.com/posener/complete"
)

//File is a string flag.Value which also performs
//an os.Stat on itself during Set, confirming it
//references a file
type File string

func (f File) String() string {
	return string(f)
}

func (f *File) Set(s string) error {
	if info, err := os.Stat(s); os.IsNotExist(err) {
		return fmt.Errorf("'%s' does not exist", s)
	} else if err != nil {
		return err
	} else if info.IsDir() {
		return fmt.Errorf("'%s' is a directory", s)
	}
	*f = File(s)
	return nil
}

var filesPredictor complete.Predictor

func (f *File) Predict(args complete.Args) []string {
	if filesPredictor == nil {
		filesPredictor = complete.PredictFiles("*")
	}
	return filesPredictor.Predict(args)
}

//Dir is a string flag.Value which also performs
//an os.Stat on itself during Set, confirming it
//references a directory
type Dir string

func (d Dir) String() string {
	return string(d)
}

func (d *Dir) Set(s string) error {
	if info, err := os.Stat(s); os.IsNotExist(err) {
		return fmt.Errorf("'%s' does not exist", s)
	} else if err != nil {
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("'%s' is a file", s)
	}
	*d = Dir(s)
	return nil
}

var dirsPredictor complete.Predictor

func (d *Dir) Predict(args complete.Args) []string {
	if dirsPredictor == nil {
		dirsPredictor = complete.PredictDirs("*")
	}
	return dirsPredictor.Predict(args)
}
