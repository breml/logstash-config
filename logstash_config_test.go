package config

import (
	"fmt"
	"strings"
	"testing"
)

func TestParserIdentic(t *testing.T) {
	cases := []struct {
		input string
	}{
		// Empty file => does not work
		// {
		// 	input: ``,
		// },

		// Single PluginSection
		{
			input: `input {
  
}
`,
		},

		// All PluginSections empty
		{
			input: `input {
  
}
filter {
  
}
output {
  
}
`,
		},

		// Plugin without attributes
		{
			input: `input {
  stdin {
    
  }
}
`,
		},

		// Multiple plugins
		{
			input: `input {
  stdin {
    
  }
  file {
    
  }
}
filter {
  mutate {
    
  }
  mutate {
    
  }
  mutate {
    
  }
}
output {
  stdout {
    
  }
}
`,
		},

		// Plugin with all attribute types
		{
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
	}

	for _, test := range cases {
		got, err := ParseReader("test", strings.NewReader(test.input))
		if err != nil {
			t.Fatalf("Expected to parse without error: %s, input:\n|%s|", err, test.input)
		}
		if test.input != fmt.Sprintf("%v", got) {
			t.Errorf("Expected parsed input to print the same as input, input:\n|%s|\n\nparsed:\n|%v|", test.input, got)
		}
	}
}

func TestParser(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input: `input { stdin {} }`,
			expected: `input {
  stdin {
    
  }
}
`,
		},
		{
			input: `input { 
  # Comment
  stdin {
    # Comment
  }
}`,
			expected: `input {
  stdin {
    
  }
}
`,
		},
	}

	for _, test := range cases {
		got, err := ParseReader("test", strings.NewReader(test.input))
		if err != nil {
			t.Fatalf("Expected to parse without error, input:\n|%s|", test.input)
		}
		if test.expected != fmt.Sprintf("%v", got) {
			t.Errorf("Expected output does not match parsed output, expected:\n|%s|\n\nparsed:\n|%v|", test.expected, got)
		}
	}
}
