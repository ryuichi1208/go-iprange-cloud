package iprange

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type AWS struct {
	Synctoken  string `json:"syncToken"`
	Createdate string `json:"createDate"`
	Prefixes   []struct {
		IPPrefix           string `json:"ip_prefix,omitempty"`
		Region             string `json:"region"`
		Service            string `json:"service"`
		NetworkBorderGroup string `json:"network_border_group"`
		Ipv6Prefix         string `json:"ipv6_prefix,omitempty"`
	} `json:"prefixes"`
}

type GCP struct {
	Synctoken    string `json:"syncToken"`
	Creationtime string `json:"creationTime"`
	Prefixes     []struct {
		Ipv4Prefix string `json:"ipv4Prefix"`
		Service    string `json:"service"`
		Scope      string `json:"scope"`
	} `json:"prefixes"`
}

type Azure struct {
	Changenumber int    `json:"changeNumber"`
	Cloud        string `json:"cloud"`
	Values       []struct {
		Name       string `json:"name"`
		ID         string `json:"id"`
		Properties struct {
			Changenumber    int      `json:"changeNumber"`
			Region          string   `json:"region"`
			Regionid        int      `json:"regionId"`
			Platform        string   `json:"platform"`
			Systemservice   string   `json:"systemService"`
			Addressprefixes []string `json:"addressPrefixes"`
			Networkfeatures []string `json:"networkFeatures"`
		} `json:"properties"`
	} `json:"values"`
}

type Oracle struct {
	LastUpdatedTimestamp string `json:"last_updated_timestamp"`
	Regions              []struct {
		Region string `json:"region"`
		Cidrs  []struct {
			Cidr string   `json:"cidr"`
			Tags []string `json:"tags"`
		} `json:"cidrs"`
	} `json:"regions"`
}

type GitHub struct {
	VerifiablePasswordAuthentication bool `json:"verifiable_password_authentication"`
	SSHKeyFingerprints               struct {
		Sha256Rsa     string `json:"SHA256_RSA"`
		Sha256Ecdsa   string `json:"SHA256_ECDSA"`
		Sha256Ed25519 string `json:"SHA256_ED25519"`
	} `json:"ssh_key_fingerprints"`
	SSHKeys    []string `json:"ssh_keys"`
	Hooks      []string `json:"hooks"`
	Web        []string `json:"web"`
	API        []string `json:"api"`
	Git        []string `json:"git"`
	Packages   []string `json:"packages"`
	Pages      []string `json:"pages"`
	Importer   []string `json:"importer"`
	Actions    []string `json:"actions"`
	Dependabot []string `json:"dependabot"`
}

func doHttpRequest(url string) string {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("Not 200")
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func getAwsRange(body string) {
	var aws AWS
	json.Unmarshal([]byte(body), &aws)
	for _, prefix := range aws.Prefixes {
		fmt.Println(prefix.IPPrefix)
	}
}

func getGoogleRange(body string) {
	var gcp GCP
	json.Unmarshal([]byte(body), &gcp)
	var ips string
	for _, prefix := range gcp.Prefixes {
		ips = prefix.Ipv4Prefix
		if len(ips) > 0 {
			fmt.Println(prefix.Ipv4Prefix)
		}
	}
}

func getAzureRange(body string) {
	var azure Azure
	json.Unmarshal([]byte(body), &azure)
	for _, prefix := range azure.Values {
		for _, p := range prefix.Properties.Addressprefixes {
			fmt.Println(p)
		}
	}
}

func getOracleRange(body string) {
	var oracle Oracle
	json.Unmarshal([]byte(body), &oracle)
	for _, prefix := range oracle.Regions {
		for _, cidr := range prefix.Cidrs {
			fmt.Println(cidr.Cidr)
		}
	}
}

func getGitHubRange(body string) {
	var github GitHub
	json.Unmarshal([]byte(body), &github)
	for _, prefix := range github.Actions {
		fmt.Println(prefix)
	}
}

func getIpRangeFromFile(provider string) error {
	fileName := fmt.Sprintf("./data/%s.txt", provider)
	bytes, err := ioutil.ReadFile(fileName)
	fmt.Println(string(bytes))
	return err
}

func Run(args []string) {
	for _, p := range args[1:] {
		switch p {
		case "heroku":
			fallthrough
		case "aws":
			getAwsRange(doHttpRequest("https://ip-ranges.amazonaws.com/ip-ranges.json"))
		case "gcp":
			getGoogleRange(doHttpRequest("https://www.gstatic.com/ipranges/cloud.json"))
		case "googlebot":
			getGoogleRange(doHttpRequest("https://developers.google.com/search/apis/ipranges/googlebot.json"))
		case "azure":
			getAzureRange(doHttpRequest("https://download.microsoft.com/download/7/1/D/71D86715-5596-4529-9B13-DA13A5DE5B63/ServiceTags_Public_20220207.json"))
		case "oracle":
			getOracleRange(doHttpRequest("https://docs.oracle.com/en-us/iaas/tools/public_ip_ranges.json"))
		case "github":
			getGitHubRange(doHttpRequest("https://api.github.com/meta"))
		default:
			if err := getIpRangeFromFile(p); err != nil {
				fmt.Println("Not Found Provider")
				fmt.Println("    : heroku, aws, gcp, googlebot, oracle, github")
			}
		}
	}
}
