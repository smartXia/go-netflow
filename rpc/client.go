package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	serverSuccessCode = 0
)

func RegisterDevice(ctx context.Context, urlProvider UrlProvider, device DeviceInfo) (*RegisterRsp, error) {
	body := map[string]interface{}{
		"data": &device,
	}
	var rspData RegisterRsp
	url, headers := urlProvider.Get("/api/common/v1/register")
	err := executeRequest(ctx, "RegisterDevice", "POST", url, headers, &body, &rspData)
	if err != nil {
		return nil, err
	}
	return &rspData, nil
}

type Client interface {
	ReportDeviceInfo(ctx context.Context, device DeviceInfo) error

	GetDialingInfo(ctx context.Context) ([]NetCardInfo, error)
	ReportNetCardInfo(ctx context.Context, netCards Netcards) error

	ReportHeartbeat(ctx context.Context) (*HeartbeatRsp, error)
	ReportMonitorInfo(ctx context.Context, monitorInfos []MonitorInfo) error

	ReportTest(ctx context.Context, networkTest []NetworkTestResult, diskResults []DiskResults) error

	ReportDeploymentStatus(ctx context.Context, deploymentStatusResult DeploymentStatusResult) error
}

func CreateRpcClient(urlProvider UrlProvider) Client {
	return &rpcClientImpl{urlProvider: urlProvider}
}

type rpcClientImpl struct {
	urlProvider UrlProvider
}

func (c *rpcClientImpl) ReportDeploymentStatus(ctx context.Context, deploymentStatusResult DeploymentStatusResult) error {
	body := map[string]interface{}{
		"data": &deploymentStatusResult,
	}
	url, headers := c.urlProvider.Get("/api/common/v1/update/deployment/status")
	return executeRequest(ctx, "ReportDeploymentStatus", "POST", url, headers, body, nil)
}

func (c *rpcClientImpl) ReportDeviceInfo(ctx context.Context, device DeviceInfo) error {
	body := map[string]interface{}{
		"data": &device,
	}
	url, headers := c.urlProvider.Get("/api/common/v1/sync/hardware")
	return executeRequest(ctx, "ReportDeviceInfo", "POST", url, headers, body, nil)
}

func (c *rpcClientImpl) ReportMonitorInfo(ctx context.Context, monitorInfos []MonitorInfo) error {
	body := map[string]interface{}{
		"data": &monitorInfos,
	}

	url, headers := c.urlProvider.Get("/api/common/v1/nethogs/monitor")
	return executeRequest(ctx, "ReportMonitorInfo", "POST", url, headers, body, nil)
}

func (c *rpcClientImpl) GetDialingInfo(ctx context.Context) ([]NetCardInfo, error) {
	url, headers := c.urlProvider.Get("/api/common/v1/dialinginfo")

	var netCards []NetCardInfo
	err := executeRequest(ctx, "GetDialingInfo", "GET", url, headers, nil, &netCards)
	return netCards, err
}

func (c *rpcClientImpl) ReportNetCardInfo(ctx context.Context, netCards Netcards) error {
	body := map[string]interface{}{
		"modifyTime": netCards.ModifyTime,
		"netCards":   netCards.Netcards,
	}
	url, headers := c.urlProvider.Get("/api/common/v1/sync/network")
	return executeRequest(ctx, "ReportNetCardInfo", "POST", url, headers, body, nil)
}

func (c *rpcClientImpl) ReportHeartbeat(ctx context.Context) (*HeartbeatRsp, error) {
	var rspData HeartbeatRsp
	url, headers := c.urlProvider.Get("/api/common/v1/heartbeat")
	err := executeRequest(ctx, "ReportHeartbeat", "POST", url, headers, []byte("{}"), &rspData)
	if err != nil {
		return nil, err
	}

	return &rspData, nil
}

func (c *rpcClientImpl) ReportTest(ctx context.Context, networkTest []NetworkTestResult, diskResults []DiskResults) error {
	body := map[string]interface{}{
		"networkResults": networkTest,
		"diskResults":    diskResults,
	}
	url, headers := c.urlProvider.Get("/api/common/v1/test/result")
	return executeRequest(ctx, "ReportTest", "POST", url, headers, body, nil)
}

func executeRequest(ctx context.Context, apiName string, method string, url string,
	headers map[string]string, body interface{}, rspDataRef interface{}) error {
	_, _, err := doExecuteRequest(ctx, method, url, headers, body, rspDataRef)

	//logFunc := log.GetLogger().Debug
	//logFields := []zap.Field{
	//	zap.String("url", url),
	//	zap.Any("headers", headers),
	//}
	//if rawRequest != "" || rawResponse != "" {
	//	fmt.Printf("\n>>>>>>>>>>>>  " + apiName + "  >>>>>>>>>>>>>>>\n")
	//	fmt.Printf(rawRequest)
	//	fmt.Printf("\n--------------------------------------\n")
	//	fmt.Printf(rawResponse)
	//	fmt.Printf("\n<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n")
	//}
	//if err != nil {
	//	logFunc = log.GetLogger().Error
	//	logFields = append(logFields, zap.Error(err))
	//} else {
	//	logFields = append(logFields, zap.Any("response", rspDataRef))
	//}
	//logFunc("executed rpc request: "+apiName, logFields...)

	return err
}

func doExecuteRequest(ctx context.Context, method string, url string, headers map[string]string,
	body interface{}, rspDataRef interface{}) (string, string, error) {
	var bodyReader io.Reader = nil
	if body != nil {
		if _, ok := body.([]byte); !ok {
			bytes, _ := json.Marshal(body)
			body = bytes
		}
		bodyReader = bytes.NewBuffer(body.([]byte))
	}
	httpReq, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return "", "", err
	}

	// avoid using `httpReq.Header.Add(k, v)` here as it would canonicalize
	// the header key, which the server may not accept it.
	if body != nil {
		httpReq.Header["Content-Type"] = []string{"application/json"}
	}
	for k, v := range headers {
		httpReq.Header[k] = []string{v}
	}

	var rawRequest string = ""
	//if log.GetLogger().Level() == zapcore.DebugLevel {
	//	b, _ := httputil.DumpRequest(httpReq, true)
	//	rawRequest = string(b)
	//}

	rsp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return rawRequest, "", err
	}
	defer rsp.Body.Close()

	rawResponse, err := parseResponse(rsp, rspDataRef)
	return rawRequest, rawResponse, err
}

func parseResponse(rsp *http.Response, dataRef interface{}) (string, error) {
	// We need to do the dump before ReadAll (dump would keep the rawResponse)
	var rawResponse string = ""
	//if log.GetLogger().Level() == zapcore.DebugLevel {
	//	b, _ := httputil.DumpResponse(rsp, true)
	//	rawResponse = string(b)
	//}

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	if rsp.StatusCode >= 200 && rsp.StatusCode < 300 {
		if dataRef != nil {
			rsp := ServerResponse{Data: dataRef}
			if err := json.Unmarshal(buf, &rsp); err != nil {
				return rawResponse, err
			}
			if rsp.Code != serverSuccessCode {
				return rawResponse, &ServerError{Code: rsp.Code, Message: rsp.Message, EventID: rsp.EventID}
			}
			return rawResponse, nil
		}

		var rsp ServerError
		if err := json.Unmarshal(buf, &rsp); err != nil {
			return rawResponse, err
		}
		if rsp.Code != serverSuccessCode {
			return rawResponse, &rsp
		}
		return rawResponse, nil
	}
	return rawResponse, &HttpError{Code: rsp.StatusCode, Message: string(buf)}
}

type ServerResponse struct {
	Code        int         `json:"code"`
	CurrentTime int64       `json:"current_time"`
	EventID     string      `json:"event_id"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data"`
}
