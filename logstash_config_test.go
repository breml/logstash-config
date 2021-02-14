package config_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	. "github.com/breml/logstash-config"
	"github.com/breml/logstash-config/ast"
)

func ExampleParseReader() {
	logstashConfig := `filter {
    mutate {
      add_tag => [ "tag" ]
    }
  }`
	got, err := ParseReader("example.conf", strings.NewReader(logstashConfig))
	if err != nil {
		log.Fatalf("Parse error: %s\n", err)
	}

	// Output: filter {
	//   mutate {
	//     add_tag => [ "tag" ]
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
			name: "Plugin with all attribte types",
			input: `input {
  stdin {
    doublequotedstring => "doublequotedstring with escaped \" "
    singlequotedstring => 'singlequotedstring with escaped \' '
    "doublequotedkey" => value
    'singlequotedkey' => value
    bareword => bareword
    intnumber => 3
    floatnumber => 3.1415
    arrayvalue => [ bareword, "doublequotedstring", 'singlequotedstring', 3, 3.1415 ]
    hashvalue => {
      doublequotedstring => "doublequotedstring"
      singlequotedstring => 'singlequotedstring'
      bareword => bareword
      intnumber => 3
      arrayvalue => [ bareword, "doublequotedstring", 'singlequotedstring', 3, 3.1415 ]
      subhashvalue => {
        subhashvaluestring => value
      }
    }
    codec => rubydebug {
      string => "a string"
    }
  }
}
`,
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
  if ! ("true" == "true") {
    plugin {}
  }
}
`,
		},
		{
			name: "Negative Selector for value in subfield",
			input: `filter {
  if ! [field][subfield] {
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
    value => [  ]
  }
}
`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseReader("test", strings.NewReader(test.input))
			if err != nil {
				t.Fatalf("Expected to parse without error: %s, input:\n|%s|", err, test.input)
			}
			if test.input != fmt.Sprintf("%v", got) {
				t.Errorf("Expected parsed input to print the same as input, input:\n|%s|\n\nparsed:\n|%v|", test.input, got)
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
  stdin {}
}
`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseReader("test", strings.NewReader(test.input))
			if err != nil {
				t.Fatalf("Expected to parse without error: %s, input:\n|%s|", err, test.input)
			}
			if test.expected != fmt.Sprintf("%v", got) {
				t.Errorf("Expected output does not match parsed output, expected:\n|%s|\n\nparsed:\n|%v|", test.expected, got)
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
			name: "missing closing curly backet (pluginsection)",
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
  if ! ( 1 == 1 {
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
			name: "missing closing squre bracket",
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
	cases := []struct {
		name  string
		input string
	}{
		{
			name: "comments everywhere",
			input: `# comment
filter # exceptional_comment
{
  # comment
  mutate # exceptional_comment
  {
    # comment
    id # exceptional_comment
    => # exceptional_comment
    comment
    # comment 1
    add_tag # exceptional_comment
    => # exceptional_comment
    [
      # comment
      "value",
      # comment
      "value2"
      # comment 2
    ]
    # comment
    add_field # exceptional_comment
    => # exceptional_comment
    {
      # comment
      "key1" => "value"
      "key2" => "value2"
      # comment
    }
    # comment
  }
  # comment
  if # exceptional_comment
  ( # exceptional_comment
  "true" # exceptional_comment
  == # exceptional_comment
  "false"
  # exceptional_comment
  ) # exceptional_comment
  and # exceptional_comment
  ! # exceptional_comment
  [foobar] # exceptional_comment
  or # exceptional_comment
  "foo" # exceptional_comment
  not # exceptional_comment
  in # exceptional_comment
  [bar]
  { # comment
    mutate {}
    # comment
  } # comment
  else # exceptional_comment
  if # exceptional_comment
  [field] # exceptional_comment
  =~ # exceptional_comment
  /regex/ # exceptional_comment
  and # exceptional_comment
  ! # exceptional_comment
  ( # exceptional_comment
  1 # exceptional_comment
  < # exceptional_comment
  2 # exceptional_comment
  ) # exceptional_comment
  { # comment
    mutate {}
    # comment
  } # comment
  else # exceptional_comment
  { # comment
    mutate {}
    # comment
  } # comment
}
# comment
`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseReader(
				"test",
				strings.NewReader(test.input),
				ExceptionalCommentsWarning(true),
			)
			if err != nil {
				t.Fatalf("Expected to parse without error: %s, input:\n|%s|", err, test.input)
			}
			config, ok := got.(ast.Config)
			if !ok {
				t.Fatalf("Expected to parse to Config")
			}
			if strings.Count(test.input, "exceptional_comment") != len(config.Warnings) {
				for _, line := range config.Warnings {
					t.Log("line", line)
				}
				t.Fatalf("Expected %d warnings, got %d", strings.Count(test.input, "exceptional_comment"), len(config.Warnings))
			}
		})
	}
}
