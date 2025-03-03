package commands

import "gitlab.ozon.dev/timofey15g/homework/internal/models"

func ParseArgs(args []string) (map[string]string, error) {
	result := make(map[string]string)

	for i := 0; i+1 < len(args); i++ {
		if args[i][0] == '-' {
			_, exists := result[args[i][1:]]
			if exists {
				return nil, models.ErrorInvalidOptionalArgs
			}
			result[args[i][1:]] = args[i+1]
		}
	}

	return result, nil
}
