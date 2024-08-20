package logging

func LogVersion(version, date, commit string) {
	if err := InitializeZapLogger("info"); err != nil {
		panic(err)
	}

	logField := func(name, value string) {
		if value == "" {
			value = "N/A"
		}
		Log.Infof("%s: %s", name, value)
	}
	logField("Build version", version)
	logField("Build date", date)
	logField("Build commit", commit)
}
