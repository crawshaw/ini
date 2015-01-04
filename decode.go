package ini

// Config represents an INI configuration loaded into memory.
type Config map[string]map[string]string

const (
	// Default is the default section value.
	Default = ""
)

func Decode(conf string) (Config, error) {
	return nil, nil
}
