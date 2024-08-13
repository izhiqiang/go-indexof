package view

import (
	"html/template"
	"indexof/config"
	"net/http"
)

func FetchIndexOf(w http.ResponseWriter, data any) error {
	name := "indexof"
	w.WriteHeader(http.StatusOK)
	if config.Global.Debug {
		fetchTemplates(w, name, data)
		return nil
	}
	w.Header().Set("Content-Type", "text/html")
	return template.Must(template.New(name).Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <meta name="viewport" content="width=device-width">
    <title>Index of {{.IndexOf}}</title>
    <style>
        body, html {
            background: #fff;
            font-family: "Bitstream Vera Sans", "Lucida Grande", "Lucida Sans Unicode", Lucidux, Verdana, Lucida, sans-serif;
        }
        tr:nth-child(even) {
            background: #f4f4f4;
        }

        th, td {
            padding: 0.1em 0.5em;
        }

        th {
            text-align: left;
            font-weight: bold;
            background: #eee;
            border-bottom: 1px solid #aaa;
        }

        #list {
            border: 1px solid #aaa;
            width: 100%;
        }

        a {
            color: #a33;
        }

        a:hover {
            color: #e33;
        }
    </style>
</head>
<body>
<h1>Index of {{.IndexOf}}</h1>

<table id="list">
    <thead>
    <tr>
        <th style="width:55%">File Name</th>
        <th style="width:20%">File Size</th>
        <th style="width:25%">Date</th>
    </tr>
    </thead>
    <tbody>
    <tbody>
    {{range .PathInfos}}
    <tr>
        <td class="link"><a href="{{.Name}}{{if .IsDir}}/{{end}}" title="{{.Name}}">{{.Name}}</a></td>
        <td class="size">
            {{if .IsDir}}-{{else}}{{.Size}}{{end}}
        </td>
        <td class="date">{{.Data.Format "2006-01-02 15:04:05"}}</td>
    </tr>
    {{end}}
    </tbody>
</table>
<footer>
    <p align=center>
        <a href="https://go.dev/" target="_blank">{{.GoVersion}}</a>
        - <a href="https://github.com/izhiqiang/go-indexof" target="_blank">开源地址</a>
        - <a href="https://www.zhiqiang.wang/" target="_blank">关于作者</a>
    </p>
</footer>
</body>
</html>
`)).Execute(w, data)
}
