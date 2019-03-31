package foo

type App struct {
	//internal config
	Ping string `type:"flag"`
	Pong string `type:"flag"`
	Zip  int    `type:"flag"`
	Zop  int    `type:"flag"`
	//internal state
	bar  int
	bazz int
}

func (f *App) Run() {
	f.bar = 42 + f.Zip
	f.bazz = 21 + f.Zop
	println("Foo is running...")
}
