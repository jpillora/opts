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

func (a *App) Run() {
	a.bar = 42 + a.Zip
	a.bazz = 21 + a.Zop
	println("App is running: %+v", a)
}
