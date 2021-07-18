var AllTableNames = []string {
	{{range $table := .Tables -}}
    "{{$table.Name}}",
	{{end -}}
}
