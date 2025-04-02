package markdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Split_MarkdownText(t *testing.T) {
	t.Parallel()
	type testCase struct {
		markdown      string
		expectedParts []string
	}
	testCases := []testCase{
		{
			markdown: `## First header: h2
Some content below the first h2.
## Second header: h2
### Third header: h3

- This is a list item of bullet type.
- This is another list item.

 *Everything* is going according to **plan**.

# Fourth header: h1
Some content below the first h1.
## Fifth header: h2
#### Sixth header: h4

Some content below h1>h2>h4.`,
			expectedParts: []string{
				"## First header: h2\nSome content below the first h2.\n## Second header: h2\n### Third header: h3\n\n- This is a list item of bullet type.\n- This is another list item.\n\n",
				" *Everything* is going according to **plan**.\n\n# Fourth header: h1\nSome content below the first h1.\n## Fifth header: h2\n#### Sixth header: h4\n\nSome content below h1>h2>h4.",
			},
		},
	}

	for _, tc := range testCases {
		parts := Split(tc.markdown, 200)
		assert.Equal(t, tc.expectedParts, parts)
	}
}

func Test_Split_MarkdownCode(t *testing.T) {
	t.Parallel()
	type testCase struct {
		markdown      string
		expectedParts []string
	}
	testCases := []testCase{
		{
			markdown: "## First header: h2\nSome content below the first h2.\n## Second header: h2\n### Third header: h3\n\n- This is a list item of bullet type.\n```go\npackage main\n\nimport \"fmt\"\n\nfunc helloWorld() {\n    fmt.Println(\"Hello, World!\")\n}\n\nfunc main() {\n    helloWorld()\n}\n\n``` *Everything* is going according to **plan**.\n\n# Fourth header: h1\nSome content below the first h1.\n## Fifth header: h2\n#### Sixth header: h4\n\nSome content below h1>h2>h4.",
			expectedParts: []string{
				"## First header: h2\nSome content below the first h2.\n## Second header: h2\n### Third header: h3\n\n- This is a list item of bullet type.\n```go\npackage main\n\nimport \"fmt\"\n\nfunc helloWorld() {\n\n```\n",
				"```go\n    fmt.Println(\"Hello, World!\")\n}\n\nfunc main() {\n    helloWorld()\n}\n\n``` *Everything* is going according to **plan**.\n\n# Fourth header: h1\nSome content below the first h1.\n## Fifth header: h2\n",
				"#### Sixth header: h4\n\nSome content below h1>h2>h4.",
			},
		},
	}

	for _, tc := range testCases {
		parts := Split(tc.markdown, 200)
		assert.Equal(t, tc.expectedParts, parts)
	}
}
