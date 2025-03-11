package main

type MyApp struct {
	metadata MyAppMetadata
	spec     MyAppSpec
}

type MyAppMetadata struct {
	name      string
	namespace string
}

type MyAppSpec struct {
	image       string
	replicas    int64
	appSelector string
}
