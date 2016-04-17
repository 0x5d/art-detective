package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
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
	flag.Usage = printUsage
	flag.Parse()

	// Validate subject.
	if !supportedSubjects.Contains(*subject) {
		fmt.Printf("Subject \"%s\" not supported.\n\n", *subject)
		flag.Usage()
		os.Exit(1)
	}

	msg := "Investigating subject \"" + *subject + "\""
	if len(*id) != 0 {
		msg += " with id \"" + *id + "\""
	}
	msg += "."
	log.Println(msg)
	accessToken, _ := artsy.GetAccessToken(os.Getenv("ARTSY_CLIENT_ID"), os.Getenv("ARTSY_CLIENT_SECRET"))

	res, _ := artsy.Get(accessToken, *subject, *id)

	if len(*field) != 0 {
		var err error
		res, err = getField(res, strings.Split(*field, "."), accessToken)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(1)
		}
	}
	jsonRes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
	fmt.Println(string(jsonRes))
}

func getField(m map[string]interface{}, fieldPath []string, token string) (map[string]interface{}, error) {
	if len(fieldPath) == 1 {
		return artsy.Do("GET", token, m[fieldPath[0]].(string), nil)
	}
	return getField(m[fieldPath[0]].(map[string]interface{}), fieldPath[1:], token)
}

func printUsage() {
	msg := `art-detective

An artsy API client that follows links in responses.
Specify an endpoint and an optional ID for a specific resource, and art-detective will print the response.
Additionally, you may specify the path to an url field and art-detective will retrieve it too.

For example, to get similar artists to Andy Warhol:
  art-detective -subject artists -id 4d8b92b34eb68a1b2c0003f4 -field _links.similar_artists.href

Usage: art-detective [options]
Options:`
	fmt.Println(msg)
	flag.PrintDefaults()
}
