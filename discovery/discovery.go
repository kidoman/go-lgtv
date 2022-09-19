package discovery

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/kidoman/go-lgtv/control"
	iputil "github.com/kidoman/go-lgtv/util/ip"
	xmlutil "github.com/kidoman/go-lgtv/util/xml"
)

const openPort = 1426

var (
	// ErrNoTVFound indicates that no TV could be found
	ErrNoTVFound  = errors.New("Failed to find a TV")
	errNoResponse = errors.New("No response from IP address")
)

// Discover searches for LG TVs running version 3.5 of WebOS.
//
// This version has port 1042 open to GET requests to display details
// of the TV. gatewayIP is the IP address of the default gateway on the local network.
//
// This has only been tested using the C7V 2017 model.
func Discover(gatewayIP string) (*control.LgTv, error) {
	// Convert the provided string in to an IP address
	gwIP, err := iputil.ParseIP(gatewayIP)
	if err != nil {
		return nil, err
	}

	// Iterate over all possible local IP addresses (based on a single gateway setup
	for i := 0; i < 256; i++ {
		gwIP[3] = byte(i)
		found, _ := pingIP(gwIP)
		if found {
			return control.NewTV(string(gwIP), "", "")
		}
	}

	return nil, ErrNoTVFound
}

func pingIP(ip net.IP) (bool, error) {
	timeout := time.Duration(500 * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}
	fmt.Printf("http://%v:%v", ip, openPort)
	resp, _ := client.Get(fmt.Sprintf("http://%v:%v", ip, openPort))
	if resp == nil {
		return false, errNoResponse
	}

	body, _ := ioutil.ReadAll(resp.Body)

	node, err := xmlutil.FindXMLNode(string(body), "modelName")
	if err != nil {
		if err != xmlutil.ErrNodeNotFound {
			return false, err
		}

		return false, nil
	}

	if string(node.Content) == "LG TV" {
		return true, nil
	}
	return false, nil
}
