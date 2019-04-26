package backend

// ShuffleTest --
type ShuffleTest struct {
	Title     string             `yaml:"title"`
	Summary   string             `yaml:"summary"`
	TestSuite string             `yaml:"test_suite"`
	Fork      string             `yaml:"fork"`
	Version   string             `yaml:"version"`
	TestCases []*ShuffleTestCase `yaml:"test_cases"`
}

// ShuffleTestCase --
type ShuffleTestCase struct {
	Input  ShuffleInput `yaml:"input,flow"`
	Output [][]uint64   `yaml:"output,flow"`
	Seed   string
}

// ShuffleInput --
type ShuffleInput struct {
	Epoch      uint64       `yaml:"epoch"`
	Validators []Validators `yaml:"validators"`
}

// Validators --
type Validators struct {
	ActivationEpoch uint64 `yaml:"activation_epoch"`
	ExitEpoch       uint64 `yaml:"exit_epoch"`
	OriginalIndex   int    `yaml:"original_index"`
}
