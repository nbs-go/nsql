package parse_test

import (
	"fmt"
	"github.com/nbs-go/nsql/parse"
	"github.com/nbs-go/nsql/test_utils"
	"log"
	"testing"
	"time"
)

func TestIntArgs(t *testing.T) {
	args := parse.IntArgs("9,1,2,a,b,c,2")
	test_utils.CompareInterfaceArray(t, "SAME ARGS", args, []interface{}{9, 1, 2})
}

func TestTimeArgNow(t *testing.T) {
	args := parse.TimeArg("now", "")
	test_utils.CompareInt(t, "NOW EXISTS", len(args), 1)
}

func TestTimeArgEpoch(t *testing.T) {
	ep := time.Now()
	args := parse.TimeArg(fmt.Sprintf("%d", ep.Unix()), "")
	exp := []interface{}{time.Unix(ep.Unix(), 0)}
	test_utils.CompareInterfaceArray(t, "EPOCH", args, exp)
}

func TestTimeArgLayoutDefault(t *testing.T) {
	// Create actual value
	v := "2022-01-01 00:00:00"
	args := parse.TimeArg(v, "")
	// Create expected
	tExp, _ := time.Parse(parse.DefaultTimeLayout, v)
	// Change test case
	actual := fmt.Sprintf("%s", args[0])
	expected := fmt.Sprintf("%s", tExp)
	if actual != expected {
		t.Errorf("unexpected parsed time. Expected = %s, Actual = %s", expected, actual)
	}
}

func TestTimeArgLayoutCustom(t *testing.T) {
	// Create actual value
	v := "2022-12-31T23:59:59+07:00"
	args := parse.TimeArg(v, time.RFC3339)

	// Create expected
	expT, err := time.Parse(time.RFC3339, v)
	if err != nil {
		log.Fatalf("cannot create expectation value. Error = %s", err)
	}
	exp := []interface{}{expT}

	test_utils.CompareInterfaceArray(t, "CUSTOM LAYOUT", args, exp)
}

func TestTimeArgError(t *testing.T) {
	// Create actual value
	v := "2022-12-31 23:59:59 07:00"
	args := parse.TimeArg(v, time.RFC3339)
	test_utils.CompareInt(t, "ERROR LAYOUT", len(args), 0)
}
