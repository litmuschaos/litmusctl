package utils

var (
	// CompatibilityMatrix to store the compatible versions of litmusctl and ChaosCenter
	CompatibilityMatrix = map[string][]string{
		"0.6.0":  {"2.2.0", "2.3.0"},
		"0.7.0":  {"2.4.0", "2.5.0", "2.6.0", "2.7.0", "2.8.0"},
		"0.8.0":  {"2.4.0", "2.5.0", "2.6.0", "2.7.0", "2.8.0"},
		"0.9.0":  {"2.4.0", "2.5.0", "2.6.0", "2.7.0", "2.8.0"},
		"0.10.0": {"2.9.0", "2.10.0", "2.11.0", "2.12.0", "2.13.0", "2.14.0", "3.0-beta1", "3.0.0-beta2", "3.0.0-beta3", "3.0.0-beta4", "3.0.0-beta5", "3.0.0-beta6", "3.0.0-beta7", "3.0.0-beta8"},
		"0.11.0": {"2.9.0", "2.10.0", "2.11.0", "2.12.0", "2.13.0", "2.14.0", "3.0-beta1", "3.0.0-beta2", "3.0.0-beta3", "3.0.0-beta4", "3.0.0-beta5", "3.0.0-beta6", "3.0.0-beta7", "3.0.0-beta8"},
		"0.12.0": {"2.9.0", "2.10.0", "2.11.0", "2.12.0", "2.13.0", "2.14.0", "3.0-beta1", "3.0.0-beta2", "3.0.0-beta3", "3.0.0-beta4", "3.0.0-beta5", "3.0.0-beta6", "3.0.0-beta7", "3.0.0-beta8"},
		"0.13.0": {"2.9.0", "2.10.0", "2.11.0", "2.12.0", "2.13.0", "2.14.0", "3.0-beta1", "3.0.0-beta2", "3.0.0-beta3", "3.0.0-beta4", "3.0.0-beta5", "3.0.0-beta6", "3.0.0-beta7", "3.0.0-beta8"},
		"0.14.0": {"2.9.0", "2.10.0", "2.11.0", "2.12.0", "2.13.0", "2.14.0", "3.0-beta1", "3.0.0-beta2", "3.0.0-beta3", "3.0.0-beta4", "3.0.0-beta5", "3.0.0-beta6", "3.0.0-beta7", "3.0.0-beta8"},
		"0.15.0": {"2.9.0", "2.10.0", "2.11.0", "2.12.0", "2.13.0", "2.14.0", "3.0-beta1", "3.0.0-beta2", "3.0.0-beta3", "3.0.0-beta4", "3.0.0-beta5", "3.0.0-beta6", "3.0.0-beta7", "3.0.0-beta8"},
		"0.16.0": {"2.9.0", "2.10.0", "2.11.0", "2.12.0", "2.13.0", "2.14.0", "3.0-beta1", "3.0.0-beta2", "3.0.0-beta3", "3.0.0-beta4", "3.0.0-beta5", "3.0.0-beta6", "3.0.0-beta7", "3.0.0-beta8"},
		"0.17.0": {"2.9.0", "2.10.0", "2.11.0", "2.12.0", "2.13.0", "2.14.0", "3.0-beta1", "3.0.0-beta2", "3.0.0-beta3", "3.0.0-beta4", "3.0.0-beta5", "3.0.0-beta6", "3.0.0-beta7", "3.0.0-beta8"},
		"0.18.0": {"2.9.0", "2.10.0", "2.11.0", "2.12.0", "2.13.0", "2.14.0", "3.0-beta1", "3.0.0-beta2", "3.0.0-beta3", "3.0.0-beta4", "3.0.0-beta5", "3.0.0-beta6", "3.0.0-beta7", "3.0.0-beta8"},
		"0.19.0": {"2.9.0", "2.10.0", "2.11.0", "2.12.0", "2.13.0", "2.14.0", "3.0-beta1", "3.0.0-beta2", "3.0.0-beta3", "3.0.0-beta4", "3.0.0-beta5", "3.0.0-beta6", "3.0.0-beta7", "3.0.0-beta8"},
		"0.20.0": {"2.9.0", "2.10.0", "2.11.0", "2.12.0", "2.13.0", "2.14.0", "3.0-beta1", "3.0.0-beta2", "3.0.0-beta3", "3.0.0-beta4", "3.0.0-beta5", "3.0.0-beta6", "3.0.0-beta7", "3.0.0-beta8"},
		"0.21.0": {"2.9.0", "2.10.0", "2.11.0", "2.12.0", "2.13.0", "2.14.0", "3.0-beta1", "3.0.0-beta2", "3.0.0-beta3", "3.0.0-beta4", "3.0.0-beta5", "3.0.0-beta6", "3.0.0-beta7", "3.0.0-beta8"},
		"0.22.0": {"2.9.0", "2.10.0", "2.11.0", "2.12.0", "2.13.0", "2.14.0", "3.0-beta1", "3.0.0-beta2", "3.0.0-beta3", "3.0.0-beta4", "3.0.0-beta5", "3.0.0-beta6", "3.0.0-beta7", "3.0.0-beta8"},
		"0.23.0": {"3.0.0-beta9", "3.0.0-beta10", "3.0.0-beta11", "3.0.0-beta12"},
		"0.24.0": {"3.0.0-beta9", "3.0.0-beta10", "3.0.0-beta11", "3.0.0-beta12"},
		"1.0.0":  {"3.0.0", "3.1.0"},
		"1.1.0":  {"3.0.0", "3.1.0", "3.2.0"},
		"1.2.0":  {"3.0.0", "3.1.0", "3.2.0", "3.3.0"},
		"1.3.0":  {"3.0.0", "3.1.0", "3.2.0", "3.3.0", "3.4.0"},
		"1.4.0":  {"3.0.0", "3.1.0", "3.2.0", "3.3.0", "3.4.0", "3.5.0"},
		"1.5.0":  {"3.0.0", "3.1.0", "3.2.0", "3.3.0", "3.4.0", "3.5.0", "3.6.0"},
		"1.6.0":  {"3.0.0", "3.1.0", "3.2.0", "3.3.0", "3.4.0", "3.5.0", "3.6.0", "3.6.1", "3.7.0"},
		"1.7.0":  {"3.0.0", "3.1.0", "3.2.0", "3.3.0", "3.4.0", "3.5.0", "3.6.0", "3.6.1", "3.7.0", "3.8.0"},
		"1.8.0":  {"3.0.0", "3.1.0", "3.2.0", "3.3.0", "3.4.0", "3.5.0", "3.6.0", "3.6.1", "3.7.0", "3.8.0", "3.9.0", "3.9.1"},
	}
)
