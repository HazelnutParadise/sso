package utils

import (
	"encoding/base64"
	"strings"
	"time"
)

func PtrString(s string) *string {
	return &s
}

func Deref[T any](s *T) T {
	if s == nil {
		var zero T
		return zero
	}
	return *s
}

func FormatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func BlobToBase64(blob []byte) string {
	if blob == nil {
		return ""
	}
	return "data:" + DetectImageMimeType(blob) + ";base64," + base64.StdEncoding.EncodeToString(blob)
}

func Base64ToBlob(base64Str string) *[]byte {
	if base64Str == "" {
		return nil
	}
	// 支援 data URI scheme
	if len(base64Str) > 5 && base64Str[:5] == "data:" {
		// 去除前綴
		parts := strings.SplitN(base64Str, ";base64,", 2)
		if len(parts) != 2 {
			return nil
		}
		base64Str = parts[1]
	}
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil
	}
	return &data
}

func DetectImageMimeType(blob []byte) string {
	if len(blob) >= 8 && blob[0] == 0x89 && blob[1] == 0x50 && blob[2] == 0x4E && blob[3] == 0x47 {
		return "image/png"
	}
	if len(blob) >= 3 && blob[0] == 0xFF && blob[1] == 0xD8 && blob[2] == 0xFF {
		return "image/jpeg"
	}
	if len(blob) >= 6 && blob[0] == 0x47 && blob[1] == 0x49 && blob[2] == 0x46 && blob[3] == 0x38 {
		return "image/gif"
	}
	// webp: "RIFF"...."WEBP"
	if len(blob) >= 12 && blob[0] == 0x52 && blob[1] == 0x49 && blob[2] == 0x46 && blob[3] == 0x46 &&
		blob[8] == 0x57 && blob[9] == 0x45 && blob[10] == 0x42 && blob[11] == 0x50 {
		return "image/webp"
	}
	// heif/heic: "ftypheic" or "ftypheix" or "ftyphevc" or "ftyphevx" or "ftypmif1" or "ftypmsf1"
	if len(blob) >= 12 && blob[4] == 0x66 && blob[5] == 0x74 && blob[6] == 0x79 && blob[7] == 0x70 {
		brand := string(blob[8:12])
		switch brand {
		case "heic", "heix", "hevc", "hevx", "mif1", "msf1":
			return "image/heif"
		}
	}
	return "application/octet-stream"
}
