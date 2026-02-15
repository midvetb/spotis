package deezer

import (
	"bytes"
	"context"
	"crypto/cipher"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/crypto/blowfish"
)

const (
	chunkSize = 2048
	quality   = "MP3_128"
)

var iv = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}

func GetDecryptionKey(secret, songID string) []byte {
	hash := md5.Sum([]byte(songID))
	hashHex := fmt.Sprintf("%x", hash)

	key := []byte(secret)
	for i := 0; i < len(hash); i++ {
		key[i] = key[i] ^ hashHex[i] ^ hashHex[i+16]
	}

	return key
}

func Decrypt(data, key []byte) ([]byte, error) {
	block, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(data))
	mode.CryptBlocks(decrypted, data)

	return decrypted, nil
}

func DownloadTrack(ctx context.Context, session *Session, trackURL string, song *Song) (*bytes.Buffer, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", trackURL, nil)

	streamingClient := *session.HttpClient
	streamingClient.Timeout = 0

	resp, err := streamingClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download: status %d", resp.StatusCode)
	}

	key := GetDecryptionKey(os.Getenv("DEEZER_SECRET"), song.ID)
	buffer := make([]byte, chunkSize)
	chunkIndex := 0

	buf := new(bytes.Buffer)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		n, err := io.ReadFull(resp.Body, buffer)
		isEOF := false

		if err != nil {
			if err == io.EOF {
				break
			}
			if err == io.ErrUnexpectedEOF {
				isEOF = true
			} else {
				return nil, err
			}
		}

		if n == 0 {
			break
		}

		if chunkIndex%3 == 0 && n == chunkSize {
			decrypted, err := Decrypt(buffer[:n], key)
			if err != nil {
				return nil, err
			}
			buf.Write(decrypted)
		} else {
			buf.Write(buffer[:n])
		}

		if isEOF {
			break
		}

		chunkIndex++
	}

	return buf, nil
}

func DownloadTrackFromURL(ctx context.Context, trackURL string) (*bytes.Buffer, string, error) {
	trackID := extractTrackID(trackURL)
	if trackID == "" {
		return nil, "", fmt.Errorf("invalid track URL")
	}

	session, err := Authenticate(ctx, os.Getenv("DEEZER_ARL"))
	if err != nil {
		return nil, "", fmt.Errorf("auth failed: %w", err)
	}

	song, err := FetchTrack(ctx, session, trackID)
	if err != nil {
		return nil, "", fmt.Errorf("fetch track failed: %w", err)
	}

	fmt.Printf("Downloading: %s - %s\n", song.Artist, song.Title)

	media, err := FetchMediaURL(ctx, session, song, quality)
	if err != nil {
		return nil, "", fmt.Errorf("fetch media failed: %w", err)
	}

	fileName := fmt.Sprintf("%s - %s.mp3", song.Artist, song.Title)

	buf, err := DownloadTrack(ctx, session, media.GetURL(), song)
	if err != nil {
		return nil, "", fmt.Errorf("download failed: %w", err)
	}

	fmt.Printf("✓ Downloaded to memory: %s\n", fileName)
	return buf, fileName, nil
}

func extractTrackID(url string) string {
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] < '0' || url[i] > '9' {
			return url[i+1:]
		}
	}
	return ""
}
