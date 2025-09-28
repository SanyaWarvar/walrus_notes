package restclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io"
	"log"
	"wn/pkg/applogger"
	"net/http"
	"os"
)

type RestClient interface {
	MakeRequest(ctx context.Context, req *http.Request) ([]byte, int, error)
}

type Client struct {
	client            *http.Client
	logger            applogger.Logger
	requestLogEnabled bool
	requestWithBody   bool
}

func (rc *Client) WithCert(certPath, keyPath string) *Client {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatal(err)
	}

	caCert, err := os.ReadFile(certPath)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	rc.client.Transport = otelhttp.NewTransport(
		&http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				RootCAs:            caCertPool,
				Certificates:       []tls.Certificate{cert},
			},
		},
	)

	return rc
}

func NewRestClient(logger applogger.Logger, requestLogEnabled, requestWithBody bool) *Client {
	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	return &Client{
		client:            client,
		logger:            logger,
		requestLogEnabled: requestLogEnabled,
		requestWithBody:   requestWithBody,
	}
}

func (rc *Client) MakeRequest(ctx context.Context, req *http.Request) ([]byte, int, error) {
	url := getURLFromUrl(req)
	if rc.requestLogEnabled {
		if rc.requestWithBody {
			requestBody := ""
			if req.Body != nil {
				body, err := io.ReadAll(req.Body)
				if err != nil {
					return nil, -1, err
				}
				requestBody = string(body)
				req.Body = io.NopCloser(bytes.NewBuffer(body))
			}

			reqMsgDebug := fmt.Sprintf("URL: %s, METHOD: %s, QUERY_PARAM: %s, HEADERS: %v, BODY: %s",
				url, req.Method, req.URL.RawQuery, req.Header, requestBody)
			rc.logger.WithCtx(ctx).Infof(reqMsgDebug)
		} else {
			rc.logger.WithCtx(ctx).Infof(fmt.Sprintf("URL: %s, METHOD: %s", url, req.Method))
		}
	}

	res, err := rc.client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, -1, err
	}

	if rc.requestLogEnabled {
		if rc.requestWithBody {
			rc.logger.WithCtx(ctx).Infof(fmt.Sprintf("URL: %s, STATUS: %d, BODY: %s", url, res.StatusCode, string(respBody)))
		} else {
			rc.logger.WithCtx(ctx).Infof(fmt.Sprintf("URL: %s, STATUS: %d", url, res.StatusCode))
		}
	}

	return respBody, res.StatusCode, nil
}

func getURLFromUrl(req *http.Request) string {
	return req.URL.Scheme + "://" + req.URL.Host + req.URL.Path
}
