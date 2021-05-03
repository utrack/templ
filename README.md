# templ

* A strongly typed HTML templating language that compiles to Go code, and has great developer tooling.

## Current state

This is alpha software, it will change a lot, and often. There's no guarantees of correctness or that APIs won't change at the moment.

The `templ compile` command works, and the `templ lsp` is in a basic "Hello World" state right now.

If you're keen to see Go be practical for Web projects, see "Help needed" for where the project needs your help.

## Features

The language compiles to Go syntax, some sections of the template (e.g. `package`, `import`, `if`, `for` and `switch` statements) are directly as Go expressions in the compiled output, while HTML elements are converted to Go code that renders their output.

The project is in alpha stage at present, but the aim is to provide a core set of functionality.

* `templ compile` compiles `*.templ` files into Go files.
* `templ lsp` provides a Language Server to support IDE integrations. The compile command generates a sourcemap which maps from the `*.templ` files to the compiled Go file. This enables the `templ` LSP to use the Go language `gopls` language server as is, providing a thin shim to do the source remapping. This is used to provide autocomplete for template variables and functions.
* `templ fmt` formats the template file by parsing it and rewriting it out.

### Help needed

The project is looking for help with:

* Writing the `fmt` tool.
* Adding support for switch/case statements based on the existing statements.
* Adding an integration test suite that, in each subdirectory, is an templ file, a `main.go` file that will render it, and an exected output file. The integration test will compile the `templ` files, run the `main.go` file and compare the expected vs the actual output.
* Adding features to the Language Server implementation, it's just at "Hello World!" stage at the moment. It needs to be able to do definition (should be easiest, because the `gopls` CLI supports it) and then autocomplete.
* Writing a VS Code plugin that uses the LSP support.
* Examples and testing of the tools.
* Adding a `hot` option to the compiler that recompiles the `*.templ` files when they change on disk. This could be achieved by documenting and making it easy to use external tools such as `ag`, ripgrep (`rg`) and `entr` in the short term.
* Writing documentation of the components.
* Writing a blog post that demonstrates using the tool to build a form-based Web application.
* Testing, it's alpha stage right now.
* An example of a web-based UI component library would be very useful, a more advanced version of the integration test suite, thatwould be a Go web server that runs the compiled `templ` file along with example JSON payloads that match the expected data structure types and renders the content - a UI playground. If it could do hot-reload, amazing.
* Low priority, but I'm thinking of developing a CSS-in-Go implementation to work in parallel. This might take the form of a pre-processor which would collect all "style" attributes of elements and automatically calculate a minimum set of CSS classes that could be created and applied to the elements - but a first pass could just be a way to define CSS classes in Go to allow the use of CSS variables.

Please get in touch if you're interested in building a feature as I don't want people to spend time on something that's already being worked on, or ends up being a waste of their time because it can't be integrated.

## Language example

```templ
{% package templ %}

{% import "strings" %}
{% import strs "strings" %}

{% templ RenderAddress(addr Address) %}
	<div>{%= addr.Address1 %}</div>
	<div>{%= addr.Address2 %}</div>
	<div>{%= addr.Address3 %}</div>
	<div>{%= addr.Address4 %}</div>
{% endtempl %}

{% templ Render(p Person) %}
   <div>
     <div>{%= p.Name() %}</div>
     <a href={%= p.URL %}>{%= strings.ToUpper(p.Name()) %}</a>
     <div>
         {% if p.Type == "test" %}
            <span>{%= "Test user" %}</span>
         {% else %}
	    <span>{%= "Not test user" %}</span>
         {% endif %}
         {% for _, v := range p.Addresses %}
            {% call RenderAddress(v) %}
         {% endfor %}
     </div>
   </div>
{% endtempl %}
```

Will compile to Go code similar to:

```go
// Code generated by templ DO NOT EDIT.

package templ

import "html"
import "context"
import "io"
import "strings"
import strs "strings"

func RenderAddress(ctx context.Context, w io.Writer, addr Address) (err error) {
	_, err = io.WriteString(w, "<div>")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, html.EscapeString(addr.Address1))
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, "</div>")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, "<div>")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, html.EscapeString(addr.Address2))
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, "</div>")
	if err != nil {
		return err
	}
	// You get the idea...
	return err
}

func Render(ctx context.Context, w io.Writer, p Person) (err error) {
	_, err = io.WriteString(w, "<div>")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, "<div>")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, html.EscapeString(p.Name()))
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, "</div>")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, "<a")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, " href=")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, "\"")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, html.EscapeString(p.URL))
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, "\"")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, ">")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, html.EscapeString(strings.ToUpper(p.Name())))
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, "</a>")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, "<div>")
	if err != nil {
		return err
	}
	if p.Type == "test" {
		_, err = io.WriteString(w, "<span>")
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, html.EscapeString("Test user"))
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, "</span>")
		if err != nil {
			return err
		}
	} else {
		_, err = io.WriteString(w, "<span>")
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, html.EscapeString("Not test user"))
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, "</span>")
		if err != nil {
			return err
		}
	}
	for _, v := range p.Addresses {
		err = Address(ctx, w, v)
		if err != nil {
			return err
		}
	}
	_, err = io.WriteString(w, "</div>")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, "</div>")
	if err != nil {
		return err
	}
	return err
}
```

# IDE Support

## vim / neovim

A vim / neovim plugin is available from https://github.com/Joe-Davidson1802/templ.vim which adds syntax highlighting.

Neovim 5 supports Language Servers directly. For the moment, I'm using https://github.com/neoclide/coc.nvim to test the language server after using Joe-Davidson1802's plugin to set the language type:

```json
{
  "languageserver": {
    "templ": {
      "command": "templ",
      "args": ["lsp"],
      "filetypes": ["templ"]
    }
}
```

To add extensive debug information, you can include additional args to the LSP, like this:

```json
{
  "languageserver": {
    "templ": {
      "command": "templ",
      "args": ["lsp",
        "--log", "/Users/adrian/github.com/a-h/templ/cmd/lsp/templ-log.txt", 
	"--goplsLog", "/Users/adrian/github.com/a-h/templ/cmd/lsp/gopls-log.txt",
	"--goplsRPCTrace", "true"
      ],
      "filetypes": ["templ"]
    }
}
```

## vscode

Yes please, talk to me about it!

# Development

## Local builds

To build a local version you can use the `go build` tool:

```
cd cmd
go build -o templ
```

## Testing

Unit tests use the `go test` tool:

```
go test ./...
```

## Release testing

This project uses https://github.com/goreleaser/goreleaser to build the command line binary and deploy it to Github. You will need to install this to test releases.

```
make build-snapshot
```

The binaries are created by me and signed by my GPG key. You can verify with my key https://adrianhesketh.com/a-h.gpg

# Inspiration

Doesn't this look like a lot like https://github.com/valyala/quicktemplate ?

Yes, yes it does. I looked at the landscape of Go templating languages before I started writing code and my initial plan was to improve the IDE support of quicktemplate, see https://github.com/valyala/quicktemplate/issues/80

The package author didn't respond (hey, we're all busy), and looking through the code, I realised that it would be hard to modify what's there to have the concept of source mapping, mostly because there's no internal object model of the language, it reads and emits code in one go.

It's also a really feature rich project, with all sorts of formatters, and support for various languages (JSON etc.), so I borrowed some syntax ideas, but left the code. If `valyala` is up for it, I'd be happy to help integrate the ideas from here. I just want Go to have a templating language with great IDE support.
