package gen

import (
    "math/rand"
    "time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateSerial(length int) string {
    seed := rand.NewSource(time.Now().UnixNano())
    random := rand.New(seed)
    serial := make([]byte, length)
    for i := range serial {
        serial[i] = charset[random.Intn(len(charset))]
    }
    return string(serial)
}