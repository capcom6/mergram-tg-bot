package mermaidinkfx

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const BaseURL = "https://mermaid.ink"

type Client struct {
	client *http.Client
}

func NewClient() (*Client, error) {
	return &Client{http.DefaultClient}, nil
}

func (c *Client) Render(ctx context.Context, diagram string) ([]byte, error) {
	input := Request{
		Code: diagram,
	}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	encoded, err := Encode(inputBytes)
	if err != nil {
		return nil, fmt.Errorf("encode diagram: %w", err)
	}

	url := fmt.Sprintf("%s/img/pako:%s", BaseURL, encoded)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get diagram: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: unexpected status code: %d", ErrAPIError, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return data, nil
}

// Encode takes a string and returns an encoded string in deflate + base64 format.
func Encode(input []byte) (string, error) {
	var buffer bytes.Buffer
	writer := zlib.NewWriter(&buffer)

	if _, err := writer.Write(input); err != nil {
		return "", fmt.Errorf("write to zlib writer: %w", err)
	}

	if err := errors.Join(writer.Flush(), writer.Close()); err != nil {
		return "", fmt.Errorf("close zlib writer: %w", err)
	}

	result := base64.URLEncoding.EncodeToString(buffer.Bytes())
	return result, nil
}
