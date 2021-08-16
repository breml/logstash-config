package config_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	. "github.com/breml/logstash-config"
	"github.com/breml/logstash-config/ast"
)

func ExampleParseReader() {
	logstashConfig := `filter {
    mutate {
      add_tag => [
        "tag"
      ]
    }
  }`
	got, err := ParseReader("example.conf", strings.NewReader(logstashConfig))
	if err != nil {
		log.Fatalf("Parse error: %s\n", err)
	}

	// Output: filter {
	//   mutate {
	//     add_tag => [
	//       "tag"
	//     ]
	//   }
	//}
	fmt.Println(got)
}

func TestParserIdentic(t *testing.T) {
	cases := []struct {
		name string

		input string
	}{
		{
			name:  "Empty file",
			input: ``,
		},
		{
			name: "Single PluginSection",
			input: `input {}
`,
		},
		{
			name: "All PluginSections empty",
			input: `input {}
filter {}
output {}
`,
		},
		{
			name: "Plugin without attributes",
			input: `input {
  stdin {}
}
`,
		},
		{
			name: "Multiple plugins",
			input: `input {
  stdin {}

  file {}
}
filter {
  mutate {}

  mutate {}

  mutate {}
}
output {
  stdout {}
}
`,
		},
		{
			name:  "Plugin with all attribute types",
			input: ``,
		},
		{
			name: "Simple if (without else) branch",
			input: `filter {
  if 1 == 1 {
    date {}
  }
}
`,
		},
		{
			name: "if with else-if and a final else branch",
			input: `filter {
  if 1 == 1 {
    date {}
  } else if 1 == 1 {
    date {}
  } else {
    date {}
  }
}
`,
		},
		{
			name: "if with multiple else-if and a final else branch, multiple plugins in each branch",
			input: `filter {
  if 1 == 1 {
    date {}

    date {}
  } else if 1 == 1 {
    date {}

    date {}

    date {}
  } else if 1 == 1 {
    date {}

    date {}

    date {}
  } else {
    date {}

    date {}

    date {}
  }
}
`,
		},
		{
			name: "if with multiple else-if and a final else branch, test for different condition types",
			input: `filter {
  if 1 != 1 {
    date {}
  } else if 1 <= 1 {
    date {}
  } else if 1 >= 1 {
    date {}
  } else if 1 < 1 {
    date {}
  } else if 1 > 1 {
    date {}
  }
}
`,
		},
		{
			name: "if with multiple compare operators",
			input: `filter {
  if "true" == "true" and "true1" == "true1" or "true2" == "true2" nand "true3" == "true3" xor "true4" == "true4" {
    plugin {}
  }
}
`,
		},
		{
			name: "Condition in parentheses",
			input: `filter {
  if ("tag" in [tags]) {
    plugin {}
  }
}
`,
		},
		{
			name: "Multiple conditions in parentheses",
			input: `filter {
  if ("tag" in [tags] or ("true" == "true" and 1 == 1)) {
    plugin {}
  }
}
`,
		},
		{
			name: "Negative condition",
			input: `filter {
  if !("true" == "true") {
    plugin {}
  }
}
`,
		},
		{
			name: "Negative Selector for value in subfield",
			input: `filter {
  if ![field][subfield] {
    plugin {}
  }
}
`,
		},
		{
			name: "InExpression",
			input: `filter {
  if "tag" in [tags] {
    plugin {}
  }
}
`,
		},
		{
			name: "NotInExpression",
			input: `filter {
  if "tag" not in [field][subfield] {
    plugin {}
  }
}
`,
		},
		{
			name: "RegexpExpression (Match)",
			input: `filter {
  if [field] =~ /.*/ {
    plugin {}
  }
}
`,
		},
		{
			name: "RegexpExpression (Not Match)",
			input: `filter {
  if [field] !~ /.*/ {
    plugin {}
  }
}
`,
		},
		{
			name: "Rvalue",
			input: `filter {
  if "string" or 10 or [field][subfield] or /.*/ {
    plugin {}
  }
}
`,
		},
		{
			name: "Empty array",
			input: `filter {
  plugin {
    value => []
  }
}
`,
		},
		{
			name: "Multiple filter sections",
			input: `filter {}

filter {}

filter {}
`,
		},

		// Comments
		{
			name: "plugin section comment",
			input: `# Comment
filter {}
`,
		},
		{
			name: "file header and plugin section comment",
			input: `# Comment

# Comment 2

# Comment 3
filter {}
`,
		},
		{
			name: "plugin section comment with whitespace after",
			input: `# Comment

filter {}
`,
		},
		{
			name: "file footer comment",
			input: `filter {}
output {}

# Comment after
`,
		},
		{
			name: "foobar",
			input: `input {
  stdin {
    # Comment
    codec => rubydebug {
      # Comment
      string => "a string"

      # Comment
    }

    # Comment
  }
}
`,
		},
		{
			name: "Empty array with comment",
			input: `input {
  stdin {
    arrayvalue => [
      # Comment
    ]
  }
}
`,
		},
		{
			name: "Empty hash with comment",
			input: `input {
  stdin {
    hashvalue => {
      # Comment
    }
  }
}
`,
		},
		{
			name: "comment only if, else-if and else",
			input: `filter {
  if 1 == 1 {
    # Comment
  } else if 1 == 1 {
    # Comment
  } else {
    # Comment
  }
}
`,
		},
		// https://github.com/magnusbaeck/logstash-filter-verifier/issues/104
		{
			name: "large numbers with and without precission",
			input: `filter {
  mutate {
    add_field => {
      "largeint" => 2000000
      "largeint_negative" => -2000000
      "int32max" => 2147483647
      "int32max_negative" => -2147483647
      "float_with_5_precision" => 0.00001
      "float_with_5_precision_negative" => -0.00001
      "float_with_10_precision" => 0.0000000001
      "float_with_10_precision_negative" => -0.0000000001
      "largefloat_with_10_precision" => 2000000.0000000005
      "largefloat_with_10_precision_negative" => -2000000.0000000005
      "int32max_with_10_precision" => 2147483647.0000009537
      "int32max_with_10_precision_negative" => -2147483647.0000009537
    }
  }
}
`,
		},
		{
			name: "multiline string",
			input: `filter {
  mutate {
    add_field => {
      "largeint" => "a string
with multiple
lines"
    }
  }
}
`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got1, err := ParseReader("test", strings.NewReader(test.input))
			if err != nil {
				t.Fatalf("Expected to parse without error: %s, input:\n%s", err, test.input)
			}
			got := fmt.Sprintf("%v", got1)
			if test.input != got {
				t.Errorf("Expected parsed input to print the same as input:\n%s", printDiff(test.input, got))
			}
		})
	}
}

func TestParserIdenticFile(t *testing.T) {
	cases := []string{
		"plugin_with_all_attribute_types",
		"plugin_with_all_attribute_types_with_comments",
		"if_else-if_and_else_branch_with_comments",
	}

	for _, test := range cases {
		t.Run(test, func(t *testing.T) {
			inputFilename := "testdata/identic/" + test + ".conf"
			got1, err := ParseFile(inputFilename)
			if err != nil {
				t.Fatalf("Expected %q to parse without error: %v", test, err)
			}
			expected, err := ioutil.ReadFile(inputFilename)
			if err != nil {
				t.Fatalf("Unable to read file %q: %v", inputFilename, err)
			}
			got := fmt.Sprintf("%v", got1)
			if string(expected) != got {
				t.Errorf("Expected %q parsed input to print the same as input:\n%s", inputFilename, printDiff(string(expected), got))
			}
		})
	}
}

func TestParser(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Whitespace, tab and newlines only",
			input: `
      
			`,
			expected: ``,
		},
		{
			name:     "Single comment (one line without newline)",
			input:    `# comment only`,
			expected: ``,
		},
		{
			name: "Comment surrounded by empty lines",
			input: `
# comment only
`,
			expected: ``,
		},
		{
			name:  "Reformat plugins",
			input: `input { stdin {} }`,
			expected: `input {
  stdin {}
}
`,
		},
		{
			name: "Comments and whitespace",
			input: `input { 
  # Comment
  stdin {

    # Comment

  }
}`,
			expected: `input {
  # Comment
  stdin {
    # Comment
  }
}
`,
		},
		{
			name: "Multiple filter sections without empty lines",
			input: `filter {}
filter {}
filter {}
`,
			expected: `filter {}

filter {}

filter {}
`,
		},

		// Comment
		{
			name: "plugin section comment with whitespace before",
			input: `

# Comment
filter {}
`,
			expected: `# Comment
filter {}
`,
		},
		{
			name: "plugin section comment with whitespace before and after",
			input: `

# Comment

filter {}
`,
			expected: `# Comment

filter {}
`,
		},
		{
			name: "file header, footer and plugin section comments with whitespace before and after",
			input: `# Pre Filter comment

# Filter comment
filter {}

# Input comment
input {}

# File footer comment
`,
			expected: `# Input comment
input {}
# Pre Filter comment

# Filter comment
filter {}

# File footer comment
`,
		},
		{
			name: "file footer comment without spaceBefore",
			input: `filter {}
# Comment after
`,
			expected: `filter {}

# Comment after
`,
		},
		{
			name: "pluginSection footer comment with and without spaceBefore",
			input: `filter {
  plugin {}

  # pluginSection footer comment
}

filter {
  plugin {}
  # pluginSection footer comment
}
`,
			expected: `filter {
  plugin {}

  # pluginSection footer comment
}

filter {
  plugin {}

  # pluginSection footer comment
}
`,
		},
		{
			name: "pluginSection footer comment empty block with and without spaceBefore",
			input: `filter {

  # pluginSection footer comment
}

filter {
  # pluginSection footer comment
}
`,
			expected: `filter {
  # pluginSection footer comment
}

filter {
  # pluginSection footer comment
}
`,
		},
		{
			name: "plugins without comments",
			input: `filter {
  plugin {}
  otherplugin {}

  thirdplugin {}
}
`,
			expected: `filter {
  plugin {}

  otherplugin {}

  thirdplugin {}
}
`,
		},
		{
			name: "plugins with comments",
			input: `filter {
  # plugin comment
  plugin {}
  # otherplugin comment
  otherplugin {}

  # third plugin
  thirdplugin {}
}
`,
			expected: `filter {
  # plugin comment
  plugin {}

  # otherplugin comment
  otherplugin {}

  # third plugin
  thirdplugin {}
}
`,
		},
		{
			name: "plugin with multiple comments",
			input: `filter {
  # additional comment

  # plugin comment
# multiple lines
  plugin {}
  # footer comment
}
`,
			expected: `filter {
  # additional comment

  # plugin comment
  # multiple lines
  plugin {}

  # footer comment
}
`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got1, err := ParseReader("test", strings.NewReader(test.input))
			if err != nil {
				t.Fatalf("Expected to parse without error: %s, input:\n|%s|", err, test.input)
			}
			got := fmt.Sprintf("%v", got1)
			if test.expected != got {
				t.Errorf("Expected output does not match parsed output:\n%s", printDiff(test.expected, got))
			}
		})
	}
}

func TestParserFile(t *testing.T) {
	cases := []string{
		"comments_everywhere",
	}

	for _, test := range cases {
		t.Run(test, func(t *testing.T) {
			inputFilename := "testdata/parser/" + test + ".conf"
			expectedFilename := "testdata/parser/" + test + ".expected.conf"
			got1, err := ParseFile(inputFilename)
			if err != nil {
				t.Fatalf("Expected %q to parse without error: %v", test, err)
			}
			expected, err := ioutil.ReadFile(expectedFilename)
			if err != nil {
				t.Fatalf("Unable to read file %q: %v", expectedFilename, err)
			}
			got := fmt.Sprintf("%v", got1)
			if string(expected) != got {
				fmt.Println(got)
				t.Errorf("Expected %q parsed input to print the same as input:\n%s", inputFilename, printDiff(string(expected), got))
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name: "misspelled plugin section",
			input: `FILTER {
  if 1 == 1 {
    plugin{}
  }
}`,
			expectedError: `expect plugin type`,
		},
		{
			name: "missing closing curly bracket (pluginsection)",
			input: `filter {
  if 1 == 1 {
    plugin{}
  }
`,
			expectedError: `expect closing curly bracket`,
		},
		{
			name: "missing closing curly bracket (plugin)",
			input: `filter {
  plugin {}
  plugin2 {
  plugin3 {}
}`,
			expectedError: `expect closing curly bracket`,
		},
		{
			name: "invalid value",
			input: `filter {
  plugin {
    value => #invalid#
  }
}`,
			expectedError: `expect closing curly bracket`,
		},
		{
			name: "invalid array value",
			input: `filter {
  plugin {
    value => [ #invalid# ]
  }
}`,
			// Upstream grammar from Logstash is wrong: https://github.com/elastic/logstash/issues/6580
			// expectedError: `invalid array value`,
			expectedError: `expect closing square bracket`,
		},
		{
			name: "missing closing double quotes",
			input: `filter {
  plugin {
    value => "invalid
  }
}`,
			expectedError: `expect closing double quotes (")`,
		},
		{
			name: "missing closing single quotes",
			input: `filter {
  plugin {
    value => 'invalid
  }
}`,
			expectedError: `expect closing single quote (')`,
		},
		{
			name: "missing closing double slash for regexp",
			input: `filter {
  if [field] =~ /.*
    plugin {}
  }
}`,
			expectedError: `expect closing slash (/) for regexp`,
		},
		{
			name: "missing closing squre bracket",
			input: `filter {
  plugin {
    value => [ "bar" 
  }
}`,
			expectedError: `expect closing square bracket`,
		},
		{
			name: "missing closing curly bracket (value)",
			input: `filter {
  plugin {
    value => { "bar" => "bar"
  }
}`,
			expectedError: `expect closing curly bracket`,
		},
		{
			name: "missing closing curly bracket (if)",
			input: `filter {
  if 1 == 1 {
    plugin {}
  else {
    plugin2 {}
  }
}`,
			expectedError: `expect closing curly bracket`,
		},
		{
			name: "missing closing curly bracket (else)",
			input: `filter {
  if 1 == 1 {
    plugin {}
  } else {
    plugin2 {}
}`,
			expectedError: `expect closing curly bracket`,
		},
		{
			name: "missing closing curly bracket (else if)",
			input: `filter {
  if 1 == 1 {
    plugin {}
  } else if 2 == 2 {
    plugin2 {}
  else {
    plugin3 {}
  }
}`,
			expectedError: `expect closing curly bracket`,
		},
		// 		{
		// 			name: "not valid if condition"
		// 			input: `filter {
		//   if this is not valid {
		//     plugin {}
		//   }
		// }`,
		// 			expectedError: `xxxxxxxxxxxxx`,
		// 		},
		{
			name: "missing closing parenthesis",
			input: `filter {
  if !( 1 == 1 {
    plugin {}
  }
}`,
			expectedError: `expect closing parenthesis`,
		},
		{
			name: "invalid value for expression",
			input: `filter {
  if [test] == #test# {
    plugin {}
  }
}`,
			expectedError: ` invalid value for expression`,
		},
		{
			name: "invalid boolean operator",
			input: `filter {
  if 1 == 1 nor 2 == 2 {
    plugin{}
  }
}`,
			expectedError: `expect boolean operator`,
		},
		{
			name: "invalid comparison operator",
			input: `filter {
  if [field] ?~ /.*/ {
    plugin{}
  }
}`,
			expectedError: `expect regexp comparison operator`,
		},
		{
			name: "invalid compare operator",
			input: `filter {
  if "test" =! "test" {
    plugin{}
  }
}`,
			expectedError: `expect compare operator`,
		},
		{
			name: "invalid in operator",
			input: `filter {
  if "test" on [field] {
    plugin{}
  }
}`,
			expectedError: `expect in operator`,
		},
		{
			name: "invalid not in operator",
			input: `filter {
  if "test" no in [field] {
    plugin{}
  }
}`,
			expectedError: `expect not in operator`,
		},
		{
			name: "missing closing square bracket",
			input: `filter {
  if "test" in [field][subfield {
    plugin{}
  }
}`,
			expectedError: `expect closing square bracket`,
		},
	}

	var errMsg string
	var hasErr bool
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := ParseReader("parse errors", strings.NewReader(test.input))
			if err == nil {
				t.Errorf("Expected parsing to fail with error: %s, input: %s", test.expectedError, test.input)
			} else {
				if errMsg, hasErr = GetFarthestFailure(); !hasErr {
					errMsg = err.Error()
				}
				if !strings.Contains(errMsg, test.expectedError) {
					t.Errorf("Expected parsing to fail with error containing: %s, got error: %s, input: %s", test.expectedError, errMsg, test.input)
				}
			}
		})
	}
}

func TestParseExceptionalComments(t *testing.T) {
	cases := []string{
		"comments_everywhere",
	}

	for _, test := range cases {
		t.Run(test, func(t *testing.T) {
			inputFilename := "testdata/exceptional_comments/" + test + ".conf"
			got, err := ParseFile(
				inputFilename,
				ExceptionalCommentsWarning(true),
			)
			if err != nil {
				t.Fatalf("Expected %q to parse without error: %v", test, err)
			}
			config, ok := got.(ast.Config)
			if !ok {
				t.Fatalf("Expected to parse to Config")
			}

			body, err := ioutil.ReadFile(inputFilename)
			if err != nil {
				t.Fatalf("Unable to read file %q: %v", inputFilename, err)
			}
			exceptionalCommentCount := strings.Count(string(body), "exceptional_comment")
			if exceptionalCommentCount != len(config.Warnings) {
				for _, line := range config.Warnings {
					t.Log("line", line)
				}
				t.Fatalf("Expected %d warnings, got %d", exceptionalCommentCount, len(config.Warnings))
			}
		})
	}
}

func TestParseIgnoreComments(t *testing.T) {
	cases := []string{
		"comments_everywhere",
	}

	for _, test := range cases {
		t.Run(test, func(t *testing.T) {
			inputFilename := "testdata/ignore_comments/" + test + ".conf"
			expectedFilename := "testdata/ignore_comments/" + test + ".expected.conf"
			got1, err := ParseFile(
				inputFilename,
				IgnoreComments(true),
			)
			if err != nil {
				t.Fatalf("Expected %q to parse without error: %v", test, err)
			}
			expected, err := ioutil.ReadFile(expectedFilename)
			if err != nil {
				t.Fatalf("Unable to read file %q: %v", expectedFilename, err)
			}
			got := fmt.Sprintf("%v", got1)
			if string(expected) != got {
				t.Errorf("Expected %q parsed input to print the same as input:\n%s", inputFilename, printDiff(string(expected), got))
			}
		})
	}
}
