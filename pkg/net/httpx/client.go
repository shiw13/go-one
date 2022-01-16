package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"go-one/pkg/logger"
	"go-one/pkg/util/bytesconv"

	"go.uber.org/zap"
)

type ClientConfig struct {
	Timeout      time.Duration
	DebugEnabled bool
}

type Client struct {
	cfg       *ClientConfig
	rawClient http.Client
}

func NewClient(cfg ClientConfig) *Client {
	c := &Client{
		cfg: &cfg,
		rawClient: http.Client{
			Timeout: cfg.Timeout,
		},
	}

	return c
}

func NewJSONRequest(method, url string, header http.Header, jo interface{}) (*http.Request, error) {
	body, err := json.Marshal(jo)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = header
	}
	req.Header.Set(HeaderContentType, ContentTypeJSON)

	return req, nil
}

func (c *Client) JSON(ctx context.Context, req *http.Request, res interface{}) error {
	var bs []byte
	var err error

	if bs, err = c.Send(ctx, req); err != nil {
		return err
	}

	if res != nil {
		if err = json.Unmarshal(bs, res); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) Send(ctx context.Context, req *http.Request) ([]byte, error) {
	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	var bs []byte
	bs, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	const logPrintMaxLen = 10 * 1024

	var reqBodyStr string
	if c.cfg.DebugEnabled && req.Body != nil {
		if bs, e := io.ReadAll(req.Body); e == nil {
			req.Body = io.NopCloser(bytes.NewBuffer(bs))
			if len(bs) < logPrintMaxLen {
				reqBodyStr = bytesconv.UnsafeBytesToString(bs)
			} else {
				reqBodyStr = "body len over limit"
			}
		}
	}

	begin := time.Now()
	resp, err := c.rawClient.Do(req.WithContext(ctx))
	latency := time.Since(begin)

	if err != nil {
		if c.cfg.DebugEnabled {
			logger.Zap().Info("http client",
				zap.String(logger.HTTPMethod, req.Method),
				zap.String(logger.URLPath, req.URL.Path),
				zap.String(logger.URLQuery, req.URL.Query().Encode()),
				zap.String(logger.DstAddr, req.Host),
				zap.String(logger.HTTPReqBody, reqBodyStr),
				zap.Int64("latency", latency.Milliseconds()),
			)
		}
		return nil, err
	}

	var resBodyStr string
	if c.cfg.DebugEnabled {
		if resp.Body != nil {
			if bs, e := io.ReadAll(resp.Body); e == nil {
				resp.Body = io.NopCloser(bytes.NewBuffer(bs))
				if len(bs) < logPrintMaxLen {
					resBodyStr = bytesconv.UnsafeBytesToString(bs)
				} else {
					resBodyStr = "body len over limit"
				}
			}
		}

		logger.Zap().Info("http client",
			zap.String(logger.HTTPMethod, req.Method),
			zap.String(logger.URLPath, req.URL.Path),
			zap.String(logger.URLQuery, req.URL.Query().Encode()),
			zap.String(logger.DstAddr, req.Host),
			zap.String(logger.HTTPReqBody, reqBodyStr),
			zap.Int(logger.HTTPStatusCode, resp.StatusCode),
			zap.String(logger.HTTPResBody, resBodyStr),
			zap.Int64("latency", latency.Milliseconds()),
		)
	}

	return resp, nil
}

func (c *Client) GetJSON(ctx context.Context, url string, params url.Values, header http.Header, res interface{}) error {
	if params != nil {
		url = url + "?" + params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	if header != nil {
		req.Header = header
	}

	return c.JSON(ctx, req, res)
}

func (c *Client) PostJSON(ctx context.Context, url string, params url.Values, header http.Header,
	reqJo interface{}, res interface{}) (err error) {
	if params != nil {
		url = url + "?" + params.Encode()
	}

	var req *http.Request
	if req, err = NewJSONRequest(http.MethodPost, url, header, reqJo); err != nil {
		return err
	}

	return c.JSON(ctx, req, res)
}
