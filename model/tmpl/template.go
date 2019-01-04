package tmpl

import (
	"log"
	"strings"
	"text/template"
)

var rootTemplate *template.Template

type TemplateParam struct {
	Table  string
	Suffix string
	Equals []string
}

func init() {
	var err error
	rootTemplate = template.New("root")
	rootTemplate.Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"indent": func(num int, name string, data interface{}) (string, error) {
			sb := &strings.Builder{}
			err := rootTemplate.ExecuteTemplate(sb, name, data)
			if err != nil {
				return "", err
			}
			strs := strings.Split(sb.String(), "\n")
			sb.Reset()
			for n, str := range strs[:len(strs)-1] {
				for i := 0; i < num; i++ {
					sb.WriteRune(' ')
				}
				sb.WriteString(str)
				if n+1 < len(strs)-1 {
					sb.WriteRune('\n')
				}
			}
			return sb.String(), nil
		},
	})

	_, err = rootTemplate.New("select").Parse(`
{{- $root := . -}}
select * from {{$root.Table}}{{$root.Suffix}} where
{{- range $index, $field := $root.Equals }}
  {{$root.Table}}_{{$field}} = :{{$root.Table}}_{{$field}}{{if lt (add $index 1) (len $root.Equals)}} and{{end}}
{{- end }}
order by {{$root.Table}}_id desc
`)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = rootTemplate.New("selectInsert").Parse(`
{{- $root := . -}}
select * from {{$root.Table}}{{$root.Suffix}} where
{{- range $index, $field := $root.Equals }}
  {{$root.Table}}_{{$field}} = :{{$root.Table}}_{{$field}} and
{{- end }}
  {{$root.Table}}_deleted is null
`)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = rootTemplate.New("insert").Parse(`
{{- $root := . -}}
with cte as (
  insert into {{ $root.Table }}{{$root.Suffix}}( 
{{- range $index, $field := $root.Equals }}
    {{$root.Table}}_{{$field}}{{if lt (add $index 1) (len $root.Equals)}},{{end}}
{{- end  }}
  ) select
{{- range $index, $field := $root.Equals }}
      :{{$root.Table}}_{{$field}}{{if lt (add $index 1) (len $root.Equals)}},{{end}}
{{- end}}
    where not exists (
{{indent 6 "selectInsert" $root}}
    )
  returning *
) select * from cte
  union
{{indent 2 "selectInsert" $root}}
`)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = rootTemplate.New("delete").Parse(`
{{- $root := . -}}
update {{$root.Table}}{{$root.Suffix}} set
  {{$root.Table}}_deleted = now()
where
  {{$root.Table}}_id in (
    select {{$root.Table}}_id from {{$root.Table}}_last where
{{- range $index, $field := $root.Equals }}
      {{$root.Table}}_{{$field}} = :{{$root.Table}}_{{$field}}{{if lt (add $index 1) (len $root.Equals)}} and{{end}}
{{- end }}
  )
`)
	if err != nil {
		log.Fatalln(err)
	}
}

func SelectTemplate(table, suffix string, equals ...string) string {
	sb := &strings.Builder{}
	err := rootTemplate.ExecuteTemplate(sb, "select", &TemplateParam{
		Table:  table,
		Suffix: suffix,
		Equals: equals,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return sb.String()
}

func InsertTemplate(table, suffix string, equals ...string) string {
	sb := &strings.Builder{}
	err := rootTemplate.ExecuteTemplate(sb, "insert", &TemplateParam{
		Table:  table,
		Suffix: suffix,
		Equals: equals,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return sb.String()
}

func DeleteTemplate(table, suffix string, equals ...string) string {
	sb := &strings.Builder{}
	err := rootTemplate.ExecuteTemplate(sb, "delete", &TemplateParam{
		Table:  table,
		Suffix: suffix,
		Equals: equals,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return sb.String()
}
