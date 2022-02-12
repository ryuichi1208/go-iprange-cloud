package iprange

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestAws(t *testing.T) {
	provider := "aws"
	fileName := fmt.Sprintf("../data/%s.json", provider)
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	getAwsRange(string(bytes))
}
