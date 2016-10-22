package hostinfo

import (
	"net"
	"net/http"
	"net/url"
	"encoding/json"
	"os"
	"bytes"
	"strings"
	"strconv"
)

//RequestInfoStruct Structure for Request Info
type RequestInfoStruct struct {
	Tag        string
	URL        url.URL
	Method     string
	Host       string
	RequestURI string
	Protocol   string
	Header     map[string]string
	RemoteAddr string
}

// RequestInfo Function Returns Details Retreived from the Request
func RequestInfo(r *http.Request) string {

	var rInfo RequestInfoStruct

	rInfo.Tag = "WFE: Server Info Page"
	rInfo.URL = *r.URL
	rInfo.Method = r.Method
	rInfo.Host = r.Host
	rInfo.RequestURI = r.RequestURI
	rInfo.Protocol = r.Proto
	rInfo.Header = make(map[string]string)
	for k, v := range r.Header {
		rInfo.Header[k] = v[0]
	}
	rInfo.RemoteAddr = r.RemoteAddr

	uglyJSON, _ := json.Marshal(rInfo)
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, uglyJSON, "", "\t")

	return prettyJSON.String()
}

//HostInfo Function Returns Details Retreived from the Request
func HostInfo() string {
	type serverInfo struct {
		OSHostName    string
		OSEnvironment map[string]string
		NETAddrs      map[string]string
		NETInts       map[string]map[string]string
	}

	var sInfo serverInfo

	sInfo.OSHostName, _ = os.Hostname()
	sInfo.OSEnvironment = make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		sInfo.OSEnvironment[pair[0]] = pair[1]
	}
	netAddrs, _ := net.InterfaceAddrs()
	sInfo.NETAddrs = make(map[string]string)
	for x, aa := range netAddrs {
		sInfo.NETAddrs[aa.Network()+strconv.Itoa(x)] = aa.String()
	}
	interfaces, _ := net.Interfaces()
	sInfo.NETInts = make(map[string]map[string]string)
	for _, aa := range interfaces {
		y := make(map[string]string)
		y["Index"] = strconv.Itoa(aa.Index)
		y["MTU"] = strconv.Itoa(aa.MTU)
		y["HardwareAddr"] = aa.HardwareAddr.String()
		y["Flags"] = aa.Flags.String()
		sInfo.NETInts[aa.Name] = y
	}
	uglyJSON, _ := json.Marshal(sInfo)
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, uglyJSON, "", "\t")

	return prettyJSON.String()
}

//LookupService Function Returns the Status of Hosts providing the Service
func LookupService(service string) (map[string]string, error) {
	checkService, err := net.LookupIP(service)
	if err != nil {
		return nil, err
	}
	hostStatus := make(map[string]string)
	for x, aa := range checkService {
		hostStatus[strconv.Itoa(x+1)] = aa.String()
	}
	return hostStatus, nil
}

//CheckServices Checks Services
func CheckServices(hosts ...string) string {
	serviceStatus := make(map[string]map[string]string)
	for _, v := range hosts {
		serviceStatus[v], _ = LookupService(v)
	}
	uglyJSON, _ := json.Marshal(serviceStatus)
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, uglyJSON, "", "\t")

	return prettyJSON.String()
}
