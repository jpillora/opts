package foo

type App struct {
	//configurable fields
	Ping string
	Pong string
	Zip  int
	Zop  int
	//internal state
	bar  int
	bazz int
}

func (f *App) Run() {
	f.bar = 42 + f.Zip
	f.bazz = 21 + f.Zop
	println("Foo is running...")
}
