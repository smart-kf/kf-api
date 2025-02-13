package website

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFS(t *testing.T) {
	f, err := FS.Open("views/domain.html")
	require.Nil(t, err)
	data, _ := ioutil.ReadAll(f)
	f.Close()
	fmt.Println(string(data))
}
