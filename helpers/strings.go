package helpers

import (
    "strconv"
)

// AtoI converts alpha to integer ignoring errors
func AtoI(original string) int {
    converted, err := strconv.Atoi(original)
    if err != nil {
        return 0
    }

    return converted
}
