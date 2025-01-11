package validation

var validTones = map[string]bool{
    "nice":   true,
    "normal": true,
    "snarky": true,
}

func IsValidTone(tone string) bool {
    return validTones[tone]
}