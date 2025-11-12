package manager

import (
	"fmt"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/inhere/kite-go/pkg/xenv/models"
)

func TestTomlUpdater(t *testing.T) {
	up := NewTomlUpdater()
	up.SetContents([]byte(`
paths = []

[sdks]
go = "1.22"
node = "20"

[envs]

[tools]
`))

	// no-change version
	t.Run("no-change version", func(t *testing.T) {
		bs := up.Build(&models.ActivityState{
			SDKs: map[string]string{
				"go":   "1.22",
				"node": "20",
			},
		}).LastContents()
		s := string(bs)
		fmt.Println(s)
		assert.Equal(t, `
paths = [
]

[sdks]
go = "1.22"
node = "20"

[envs]
[tools]
`, s)
	})

	// change version
	t.Run("change version", func(t *testing.T) {
		bs := up.Build(&models.ActivityState{
			SDKs: map[string]string{
				"go":   "1.23",
				"node": "21",
				"java": "17",
			},
		}).LastContents()
		s := string(bs)
		fmt.Println(s)
		assert.Equal(t, `
paths = [
]

[sdks]
go = "1.23"
node = "21"
java = "17"

[envs]
[tools]
`, s)
	})
}
