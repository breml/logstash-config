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
		// Simple if (without else) branch
		{
			input: `filter {
  if 1 == 1 {
    date {
      
    }
  }
}
`,
		},
		// if with else-if and a final else branch
		{
			input: `filter {
  if 1 == 1 {
    date {
      
    }
  } else if 1 == 1 {
    date {
      
    }
  } else {
    date {
      
    }
  }
}
`,
		},
		// if with multiple else-if and a final else branch
		// multiple plugins in each branch
		{
			input: `filter {
  if 1 == 1 {
    date {
      
    }
    date {
      
    }
  } else if 1 == 1 {
    date {
      
    }
    date {
      
    }
    date {
      
    }
  } else if 1 == 1 {
    date {
      
    }
    date {
      
    }
    date {
      
    }
  } else {
    date {
      
    }
    date {
      
    }
    date {
      
    }
  }
}
`,
		},
		// if with multiple else-if and a final else branch
		// test for different condition types
		{
			input: `filter {
  if 1 != 1 {
    date {
      
    }
  } else if 1 <= 1 {
    date {
      
    }
  } else if 1 >= 1 {
    date {
      
    }
  } else if 1 < 1 {
    date {
      
    }
  } else if 1 > 1 {
    date {
      
    }
  }
}
`,
		},
		// if with multiple compare operators
		{
			input: `filter {
  if "true" == "true" and "true1" == "true1" or "true2" == "true2" nand "true3" == "true3" xor "true4" == "true4" {
    plugin {
      
    }
  }
}
`,
		},
		// Condition in parentheses
		{
			input: `filter {
  if ("tag" in [tags]) {
    plugin {
      
    }
  }
}
`,
		},
		// Multiple conditions in parentheses
		{
			input: `filter {
  if ("tag" in [tags] or ("true" == "true" and 1 == 1)) {
    plugin {
      
    }
  }
}
`,
		},
		// Negative condition
		{
			input: `filter {
  if ! ("true" == "true") {
    plugin {
      
    }
  }
}
`,
		},
		// Negative Selector for value in subfield
		{
			input: `filter {
  if ! [field][subfield] {
    plugin {
      
    }
  }
}
`,
		},
		// InExpression
		{
			input: `filter {
  if "tag" in [tags] {
    plugin {
      
    }
  }
}
`,
		},
		// NotInExpression
		{
			input: `filter {
  if "tag" not in [field][subfield] {
    plugin {
      
    }
  }
}
`,
		},
		// RegexpExpression (Match)
		{
			input: `filter {
  if [field] =~ /.*/ {
    plugin {
      
    }
  }
}
`,
		},
		// RegexpExpression (Not Match)
		{
			input: `filter {
  if [field] !~ /.*/ {
    plugin {
      
    }
  }
}
`,
		},
		// Rvalue
		{
			input: `filter {
  if "string" or 10 or [field][subfield] or /.*/ {
    plugin {
      
    }
  }
}
`,
		},
		//
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
