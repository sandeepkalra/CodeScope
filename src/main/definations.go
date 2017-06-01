package main

type Symbols map[string][]Occourance

type Occourance struct {
	OccouranceCountInThisLine int
	FileName                  string
	LineNumber                int
	LineContent               string
}

/* ------- json --- cfg -------------- */

type Language struct {
	Extensions []string `json:"extensions"`
	Keywords   []string `json:"keywords"`
	Whitespace []string `json:"whitespace"`
	Operators  []string `json:"operators"`
	IgnoreDirs []string `json:"ignore_dirs"`
}

type JsonPkg struct {
	Languages []Language `json:"languages"`
}
