package main

/////
// A Go Client for NWS API queries.  Only active alerts and count enpoints.
// v1.1
// @Jinxd  (MIT License)
/////

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

type NWSAlert struct {
	Context  interface{}    `json:"@context"`
	Type     string         `json:"type"`
	Features []FeaturesItem `json:"features"`
	Title    string         `json:"title"`
	Updated  string         `json:"updated"`
}

type FeaturesItem struct {
	Id         string `json:"id"`
	Type       string
	Geometry   interface{}
	Properties struct {
		Type          string      `json:"@type"`
		Id            string      `json:"@id"`
		AreaDesc      string      `json:"areaDesc"`
		GeoCode       interface{} `json:"geocode"`
		AffectedZones interface{} `json:"affectedzones"`
		Sent          string      `json:"sent"`
		Effective     string      `json:"effective"`
		Onset         string      `json:"onset"`
		Expires       string      `json."expires"`
		Ends          string      `json:"ends"`
		MessageType   string      `json:"messageType"`
		Category      string      `json:"category"`
		Certainty     string      `json:"certainty"`
		Urgency       string      `json:"urgency"`
		Event         string      `json:"event"`
		Sender        string      `json:"sender"`
		SenderName    string      `json:"senderName"`
		Headline      string      `json:"headline"`
		Description   string      `json:"Description"`
		Response      string      `json:"response"`
		Instruction   string      `json:"instruction"`
		Paramaters    interface{} `json:"parameters"`
	} `json:"properties"`
}

// Command options
type Options struct {
	status       string // actual,exercise,system,test,draft, arr
	message_type string // alert, update, cancel, arr
	event        string // Event Name, arr
	code         string // Event Code, arr					//
	area         string // State/territory code or marine area,arr
	point        string // Point (lat,long)
	region       string // Marine Region Code (AL,AT,GL,PA,PI), arr
	region_type  string // land, marine
	zone         string // Zone ID (forecast only), arr
	urgency      string // Immediate, Expected,Future,Past,Uknown, arr
	severity     string // Extreme, Severe, Moderate, Minor, Unknown, arr
	certainty    string // Observed,Likely,Possible,Unlikely,Unknown, arr
	limit        int    // Query limit
	service      string // API service
}

// Result Structure for /active/count
type NWSAlertCount struct {
	Total   int            `json:"total"`
	Land    int            `json:"land"`
	Marine  int            `json:"marine"`
	Regions map[string]int `json:"regions"`
	Areas   map[string]int `json:"areas"`
	Zones   map[string]int `json:"zones"`
}

type params map[string]interface{}

//type doFunc func(req *http.Request) (*http.Response, error)

type Client struct {
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
	Logger     *log.Logger
	Debug      bool
}

func NewClient() *Client {
	return &Client{
		BaseURL:    "https://api.weather.gov",
		UserAgent:  "go-nws",
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, "Go-nws", log.LstdFlags),
		Debug:      false,
	}
}

func (c *Client) debug(format string, v ...interface{}) {
	if c.Debug {
		c.Logger.Printf(format, v...)
	}
}

type request struct {
	method   string
	fullURL  string
	endpoint string
	query    url.Values
	form     url.Values
	header   http.Header
	body     io.Reader
}

func init() {
	flag.StringVar(&opts.area, "area", "", "Specify area (AR,AH,CA,FL,....).")
	flag.StringVar(&opts.certainty, "c", "", "Specify certainty (Observed,Likely,Possible,Unlikely)")
	flag.StringVar(&opts.severity, "s", "", "Specify severity (Extreme,Severe,Moderate,Minor,Unknown)")
	flag.StringVar(&opts.service, "x", "", "Specify api endpoint (active,count)")
	flag.StringVar(&opts.status, "t", "", "Specify status (actual,exercise,system,test,draft)")
	flag.StringVar(&opts.region, "r", "", "Specify Marine Region code (AL,AT,GL,PA,PI).")
	flag.StringVar(&opts.region_type, "rt", "", "Specify region type (land, marine).")
	flag.StringVar(&opts.event, "e", "", "Specify Event Name.")
	flag.StringVar(&opts.zone, "z", "", "Specify Zone.")
	flag.StringVar(&opts.urgency, "u", "", "Specify Urgency (Immediate,Expected,Future,Past,Unknown)")
}

// Map command line options to parameters
func (o *Options) SetParams(p *params) {
	if opts.urgency != "" {
		(*p)["urgency"] = opts.urgency
	}
	if opts.area != "" {
		(*p)["area"] = strings.ToUpper(opts.area)
	}
	if opts.severity != "" {
		(*p)["severity"] = opts.severity
	}
	if opts.certainty != "" {
		(*p)["certainty"] = opts.certainty
	}
	if opts.region != "" {
		(*p)["region"] = opts.region
	}
	if opts.region_type != "" {
		(*p)["region_type"] = opts.region_type
	}
	if opts.event != "" {
		(*p)["event"] = opts.event
	}
	if opts.zone != "" {
		(*p)["zone"] = opts.zone
	}
}

func Exists(s string, a []string) bool {
	for _, v := range a {
		if s == v {
			return true
		}
	}
	return false
}

var line = func(s string, i int) *string {
	l := fmt.Sprintf(strings.Repeat(s, i))
	return &l
}

func (data *NWSAlertCount) PrintReport() {

	l := *line("-", 40)
	fmt.Printf("%s\nTotal Alert Count Report\n%s\n", l, l)
	fmt.Printf("Total  : %d\n", data.Total)
	fmt.Printf("Land   : %d\n", data.Land)
	fmt.Printf("Marine : %d\n", data.Marine)

	// Sort the keys by area name
	keys := make([]string, 0, len(data.Areas))
	for k := range data.Areas {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	b := bytes.Buffer{}
	var c int
	fmt.Fprintf(&b, "%s\nAlerts per Area\n%s\n", l, l)
	for _, k := range keys {
		c++
		fmt.Fprintf(&b, " %s : %3d ", k, data.Areas[k])
		if c%4 == 0 {
			fmt.Fprintf(&b, "%s", "\n")
		}
	}
	fmt.Println(b.String())
}

func (data *NWSAlert) PrintReport() {

	b := bytes.Buffer{}
	cnt := 0
	sep := *line("=", 50)
	for _, v := range data.Features {
		fmt.Fprintf(&b, "%s\nEvent    : %s\n", sep, v.Properties.Event)
		fmt.Fprintf(&b, "Headline : %s\n", v.Properties.Headline)
		fmt.Fprintf(&b, "Category : %s\n", v.Properties.Category)
		fmt.Fprintf(&b, "Msgtype  : %s\n", v.Properties.MessageType)
		fmt.Fprintf(&b, "Urgency  : %s\n", v.Properties.Urgency)
		fmt.Fprintf(&b, "Certainty: %s\n", v.Properties.Certainty)
		fmt.Fprintf(&b, "Type     : %s\n", v.Properties.Type)
		fmt.Fprintf(&b, "Sent     : %s\n", v.Properties.Sent)
		fmt.Fprintf(&b, "Effective: %s\n", v.Properties.Effective)
		fmt.Fprintf(&b, "Onset    : %s\n", v.Properties.Onset)
		fmt.Fprintf(&b, "Expires  : %s\n", v.Properties.Expires)
		fmt.Fprintf(&b, "Sender   : %s (%s)\n", v.Properties.SenderName, v.Properties.Sender)
		fmt.Fprintf(&b, "AreaDesc : %s\n", v.Properties.AreaDesc)
		fmt.Fprintf(&b, "Description  :\n%s\n", v.Properties.Description)
		if len(v.Properties.Instruction) > 0 {
			fmt.Fprintf(&b, "%s\n", *line("--", 25))
			fmt.Fprintf(&b, "Instructions : %s\n", v.Properties.Instruction)
		}
		fmt.Fprintf(&b, "%s\n", sep)
		cnt++

	}
	fmt.Println(b.String())
	fmt.Printf("%d Alerts listed.\n", cnt)
}

func (c *Client) callAPI(r *request, p *params) (data []byte, err error) {

	// Build Query String
	query := url.Values{}
	var urlstring string
	if p != nil {
		for k, v := range *p {
			query.Set(k, fmt.Sprintf("%v", v))
		}
		urlstring = query.Encode()
	}
	fullURL, _ := url.JoinPath(c.BaseURL, r.endpoint)
	if urlstring != "" {
		fullURL = fmt.Sprintf("%s?%s", fullURL, urlstring)
	}
	c.debug("FullURL: %s", fullURL)

	// Build new request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	c.debug("Request: %#v", r)
	if r.method != "" {
		req.Method = r.method
	}
	if r.header != nil {
		req.Header = r.header
	}
	req.Header.Set("UserAgent", "go-nwsclient")
	req.Header.Set("Accept", "application/geo+json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	c.debug("response: %#v", resp)
	c.debug("response body: %s", string(body))
	c.debug("response status code: %d", resp.StatusCode)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode > http.StatusOK {
		if resp.StatusCode == 503 {
			log.Println("The API is probably too busy, try again in a few.")
		}
		log.Fatalf("Error: API returned: %s\n", resp.Status)
		/*              apiErr := new(APIError)
		                e := json.Unmarshal(data, apiErr)
		                if e != nil {
		                        c.debug("failed to unmarshal json: %s", e)
		                }
		                return nil, apiErr
		*/
	}
	return body, nil
}

// alert/count endpoint
func (c *Client) NWSAlertCount(p *params) (*NWSAlertCount, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/alerts/active/count",
	}
	res, err := c.callAPI(r, nil)
	if err != nil {
		return nil, err
	}
	j := new(NWSAlertCount)
	err = json.Unmarshal(res, &j)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// alert/active endpoint
func (c *Client) NWSAlertActive(p *params) (*NWSAlert, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/alerts/active",
	}
	res, err := c.callAPI(r, p)
	// Unmarshal response
	data := NWSAlert{}
	err = json.Unmarshal(res, &data)
	if err != nil {
		log.Fatal(err)
	}
	return &data, nil

}

var opts Options
var MarineRegionCodes = [6]string{"AL", "AT", "GL", "GM", "PA", "PI"}
var RegionCodes = [6]string{"AR", "CR", "ER", "PR", "SR", "WR"}
var MarineAreaCodes = [15]string{
	"AM", "AN", "GM", "LC", "LE", "LH", "LM", "LO",
	"LS", "PH", "PK", "PM", "PS", "PZ", "SL",
}

func main() {

	c := NewClient()
	flag.Parse()

	c.Debug = false
	c.debug("Options %+v", opts)

	parms := params{}
	opts.SetParams(&parms)

	switch opts.service {
	case "count":
		data, err := c.NWSAlertCount(&parms)
		if err != nil {
			panic(err)
		}
		data.PrintReport()
	case "alerts":
		data, err := c.NWSAlertActive(&parms)
		if err != nil {
			panic(err)
		}
		data.PrintReport()
	default:
		fmt.Println("Go-nws: Specify report type -x <alerts,count> -h for more info and options.")
	}
}
