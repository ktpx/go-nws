package main

//
// A Go client for NWS API queries.  Only active alerts.
// v1.0
//

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

type params map[string]interface{}
type doFunc func(req *http.Request) (*http.Response, error)

type Client struct {
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
	Logger     *log.Logger
	Debug      bool
	do         doFunc
}

func NewClient() *Client {
	return &Client{
		BaseURL:    "https://api.weather.gov/alerts",
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
	flag.StringVar(&opts.certainty, "c", "", "Specify certainty (Observed,Likely,Possible,Unlikely)")
	flag.StringVar(&opts.severity, "s", "", "Specify severity (Extreme,Severe,Moderate,Minor,Unknown)")
	flag.StringVar(&opts.service, "x", "", "Specify api service (unused)")
	flag.StringVar(&opts.status, "t", "", "Specify status (actual,exercise,system,test,draft)")
	flag.StringVar(&opts.region, "r", "", "Specify Marine Region code (AL,AT,GL,PA,PI).")
	flag.StringVar(&opts.region_type, "rt", "", "Specify region type (land, marine).")
	flag.StringVar(&opts.area, "area", "", "Specify area (AR,AH,CA,FL,....).")
	flag.StringVar(&opts.event, "e", "", "Specify Event Name.")
	flag.StringVar(&opts.zone, "z", "", "Specify Zone.")
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

var line = func(s string, i int) *string {
	l := fmt.Sprintf(strings.Repeat(s, i))
	return &l
}

func (data *NWSAlert) PrintReport() {

	b := bytes.Buffer{}
	cnt := 0
	for _, v := range data.Features {
		fmt.Fprintf(&b, "Event    : %s\n", v.Properties.Event)
		fmt.Fprintf(&b, "Headline : %s\n", v.Properties.Headline)
		fmt.Fprintf(&b, "Category : %s\n", v.Properties.Category)
		fmt.Fprintf(&b, "Urgency  : %s\n", v.Properties.Urgency)
		fmt.Fprintf(&b, "Type     : %s\n", v.Properties.Type)
		fmt.Fprintf(&b, "Sent     : %s\n", v.Properties.Sent)
		fmt.Fprintf(&b, "Effective: %s\n", v.Properties.Effective)
		fmt.Fprintf(&b, "Onset    : %s\n", v.Properties.Onset)
		fmt.Fprintf(&b, "Expires  : %s\n", v.Properties.Expires)
		fmt.Fprintf(&b, "Sender   : %s (%s)\n", v.Properties.SenderName, v.Properties.Sender)
		fmt.Fprintf(&b, "Msgtype  : %s\n", v.Properties.MessageType)
		fmt.Fprintf(&b, "Desc     :\n%s\n", v.Properties.Description)
		if len(v.Properties.Instruction) > 0 {
			fmt.Fprintf(&b, "%s\n", *line("-", 50))
			fmt.Fprintf(&b, "Instructions : %s\n", v.Properties.Instruction)
		}
		fmt.Fprintf(&b, "%s\n", *line("=", 80))
		cnt++

	}
	fmt.Println(b.String())
	fmt.Printf("%d Alerts listed.\n", cnt)
}

var baseURL = "https://api.weather.gov/alerts/active"
var opts Options

func main() {

	c := NewClient()
	parms := params{}

	flag.Parse()
	opts.SetParams(&parms)
	c.Debug = false
	c.debug("Options %v", opts)
	query := url.Values{}
	for k, v := range parms {
		query.Set(k, fmt.Sprintf("%v", v))
	}
	urlstring := query.Encode()
	fullURL := fmt.Sprintf("%s?%s", baseURL, urlstring)

	c.debug("FullURL: %s", fullURL)

	// Build new request
	r, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	c.debug("Request: %#v", r)
	r.Header.Set("UserAgent", "go-nwsclient")
	resp, err := c.HTTPClient.Do(r)
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
		log.Fatalf("Error: API returned: %s\n", resp.Status)
	}
	// Unmarshal response
	data := NWSAlert{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	data.PrintReport()
}
