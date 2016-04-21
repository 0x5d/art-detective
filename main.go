package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/castillobg/art-detective/artsy"
	mapset "github.com/deckarep/golang-set"
)

var supportedSubjects = mapset.NewSetFromSlice([]interface{}{"artworks", "artists"})

func main() {
	// Init flags.
	subject := flag.String(
		"subject",
		"artworks",
		"(Required) The subject to investigate. Default is \"artwork\". Supported subjects are \"artworks\" and \"artists\"",
	)
	id := flag.String(
		"id",
		"",
		"(Optional) The id of the subject. If set, the subject with the given id will be retrieved, instead of the whole list.",
	)
	field := flag.String(
		"field",
		"",
		"A URL field in the body response that art-detective should follow. If not set, the initial raw response will be printed.",
	)
	help := flag.Bool("h", false, "Print the usage")
	flag.Usage = printUsage
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Validate subject.
	if !supportedSubjects.Contains(*subject) {
		fmt.Printf("Subject \"%s\" not supported.\n\n", *subject)
		flag.Usage()
		os.Exit(1)
	}

	err := investigate(*subject, *id, *field)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func investigate(subject, id, field string) error {
	msg := "Investigating subject \"" + subject + "\""
	if len(id) != 0 {
		msg += " with id \"" + id + "\""
	}
	msg += "."
	fmt.Println(msg)
	clientID := os.Getenv("ARTSY_CLIENT_ID")
	clientSecret := os.Getenv("ARTSY_CLIENT_SECRET")
	if len(clientID) == 0 || len(clientSecret) == 0 {
		return errors.New("Please set the ARTSY_CLIENT_ID and ARTSY_CLIENT_SECRET environment variables.")
	}
	accessToken, _ := artsy.GetAccessToken(clientID, clientSecret)

	res, _ := artsy.Get(accessToken, subject, id)

	if len(field) != 0 {
		var err error
		res, err = getField(res, strings.Split(field, "."), accessToken)
		if err != nil {
			return err
		}
	}
	jsonRes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonRes))
	return nil
}

func getField(m map[string]interface{}, fieldPath []string, token string) (map[string]interface{}, error) {
	if len(fieldPath) == 1 {
		return artsy.Do("GET", token, m[fieldPath[0]].(string), nil)
	}
	return getField(m[fieldPath[0]].(map[string]interface{}), fieldPath[1:], token)
}

func printUsage() {
	msg := `art-detective

An Artsy API client that follows links in responses.
Specify an endpoint and an optional ID for a specific resource, and art-detective will print the response.
Additionally, you may specify the path to an url field and art-detective will retrieve it too.

For example, to get similar artists to Andy Warhol:
  art-detective -subject artists -id 4d8b92b34eb68a1b2c0003f4 -field _links.similar_artists.href

Usage: art-detective [options]
Options:`
	fmt.Println(msg)
	flag.PrintDefaults()
}
